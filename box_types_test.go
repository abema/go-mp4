package mp4

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestElstMarshal(t *testing.T) {
	testCases := []struct {
		name string
		elst Elst
		want []byte
	}{
		{
			name: "version 0",
			elst: Elst{
				FullBox: FullBox{
					Version: 0,
					Flags:   [3]byte{0x00, 0x00, 0x00},
				},
				EntryCount: 2,
				Entries: []ElstEntry{
					{
						SegmentDurationV0: 0x0100000a,
						MediaTimeV0:       0x0100000b,
						MediaRateInteger:  0x010c,
						MediaRateFraction: 0x010d,
					}, {
						SegmentDurationV0: 0x0200000a,
						MediaTimeV0:       0x0200000b,
						MediaRateInteger:  0x020c,
						MediaRateFraction: 0x020d,
					},
				},
			},
			want: []byte{
				0,                // version
				0x00, 0x00, 0x00, // flags
				0x00, 0x00, 0x00, 0x02, // entry count
				0x01, 0x00, 0x00, 0x0a, // segment duration v0
				0x01, 0x00, 0x00, 0x0b, // media time v0
				0x01, 0x0c, // media rate integer
				0x01, 0x0d, // media rate fraction
				0x02, 0x00, 0x00, 0x0a, // segment duration v0
				0x02, 0x00, 0x00, 0x0b, // media time v0
				0x02, 0x0c, // media rate integer
				0x02, 0x0d, // media rate fraction
			},
		},
		{
			name: "version 1",
			elst: Elst{
				FullBox: FullBox{
					Version: 1,
					Flags:   [3]byte{0x00, 0x00, 0x00},
				},
				EntryCount: 2,
				Entries: []ElstEntry{
					{
						SegmentDurationV1: 0x010000000000000a,
						MediaTimeV1:       0x010000000000000b,
						MediaRateInteger:  0x010c,
						MediaRateFraction: 0x010d,
					}, {
						SegmentDurationV1: 0x020000000000000a,
						MediaTimeV1:       0x020000000000000b,
						MediaRateInteger:  0x020c,
						MediaRateFraction: 0x020d,
					},
				},
			},
			want: []byte{
				1,                // version
				0x00, 0x00, 0x00, // flags
				0x00, 0x00, 0x00, 0x02, // entry count
				0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0a, // segment duration v1
				0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0b, // media time v1
				0x01, 0x0c, // media rate integer
				0x01, 0x0d, // media rate fraction
				0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0a, // segment duration v1
				0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0b, // media time v1
				0x02, 0x0c, // media rate integer
				0x02, 0x0d, // media rate fraction
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			n, err := Marshal(buf, &c.elst)
			require.NoError(t, err)
			assert.Equal(t, uint64(len(c.want)), n)
			assert.Equal(t, c.want, buf.Bytes())
		})
	}
}

func TestEmsgMarshal(t *testing.T) {
	emsg := Emsg{
		FullBox: FullBox{
			Version: 0,
			Flags:   [3]byte{0x00, 0x00, 0x00},
		},
		SchemeIdUri:           "urn:test",
		Value:                 "foo",
		Timescale:             1000,
		PresentationTimeDelta: 123,
		EventDuration:         3000,
		Id:                    0xabcd,
		MessageData:           []byte("abema"),
	}
	expect := []byte{
		0,                // version
		0x00, 0x00, 0x00, // flags
		0x75, 0x72, 0x6e, 0x3a, 0x74, 0x65, 0x73, 0x74, 0x00, // scheme id uri
		0x66, 0x6f, 0x6f, 0x00, // value
		0x00, 0x00, 0x03, 0xe8, // timescale
		0x00, 0x00, 0x00, 0x7b, // presentation time delta
		0x00, 0x00, 0x0b, 0xb8, // event duration
		0x00, 0x00, 0xab, 0xcd, // id
		0x61, 0x62, 0x65, 0x6d, 0x61, // message data
	}

	buf := bytes.NewBuffer(nil)
	n, err := Marshal(buf, &emsg)
	require.NoError(t, err)
	assert.Equal(t, uint64(len(expect)), n)
	assert.Equal(t, expect, buf.Bytes())
}

