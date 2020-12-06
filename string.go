package mp4

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"

	"github.com/abema/go-mp4/util"
)

type stringifier struct {
	buf    *bytes.Buffer
	src    IImmutableBox
	indent string
	ctx    Context
}

func Stringify(src IImmutableBox, ctx Context) (string, error) {
	return StringifyWithIndent(src, "", ctx)
}

func StringifyWithIndent(src IImmutableBox, indent string, ctx Context) (string, error) {
	boxDef := src.GetType().getBoxDef(ctx)
	if boxDef == nil {
		return "", ErrBoxInfoNotFound
	}

	v := reflect.ValueOf(src).Elem()

	m := &stringifier{
		buf:    bytes.NewBuffer(nil),
		src:    src,
		indent: indent,
		ctx:    ctx,
	}

	err := m.stringifyStruct(boxDef.fields, v, 0, true)
	if err != nil {
		return "", err
	}

	return m.buf.String(), nil
}

func (m *stringifier) stringify(f *field, v reflect.Value, config fieldConfig, depth int) error {
	switch v.Type().Kind() {
	case reflect.Ptr:
		return m.stringifyPtr(f, v, config, depth)
	case reflect.Struct:
		return m.stringifyStruct(f.children, v, depth, config.extend)
	case reflect.Array:
		return m.stringifyArray(f, v, config, depth)
	case reflect.Slice:
		return m.stringifySlice(f, v, config, depth)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return m.stringifyInt(v, config, depth)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return m.stringifyUint(v, config, depth)
	case reflect.Bool:
		return m.stringifyBool(v, config, depth)
	case reflect.String:
		return m.stringifyString(v, config, depth)
	default:
		return fmt.Errorf("unsupported type: %s", v.Type().Kind())
	}
}

func (m *stringifier) stringifyPtr(f *field, v reflect.Value, config fieldConfig, depth int) error {
	return m.stringify(f, v.Elem(), config, depth)
}

func (m *stringifier) stringifyStruct(fs []*field, v reflect.Value, depth int, extended bool) error {
	if !extended {
		m.buf.WriteString("{")
		if m.indent != "" {
			m.buf.WriteString("\n")
		}
		depth++
	}

	for _, f := range fs {
		sf, _ := v.Type().FieldByName(f.name)
		tagStr, _ := sf.Tag.Lookup("mp4")
		config, err := readFieldConfig(m.src, v, f.name, parseFieldTag(tagStr), m.ctx)
		if err != nil {
			return err
		}

		if !isTargetField(m.src, config, m.ctx) {
			continue
		}

		if config.cnst != "" || config.hidden {
			continue
		}

		if !config.extend {
			if m.indent != "" {
				writeIndent(m.buf, m.indent, depth+1)
			} else if m.buf.Len() != 0 && m.buf.Bytes()[m.buf.Len()-1] != '{' {
				m.buf.WriteString(" ")
			}
			m.buf.WriteString(f.name)
			m.buf.WriteString("=")
		}

		str, ok := config.cfo.StringifyField(f.name, m.indent, depth+1, m.ctx)
		if ok {
			m.buf.WriteString(str)
			if !config.extend && m.indent != "" {
				m.buf.WriteString("\n")
			}
			continue
		}

		if f.name == "Version" {
			m.buf.WriteString(strconv.Itoa(int(m.src.GetVersion())))
		} else if f.name == "Flags" {
			fmt.Fprintf(m.buf, "0x%06x", m.src.GetFlags())
		} else {
			err = m.stringify(f, v.FieldByName(f.name), config, depth)
			if err != nil {
				return err
			}
		}

		if !config.extend && m.indent != "" {
			m.buf.WriteString("\n")
		}
	}

	if !extended {
		if m.indent != "" {
			writeIndent(m.buf, m.indent, depth)
		}
		m.buf.WriteString("}")
	}

	return nil
}

func (m *stringifier) stringifyArray(f *field, v reflect.Value, config fieldConfig, depth int) error {
	begin, sep, end := "[", ", ", "]"
	if config.str || config.iso639_2 {
		begin, sep, end = "\"", "", "\""
	} else if config.uuid {
		begin, sep, end = "", "", ""
	}

	m.buf.WriteString(begin)

	m2 := *m
	if config.str {
		m2.buf = bytes.NewBuffer(nil)
	}
	size := v.Type().Size()
	for i := 0; i < int(size)/int(v.Type().Elem().Size()); i++ {
		if i != 0 {
			m2.buf.WriteString(sep)
		}

		if err := m2.stringify(f, v.Index(i), config, depth+1); err != nil {
			return err
		}

		if config.uuid && (i == 3 || i == 5 || i == 7 || i == 9) {
			m.buf.WriteString("-")
		}
	}
	if config.str {
		m.buf.WriteString(util.EscapeUnprintables(m2.buf.String()))
	}

	m.buf.WriteString(end)

	return nil
}

func (m *stringifier) stringifySlice(f *field, v reflect.Value, config fieldConfig, depth int) error {
	begin, sep, end := "[", ", ", "]"
	if config.str || config.iso639_2 {
		begin, sep, end = "\"", "", "\""
	}

	m.buf.WriteString(begin)

	m2 := *m
	if config.str {
		m2.buf = bytes.NewBuffer(nil)
	}
	for i := 0; i < v.Len(); i++ {
		if config.length != LengthUnlimited && uint(i) >= config.length {
			break
		}

		if i != 0 {
			m2.buf.WriteString(sep)
		}

		if err := m2.stringify(f, v.Index(i), config, depth+1); err != nil {
			return err
		}
	}
	if config.str {
		m.buf.WriteString(util.EscapeUnprintables(m2.buf.String()))
	}

	m.buf.WriteString(end)

	return nil
}

func (m *stringifier) stringifyInt(v reflect.Value, config fieldConfig, depth int) error {
	if config.hex {
		val := v.Int()
		if val >= 0 {
			m.buf.WriteString("0x")
			m.buf.WriteString(strconv.FormatInt(val, 16))
		} else {
			m.buf.WriteString("-0x")
			m.buf.WriteString(strconv.FormatInt(-val, 16))
		}
	} else {
		m.buf.WriteString(strconv.FormatInt(v.Int(), 10))
	}
	return nil
}

func (m *stringifier) stringifyUint(v reflect.Value, config fieldConfig, depth int) error {
	if config.iso639_2 {
		m.buf.WriteString(string([]byte{byte(v.Uint() + 0x60)}))
	} else if config.uuid {
		fmt.Fprintf(m.buf, "%02x", v.Uint())
	} else if config.str {
		m.buf.WriteString(string([]byte{byte(v.Uint())}))
	} else if config.hex || (!config.dec && v.Type().Kind() == reflect.Uint8) || v.Type().Kind() == reflect.Uintptr {
		m.buf.WriteString("0x")
		m.buf.WriteString(strconv.FormatUint(v.Uint(), 16))
	} else {
		m.buf.WriteString(strconv.FormatUint(v.Uint(), 10))
	}

	return nil
}

func (m *stringifier) stringifyBool(v reflect.Value, config fieldConfig, depth int) error {
	m.buf.WriteString(strconv.FormatBool(v.Bool()))

	return nil
}

func (m *stringifier) stringifyString(v reflect.Value, config fieldConfig, depth int) error {
	m.buf.WriteString("\"")
	m.buf.WriteString(util.EscapeUnprintables(v.String()))
	m.buf.WriteString("\"")

	return nil
}

func writeIndent(w io.Writer, indent string, depth int) {
	for i := 0; i < depth; i++ {
		io.WriteString(w, indent)
	}
}
