package mp4

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"strconv"
	"strings"
)

const (
	anyVersion      = math.MaxUint8
	lengthUnlimited = math.MaxUint32
)

var ErrUnsupportedBoxVersion = errors.New("unsupported box version")

type marshaller struct {
	writer io.Writer
	octet  byte
	width  uint
	wbytes uint64
	src    IImmutableBox
}

func Marshal(w io.Writer, src IImmutableBox) (n uint64, err error) {
	t := reflect.TypeOf(src).Elem()
	v := reflect.ValueOf(src).Elem()

	m := &marshaller{
		writer: w,
		src:    src,
	}

	if err := m.marshalStruct(t, v); err != nil {
		return 0, err
	}

	if m.width != 0 {
		return 0, fmt.Errorf("box size is not multiple of 8 bits: type=%s, width=%d", src.GetType().String(), m.width)
	}

	return m.wbytes, nil
}

func (m *marshaller) marshal(t reflect.Type, v reflect.Value, config fieldConfig) error {
	switch t.Kind() {
	case reflect.Ptr:
		return m.marshalPtr(t, v, config)
	case reflect.Struct:
		return m.marshalStruct(t, v)
	case reflect.Array:
		return m.marshalArray(t, v, config)
	case reflect.Slice:
		return m.marshalSlice(t, v, config)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return m.marshalInt(t, v, config)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return m.marshalUint(t, v, config)
	case reflect.Bool:
		return m.marshalBool(t, v, config)
	case reflect.String:
		return m.marshalString(t, v, config)
	default:
		return fmt.Errorf("unsupported type: %s", t.Kind())
	}
}

func (m *marshaller) marshalPtr(t reflect.Type, v reflect.Value, config fieldConfig) error {
	return m.marshal(t.Elem(), v.Elem(), config)
}