func TestEmsgUnmarshal(t *testing.T) {
	data := []byte{
		0,                // version
		0x00, 0x00, 0x00, // flags
		0x75, 0x72, 0x6e, 0x3a, 0x74, 0x65, 0x73, 0x74, 0x00, // scheme id uri
		0x66, 0x6f, 0x6f, 0x00, // value
		0x00, 0x00, 0x03, 0xe8, // timescale
		0x00, 0x00, 0x00, 0x7b, // presentation time delta
		0x00, 0x00, 0x0b, 0xb8, // event duration
		0x00, 0x00, 0xab, 0xcd, // id
		0x61, 0x62, 0x65, 0x6d, 0x61, // message data
	}

	buf := bytes.NewReader(data)
	emsg := Emsg{}
	n, err := Unmarshal(buf, uint64(buf.Len()), &emsg)
	require.NoError(t, err)
	assert.Equal(t, uint64(len(data)), n)
	assert.Equal(t, uint8(0), emsg.Version)
	assert.Equal(t, uint32(0x000000), emsg.GetFlags())
	assert.Equal(t, "urn:test", emsg.SchemeIdUri)
	assert.Equal(t, "foo", emsg.Value)
	assert.Equal(t, uint32(1000), emsg.Timescale)
	assert.Equal(t, uint32(123), emsg.PresentationTimeDelta)
	assert.Equal(t, uint32(3000), emsg.EventDuration)
	assert.Equal(t, uint32(0xabcd), emsg.Id)
	assert.Equal(t, []byte("abema"), emsg.MessageData)

	buf = bytes.NewReader(data)
	result, n, err := UnmarshalAny(buf, BoxTypeEmsg(), uint64(len(data)))
	require.NoError(t, err)
	assert.Equal(t, uint64(len(data)), n)
	pemsg, ok := result.(*Emsg)
	require.True(t, ok)
	assert.Equal(t, uint8(0), pemsg.Version)
	assert.Equal(t, uint32(0x000000), pemsg.GetFlags())
	assert.Equal(t, "urn:test", pemsg.SchemeIdUri)
	assert.Equal(t, "foo", pemsg.Value)
	assert.Equal(t, uint32(1000), pemsg.Timescale)
	assert.Equal(t, uint32(123), pemsg.PresentationTimeDelta)
	assert.Equal(t, uint32(3000), pemsg.EventDuration)
	assert.Equal(t, uint32(0xabcd), pemsg.Id)
	assert.Equal(t, []byte("abema"), pemsg.MessageData)
}

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

