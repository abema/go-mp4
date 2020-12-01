package mp4

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"

	"github.com/abema/go-mp4/bitio"
)

const (
	anyVersion = math.MaxUint8
)

var ErrUnsupportedBoxVersion = errors.New("unsupported box version")

type marshaller struct {
	writer bitio.Writer
	wbits  uint64
	src    IImmutableBox
	ctx    Context
}

func Marshal(w io.Writer, src IImmutableBox, ctx Context) (n uint64, err error) {
	t := reflect.TypeOf(src).Elem()
	v := reflect.ValueOf(src).Elem()

	m := &marshaller{
		writer: bitio.NewWriter(w),
		src:    src,
		ctx:    ctx,
	}

	if err := m.marshalStruct(t, v); err != nil {
		return 0, err
	}

	if m.wbits%8 != 0 {
		return 0, fmt.Errorf("box size is not multiple of 8 bits: type=%s, bits=%d", src.GetType().String(), m.wbits)
	}

	return m.wbits / 8, nil
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
		config, err := readFieldConfig(m.src, v, f.Name, parseFieldTag(tagStr), m.ctx)
		if err != nil {
			return err
		}

		if !isTargetField(m.src, config, m.ctx) {
			continue
		}

		wbits, override, err := config.cfo.OnWriteField(f.Name, m.writer, m.ctx)
		if err != nil {
			return err
		}
		m.wbits += wbits
		if override {
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
	if config.length != LengthUnlimited {
		if length < uint64(config.length) {
			return fmt.Errorf("the slice has too few elements: required=%d actual=%d", config.length, length)
		}
		length = uint64(config.length)
	}

	elemType := t.Elem()
	if elemType.Kind() == reflect.Uint8 && config.size == 8 && m.wbits%8 == 0 {
		if _, err := io.CopyN(m.writer, bytes.NewBuffer(v.Bytes()), int64(length)); err != nil {
			return err
		}
		m.wbits += length * 8
		return nil
	}

	for i := 0; i < int(length); i++ {
		m.marshal(t.Elem(), v.Index(i), config)
	}
	return nil
}

func (m *marshaller) marshalInt(t reflect.Type, v reflect.Value, config fieldConfig) error {
	signed := v.Int()

	if config.varint {
		return errors.New("signed varint is unsupported")
	}

	signBit := signed < 0
	val := uint64(signed)
	for i := uint(0); i < config.size; i += 8 {
		v := val
		size := uint(8)
		if config.size > i+8 {
			v = v >> (config.size - (i + 8))
		} else if config.size < i+8 {
			size = config.size - i
		}

		// set sign bit
		if i == 0 {
			if signBit {
				v |= 0x1 << (size - 1)
			} else {
				v &= 0x1<<(size-1) - 1
			}
		}

		if err := m.writer.WriteBits([]byte{byte(v)}, size); err != nil {
			return err
		}
		m.wbits += uint64(size)
	}

	return nil
}

func (m *marshaller) marshalUint(t reflect.Type, v reflect.Value, config fieldConfig) error {
	val := v.Uint()

	if config.varint {
		m.writeUvarint(val)
		return nil
	}

	for i := uint(0); i < config.size; i += 8 {
		v := val
		size := uint(8)
		if config.size > i+8 {
			v = v >> (config.size - (i + 8))
		} else if config.size < i+8 {
			size = config.size - i
		}
		if err := m.writer.WriteBits([]byte{byte(v)}, size); err != nil {
			return err
		}
		m.wbits += uint64(size)
	}

	return nil
}

func (m *marshaller) marshalBool(t reflect.Type, v reflect.Value, config fieldConfig) error {
	var val byte
	if v.Bool() {
		val = 0xff
	} else {
		val = 0x00
	}
	if err := m.writer.WriteBits([]byte{val}, config.size); err != nil {
		return err
	}
	m.wbits += uint64(config.size)
	return nil
}

func (m *marshaller) marshalString(t reflect.Type, v reflect.Value, config fieldConfig) error {
	data := []byte(v.String())
	for _, b := range data {
		if err := m.writer.WriteBits([]byte{b}, 8); err != nil {
			return err
		}
		m.wbits += 8
	}
	// null character
	if err := m.writer.WriteBits([]byte{0x00}, 8); err != nil {
		return err
	}
	m.wbits += 8
	return nil
}

func (m *marshaller) writeUvarint(u uint64) error {
	if u == 0 {
		if err := m.writer.WriteBits([]byte{0}, 8); err != nil {
			return err
		}
		m.wbits += 8
		return nil
	}
	for i := 63; i >= 0; i -= 7 {
		if u>>uint(i) != 0 {
			data := byte(u>>uint(i)) & 0x7f
			if i != 0 {
				data |= 0x80
			}
			if err := m.writer.WriteBits([]byte{data}, 8); err != nil {
				return err
			}
			m.wbits += 8
		}
	}
	return nil
}

type unmarshaller struct {
	reader bitio.ReadSeeker
	dst    IBox
	size   uint64
	rbits  uint64
	ctx    Context
}

func UnmarshalAny(r io.ReadSeeker, boxType BoxType, payloadSize uint64, ctx Context) (box IBox, n uint64, err error) {
	if dst, err := boxType.New(ctx); err != nil {
		return nil, 0, err
	} else {
		n, err := Unmarshal(r, payloadSize, dst, ctx)
		return dst, n, err
	}
}

func Unmarshal(r io.ReadSeeker, payloadSize uint64, dst IBox, ctx Context) (n uint64, err error) {
	t := reflect.TypeOf(dst).Elem()
	v := reflect.ValueOf(dst).Elem()

	dst.SetVersion(anyVersion)

	u := &unmarshaller{
		reader: bitio.NewReadSeeker(r),
		dst:    dst,
		size:   payloadSize,
		ctx:    ctx,
	}

	if n, override, err := dst.BeforeUnmarshal(r, payloadSize, u.ctx); err != nil {
		return 0, err
	} else if override {
		return n, nil
	} else {
		u.rbits = n * 8
	}

	sn, err := r.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	if err := u.unmarshalStruct(t, v); err != nil {
		if err == ErrUnsupportedBoxVersion {
			r.Seek(sn, io.SeekStart)
		}
		return 0, err
	}

	if u.rbits%8 != 0 {
		return 0, fmt.Errorf("box size is not multiple of 8 bits: type=%s, size=%d, bits=%d", dst.GetType().String(), u.size, u.rbits)
	}

	if u.rbits > u.size*8 {
		return 0, fmt.Errorf("overrun error: type=%s, size=%d, bits=%d", dst.GetType().String(), u.size, u.rbits)
	}

	return u.rbits / 8, nil
}

func (u *unmarshaller) unmarshal(t reflect.Type, v reflect.Value, config fieldConfig) error {
	var err error
	switch t.Kind() {
	case reflect.Ptr:
		err = u.unmarshalPtr(t, v, config)
	case reflect.Struct:
		err = u.unmarshalStructWithConfig(t, v, config)
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

func (u *unmarshaller) unmarshalStructWithConfig(t reflect.Type, v reflect.Value, config fieldConfig) error {
	if config.size != 0 && config.size%8 == 0 {
		u2 := *u
		u2.size = uint64(config.size / 8)
		u2.rbits = 0
		if err := u2.unmarshalStruct(t, v); err != nil {
			return err
		}
		u.rbits += u2.rbits
		if u2.rbits != uint64(config.size) {
			return errors.New("invalid alignment")
		}
		return nil
	}

	return u.unmarshalStruct(t, v)
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
		config, err := readFieldConfig(u.dst, v, f.Name, parseFieldTag(tagStr), u.ctx)
		if err != nil {
			return err
		}

		if !isTargetField(u.dst, config, u.ctx) {
			continue
		}

		rbits, override, err := config.cfo.OnReadField(f.Name, u.reader, u.size*8-u.rbits, u.ctx)
		if err != nil {
			return err
		}
		u.rbits += rbits
		if override {
			continue
		}

		err = u.unmarshal(ft, fv, config)
		if err != nil {
			return err
		}

		if ft == reflect.TypeOf(FullBox{}) && !u.dst.GetType().IsSupportedVersion(u.dst.GetVersion(), u.ctx) {
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

	length := uint64(config.length)
	if config.length == LengthUnlimited {
		if config.size != 0 {
			left := (u.size)*8 - u.rbits
			if left%uint64(config.size) != 0 {
				return errors.New("invalid alignment")
			}
			length = left / uint64(config.size)
		} else {
			length = 0
		}
	}

	if length > math.MaxInt32 {
		return fmt.Errorf("out of memory: requestedSize=%d", length)
	}

	if config.size != 0 && config.size%8 == 0 && u.rbits%8 == 0 && elemType.Kind() == reflect.Uint8 && config.size == 8 {
		totalSize := length * uint64(config.size) / 8
		buf := bytes.NewBuffer(make([]byte, 0, totalSize))
		if _, err := io.CopyN(buf, u.reader, int64(totalSize)); err != nil {
			return err
		}
		slice = reflect.ValueOf(buf.Bytes())
		u.rbits += uint64(totalSize) * 8

	} else {
		slice = reflect.MakeSlice(t, 0, int(length))
		for i := 0; ; i++ {
			if config.length != LengthUnlimited && uint(i) >= config.length {
				break
			}

			if config.length == LengthUnlimited && u.rbits >= u.size*8 {
				break
			}

			slice = reflect.Append(slice, reflect.Zero(elemType))

			var err error
			err = u.unmarshal(elemType, slice.Index(i), config)
			if err != nil {
				return err
			}

			if u.rbits > u.size*8 {
				return fmt.Errorf("failed to read array completely: fieldName=\"%s\"", config.name)
			}
		}
	}

	v.Set(slice)
	return nil
}

func (u *unmarshaller) unmarshalInt(t reflect.Type, v reflect.Value, config fieldConfig) error {
	if config.varint {
		return errors.New("signed varint is unsupported")
	}

	if config.size == 0 {
		return fmt.Errorf("size must not be zero: %s", config.name)
	}

	data, err := u.reader.ReadBits(config.size)
	if err != nil {
		return err
	}
	u.rbits += uint64(config.size)

	signBit := false
	if len(data) > 0 {
		signMask := byte(0x01) << ((config.size - 1) % 8)
		signBit = data[0]&signMask != 0
		if signBit {
			data[0] |= ^(signMask - 1)
		}
	}

	var val uint64
	if signBit {
		val = ^uint64(0)
	}
	for i := range data {
		val <<= 8
		val |= uint64(data[i])
	}
	v.SetInt(int64(val))
	return nil
}

func (u *unmarshaller) unmarshalUint(t reflect.Type, v reflect.Value, config fieldConfig) error {
	if config.varint {
		val, err := u.readUvarint()
		if err != nil {
			return err
		}
		v.SetUint(val)
		return nil
	}

	if config.size == 0 {
		return fmt.Errorf("size must not be zero: %s", config.name)
	}

	data, err := u.reader.ReadBits(config.size)
	if err != nil {
		return err
	}
	u.rbits += uint64(config.size)

	val := uint64(0)
	for i := range data {
		val <<= 8
		val |= uint64(data[i])
	}
	v.SetUint(val)

	return nil
}

func (u *unmarshaller) unmarshalBool(t reflect.Type, v reflect.Value, config fieldConfig) error {
	if config.size == 0 {
		return fmt.Errorf("size must not be zero: %s", config.name)
	}

	data, err := u.reader.ReadBits(config.size)
	if err != nil {
		return err
	}
	u.rbits += uint64(config.size)

	val := false
	for _, b := range data {
		val = val || (b != byte(0))
	}
	v.SetBool(val)

	return nil
}

func (u *unmarshaller) unmarshalString(t reflect.Type, v reflect.Value, config fieldConfig) error {
	switch config.strType {
	case StringType_C:
		return u.unmarshalString_C(t, v, config)
	case StringType_C_P:
		return u.unmarshalString_C_P(t, v, config)
	default:
		return fmt.Errorf("unknown string type: %d", config.strType)
	}
}

func (u *unmarshaller) unmarshalString_C(t reflect.Type, v reflect.Value, config fieldConfig) error {
	data := make([]byte, 0, 16)
	for {
		if u.rbits >= u.size*8 {
			break
		}

		c, err := u.reader.ReadBits(8)
		if err != nil {
			return err
		}
		u.rbits += 8

		if c[0] == 0 {
			break // null character
		}

		data = append(data, c[0])
	}
	v.SetString(string(data))

	return nil
}

func (u *unmarshaller) unmarshalString_C_P(t reflect.Type, v reflect.Value, config fieldConfig) error {
	if ok, err := u.tryReadPString(t, v, config); err != nil {
		return err
	} else if ok {
		return nil
	}
	return u.unmarshalString_C(t, v, config)
}

func (u *unmarshaller) tryReadPString(t reflect.Type, v reflect.Value, config fieldConfig) (ok bool, err error) {
	remainingSize := (u.size*8 - u.rbits) / 8
	if remainingSize < 2 {
		return false, nil
	}

	offset, err := u.reader.Seek(0, io.SeekCurrent)
	if err != nil {
		return false, err
	}
	defer func() {
		if err == nil && !ok {
			_, err = u.reader.Seek(offset, io.SeekStart)
		}
	}()

	buf0 := make([]byte, 1)
	if _, err := io.ReadFull(u.reader, buf0); err != nil {
		return false, err
	}
	remainingSize--
	plen := buf0[0]
	if uint64(plen) > remainingSize {
		return false, nil
	}
	buf := make([]byte, int(plen))
	if _, err := io.ReadFull(u.reader, buf); err != nil {
		return false, err
	}
	remainingSize -= uint64(plen)
	if config.cfo.IsPString(config.name, buf, remainingSize, u.ctx) {
		u.rbits += uint64(len(buf)+1) * 8
		v.SetString(string(buf))
		return true, nil
	}
	return false, nil
}

func (u *unmarshaller) readUvarint() (uint64, error) {
	var val uint64
	for {
		octet, err := u.reader.ReadBits(8)
		if err != nil {
			return 0, err
		}
		u.rbits += 8

		val = (val << 7) + uint64(octet[0]&0x7f)

		if octet[0]&0x80 == 0 {
			return val, nil
		}
	}
}
