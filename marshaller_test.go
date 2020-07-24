package mp4

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockBox struct {
	Type       BoxType
	DynSizeMap map[string]uint
	DynLenMap  map[string]uint
}

func (m *mockBox) GetType() BoxType {
	return m.Type
}

func (m *mockBox) GetFieldSize(n string) uint {
	if s, ok := m.DynSizeMap[n]; !ok {
		panic(fmt.Errorf("invalid name of dynamic-size field: %s", n))
	} else {
		return s
	}
}

func (m *mockBox) GetFieldLength(n string) uint {
	if l, ok := m.DynLenMap[n]; !ok {
		panic(fmt.Errorf("invalid name of dynamic-length field: %s", n))
	} else {
		return l
	}
}

func TestMarshal(t *testing.T) {
	type inner struct {
		Array [4]byte `mp4:"size=8,string"`
	}

	type testBox struct {
		mockBox
		FullBox `mp4:"extend"`

		// integer
		Int32  int32  `mp4:"size=32"`
		Uint32 uint32 `mp4:"size=32"`
		Int64  int64  `mp4:"size=64"`
		Uint64 uint64 `mp4:"size=64"`

		// left-justified
		Int32l   int32  `mp4:"size=29"`
		Padding0 uint8  `mp4:"size=3,const=0"`
		Uint32l  uint32 `mp4:"size=29"`
		Padding1 uint8  `mp4:"size=3,const=0"`
		Int64l   int64  `mp4:"size=59"`
		Padding2 uint8  `mp4:"size=5,const=0"`
		Uint64l  uint64 `mp4:"size=59"`
		Padding3 uint8  `mp4:"size=5,const=0"`

		// right-justified
		Padding4 uint8  `mp4:"size=3,const=0"`
		Int32r   int32  `mp4:"size=29"`
		Padding5 uint8  `mp4:"size=3,const=0"`
		Uint32r  uint32 `mp4:"size=29"`
		Padding6 uint8  `mp4:"size=5,const=0"`
		Int64r   int64  `mp4:"size=59"`
		Padding7 uint8  `mp4:"size=5,const=0"`
		Uint64r  uint64 `mp4:"size=59"`

		// varint
		Varint uint16 `mp4:"varint"`

		// string, slice, pointer, array
		String     string `mp4:"string"`
		String_C_P string `mp4:"string=c_p"`
		Bytes      []byte `mp4:"size=8,len=5"`
		Uints      []uint `mp4:"size=16,len=dynamic"`
		Ptr        *inner `mp4:"extend"`

		// bool
		Bool     bool  `mp4:"size=1"`
		Padding8 uint8 `mp4:"size=7,const=0"`

		// dynamic-size
		DynUint uint `mp4:"size=dynamic"`

		// optional
		OptUint1 uint `mp4:"size=8,opt=0x0100"`  // enabled
		OptUint2 uint `mp4:"size=8,opt=0x0200"`  // disabled
		OptUint3 uint `mp4:"size=8,nopt=0x0400"` // disabled
		OptUint4 uint `mp4:"size=8,nopt=0x0800"` // enabled
	}

	boxType := StrToBoxType("test")
	mb := mockBox{
		Type: boxType,
		DynSizeMap: map[string]uint{
			"DynUint": 24,
		},
		DynLenMap: map[string]uint{
			"Uints": 5,
		},
	}
	AddBoxDef(&testBox{mockBox: mb}, 0)

	src := testBox{
		mockBox: mb,

		FullBox: FullBox{
			Version: 0,
			Flags:   [3]byte{0x00, 0x05, 0x00},
		},

		Int32:  -0x1234567,
		Uint32: 0x1234567,
		Int64:  -0x123456789abcdef,
		Uint64: 0x123456789abcdef,

		Int32l:  -0x123456,
		Uint32l: 0x123456,
		Int64l:  -0x123456789abcd,
		Uint64l: 0x123456789abcd,

		Int32r:  -0x123456,
		Uint32r: 0x123456,
		Int64r:  -0x123456789abcd,
		Uint64r: 0x123456789abcd,

		// raw   : 0x1234=0001,0010,0011,0100b
		// varint: 0xa434=1010,0100,0011,0100b
		Varint: 0x1234,

		String:     "abema.tv",
		String_C_P: "CyberAgent, Inc.",
		Bytes:      []byte("abema"),
		Uints:      []uint{0x01, 0x02, 0x03, 0x04, 0x05},
		Ptr: &inner{
			Array: [4]byte{'h', 'o', 'g', 'e'},
		},

		Bool: true,

		DynUint: 0x123456,

		OptUint1: 0x11,
		OptUint4: 0x44,
	}

	bin := []byte{
		0,                // version
		0x00, 0x05, 0x00, // flags
		0x81, 0x23, 0x45, 0x67, // int32
		0x01, 0x23, 0x45, 0x67, // uint32
		0x81, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, // int64
		0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, // uint64
		0x80, 0x91, 0xa2, 0xb0, // int32l & padding
		0x00, 0x91, 0xa2, 0xb0, // uint32l & padding
		0x80, 0x24, 0x68, 0xAC, 0xF1, 0x35, 0x79, 0xA0, // int64l & padding
		0x00, 0x24, 0x68, 0xAC, 0xF1, 0x35, 0x79, 0xA0, // uint64l & padding
		0x10, 0x12, 0x34, 0x56, // padding & int32r
		0x00, 0x12, 0x34, 0x56, // padding & uint32r
		0x04, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, // padding & int64r
		0x00, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, // padding & uint64r
		0xa4, 0x34, // varint
		'a', 'b', 'e', 'm', 'a', '.', 't', 'v', 0, // string
		'C', 'y', 'b', 'e', 'r', 'A', 'g', 'e', 'n', 't', ',', ' ', 'I', 'n', 'c', '.', 0, // string
		'a', 'b', 'e', 'm', 'a', // bytes
		0x00, 0x01, 0x00, 0x02, 0x00, 0x03, 0x00, 0x04, 0x00, 0x05, // uints
		'h', 'o', 'g', 'e', // inner.array
		0x80,             // bool & padding
		0x12, 0x34, 0x56, // dynUint
		0x11, // optUint1
		0x44, // optUint4
	}

	// marshal
	buf := &bytes.Buffer{}
	n, err := Marshal(buf, &src)
	require.NoError(t, err)
	assert.Equal(t, uint64(len(bin)), n)
	assert.Equal(t, bin, buf.Bytes())

	// unmarshal
	dst := testBox{mockBox: mb}
	n, err = Unmarshal(bytes.NewReader(bin), uint64(len(bin)+8), &dst)
	assert.NoError(t, err)
	assert.Equal(t, uint64(len(bin)), n)
	assert.Equal(t, src, dst)
}