func TestHdlrUnmarshalHandlerName(t *testing.T) {
	testCases := []struct {
		name          string
		componentType []byte
		bytes         []byte
		want          string
	}{
		{
			name:          "NormalString",
			componentType: []byte{0x00, 0x00, 0x00, 0x00},
			bytes:         []byte("abema"),
			want:          "abema",
		},
		{
			name:          "EmptyString",
			componentType: []byte{0x00, 0x00, 0x00, 0x00},
			bytes:         nil,
			want:          "",
		},
		{
			name:          "NormalLongString",
			componentType: []byte{0x00, 0x00, 0x00, 0x00},
			bytes:         []byte(" a 1st byte equals to this length"),
			want:          " a 1st byte equals to this length",
		},
		{
			name:          "AppleQuickTimePascalString",
			componentType: []byte("mhlr"),
			bytes:         []byte{5, 'a', 'b', 'e', 'm', 'a'},
			want:          "abema",
		},
		{
			name:          "AppleQuickTimePascalStringWithEmpty",
			componentType: []byte("mhlr"),
			bytes:         []byte{0x00},
			want:          "",
		},
		{
			name:          "AppleQuickTimePascalStringLong",
			componentType: []byte("mhlr"),
			bytes:         []byte(" a 1st byte equals to this length"),
			want:          "a 1st byte equals to this length",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bin := []byte{
				0,                // version
				0x00, 0x00, 0x00, // flags
			}
			bin = append(bin, tc.componentType...)
			bin = append(bin,
				'v', 'i', 'd', 'e', // handler type
				0x00, 0x00, 0x00, 0x00, // reserved
				0x00, 0x00, 0x00, 0x00, // reserved
				0x00, 0x00, 0x00, 0x00, // reserved
			)
			bin = append(bin, tc.bytes...)

			// unmarshal
			dst := Hdlr{}
			r := bytes.NewReader(bin)
			n, err := Unmarshal(r, uint64(len(bin)), &dst)
			assert.NoError(t, err)
			assert.Equal(t, uint64(len(bin)), n)
			assert.Equal(t, [4]byte{'v', 'i', 'd', 'e'}, dst.HandlerType)
			assert.Equal(t, tc.want, dst.Name)
		})
	}
}

func TestMdhdMarshal(t *testing.T) {
	// Version 0
	mdhd := Mdhd{
		FullBox: FullBox{
			Version: 0,
			Flags:   [3]byte{0x00, 0x00, 0x00},
		},
		CreationTimeV0:     0x12345678,
		ModificationTimeV0: 0x23456789,
		TimescaleV0:        0x01020304,
		DurationV0:         0x02030405,
		Pad:                true,
		Language:           [3]byte{'j' - 0x60, 'p' - 0x60, 'n' - 0x60}, // 0x0a, 0x10, 0x0e
		PreDefined:         0,
	}
	expect := []byte{
		0,                // version
		0x00, 0x00, 0x00, // flags
		0x12, 0x34, 0x56, 0x78, // creation time
		0x23, 0x45, 0x67, 0x89, // modification time
		0x01, 0x02, 0x03, 0x04, // timescale
		0x02, 0x03, 0x04, 0x05, // duration
		0xaa, 0x0e, // pad, language (1 01010 10000 01110)
		0x00, 0x00, // pre defined
	}

	buf := bytes.NewBuffer(nil)
	n, err := Marshal(buf, &mdhd)
	require.NoError(t, err)
	assert.Equal(t, uint64(len(expect)), n)
	assert.Equal(t, expect, buf.Bytes())

	// Version 1
	mdhd = Mdhd{
		FullBox: FullBox{
			Version: 1,
			Flags:   [3]byte{0x00, 0x00, 0x00},
		},
		CreationTimeV1:     0x123456789abcdef0,
		ModificationTimeV1: 0x23456789abcdef01,
		TimescaleV1:        0x01020304,
		DurationV1:         0x0203040506070809,
		Pad:                true,
		Language:           [3]byte{'j' - 0x60, 'p' - 0x60, 'n' - 0x60}, // 0x0a, 0x10, 0x0e
		PreDefined:         0,
	}
	expect = []byte{
		1,                // version
		0x00, 0x00, 0x00, // flags
		0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0, // creation time
		0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x01, // modification time
		0x01, 0x02, 0x03, 0x04, // timescale
		0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, // duration
		0xaa, 0x0e, // pad, language (1 01010 10000 01110)
		0x00, 0x00, // pre defined
	}

	buf = bytes.NewBuffer(nil)
	n, err = Marshal(buf, &mdhd)
	require.NoError(t, err)
	assert.Equal(t, uint64(len(expect)), n)
	assert.Equal(t, expect, buf.Bytes())
}

