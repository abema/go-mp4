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
	t := reflect.TypeOf(src).Elem()
	v := reflect.ValueOf(src).Elem()

	m := &stringifier{
		buf:    bytes.NewBuffer(nil),
		src:    src,
		indent: indent,
		ctx:    ctx,
	}

	err := m.stringifyStruct(t, v, 0, true)
	if err != nil {
		return "", err
	}

	return m.buf.String(), nil
}

func (m *stringifier) stringify(t reflect.Type, v reflect.Value, config fieldConfig, depth int) error {
	switch t.Kind() {
	case reflect.Ptr:
		return m.stringifyPtr(t, v, config, depth)
	case reflect.Struct:
		return m.stringifyStruct(t, v, depth, config.extend)
	case reflect.Array:
		return m.stringifyArray(t, v, config, depth)
	case reflect.Slice:
		return m.stringifySlice(t, v, config, depth)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return m.stringifyInt(t, v, config, depth)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return m.stringifyUint(t, v, config, depth)
	case reflect.Bool:
		return m.stringifyBool(t, v, config, depth)
	case reflect.String:
		return m.stringifyString(t, v, config, depth)
	default:
		return fmt.Errorf("unsupported type: %s", t.Kind())
	}
}

func (m *stringifier) stringifyPtr(t reflect.Type, v reflect.Value, config fieldConfig, depth int) error {
	return m.stringify(t.Elem(), v.Elem(), config, depth)
}

func (m *stringifier) stringifyStruct(t reflect.Type, v reflect.Value, depth int, extended bool) error {
	if !extended {
		m.buf.WriteString("{")
		if m.indent != "" {
			m.buf.WriteString("\n")
		}
		depth++
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		ft := f.Type
		fv := v.Field(i)

		tagStr, ok := f.Tag.Lookup("mp4")
		if !ok {
			continue
		}
		config, err := readFieldConfig(m.src, v, f.Name, parseFieldTag(tagStr), m.ctx)
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
			m.buf.WriteString(f.Name)
			m.buf.WriteString("=")
		}

		str, ok := config.cfo.StringifyField(f.Name, m.indent, depth+1, m.ctx)
		if ok {
			m.buf.WriteString(str)
			if !config.extend && m.indent != "" {
				m.buf.WriteString("\n")
			}
			continue
		}

		if f.Name == "Version" {
			m.buf.WriteString(strconv.Itoa(int(m.src.GetVersion())))
		} else if f.Name == "Flags" {
			fmt.Fprintf(m.buf, "0x%06x", m.src.GetFlags())
		} else {
			err = m.stringify(ft, fv, config, depth)
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

func (m *stringifier) stringifyArray(t reflect.Type, v reflect.Value, config fieldConfig, depth int) error {
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
	size := t.Size()
	for i := 0; i < int(size)/int(t.Elem().Size()); i++ {
		if i != 0 {
			m2.buf.WriteString(sep)
		}

		if err := m2.stringify(t.Elem(), v.Index(i), config, depth+1); err != nil {
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

func (m *stringifier) stringifySlice(t reflect.Type, v reflect.Value, config fieldConfig, depth int) error {
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

		if err := m2.stringify(t.Elem(), v.Index(i), config, depth+1); err != nil {
			return err
		}
	}
	if config.str {
		m.buf.WriteString(util.EscapeUnprintables(m2.buf.String()))
	}

	m.buf.WriteString(end)

	return nil
}

func (m *stringifier) stringifyInt(t reflect.Type, v reflect.Value, config fieldConfig, depth int) error {
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

func (m *stringifier) stringifyUint(t reflect.Type, v reflect.Value, config fieldConfig, depth int) error {
	if config.iso639_2 {
		m.buf.WriteString(string([]byte{byte(v.Uint() + 0x60)}))
	} else if config.uuid {
		fmt.Fprintf(m.buf, "%02x", v.Uint())
	} else if config.str {
		m.buf.WriteString(string([]byte{byte(v.Uint())}))
	} else if config.hex || (!config.dec && t.Kind() == reflect.Uint8) || t.Kind() == reflect.Uintptr {
		m.buf.WriteString("0x")
		m.buf.WriteString(strconv.FormatUint(v.Uint(), 16))
	} else {
		m.buf.WriteString(strconv.FormatUint(v.Uint(), 10))
	}

	return nil
}

func (m *stringifier) stringifyBool(t reflect.Type, v reflect.Value, config fieldConfig, depth int) error {
	m.buf.WriteString(strconv.FormatBool(v.Bool()))

	return nil
}

func (m *stringifier) stringifyString(t reflect.Type, v reflect.Value, config fieldConfig, depth int) error {
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
