package mp4

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoxTypesISO14496_14(t *testing.T) {
	testCases := []struct {
		name string
		src  IImmutableBox
		dst  IBox
		bin  []byte
		str  string
		ctx  Context
	}{
		{
			name: "esds",
			src: &Esds{
				FullBox: FullBox{
					Version: 0,
					Flags:   [3]byte{0x00, 0x00, 0x00},
				},
				Descriptors: []Descriptor{
					{
						Tag:  ESDescrTag,
						Size: 0x1234567,
						ESDescriptor: &ESDescriptor{
							ESID:                 0x1234,
							StreamDependenceFlag: true,
							UrlFlag:              false,
							OcrStreamFlag:        true,
							StreamPriority:       0x03,
							DependsOnESID:        0x2345,
							OCRESID:              0x3456,
						},
					},
					{
						Tag:  ESDescrTag,
						Size: 0x1234567,
						ESDescriptor: &ESDescriptor{
							ESID:                 0x1234,
							StreamDependenceFlag: false,
							UrlFlag:              true,
							OcrStreamFlag:        false,
							StreamPriority:       0x03,
							URLLength:            11,
							URLString:            []byte("http://hoge"),
						},
					},
					{
						Tag:  DecoderConfigDescrTag,
						Size: 0x1234567,
						DecoderConfigDescriptor: &DecoderConfigDescriptor{
							ObjectTypeIndication: 0x12,
							StreamType:           0x15,
							UpStream:             true,
							Reserved:             false,
							BufferSizeDB:         0x123456,
							MaxBitrate:           0x12345678,
							AvgBitrate:           0x23456789,
						},
					},
					{
						Tag:  DecSpecificInfoTag,
						Size: 0x03,
						Data: []byte{0x11, 0x22, 0x33},
					},
					{
						Tag:  SLConfigDescrTag,
						Size: 0x05,
						Data: []byte{0x11, 0x22, 0x33, 0x44, 0x55},
					},
				},
			},
			dst: &Esds{},
			bin: []byte{
				0,                // version
				0x00, 0x00, 0x00, // flags
				//
				0x03,                   // tag
				0x89, 0x8d, 0x8a, 0x67, // size (varint)
				0x12, 0x34, // esid
				0xa3,       // flags & streamPriority
				0x23, 0x45, // dependsOnESID
				0x34, 0x56, // ocresid
				//
				0x03,                   // tag
				0x89, 0x8d, 0x8a, 0x67, // size (varint)
				0x12, 0x34, // esid
				0x43,                                                  // flags & streamPriority
				11,                                                    // urlLength
				'h', 't', 't', 'p', ':', '/', '/', 'h', 'o', 'g', 'e', // urlString
				//
				0x04,                   // tag
				0x89, 0x8d, 0x8a, 0x67, // size (varint)
				0x12,             // objectTypeIndication
				0x56,             // streamType & upStream & reserved
				0x12, 0x34, 0x56, // bufferSizeDB
				0x12, 0x34, 0x56, 0x78, // maxBitrate
				0x23, 0x45, 0x67, 0x89, // avgBitrate
				//
				0x05,                   // tag
				0x80, 0x80, 0x80, 0x03, // size (varint)
				0x11, 0x22, 0x33, // data
				//
				0x06,                   // tag
				0x80, 0x80, 0x80, 0x05, // size (varint)
				0x11, 0x22, 0x33, 0x44, 0x55, // data
			},
			str: `Version=0 Flags=0x000000 Descriptors=[` +
				`{Tag=ESDescr Size=19088743 ESID=4660 StreamDependenceFlag=true UrlFlag=false OcrStreamFlag=true StreamPriority=3 DependsOnESID=9029 OCRESID=13398}, ` +
				`{Tag=ESDescr Size=19088743 ESID=4660 StreamDependenceFlag=false UrlFlag=true OcrStreamFlag=false StreamPriority=3 URLLength=0xb URLString="http://hoge"}, ` +
				`{Tag=DecoderConfigDescr Size=19088743 ObjectTypeIndication=0x12 StreamType=21 UpStream=true Reserved=false BufferSizeDB=1193046 MaxBitrate=305419896 AvgBitrate=591751049}, ` +
				"{Tag=DecSpecificInfo Size=3 Data=[0x11, 0x22, 0x33]}, " +
				"{Tag=SLConfigDescr Size=5 Data=[0x11, 0x22, 0x33, 0x44, 0x55]}]",
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