func (m *marshaller) marshalStruct(t reflect.Type, v reflect.Value) error {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		ft := f.Type
		fv := v.Field(i)

		tagStr, ok := f.Tag.Lookup("mp4")
		if !ok {
			continue
		}
		config, err := readFieldConfig(m.src, v, f.Name, parseFieldTag(tagStr))
		if err != nil {
			return err
		}

		if !isTargetField(m.src, config) {
			continue
		}

		err = m.marshal(ft, fv, config)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *marshaller) marshalArray(t reflect.Type, v reflect.Value, config fieldConfig) error {
	size := t.Size()
	for i := 0; i < int(size)/int(t.Elem().Size()); i++ {
		var err error
		err = m.marshal(t.Elem(), v.Index(i), config)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *marshaller) marshalSlice(t reflect.Type, v reflect.Value, config fieldConfig) error {
	length := uint64(v.Len())
	if config.Len != lengthUnlimited {
		if length < uint64(config.Len) {
			return fmt.Errorf("the slice has too few elements: required=%d actual=%d", config.Len, length)
		}
		length = uint64(config.Len)
	}

	elemType := t.Elem()
	if elemType.Kind() == reflect.Uint8 && config.Size == 8 && m.width == 0 {
		io.CopyN(m.writer, bytes.NewBuffer(v.Bytes()), int64(length))
		m.wbytes += length
		return nil
	}

	for i := 0; i < int(length); i++ {
		m.marshal(t.Elem(), v.Index(i), config)
	}
	return nil
}

func (m *marshaller) marshalInt(t reflect.Type, v reflect.Value, config fieldConfig) error {
	signed := v.Int()

	if config.Varint {
		return errors.New("signed varint is unsupported")
	}

	signBit := signed < 0
	var val uint64
	if signBit {
		val = uint64(-signed)
	} else {
		val = uint64(signed)
	}

	for i := uint(0); i < config.Size; i += 8 {
		v := val
		size := uint(8)
		if config.Size > i+8 {
			v = v >> (config.Size - (i + 8))
		} else if config.Size < i+8 {
			v = v << ((i + 8) - config.Size)
			size = config.Size - i
		}

		// set sign bit
		if i == 0 {
			if signBit {
				v |= 0x80
			} else {
				v &= 0x7f
			}
		}

		m.write(byte(v), size)
	}

	return nil
}

func (m *marshaller) marshalUint(t reflect.Type, v reflect.Value, config fieldConfig) error {
	val := v.Uint()

	if config.Varint {
		m.writeUvarint(val)
		return nil
	}

	for i := uint(0); i < config.Size; i += 8 {
		v := val
		size := uint(8)
		if config.Size > i+8 {
			v = v >> (config.Size - (i + 8))
		} else if config.Size < i+8 {
			v = v << ((i + 8) - config.Size)
			size = config.Size - i
		}
		m.write(byte(v), size)
	}

	return nil
}

func (m *marshaller) marshalBool(t reflect.Type, v reflect.Value, config fieldConfig) error {
	var val byte
	if v.Bool() {
		val = 0x80
	} else {
		val = 0x00
	}
	m.write(val, config.Size)

	return nil
}

func (m *marshaller) marshalString(t reflect.Type, v reflect.Value, config fieldConfig) error {
	data := []byte(v.String())
	for _, b := range data {
		m.write(b, 8)
	}
	m.write(0x00, 8) // null character

	return nil
}

func (m *marshaller) writeUvarint(u uint64) {
	for i := 63; i >= 0; i -= 7 {
		if u>>uint(i) != 0 {
			data := byte(u>>uint(i)) & 0x7f
			if i != 0 {
				data |= 0x80
			}
			m.write(data, 8)
		}
	}
}

func (m *marshaller) write(data byte, size uint) {
	for i := uint(0); i < size; i++ {
		b := (data >> (7 - i)) & 0x01

		m.octet |= b << (7 - m.width)
		m.width++

		if m.width == 8 {
			m.writer.Write([]byte{m.octet})
			m.octet = 0x00
			m.width = 0
			m.wbytes++
		}
	}
}

type unmarshaller struct {
	reader io.ReadSeeker
	dst    IBox
	octet  byte
	width  uint
	size   uint64
	rbytes uint64
}

func UnmarshalAny(r io.ReadSeeker, boxType BoxType, payloadSize uint64) (box IBox, n uint64, err error) {
	if dst, err := boxType.New(); err != nil {
		return nil, 0, err
	} else {
		n, err := Unmarshal(r, payloadSize, dst)
		return dst, n, err
	}
}

func Unmarshal(r io.ReadSeeker, payloadSize uint64, dst IBox) (n uint64, err error) {
	t := reflect.TypeOf(dst).Elem()
	v := reflect.ValueOf(dst).Elem()

	dst.SetVersion(anyVersion)

	u := &unmarshaller{
		reader: r,
		dst:    dst,
		size:   payloadSize,
	}

	if n, override, err := dst.BeforeUnmarshal(r); err != nil {
		return 0, err
	} else if override {
		return n, nil
	} else {
		u.rbytes = n
	}

	if err := u.unmarshalStruct(t, v); err != nil {
		if err == ErrUnsupportedBoxVersion {
			r.Seek(-int64(u.rbytes), io.SeekCurrent)
		}
		return 0, err
	}

	if u.width != 0 {
		return 0, fmt.Errorf("box size is not multiple of 8 bits: type=%s, size=%d, readBytes=%d, width=%d", dst.GetType().String(), u.size, u.rbytes, u.width)
	}

	if u.rbytes > u.size {
		return 0, fmt.Errorf("overrun error: type=%s, size=%d, readBytes=%d, width=%d", dst.GetType().String(), u.size, u.rbytes, u.width)
	}

	// for Apple Quick Time (hdlr handlerName)
	if dst.GetType() == BoxTypeHdlr() {
		hdlr := dst.(*Hdlr)
		if err := hdlr.unmarshalHandlerName(u); err != nil {
			return 0, err
		}
	}

	return u.rbytes, nil
}

func (u *unmarshaller) unmarshal(t reflect.Type, v reflect.Value, config fieldConfig) error {
	var err error
	switch t.Kind() {
	case reflect.Ptr:
		err = u.unmarshalPtr(t, v, config)
	case reflect.Struct:
		err = u.unmarshalStruct(t, v)
	case reflect.Array:
		err = u.unmarshalArray(t, v, config)
	case reflect.Slice:
		err = u.unmarshalSlice(t, v, config)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		err = u.unmarshalInt(t, v, config)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		err = u.unmarshalUint(t, v, config)
	case reflect.Bool:
		err = u.unmarshalBool(t, v, config)
	case reflect.String:
		err = u.unmarshalString(t, v, config)
	default:
		return fmt.Errorf("unsupported type: %s", t.Kind())
	}
	return err
}

func (u *unmarshaller) unmarshalPtr(t reflect.Type, v reflect.Value, config fieldConfig) error {
	v.Set(reflect.New(t.Elem()))
	return u.unmarshal(t.Elem(), v.Elem(), config)
}

func (u *unmarshaller) unmarshalStruct(t reflect.Type, v reflect.Value) error {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		ft := f.Type
		fv := v.Field(i)

		tagStr, ok := f.Tag.Lookup("mp4")
		if !ok {
			continue
		}
		config, err := readFieldConfig(u.dst, v, f.Name, parseFieldTag(tagStr))
		if err != nil {
			return err
		}

		if !isTargetField(u.dst, config) {
			continue
		}

		err = u.unmarshal(ft, fv, config)
		if err != nil {
			return err
		}

		if ft == reflect.TypeOf(FullBox{}) && !u.dst.GetType().IsSupportedVersion(u.dst.GetVersion()) {
			return ErrUnsupportedBoxVersion
		}
	}

	return nil
}