func TestUnsupportedBoxVersionErr(t *testing.T) {
	type testBox struct {
		mockBox
		FullBox `mp4:"extend"`
	}

	boxType := StrToBoxType("test")
	mb := mockBox{
		Type: boxType,
	}
	AddBoxDef(&testBox{mockBox: mb}, 0, 1, 2)

	for _, e := range []struct {
		version byte
		enabled bool
	}{
		{version: 0, enabled: true},
		{version: 1, enabled: true},
		{version: 2, enabled: true},
		{version: 3, enabled: false},
		{version: 4, enabled: false},
	} {
		expected := testBox{
			mockBox: mb,
			FullBox: FullBox{
				Version: e.version,
				Flags:   [3]byte{0x00, 0x00, 0x00},
			},
		}

		bin := []byte{
			e.version,        // version
			0x00, 0x00, 0x00, // flags
		}

		dst := testBox{mockBox: mb}
		n, err := Unmarshal(bytes.NewReader(bin), uint64(len(bin)+8), &dst)

		if e.enabled {
			assert.NoError(t, err, "version=%d", e.version)
			assert.Equal(t, uint64(len(bin)), n, "version=%d", e.version)
			assert.Equal(t, expected, dst, "version=%d", e.version)
		} else {
			assert.Error(t, err, "version=%d", e.version)
		}
	}
}

func TestReadVarint(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		err      bool
		expected uint64
	}{
		{name: "1 byte", input: []byte{0x6c}, expected: 0x6c},
		{name: "2 bytes", input: []byte{0xac, 0x52}, expected: 0x1652},
		{name: "3 bytes", input: []byte{0xac, 0xd2, 0x43}, expected: 0xb2943},
		{name: "overrun", input: []byte{0xac, 0xd2, 0xef}, err: true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := &unmarshaller{
				reader: bytes.NewReader(tc.input),
				size:   uint64(len(tc.input)),
			}
			val, err := u.readUvarint()
			if tc.err {
				require.Error(t, err)
				return
			}
			if tc.err {
				assert.Error(t, err)
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expected, val)
		})
	}
}