func TestMdhdUnmarshal(t *testing.T) {
	// Version 0
	data := []byte{
		0,                // version
		0x00, 0x00, 0x00, // flags
		0x12, 0x34, 0x56, 0x78, // creation time
		0x23, 0x45, 0x67, 0x89, // modification time
		0x34, 0x56, 0x78, 0x9a, // timescale
		0x45, 0x67, 0x89, 0xab, // duration
		0xaa, 0x67, // pad, language (1 01010 10011 00111)
		0x00, 0x00, // pre defined
	}

	buf := bytes.NewReader(data)
	mdhd := Mdhd{}
	n, err := Unmarshal(buf, uint64(buf.Len()), &mdhd)
	require.NoError(t, err)
	assert.Equal(t, uint64(len(data)), n)
	assert.Equal(t, uint8(0), mdhd.Version)
	assert.Equal(t, uint32(0x0), mdhd.GetFlags())
	assert.Equal(t, uint32(0x12345678), mdhd.CreationTimeV0)
	assert.Equal(t, uint32(0x23456789), mdhd.ModificationTimeV0)
	assert.Equal(t, uint32(0x3456789a), mdhd.TimescaleV0)
	assert.Equal(t, uint32(0x456789ab), mdhd.DurationV0)
	assert.Equal(t, uint64(0x0), mdhd.CreationTimeV1)
	assert.Equal(t, uint64(0x0), mdhd.ModificationTimeV1)
	assert.Equal(t, uint32(0x0), mdhd.TimescaleV1)
	assert.Equal(t, uint64(0x0), mdhd.DurationV1)
	assert.Equal(t, true, mdhd.Pad)
	assert.Equal(t, byte(0x0a), mdhd.Language[0])
	assert.Equal(t, byte(0x13), mdhd.Language[1])
	assert.Equal(t, byte(0x07), mdhd.Language[2])
	assert.Equal(t, uint16(0), mdhd.PreDefined)

	// Version 1
	data = []byte{
		1,                // version
		0x00, 0x00, 0x00, // flags
		0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0, // creation time
		0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x01, // modification time
		0x34, 0x56, 0x78, 0x9a, // timescale
		0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x01, 0x23, // duration
		0xaa, 0x67, // pad, language (1 01010 10011 00111)
		0x00, 0x00, // pre defined
	}

	buf = bytes.NewReader(data)
	mdhd = Mdhd{}
	n, err = Unmarshal(buf, uint64(buf.Len()), &mdhd)
	require.NoError(t, err)
	assert.Equal(t, uint64(len(data)), n)
	assert.Equal(t, uint8(1), mdhd.Version)
	assert.Equal(t, uint32(0x0), mdhd.GetFlags())
	assert.Equal(t, uint32(0x0), mdhd.CreationTimeV0)
	assert.Equal(t, uint32(0x0), mdhd.ModificationTimeV0)
	assert.Equal(t, uint32(0x0), mdhd.TimescaleV0)
	assert.Equal(t, uint32(0x0), mdhd.DurationV0)
	assert.Equal(t, uint64(0x123456789abcdef0), mdhd.CreationTimeV1)
	assert.Equal(t, uint64(0x23456789abcdef01), mdhd.ModificationTimeV1)
	assert.Equal(t, uint32(0x3456789a), mdhd.TimescaleV1)
	assert.Equal(t, uint64(0x456789abcdef0123), mdhd.DurationV1)
	assert.Equal(t, true, mdhd.Pad)
	assert.Equal(t, byte(0x0a), mdhd.Language[0])
	assert.Equal(t, byte(0x13), mdhd.Language[1])
	assert.Equal(t, byte(0x07), mdhd.Language[2])
	assert.Equal(t, uint16(0), mdhd.PreDefined)
}