func (u *unmarshaller) unmarshalArray(t reflect.Type, v reflect.Value, config fieldConfig) error {
	size := t.Size()
	for i := 0; i < int(size)/int(t.Elem().Size()); i++ {
		var err error
		err = u.unmarshal(t.Elem(), v.Index(i), config)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *unmarshaller) unmarshalSlice(t reflect.Type, v reflect.Value, config fieldConfig) error {
	var slice reflect.Value
	elemType := t.Elem()

	length := uint64(config.Len)
	if config.Len == lengthUnlimited {
		if config.Size != 0 {
			left := (u.size-u.rbytes)*8 + uint64(u.width)
			if left%uint64(config.Size) != 0 {
				return errors.New("invalid alignment")
			}
			length = left / uint64(config.Size)
		} else {
			length = 0
		}
	}

	if length > math.MaxInt32 {
		return fmt.Errorf("out of memory: requestedSize=%d", length)
	}

	if config.Size != 0 && config.Size%8 == 0 && u.width == 0 {
		totalSize := length * uint64(config.Size) / 8
		buf := bytes.NewBuffer(make([]byte, 0, totalSize))
		if _, err := io.CopyN(buf, u.reader, int64(totalSize)); err != nil {
			return err
		}
		data := buf.Bytes()

		if elemType.Kind() == reflect.Uint8 && config.Size == 8 {
			slice = reflect.ValueOf(data)
			u.rbytes += uint64(totalSize)

		} else {
			tmpReader := bytes.NewReader(data)
			orgReader := u.reader
			u.reader = tmpReader
			defer func() {
				u.reader = orgReader
			}()

			slice = reflect.MakeSlice(t, 0, int(length))
			for i := 0; i < int(length); i++ {
				slice = reflect.Append(slice, reflect.Zero(elemType))

				var err error
				err = u.unmarshal(elemType, slice.Index(i), config)
				if err != nil {
					return err
				}
			}

			if tmpReader.Len() != 0 {
				return fmt.Errorf("unread bytes are detected: %d", tmpReader.Len())
			}
		}

	} else {
		slice = reflect.MakeSlice(t, 0, int(length))
		for i := 0; ; i++ {
			if config.Len != lengthUnlimited && uint(i) >= config.Len {
				break
			}

			if config.Len == lengthUnlimited && u.rbytes >= u.size && u.width == 0 {
				break
			}

			slice = reflect.Append(slice, reflect.Zero(elemType))

			var err error
			err = u.unmarshal(elemType, slice.Index(i), config)
			if err != nil {
				return err
			}

			if u.rbytes > u.size {
				return fmt.Errorf("failed to read array completely: fieldName=\"%s\"", config.Name)
			}
		}
	}

	v.Set(slice)
	return nil
}

func (u *unmarshaller) unmarshalInt(t reflect.Type, v reflect.Value, config fieldConfig) error {
	if config.Varint {
		return errors.New("signed varint is unsupported")
	}

	if config.Size == 0 {
		return fmt.Errorf("size must not be zero: %s", config.Name)
	}

	data, err := u.read(config.Size)
	if err != nil {
		return err
	}

	signBit := false
	if len(data) > 0 {
		mask := byte(0x01) << ((config.Size - 1) % 8)
		signBit = data[0]&mask != 0
		data[0] &= ^mask
	}

	val := uint64(0)
	for i := range data {
		val <<= 8
		val |= uint64(data[i])
	}

	if signBit {
		v.SetInt(-int64(val))
	} else {
		v.SetInt(int64(val))
	}

	return nil
}

func (u *unmarshaller) unmarshalUint(t reflect.Type, v reflect.Value, config fieldConfig) error {
	if config.Varint {
		val, err := u.readUvarint()
		if err != nil {
			return err
		}
		v.SetUint(val)
		return nil
	}

	if config.Size == 0 {
		return fmt.Errorf("size must not be zero: %s", config.Name)
	}

	data, err := u.read(config.Size)
	if err != nil {
		return err
	}

	val := uint64(0)
	for i := range data {
		val <<= 8
		val |= uint64(data[i])
	}
	v.SetUint(val)

	return nil
}

func (u *unmarshaller) unmarshalBool(t reflect.Type, v reflect.Value, config fieldConfig) error {
	if config.Size == 0 {
		return fmt.Errorf("size must not be zero: %s", config.Name)
	}

	data, err := u.read(config.Size)
	if err != nil {
		return err
	}

	val := false
	for _, b := range data {
		val = val || (b != byte(0))
	}
	v.SetBool(val)

	return nil
}

func (u *unmarshaller) unmarshalString(t reflect.Type, v reflect.Value, config fieldConfig) error {
	data := make([]byte, 0, 16)
	for {
		if u.rbytes >= u.size {
			break
		}

		c, err := u.readOctet()
		if err != nil {
			return err
		}

		if c == 0 {
			break // null character
		}

		data = append(data, c)
	}
	v.SetString(string(data))

	return nil
}

func (u *unmarshaller) readUvarint() (uint64, error) {
	var val uint64
	for {
		octet, err := u.readOctet()
		if err != nil {
			return 0, err
		}

		val = (val << 7) + uint64(octet&0x7f)

		if octet&0x80 == 0 {
			return val, nil
		}
	}
}

func (u *unmarshaller) read(size uint) ([]byte, error) {
	// return value format:
	//  |-1-byte-block-|--------------|--------------|--------------|
	//  |<-offset->|<-------------------data----------------------->|
	bytes := (size + 7) / 8
	result := make([]byte, bytes)
	offset := (bytes * 8) - (size)

	for i := uint(0); i < size; i++ {
		b, err := u.readBit()
		if err != nil {
			return nil, err
		}

		byteIdx := (offset + i) / 8
		bitIdx := 7 - (offset+i)%8
		result[byteIdx] |= b << bitIdx
	}

	return result, nil
}

func (u *unmarshaller) readBit() (byte, error) {
	if u.width == 0 {
		octet, err := u.readOctet()
		if err != nil {
			return 0, err
		}
		u.octet = octet
		u.width = 8
	}

	u.width--
	return (u.octet >> u.width) & 0x01, nil
}

func (u *unmarshaller) readOctet() (byte, error) {
	if u.width != 0 {
		return 0, errors.New("invalid alignment")
	}

	buf := make([]byte, 1)
	if _, err := u.reader.Read(buf); err != nil {
		return 0, err
	}
	u.rbytes++
	return buf[0], nil
}

type fieldConfig struct {
	Name       string
	CFO        ICustomFieldObject
	Size       uint
	Len        uint
	Varint     bool
	Version    uint8
	NVersion   uint8
	OptDynamic bool
	OptFlag    uint32
	NOptFlag   uint32
	Const      string
	Extend     bool
	Hex        bool
	String     bool
	ISO639_2   bool
}

func readFieldConfig(box IImmutableBox, parent reflect.Value, fieldName string, tag fieldTag) (config fieldConfig, err error) {
	config.Name = fieldName
	cfo, ok := parent.Addr().Interface().(ICustomFieldObject)
	if ok {
		config.CFO = cfo
	} else {
		config.CFO = box
	}

	if val, contained := tag["size"]; contained {
		if val == "dynamic" {
			config.Size = config.CFO.GetFieldSize(fieldName)
		} else {
			var size uint64
			size, err = strconv.ParseUint(val, 10, 32)
			if err != nil {
				return
			}
			config.Size = uint(size)
		}
	}

	config.Len = lengthUnlimited
	if val, contained := tag["len"]; contained {
		if val == "dynamic" {
			config.Len = config.CFO.GetFieldLength(fieldName)
		} else {
			var l uint64
			l, err = strconv.ParseUint(val, 10, 32)
			if err != nil {
				return
			}
			config.Len = uint(l)
		}
	}

	if _, contained := tag["varint"]; contained {
		config.Varint = true
	}

	config.Version = anyVersion
	if val, contained := tag["ver"]; contained {
		var ver int
		ver, err = strconv.Atoi(val)
		if err != nil {
			return
		}
		config.Version = uint8(ver)
	}

	config.NVersion = anyVersion
	if val, contained := tag["nver"]; contained {
		var ver int
		ver, err = strconv.Atoi(val)
		if err != nil {
			return
		}
		config.NVersion = uint8(ver)
	}

	if val, contained := tag["opt"]; contained {
		if val == "dynamic" {
			config.OptDynamic = true
		} else {
			var opt uint64
			if strings.HasPrefix(val, "0x") {
				opt, err = strconv.ParseUint(val[2:], 16, 32)
			} else {
				opt, err = strconv.ParseUint(val, 10, 32)
			}
			if err != nil {
				return
			}
			config.OptFlag = uint32(opt)
		}
	}

	if val, contained := tag["nopt"]; contained {
		var nopt uint64
		if strings.HasPrefix(val, "0x") {
			nopt, err = strconv.ParseUint(val[2:], 16, 32)
		} else {
			nopt, err = strconv.ParseUint(val, 10, 32)
		}
		if err != nil {
			return
		}
		config.NOptFlag = uint32(nopt)
	}

	if val, contained := tag["const"]; contained {
		config.Const = val
	}

	if _, contained := tag["extend"]; contained {
		config.Extend = true
	}

	if _, contained := tag["hex"]; contained {
		config.Hex = true
	}

	if _, contained := tag["string"]; contained {
		config.String = true
	}

	if _, contained := tag["iso639-2"]; contained {
		config.ISO639_2 = true
	}

	return
}

type fieldTag map[string]string

func parseFieldTag(str string) fieldTag {
	tag := make(map[string]string, 8)

	list := strings.Split(str, ",")
	for _, e := range list {
		kv := strings.SplitN(e, "=", 2)
		if len(kv) == 2 {
			tag[strings.Trim(kv[0], " ")] = strings.Trim(kv[1], " ")
		} else {
			tag[strings.Trim(kv[0], " ")] = ""
		}
	}

	return tag
}

func isTargetField(box IImmutableBox, config fieldConfig) bool {
	if box.GetVersion() != anyVersion {
		if config.Version != anyVersion && box.GetVersion() != config.Version {
			return false
		}

		if config.NVersion != anyVersion && box.GetVersion() == config.NVersion {
			return false
		}
	}

	if config.OptFlag != 0 && box.GetFlags()&config.OptFlag == 0 {
		return false
	}

	if config.NOptFlag != 0 && box.GetFlags()&config.NOptFlag != 0 {
		return false
	}

	if config.OptDynamic && !config.CFO.IsOptFieldEnabled(config.Name) {
		return false
	}

	return true
}
