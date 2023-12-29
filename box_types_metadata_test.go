package mp4

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoxTypesMetadata(t *testing.T) {
	testCases := []struct {
		name string
		src  IImmutableBox
		dst  IBox
		bin  []byte
		str  string
		ctx  Context
	}{
		{
			name: "ilst",
			src:  &Ilst{},
			dst:  &Ilst{},
			bin:  nil,
			str:  ``,
		},
		{
			name: "ilst meta container",
			src: &IlstMetaContainer{
				AnyTypeBox: AnyTypeBox{Type: StrToBoxType("----")},
			},
			dst: &IlstMetaContainer{
				AnyTypeBox: AnyTypeBox{Type: StrToBoxType("----")},
			},
			bin: nil,
			str: ``,
			ctx: Context{UnderIlst: true},
		},
		{
			name: "ilst data (binary)",
			src:  &Data{DataType: 0, DataLang: 0x12345678, Data: []byte("foo")},
			dst:  &Data{},
			bin: []byte{
				0x00, 0x00, 0x00, 0x00, // data type
				0x12, 0x34, 0x56, 0x78, // data lang
				0x66, 0x6f, 0x6f, // data
			},
			str: `DataType=BINARY DataLang=305419896 Data=[0x66, 0x6f, 0x6f]`,
			ctx: Context{UnderIlstMeta: true},
		},
		{
			name: "ilst data (utf8)",
			src:  &Data{DataType: 1, DataLang: 0x12345678, Data: []byte("foo")},
			dst:  &Data{},
			bin: []byte{
				0x00, 0x00, 0x00, 0x01, // data type
				0x12, 0x34, 0x56, 0x78, // data lang
				0x66, 0x6f, 0x6f, // data
			},
			str: `DataType=UTF8 DataLang=305419896 Data="foo"`,
			ctx: Context{UnderIlstMeta: true},
		},
		{
			name: "ilst data (utf8 escaped)",
			src:  &Data{DataType: 1, DataLang: 0x12345678, Data: []byte{0x00, 'f', 'o', 'o'}},
			dst:  &Data{},
			bin: []byte{
				0x00, 0x00, 0x00, 0x01, // data type
				0x12, 0x34, 0x56, 0x78, // data lang
				0x00, 0x66, 0x6f, 0x6f, // data
			},
			str: `DataType=UTF8 DataLang=305419896 Data=".foo"`,
			ctx: Context{UnderIlstMeta: true},
		},
		{
			name: "ilst data (utf16)",
			src:  &Data{DataType: 2, DataLang: 0x12345678, Data: []byte("foo")},
			dst:  &Data{},
			bin: []byte{
				0x00, 0x00, 0x00, 0x02, // data type
				0x12, 0x34, 0x56, 0x78, // data lang
				0x66, 0x6f, 0x6f, // data
			},
			str: `DataType=UTF16 DataLang=305419896 Data=[0x66, 0x6f, 0x6f]`,
			ctx: Context{UnderIlstMeta: true},
		},
		{
			name: "ilst data (mac string)",
			src:  &Data{DataType: 3, DataLang: 0x12345678, Data: []byte("foo")},
			dst:  &Data{},
			bin: []byte{
				0x00, 0x00, 0x00, 0x03, // data type
				0x12, 0x34, 0x56, 0x78, // data lang
				0x66, 0x6f, 0x6f, // data
			},
			str: `DataType=MAC_STR DataLang=305419896 Data=[0x66, 0x6f, 0x6f]`,
			ctx: Context{UnderIlstMeta: true},
		},
		{
			name: "ilst data (jpsg)",
			src:  &Data{DataType: 14, DataLang: 0x12345678, Data: []byte("foo")},
			dst:  &Data{},
			bin: []byte{
				0x00, 0x00, 0x00, 0x0e, // data type
				0x12, 0x34, 0x56, 0x78, // data lang
				0x66, 0x6f, 0x6f, // data
			},
			str: `DataType=JPEG DataLang=305419896 Data=[0x66, 0x6f, 0x6f]`,
			ctx: Context{UnderIlstMeta: true},
		},
		{
			name: "ilst data (int)",
			src:  &Data{DataType: 21, DataLang: 0x12345678, Data: []byte("foo")},
			dst:  &Data{},
			bin: []byte{
				0x00, 0x00, 0x00, 0x15, // data type
				0x12, 0x34, 0x56, 0x78, // data lang
				0x66, 0x6f, 0x6f, // data
			},
			str: `DataType=INT DataLang=305419896 Data=[0x66, 0x6f, 0x6f]`,
			ctx: Context{UnderIlstMeta: true},
		},
		{
			name: "ilst data (float32)",
			src:  &Data{DataType: 22, DataLang: 0x12345678, Data: []byte("foo")},
			dst:  &Data{},
			bin: []byte{
				0x00, 0x00, 0x00, 0x16, // data type
				0x12, 0x34, 0x56, 0x78, // data lang
				0x66, 0x6f, 0x6f, // data
			},
			str: `DataType=FLOAT32 DataLang=305419896 Data=[0x66, 0x6f, 0x6f]`,
			ctx: Context{UnderIlstMeta: true},
		},
		{
			name: "ilst data (float64)",
			src:  &Data{DataType: 23, DataLang: 0x12345678, Data: []byte("foo")},
			dst:  &Data{},
			bin: []byte{
				0x00, 0x00, 0x00, 0x17, // data type
				0x12, 0x34, 0x56, 0x78, // data lang
				0x66, 0x6f, 0x6f, // data
			},
			str: `DataType=FLOAT64 DataLang=305419896 Data=[0x66, 0x6f, 0x6f]`,
			ctx: Context{UnderIlstMeta: true},
		},
		{
			name: "ilst data (string)",
			src: &StringData{
				AnyTypeBox: AnyTypeBox{Type: StrToBoxType("mean")},
				Data:       []byte{0x00, 'f', 'o', 'o'},
			},
			dst: &StringData{
				AnyTypeBox: AnyTypeBox{Type: StrToBoxType("mean")},
			},
			bin: []byte{
				0x00, 0x66, 0x6f, 0x6f, // data
			},
			str: `Data=".foo"`,
			ctx: Context{UnderIlstFreeMeta: true},
		},
		{
			name: "ilst numbered item",
			src: &Item{
				AnyTypeBox: AnyTypeBox{Type: Uint32ToBoxType(1)},
				Version:    0,
				Flags:      [3]byte{'0'},
				ItemName:   []byte("data"),
				Data:       Data{DataType: 0, DataLang: 0x12345678, Data: []byte("foo")}},
			dst: &Item{
				AnyTypeBox: AnyTypeBox{Type: Uint32ToBoxType(1)},
			},
			bin: []byte{
				0x00,            // Version
				0x30, 0x00, 0x0, // Flags
				0x64, 0x61, 0x74, 0x61, // Item Name
				0x0, 0x0, 0x0, 0x0, // data type
				0x12, 0x34, 0x56, 0x78, // data lang
				0x66, 0x6f, 0x6f, // data
			},
			str: `Version=0 Flags=0x000000 ItemName="data" Data={DataType=BINARY DataLang=305419896 Data=[0x66, 0x6f, 0x6f]}`,
			ctx: Context{UnderIlst: true},
		},
		{
			name: "keys",
			src: &Keys{
				EntryCount: 2,
				Entries: []Key{
					{
						KeySize:      27,
						KeyNamespace: []byte("mdta"),
						KeyValue:     []byte("com.android.version"),
					},
					{
						KeySize:      25,
						KeyNamespace: []byte("mdta"),
						KeyValue:     []byte("com.android.model"),
					},
				},
			},
			dst: &Keys{},
			bin: []byte{
				0x0,           // Version
				0x0, 0x0, 0x0, // Flags
				0x0, 0x0, 0x0, 0x2, // entry count
				0x0, 0x0, 0x0, 0x1b, // entry 1 keysize
				0x6d, 0x64, 0x74, 0x61, // entry 1 key namespace
				0x63, 0x6f, 0x6d, 0x2e, 0x61, 0x6e, 0x64, 0x72, 0x6f, 0x69, 0x64, 0x2e, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, // entry 1 key value
				0x0, 0x0, 0x0, 0x19, // entry 2 keysize
				0x6d, 0x64, 0x74, 0x61, // entry 2 key namespace
				0x63, 0x6f, 0x6d, 0x2e, 0x61, 0x6e, 0x64, 0x72, 0x6f, 0x69, 0x64, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, // entry 2 key value
			},
			str: `Version=0 Flags=0x000000 EntryCount=2 Entries=[{KeySize=27 KeyNamespace="mdta" KeyValue="com.android.version"}, {KeySize=25 KeyNamespace="mdta" KeyValue="com.android.model"}]`,
			ctx: Context{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Marshal
			buf := bytes.NewBuffer(nil)
			n, err := Marshal(buf, tc.src, tc.ctx)
			require.NoError(t, err)
			assert.Equal(t, uint64(len(tc.bin)), n)
			assert.Equal(t, tc.bin, buf.Bytes())

			// Unmarshal
			r := bytes.NewReader(tc.bin)
			n, err = Unmarshal(r, uint64(len(tc.bin)), tc.dst, tc.ctx)
			require.NoError(t, err)
			assert.Equal(t, uint64(buf.Len()), n)
			assert.Equal(t, tc.src, tc.dst)
			s, err := r.Seek(0, io.SeekCurrent)
			require.NoError(t, err)
			assert.Equal(t, int64(buf.Len()), s)

			// UnmarshalAny
			dst, n, err := UnmarshalAny(bytes.NewReader(tc.bin), tc.src.GetType(), uint64(len(tc.bin)), tc.ctx)
			require.NoError(t, err)
			assert.Equal(t, uint64(buf.Len()), n)
			assert.Equal(t, tc.src, dst)
			s, err = r.Seek(0, io.SeekCurrent)
			require.NoError(t, err)
			assert.Equal(t, int64(buf.Len()), s)

			// Stringify
			str, err := Stringify(tc.src, tc.ctx)
			require.NoError(t, err)
			assert.Equal(t, tc.str, str)
		})
	}
}
