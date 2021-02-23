package mp4

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringify(t *testing.T) {
	type inner struct {
		Uint64 uint64 `mp4:"0,size=64,hex"`
	}

	type testBox struct {
		AnyTypeBox
		FullBox       `mp4:"0,extend"`
		String        string   `mp4:"1,string"`
		Int32         int32    `mp4:"2,size=32"`
		Int32Hex      int32    `mp4:"3,size=32,hex"`
		Int32HexMinus int32    `mp4:"4,size=32,hex"`
		Uint32        uint32   `mp4:"5,size=32"`
		Bytes         []byte   `mp4:"6,size=8,string"`
		Ptr           *inner   `mp4:"7"`
		PtrEx         *inner   `mp4:"8,extend"`
		Struct        inner    `mp4:"9"`
		StructEx      inner    `mp4:"10,extend"`
		Array         [7]byte  `mp4:"11,size=8,string"`
		Bool          bool     `mp4:"12,size=1"`
		UUID          [16]byte `mp4:"13,size=8,uuid"`
		NotSorted15   uint8    `mp4:"15,size=8,dec"`
		NotSorted16   uint8    `mp4:"16,size=8,dec"`
		NotSorted14   uint8    `mp4:"14,size=8,dec"`
	}
	boxType := StrToBoxType("test")
	AddAnyTypeBoxDef(&testBox{}, boxType)

	box := testBox{
		AnyTypeBox: AnyTypeBox{
			Type: boxType,
		},
		FullBox: FullBox{
			Version: 0,
			Flags:   [3]byte{0x00, 0x00, 0x00},
		},
		String:        "abema.tv",
		Int32:         -1234567890,
		Int32Hex:      0x12345678,
		Int32HexMinus: -0x12345678,
		Uint32:        1234567890,
		Bytes:         []byte{'A', 'B', 'E', 'M', 'A', 0x00, 'T', 'V'},
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
		Array:       [7]byte{'f', 'o', 'o', 0x00, 'b', 'a', 'r'},
		Bool:        true,
		UUID:        [16]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef},
		NotSorted15: 15,
		NotSorted16: 16,
		NotSorted14: 14,
	}

	str, err := StringifyWithIndent(&box, " ", Context{})
	require.NoError(t, err)
	assert.Equal(t, ` Version=0`+"\n"+
		` Flags=0x000000`+"\n"+
		` String="abema.tv"`+"\n"+
		` Int32=-1234567890`+"\n"+
		` Int32Hex=0x12345678`+"\n"+
		` Int32HexMinus=-0x12345678`+"\n"+
		` Uint32=1234567890`+"\n"+
		` Bytes="ABEMA.TV"`+"\n"+
		` Ptr={`+"\n"+
		`  Uint64=0x1234567890`+"\n"+
		` }`+"\n"+
		` Uint64=0x1234567890`+"\n"+
		` Struct={`+"\n"+
		`  Uint64=0x1234567890`+"\n"+
		` }`+"\n"+
		` Uint64=0x1234567890`+"\n"+
		` Array="foo.bar"`+"\n"+
		` Bool=true`+"\n"+
		` UUID=01234567-89ab-cdef-0123-456789abcdef`+"\n"+
		` NotSorted14=14`+"\n"+
		` NotSorted15=15`+"\n"+
		` NotSorted16=16`+"\n", str)

	str, err = Stringify(&box, Context{})
	require.NoError(t, err)
	assert.Equal(t, `Version=0 Flags=0x000000 String="abema.tv" Int32=-1234567890 Int32Hex=0x12345678 Int32HexMinus=-0x12345678 Uint32=1234567890 Bytes="ABEMA.TV" Ptr={Uint64=0x1234567890} Uint64=0x1234567890 Struct={Uint64=0x1234567890} Uint64=0x1234567890 Array="foo.bar" Bool=true UUID=01234567-89ab-cdef-0123-456789abcdef NotSorted14=14 NotSorted15=15 NotSorted16=16`, str)
}
