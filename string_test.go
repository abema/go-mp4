package mp4

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmsgStringify(t *testing.T) {
	type inner struct {
		Uint64 uint64 `mp4:"size=64,hex"`
	}

	type testBox struct {
		AnyTypeBox
		FullBox  `mp4:"extend"`
		String   string  `mp4:"string"`
		Int32    int32   `mp4:"size=32"`
		Int32Hex int32   `mp4:"size=32,hex"`
		Uint32   uint32  `mp4:"size=32"`
		Bytes    []byte  `mp4:"size=8,string"`
		Ptr      *inner  `mp4:""`
		PtrEx    *inner  `mp4:"extend"`
		Struct   inner   `mp4:""`
		StructEx inner   `mp4:"extend"`
		Array    [4]byte `mp4:"size=8,string"`
		Bool     bool    `mp4:"size=1"`
	}

	box := testBox{
		FullBox: FullBox{
			Version: 0,
			Flags:   [3]byte{0x00, 0x00, 0x00},
		},
		String:   "abema.tv",
		Int32:    -1234567890,
		Int32Hex: 0x12345678,
		Uint32:   1234567890,
		Bytes:    []byte("abema"),
		Ptr: &inner{
			Uint64: 0x1234567890,
		},
		PtrEx: &inner{
			Uint64: 0x1234567890,
		},
		Struct: inner{
			Uint64: 0x1234567890,
		},
		StructEx: inner{
			Uint64: 0x1234567890,
		},
		Array: [4]byte{'h', 'o', 'g', 'e'},
		Bool:  true,
	}

	str, err := StringifyWithIndent(&box, " ")
	require.NoError(t, err)
	assert.Equal(t, ` Version=0`+"\n"+
		` Flags=0x000000`+"\n"+
		` String="abema.tv"`+"\n"+
		` Int32=-1234567890`+"\n"+
		` Int32Hex=0x12345678`+"\n"+
		` Uint32=1234567890`+"\n"+
		` Bytes="abema"`+"\n"+
		` Ptr={`+"\n"+
		`  Uint64=0x1234567890`+"\n"+
		` }`+"\n"+
		` Uint64=0x1234567890`+"\n"+
		` Struct={`+"\n"+
		`  Uint64=0x1234567890`+"\n"+
		` }`+"\n"+
		` Uint64=0x1234567890`+"\n"+
		` Array="hoge"`+"\n"+
		` Bool=true`+"\n", str)

	str, err = Stringify(&box)
	require.NoError(t, err)
	assert.Equal(t, `Version=0 Flags=0x000000 String="abema.tv" Int32=-1234567890 Int32Hex=0x12345678 Uint32=1234567890 Bytes="abema" Ptr={Uint64=0x1234567890} Uint64=0x1234567890 Struct={Uint64=0x1234567890} Uint64=0x1234567890 Array="hoge" Bool=true`, str)
}