func TestRead(t *testing.T) {
	testCases := []struct {
		name         string
		octet        byte
		input        []byte
		width        uint
		size         uint
		err          bool
		expectedData []byte
	}{
		{name: "no width", input: []byte{0x6c, 0xa5}, size: 10, expectedData: []byte{0x01, 0xb2}},
		{name: "width 3", octet: 0x6c, input: []byte{0xa5}, width: 3, size: 10, expectedData: []byte{0x02, 0x52}},
		{name: "reach to end of box", input: []byte{0x6c, 0xa5}, size: 16, expectedData: []byte{0x6c, 0xa5}},
		{name: "overrun", input: []byte{0x6c, 0xa5}, size: 17, err: true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := &unmarshaller{
				reader: bytes.NewReader(tc.input),
				size:   uint64(len(tc.input)),
				octet:  tc.octet,
				width:  tc.width,
			}
			data, err := u.read(tc.size)
			if tc.err {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedData, data)
		})
	}
}

func TestReadBit(t *testing.T) {
	u := &unmarshaller{
		reader: bytes.NewReader([]byte{0x6c, 0xa5}),
		rbytes: 10,
	}
	outputs := []struct {
		bit    byte
		rbytes uint64
		octet  byte
	}{
		{bit: 0x00, rbytes: 11, octet: 0x6c},
		{bit: 0x01, rbytes: 11, octet: 0x6c},
		{bit: 0x01, rbytes: 11, octet: 0x6c},
		{bit: 0x00, rbytes: 11, octet: 0x6c},
		{bit: 0x01, rbytes: 11, octet: 0x6c},
		{bit: 0x01, rbytes: 11, octet: 0x6c},
		{bit: 0x00, rbytes: 11, octet: 0x6c},
		{bit: 0x00, rbytes: 11, octet: 0x6c},
		{bit: 0x01, rbytes: 12, octet: 0xa5},
		{bit: 0x00, rbytes: 12, octet: 0xa5},
		{bit: 0x01, rbytes: 12, octet: 0xa5},
		{bit: 0x00, rbytes: 12, octet: 0xa5},
		{bit: 0x00, rbytes: 12, octet: 0xa5},
		{bit: 0x01, rbytes: 12, octet: 0xa5},
		{bit: 0x00, rbytes: 12, octet: 0xa5},
		{bit: 0x01, rbytes: 12, octet: 0xa5},
	}
	for _, o := range outputs {
		bit, err := u.readBit()
		require.NoError(t, err)
		assert.Equal(t, o.bit, bit)
		assert.Equal(t, o.rbytes, u.rbytes)
		assert.Equal(t, o.octet, u.octet)
	}
	_, err := u.readBit()
	require.Error(t, err)
}

func TestReadOctet(t *testing.T) {
	u := &unmarshaller{
		reader: bytes.NewReader([]byte{0x6c, 0xa5}),
		rbytes: 10,
	}
	octet, err := u.readOctet()
	require.NoError(t, err)
	assert.Equal(t, byte(0x6c), octet)
	assert.Equal(t, uint64(11), u.rbytes)
	octet, err = u.readOctet()
	require.NoError(t, err)
	assert.Equal(t, byte(0xa5), octet)
	assert.Equal(t, uint64(12), u.rbytes)
	_, err = u.readOctet()
	require.Error(t, err)
}

func TestReadOctetInvalidAlignment(t *testing.T) {
	u := &unmarshaller{
		reader: bytes.NewReader([]byte{0x6c, 0x00}),
		width:  3,
		rbytes: 10,
	}
	_, err := u.readOctet()
	require.Error(t, err)
}

