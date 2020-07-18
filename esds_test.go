package mp4

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEsdsMarshal(t *testing.T) {
	src := Esds{
		FullBox: FullBox{
			Version: 0,
			Flags:   [3]byte{0x00, 0x00, 0x00},
		},
		Descriptors: []Descriptor{
			{
				Tag:  ESDescrTag,
				Size: 0x12345678,
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
				Size: 0x12345678,
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
				Size: 0x12345678,
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
		},
	}
	bin := []byte{
		0,                // version
		0x00, 0x00, 0x00, // flags
		//
		0x03,                         // tag
		0x81, 0x91, 0xD1, 0xAC, 0x78, // size (varint)
		0x12, 0x34, // esid
		0xa3,       // flags & streamPriority
		0x23, 0x45, // dependsOnESID
		0x34, 0x56, // ocresid
		//
		0x03,                         // tag
		0x81, 0x91, 0xD1, 0xAC, 0x78, // size (varint)
		0x12, 0x34, // esid
		0x43,                                                  // flags & streamPriority
		11,                                                    // urlLength
		'h', 't', 't', 'p', ':', '/', '/', 'h', 'o', 'g', 'e', // urlString
		//
		0x04,                         // tag
		0x81, 0x91, 0xD1, 0xAC, 0x78, // size (varint)
		0x12,             // objectTypeIndication
		0x56,             // streamType & upStream & reserved
		0x12, 0x34, 0x56, // bufferSizeDB
		0x12, 0x34, 0x56, 0x78, // maxBitrate
		0x23, 0x45, 0x67, 0x89, // avgBitrate
	}

	// marshal
	buf := bytes.NewBuffer(nil)
	n, err := Marshal(buf, &src)
	require.NoError(t, err)
	assert.Equal(t, uint64(len(bin)), n)
	assert.Equal(t, bin, buf.Bytes())

	// unmarshal
	dst := Esds{}
	n, err = Unmarshal(bytes.NewReader(bin), uint64(len(bin)), &dst)
	assert.NoError(t, err)
	assert.Equal(t, uint64(buf.Len()), n)
	assert.Equal(t, src, dst)
}

func TestEsdsStringify(t *testing.T) {
	esds := Esds{
		FullBox: FullBox{
			Version: 0,
			Flags:   [3]byte{0x00, 0x00, 0x00},
		},
		Descriptors: []Descriptor{
			{
				Tag:  ESDescrTag,
				Size: 0x12345678,
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
				Tag:  DecoderConfigDescrTag,
				Size: 0x12345678,
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
				Size: 0x12345678,
			},
			{
				Tag:  SLConfigDescrTag,
				Size: 0x12345678,
			},
		},
	}
	str, err := Stringify(&esds)
	require.NoError(t, err)
	assert.Equal(t, "Version=0 Flags=0x000000 Descriptors="+
		"[{Tag=ESDescr Size=305419896 ESID=4660 StreamDependenceFlag=true UrlFlag=false OcrStreamFlag=true StreamPriority=3 DependsOnESID=9029 OCRESID=13398},"+
		" {Tag=DecoderConfigDescr Size=305419896 ObjectTypeIndication=0x12 StreamType=21 UpStream=true Reserved=false BufferSizeDB=1193046 MaxBitrate=305419896 AvgBitrate=591751049},"+
		" {Tag=DecSpecificInfo Size=305419896 Data=[]},"+
		" {Tag=SLConfigDescr Size=305419896 Data=[]}"+
		"]", str)
}
