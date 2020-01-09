package mp4

import (
	"bytes"
	"fmt"
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
		String string `mp4:"string"`
		Bytes  []byte `mp4:"size=8,len=5"`
		Uints  []uint `mp4:"size=16,len=dynamic"`
		Ptr    *inner `mp4:"extend"`

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

		String: "abema.tv",
		Bytes:  []byte("abema"),
		Uints:  []uint{0x01, 0x02, 0x03, 0x04, 0x05},
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