func TestMetaMarshal(t *testing.T) {
	src := Meta{
		FullBox: FullBox{
			Version: 0,
			Flags:   [3]byte{0x00, 0x00, 0x00},
		},
	}
	bin := []byte{
		0,                // version
		0x00, 0x00, 0x00, // flags
	}

	// marshal
	buf := bytes.NewBuffer(nil)
	n, err := Marshal(buf, &src)
	require.NoError(t, err)
	assert.Equal(t, uint64(len(bin)), n)
	assert.Equal(t, bin, buf.Bytes())

	// unmarshal
	dst := Meta{}
	n, err = Unmarshal(bytes.NewReader(bin), uint64(len(bin)+8), &dst)
	assert.NoError(t, err)
	assert.Equal(t, uint64(buf.Len()), n)
	assert.Equal(t, src, dst)
}

func TestMetaMarshalAppleQuickTime(t *testing.T) {
	bin := []byte{
		0x00, 0x00, 0x01, 0x00, // size
		'h', 'd', 'l', 'r', // type
		0,                // version
		0x00, 0x00, 0x00, // flags
	}

	// unmarshal
	dst := Meta{}
	r := bytes.NewReader(bin)
	n, err := Unmarshal(r, uint64(len(bin)+8), &dst)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), n)
	s, _ := r.Seek(0, io.SeekCurrent)
	assert.Equal(t, int64(0), s)
	assert.Equal(t, uint8(0), dst.GetVersion())
	assert.Equal(t, uint32(0), dst.GetFlags())
}

func TestPsshStringify(t *testing.T) {
	flags := [3]byte{0x00, 0x00, 0x00}
	systemID := [16]byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
	}
	kid1 := PsshKID{KID: [16]byte{
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x10,
	}}
	kid2 := PsshKID{KID: [16]byte{
		0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
		0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x20,
	}}
	data := []byte{0x21, 0x22, 0x23, 0x24, 0x25}

	testCases := []struct {
		name string
		pssh Pssh
		want string
	}{
		{
			name: "version 0: no KIDs",
			pssh: Pssh{
				FullBox: FullBox{
					Version: 0,
					Flags:   flags,
				},
				SystemID: systemID,
				DataSize: int32(len(data)),
				Data:     data,
			},
			want: `Version=0 ` +
				`Flags=0x000000 ` +
				`SystemID="0102030405060708090a0b0c0d0e0f10" ` +
				`DataSize=5 ` +
				`Data=[0x21, 0x22, 0x23, 0x24, 0x25]`,
		},
		{
			name: "version 1: with KIDs",
			pssh: Pssh{
				FullBox: FullBox{
					Version: 1,
					Flags:   flags,
				},
				SystemID: systemID,
				KIDCount: 2,
				KIDs:     []PsshKID{kid1, kid2},
				DataSize: int32(len(data)),
				Data:     data,
			},
			want: `Version=1 ` +
				`Flags=0x000000 ` +
				`SystemID="0102030405060708090a0b0c0d0e0f10" ` +
				`KIDCount=2 ` +
				`KIDs=["1112131415161718191a1b1c1d1e1f10" "2122232425262728292a2b2c2d2e2f20"] ` +
				`DataSize=5 ` +
				`Data=[0x21, 0x22, 0x23, 0x24, 0x25]`,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			str, err := Stringify(&c.pssh)
			require.NoError(t, err)
			assert.Equal(t, c.want, str)
		})
	}
}

func TestVisualSampleEntryStringify(t *testing.T) {
	vse := VisualSampleEntry{
		SampleEntry:     SampleEntry{},
		PreDefined:      0x0101,
		PreDefined2:     [3]uint32{0x01000001, 0x01000002, 0x01000003},
		Width:           0x0102,
		Height:          0x0103,
		Horizresolution: 0x01000004,
		Vertresolution:  0x01000005,
		Reserved2:       0x01000006,
		FrameCount:      0x0104,
		Compressorname:  [32]byte{'a', 'b', 'e', 'm', 'a'},
		Depth:           0x0105,
		PreDefined3:     1001,
	}
	str, err := Stringify(&vse)
	require.NoError(t, err)
	assert.Equal(t, "DataReferenceIndex=0"+
		" PreDefined=257"+
		" PreDefined2=[16777217,"+
		" 16777218,"+
		" 16777219]"+
		" Width=258"+
		" Height=259"+
		" Horizresolution=16777220"+
		" Vertresolution=16777221"+
		" FrameCount=260"+
		" Compressorname=\"abema\""+
		" Depth=261"+
		" PreDefined3=1001", str)
}