func TestReadFieldConfig(t *testing.T) {
	box := &struct {
		mockBox
		FullBox
		ByteArray []byte
		String    string
		Int       int32
	}{
		mockBox: mockBox{
			DynSizeMap: map[string]uint{
				"ByteArray": 3,
			},
			DynLenMap: map[string]uint{
				"ByteArray": 7,
			},
		},
	}

	testCases := []struct {
		name      string
		box       IImmutableBox
		fieldName string
		fieldTag  fieldTag
		err       bool
		expected  fieldConfig
	}{
		{
			name:      "static size",
			box:       box,
			fieldName: "ByteArray",
			fieldTag:  fieldTag{"size": "8"},
			expected: fieldConfig{
				Name:     "ByteArray",
				CFO:      box,
				Size:     8,
				Len:      lengthUnlimited,
				Version:  anyVersion,
				NVersion: anyVersion,
			},
		},
		{
			name:      "invalid size",
			box:       box,
			fieldName: "ByteArray",
			fieldTag:  fieldTag{"size": "invalid"},
			err:       true,
		},
		{
			name:      "dynamic size",
			box:       box,
			fieldName: "ByteArray",
			fieldTag:  fieldTag{"size": "dynamic"},
			expected: fieldConfig{
				Name:     "ByteArray",
				CFO:      box,
				Size:     3,
				Len:      lengthUnlimited,
				Version:  anyVersion,
				NVersion: anyVersion,
			},
		},
		{
			name:      "static length",
			box:       box,
			fieldName: "ByteArray",
			fieldTag:  fieldTag{"len": "16", "size": "8"},
			expected: fieldConfig{
				Name:     "ByteArray",
				CFO:      box,
				Size:     8,
				Len:      16,
				Version:  anyVersion,
				NVersion: anyVersion,
			},
		},
		{
			name:      "invalid length",
			box:       box,
			fieldName: "ByteArray",
			fieldTag:  fieldTag{"len": "foo", "size": "8"},
			err:       true,
		},
		{
			name:      "dynamic length",
			box:       box,
			fieldName: "ByteArray",
			fieldTag:  fieldTag{"len": "dynamic", "size": "8"},
			expected: fieldConfig{
				Name:     "ByteArray",
				CFO:      box,
				Size:     8,
				Len:      7,
				Version:  anyVersion,
				NVersion: anyVersion,
			},
		},
		{
			name:      "varint",
			box:       box,
			fieldName: "Int",
			fieldTag:  fieldTag{"varint": "", "size": "13"},
			expected: fieldConfig{
				Name:     "Int",
				CFO:      box,
				Size:     13,
				Len:      lengthUnlimited,
				Version:  anyVersion,
				NVersion: anyVersion,
				Varint:   true,
			},
		},
		{
			name:      "ver 0",
			box:       box,
			fieldName: "Int",
			fieldTag:  fieldTag{"ver": "0", "size": "32"},
			expected: fieldConfig{
				Name:     "Int",
				CFO:      box,
				Size:     32,
				Len:      lengthUnlimited,
				Version:  0,
				NVersion: anyVersion,
			},
		},
		{
			name:      "ver 1",
			box:       box,
			fieldName: "Int",
			fieldTag:  fieldTag{"ver": "1", "size": "32"},
			expected: fieldConfig{
				Name:     "Int",
				CFO:      box,
				Size:     32,
				Len:      lengthUnlimited,
				Version:  1,
				NVersion: anyVersion,
			},
		},
		{
			name:      "invalid ver",
			box:       box,
			fieldName: "Int",
			fieldTag:  fieldTag{"ver": "foo", "size": "32"},
			err:       true,
		},
		{
			name:      "nver 0",
			box:       box,
			fieldName: "Int",
			fieldTag:  fieldTag{"nver": "0", "size": "32"},
			expected: fieldConfig{
				Name:     "Int",
				CFO:      box,
				Size:     32,
				Len:      lengthUnlimited,
				Version:  anyVersion,
				NVersion: 0,
			},
		},
		{
			name:      "nver 1",
			box:       box,
			fieldName: "Int",
			fieldTag:  fieldTag{"nver": "1", "size": "32"},
			expected: fieldConfig{
				Name:     "Int",
				CFO:      box,
				Size:     32,
				Len:      lengthUnlimited,
				Version:  anyVersion,
				NVersion: 1,
			},
		},
		{
			name:      "invalid nver",
			box:       box,
			fieldName: "Int",
			fieldTag:  fieldTag{"nver": "foo", "size": "32"},
			err:       true,
		},
		{
			name:      "opt dynamic",
			box:       box,
			fieldName: "String",
			fieldTag:  fieldTag{"opt": "dynamic"},
			expected: fieldConfig{
				Name:       "String",
				CFO:        box,
				Len:        lengthUnlimited,
				Version:    anyVersion,
				NVersion:   anyVersion,
				OptDynamic: true,
			},
		},
		{
			name:      "opt hex",
			box:       box,
			fieldName: "String",
			fieldTag:  fieldTag{"opt": "0x0100"},
			expected: fieldConfig{
				Name:     "String",
				CFO:      box,
				Len:      lengthUnlimited,
				Version:  anyVersion,
				NVersion: anyVersion,
				OptFlag:  0x0100,
			},
		},
		{
			name:      "opt dec",
			box:       box,
			fieldName: "String",
			fieldTag:  fieldTag{"opt": "32"},
			expected: fieldConfig{
				Name:     "String",
				CFO:      box,
				Len:      lengthUnlimited,
				Version:  anyVersion,
				NVersion: anyVersion,
				OptFlag:  0x0020,
			},
		},
		{
			name:      "invalid opt",
			box:       box,
			fieldName: "String",
			fieldTag:  fieldTag{"opt": "foo"},
			err:       true,
		},
		{
			name:      "nopt hex",
			box:       box,
			fieldName: "String",
			fieldTag:  fieldTag{"nopt": "0x0100"},
			expected: fieldConfig{
				Name:     "String",
				CFO:      box,
				Len:      lengthUnlimited,
				Version:  anyVersion,
				NVersion: anyVersion,
				NOptFlag: 0x0100,
			},
		},
		{
			name:      "nopt dec",
			box:       box,
			fieldName: "String",
			fieldTag:  fieldTag{"nopt": "32"},
			expected: fieldConfig{
				Name:     "String",
				CFO:      box,
				Len:      lengthUnlimited,
				Version:  anyVersion,
				NVersion: anyVersion,
				NOptFlag: 0x0020,
			},
		},
		{
			name:      "invalid nopt",
			box:       box,
			fieldName: "String",
			fieldTag:  fieldTag{"nopt": "foo"},
			err:       true,
		},
		{
			name:      "const",
			box:       box,
			fieldName: "Int",
			fieldTag:  fieldTag{"const": "0", "size": "32"},
			expected: fieldConfig{
				Name:     "Int",
				CFO:      box,
				Size:     32,
				Len:      lengthUnlimited,
				Version:  anyVersion,
				NVersion: anyVersion,
				Const:    "0",
			},
		},
		{
			name:      "extend",
			box:       box,
			fieldName: "FullBox",
			fieldTag:  fieldTag{"extend": ""},
			expected: fieldConfig{
				Name:     "FullBox",
				CFO:      box,
				Len:      lengthUnlimited,
				Version:  anyVersion,
				NVersion: anyVersion,
				Extend:   true,
			},
		},
		{
			name:      "hex",
			box:       box,
			fieldName: "Int",
			fieldTag:  fieldTag{"hex": "", "size": "32"},
			expected: fieldConfig{
				Name:     "Int",
				CFO:      box,
				Size:     32,
				Len:      lengthUnlimited,
				Version:  anyVersion,
				NVersion: anyVersion,
				Hex:      true,
			},
		},
		{
			name:      "string - c style",
			box:       box,
			fieldName: "String",
			fieldTag:  fieldTag{"string": ""},
			expected: fieldConfig{
				Name:       "String",
				CFO:        box,
				Len:        lengthUnlimited,
				Version:    anyVersion,
				NVersion:   anyVersion,
				String:     true,
				StringType: StringType_C,
			},
		},
		{
			name:      "string - c style or pascal style",
			box:       box,
			fieldName: "String",
			fieldTag:  fieldTag{"string": "c_p"},
			expected: fieldConfig{
				Name:       "String",
				CFO:        box,
				Len:        lengthUnlimited,
				Version:    anyVersion,
				NVersion:   anyVersion,
				String:     true,
				StringType: StringType_C_P,
			},
		},
		{
			name:      "iso639-2",
			box:       box,
			fieldName: "String",
			fieldTag:  fieldTag{"iso639-2": ""},
			expected: fieldConfig{
				Name:     "String",
				CFO:      box,
				Len:      lengthUnlimited,
				Version:  anyVersion,
				NVersion: anyVersion,
				ISO639_2: true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(tc.box).Elem()
			config, err := readFieldConfig(tc.box, v, tc.fieldName, tc.fieldTag)
			if tc.err {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expected, config)
		})
	}
}