func TestTfhdMarshal(t *testing.T) {
	tfhd := Tfhd{
		FullBox: FullBox{
			Version: 0,
			Flags:   [3]byte{0x00, 0x00, 0x00},
		},
		TrackID:                0x08404649,
		BaseDataOffset:         0x0123456789abcdef,
		SampleDescriptionIndex: 0x12345678,
		DefaultSampleDuration:  0x23456789,
		DefaultSampleSize:      0x3456789a,
		DefaultSampleFlags:     0x456789ab,
	}
	expect := []byte{
		0,                // version
		0x00, 0x00, 0x00, // flags
		0x08, 0x40, 0x46, 0x49, // track ID
	}

	buf := bytes.NewBuffer(nil)
	n, err := Marshal(buf, &tfhd)
	require.NoError(t, err)
	assert.Equal(t, uint64(len(expect)), n)
	assert.Equal(t, expect, buf.Bytes())

	tfhd.SetFlags(TfhdBaseDataOffsetPresent | TfhdDefaultSampleDurationPresent)
	expect = []byte{
		0,                // version
		0x00, 0x00, 0x09, // flags (0000 0000 1001)
		0x08, 0x40, 0x46, 0x49, // track ID
		0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
		0x23, 0x45, 0x67, 0x89,
	}

	buf = bytes.NewBuffer(nil)
	n, err = Marshal(buf, &tfhd)
	require.NoError(t, err)
	assert.Equal(t, uint64(len(expect)), n)
	assert.Equal(t, expect, buf.Bytes())
}

func TestTfhdUnmarshal(t *testing.T) {
	data := []byte{
		0,                // version
		0x00, 0x00, 0x00, // flags (0000 0000 1001)
		0x08, 0x40, 0x46, 0x49, // track ID
	}

	buf := bytes.NewReader(data)
	tfhd := Tfhd{}
	n, err := Unmarshal(buf, uint64(buf.Len()), &tfhd)
	require.NoError(t, err)
	assert.Equal(t, uint64(len(data)), n)
	assert.Equal(t, uint8(0), tfhd.Version)
	assert.Equal(t, uint32(0x00), tfhd.GetFlags())
	assert.Equal(t, uint32(0x08404649), tfhd.TrackID)

	data = []byte{
		0,                // version
		0x00, 0x00, 0x09, // flags (0000 0000 1001)
		0x08, 0x40, 0x46, 0x49, // track ID
		0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
		0x23, 0x45, 0x67, 0x89,
	}

	buf = bytes.NewReader(data)
	tfhd = Tfhd{}
	n, err = Unmarshal(buf, uint64(buf.Len()), &tfhd)
	require.NoError(t, err)
	assert.Equal(t, uint64(len(data)), n)
	assert.Equal(t, uint8(0), tfhd.Version)
	assert.Equal(t, uint32(0x09), tfhd.GetFlags())
	assert.Equal(t, uint32(0x08404649), tfhd.TrackID)
	assert.Equal(t, uint64(0x0123456789abcdef), tfhd.BaseDataOffset)
	assert.Equal(t, uint32(0x23456789), tfhd.DefaultSampleDuration)
}

func TestTkhdMarshal(t *testing.T) {
	// Version 0
	tkhd := Tkhd{
		FullBox: FullBox{
			Version: 0,
			Flags:   [3]byte{0x00, 0x00, 0x00},
		},
		CreationTimeV0:     0x01234567,
		ModificationTimeV0: 0x12345678,
		TrackIDV0:          0x23456789,
		ReservedV0:         0x3456789a,
		DurationV0:         0x456789ab,
		Reserved:           [2]uint32{0, 0},
		Layer:              23456,  // 0x5ba0
		AlternateGroup:     -23456, // 0xdba0
		Volume:             0x0100,
		Reserved2:          0,
		Matrix: [9]int32{
			0x00010000, 0, 0,
			0, 0x00010000, 0,
			0, 0, 0x40000000,
		},
		Width:  0x56789abc,
		Height: 0x6789abcd,
	}
	expect := []byte{
		0,                // version
		0x00, 0x00, 0x00, // flags
		0x01, 0x23, 0x45, 0x67, // creation time
		0x12, 0x34, 0x56, 0x78, // modification time
		0x23, 0x45, 0x67, 0x89, // track ID
		0x34, 0x56, 0x78, 0x9a, // reserved
		0x45, 0x67, 0x89, 0xab, // duration
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // reserved
		0x5b, 0xa0, // layer
		0xdb, 0xa0, // alternate group
		0x01, 0x00, // volume
		0x00, 0x00, // reserved
		0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00, // matrix
		0x56, 0x78, 0x9a, 0xbc, // width
		0x67, 0x89, 0xab, 0xcd, // height
	}

	buf := bytes.NewBuffer(nil)
	n, err := Marshal(buf, &tkhd)
	require.NoError(t, err)
	assert.Equal(t, uint64(len(expect)), n)
	assert.Equal(t, expect, buf.Bytes())
}

func TestTkhdUnmarshal(t *testing.T) {
	// Version 0
	data := []byte{
		0,                // version
		0x00, 0x00, 0x00, // flags (0000 0000 1001)
		0x01, 0x23, 0x45, 0x67, // creation time
		0x12, 0x34, 0x56, 0x78, // modification time
		0x23, 0x45, 0x67, 0x89, // track ID
		0x34, 0x56, 0x78, 0x9a, // reserved
		0x45, 0x67, 0x89, 0xab, // duration
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // reserved
		0x5b, 0xa0, // layer
		0xdb, 0xa0, // alternate group
		0x01, 0x00, // volume
		0x00, 0x00, // reserved
		0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00, // matrix
		0x56, 0x78, 0x9a, 0xbc, // width
		0x67, 0x89, 0xab, 0xcd, // height
	}

	buf := bytes.NewReader(data)
	tkhd := Tkhd{}
	n, err := Unmarshal(buf, uint64(buf.Len()), &tkhd)
	require.NoError(t, err)
	assert.Equal(t, uint64(len(data)), n)
	assert.Equal(t, uint8(0), tkhd.Version)
	assert.Equal(t, uint32(0x00), tkhd.GetFlags())
	assert.Equal(t, uint32(0x01234567), tkhd.CreationTimeV0)
	assert.Equal(t, uint32(0x12345678), tkhd.ModificationTimeV0)
	assert.Equal(t, uint32(0x23456789), tkhd.TrackIDV0)
	assert.Equal(t, uint32(0x3456789a), tkhd.ReservedV0)
	assert.Equal(t, uint32(0x456789ab), tkhd.DurationV0)
	assert.Equal(t, [2]uint32{0, 0}, tkhd.Reserved)
	assert.Equal(t, int16(23456), tkhd.Layer)
	assert.Equal(t, int16(-23456), tkhd.AlternateGroup)
	assert.Equal(t, int16(0x0100), tkhd.Volume)
	assert.Equal(t, uint16(0), tkhd.Reserved2)
	assert.Equal(t, [9]int32{
		0x00010000, 0, 0,
		0, 0x00010000, 0,
		0, 0, 0x40000000,
	}, tkhd.Matrix)
	assert.Equal(t, uint32(0x56789abc), tkhd.Width)
	assert.Equal(t, uint32(0x6789abcd), tkhd.Height)
}
