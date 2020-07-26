package mp4

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoxTypes(t *testing.T) {
	testCases := []struct {
		name string
		src  IImmutableBox
		dst  IBox
		bin  []byte
		str  string
	}{
		{
			name: "elst: version 0",
			src: &Elst{
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
			dst: &Elst{},
			bin: []byte{
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
			str: `Version=0 Flags=0x000000 EntryCount=2 ` +
				`Entries=[{SegmentDurationV0=16777226 MediaTimeV0=16777227 MediaRateInteger=268}, ` +
				`{SegmentDurationV0=33554442 MediaTimeV0=33554443 MediaRateInteger=524}]`,
		},
		{
			name: "elst: version 1",
			src: &Elst{
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
			dst: &Elst{},
			bin: []byte{
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
			str: `Version=1 Flags=0x000000 EntryCount=2 ` +
				`Entries=[{SegmentDurationV1=72057594037927946 MediaTimeV1=72057594037927947 MediaRateInteger=268}, ` +
				`{SegmentDurationV1=144115188075855882 MediaTimeV1=144115188075855883 MediaRateInteger=524}]`,
		},
		{
			name: "emsg",
			src: &Emsg{
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
			},
			dst: &Emsg{},
			bin: []byte{
				0,                // version
				0x00, 0x00, 0x00, // flags
				0x75, 0x72, 0x6e, 0x3a, 0x74, 0x65, 0x73, 0x74, 0x00, // scheme id uri
				0x66, 0x6f, 0x6f, 0x00, // value
				0x00, 0x00, 0x03, 0xe8, // timescale
				0x00, 0x00, 0x00, 0x7b, // presentation time delta
				0x00, 0x00, 0x0b, 0xb8, // event duration
				0x00, 0x00, 0xab, 0xcd, // id
				0x61, 0x62, 0x65, 0x6d, 0x61, // message data
			},
			str: `Version=0 Flags=0x000000 ` +
				`SchemeIdUri="urn:test" ` +
				`Value="foo" ` +
				`Timescale=1000 ` +
				`PresentationTimeDelta=123 ` +
				`EventDuration=3000 ` +
				`Id=43981 ` +
				`MessageData="abema"`,
		},
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
				//
				0x05,             // tag
				0x03,             // size (varint)
				0x11, 0x22, 0x33, // data
				//
				0x06,                         // tag
				0x05,                         // size (varint)
				0x11, 0x22, 0x33, 0x44, 0x55, // data
			},
			str: `Version=0 Flags=0x000000 Descriptors=[` +
				`{Tag=ESDescr Size=305419896 ESID=4660 StreamDependenceFlag=true UrlFlag=false OcrStreamFlag=true StreamPriority=3 DependsOnESID=9029 OCRESID=13398}, ` +
				`{Tag=ESDescr Size=305419896 ESID=4660 StreamDependenceFlag=false UrlFlag=true OcrStreamFlag=false StreamPriority=3 URLLength=0xb URLString="http://hoge"}, ` +
				`{Tag=DecoderConfigDescr Size=305419896 ObjectTypeIndication=0x12 StreamType=21 UpStream=true Reserved=false BufferSizeDB=1193046 MaxBitrate=305419896 AvgBitrate=591751049}, ` +
				"{Tag=DecSpecificInfo Size=3 Data=[0x11, 0x22, 0x33]}, " +
				"{Tag=SLConfigDescr Size=5 Data=[0x11, 0x22, 0x33, 0x44, 0x55]}]",
		},
		{
			name: "mdhd: version 0",
			src: &Mdhd{
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
			},
			dst: &Mdhd{},
			bin: []byte{
				0,                // version
				0x00, 0x00, 0x00, // flags
				0x12, 0x34, 0x56, 0x78, // creation time
				0x23, 0x45, 0x67, 0x89, // modification time
				0x01, 0x02, 0x03, 0x04, // timescale
				0x02, 0x03, 0x04, 0x05, // duration
				0xaa, 0x0e, // pad, language (1 01010 10000 01110)
				0x00, 0x00, // pre defined
			},
			str: `Version=0 Flags=0x000000 ` +
				`CreationTimeV0=305419896 ` +
				`ModificationTimeV0=591751049 ` +
				`TimescaleV0=16909060 ` +
				`DurationV0=33752069 ` +
				`Pad=true ` +
				`Language="jpn" ` +
				`PreDefined=0`,
		},
		{
			name: "mdhd: version 1",
			src: &Mdhd{
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
			},
			dst: &Mdhd{},
			bin: []byte{
				1,                // version
				0x00, 0x00, 0x00, // flags
				0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0, // creation time
				0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x01, // modification time
				0x01, 0x02, 0x03, 0x04, // timescale
				0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, // duration
				0xaa, 0x0e, // pad, language (1 01010 10000 01110)
				0x00, 0x00, // pre defined
			},
			str: `Version=1 Flags=0x000000 ` +
				`CreationTimeV1=1311768467463790320 ` +
				`ModificationTimeV1=2541551405711093505 ` +
				`TimescaleV1=16909060 ` +
				`DurationV1=144964032628459529 ` +
				`Pad=true ` +
				`Language="jpn" ` +
				`PreDefined=0`,
		},
		{
			name: "meta",
			src: &Meta{
				FullBox: FullBox{
					Version: 0,
					Flags:   [3]byte{0x00, 0x00, 0x00},
				},
			},
			dst: &Meta{},
			bin: []byte{
				0,                // version
				0x00, 0x00, 0x00, // flags
			},
			str: `Version=0 Flags=0x000000`,
		},
		{
			name: "pssh: version 0: no KIDs",
			src: &Pssh{
				FullBox: FullBox{
					Version: 0,
					Flags:   [3]byte{0x00, 0x00, 0x00},
				},
				SystemID: [16]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10},
				DataSize: 5,
				Data:     []byte{0x21, 0x22, 0x23, 0x24, 0x25},
			},
			dst: &Pssh{},
			bin: []byte{
				0,                // version
				0x00, 0x00, 0x00, // flags
				// system ID
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
				0x00, 0x00, 0x00, 0x05, // data size
				0x21, 0x22, 0x23, 0x24, 0x25, // data
			},
			str: `Version=0 Flags=0x000000 ` +
				`SystemID="0102030405060708090a0b0c0d0e0f10" ` +
				`DataSize=5 ` +
				`Data=[0x21, 0x22, 0x23, 0x24, 0x25]`,
		},
		{
			name: "pssh: version 1: with KIDs",
			src: &Pssh{
				FullBox: FullBox{
					Version: 1,
					Flags:   [3]byte{0x00, 0x00, 0x00},
				},
				SystemID: [16]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10},
				KIDCount: 2,
				KIDs: []PsshKID{
					{KID: [16]byte{0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x10}},
					{KID: [16]byte{0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x20}},
				},
				DataSize: 5,
				Data:     []byte{0x21, 0x22, 0x23, 0x24, 0x25},
			},
			dst: &Pssh{},
			bin: []byte{
				1,                // version
				0x00, 0x00, 0x00, // flags
				// system ID
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
				0x00, 0x00, 0x00, 0x02, // KID count
				// KIDs
				0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x10,
				0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x20,
				0x00, 0x00, 0x00, 0x05, // data size
				0x21, 0x22, 0x23, 0x24, 0x25, // data
			},
			str: `Version=1 Flags=0x000000 ` +
				`SystemID="0102030405060708090a0b0c0d0e0f10" ` +
				`KIDCount=2 ` +
				`KIDs=["1112131415161718191a1b1c1d1e1f10" "2122232425262728292a2b2c2d2e2f20"] ` +
				`DataSize=5 ` +
				`Data=[0x21, 0x22, 0x23, 0x24, 0x25]`,
		},
		{
			name: "sgpd: version 0",
			src: &Sgpd{
				FullBox: FullBox{
					Version: 0,
					Flags:   [3]byte{0x00, 0x00, 0x00},
				},
				GroupingType: [4]byte{'r', 'o', 'l', 'l'},
				EntryCount:   2,
				Unsupported:  []byte{0x11, 0x22, 0x33, 0x44},
			},
			dst: &Sgpd{},
			bin: []byte{
				0,                // version
				0x00, 0x00, 0x00, // flags
				'r', 'o', 'l', 'l', // grouping type
				0x00, 0x00, 0x00, 0x02, // entry count
				0x11, 0x22, 0x33, 0x44, // unsupported
			},
			str: `Version=0 Flags=0x000000 ` +
				`GroupingType="roll" ` +
				`EntryCount=2 ` +
				`Unsupported=[0x11, 0x22, 0x33, 0x44]`,
		},
		{
			name: "sgpd: version 1 roll",
			src: &Sgpd{
				FullBox: FullBox{
					Version: 1,
					Flags:   [3]byte{0x00, 0x00, 0x00},
				},
				GroupingType:  [4]byte{'r', 'o', 'l', 'l'},
				DefaultLength: 2,
				EntryCount:    2,
				RollDistances: []int16{0x1111, 0x2222},
				Unsupported:   []byte{},
			},
			dst: &Sgpd{},
			bin: []byte{
				1,                // version
				0x00, 0x00, 0x00, // flags
				'r', 'o', 'l', 'l', // grouping type
				0x00, 0x00, 0x00, 0x02, // default length
				0x00, 0x00, 0x00, 0x02, // entry count
				0x11, 0x11, 0x22, 0x22, // unsupported
			},
			str: `Version=1 Flags=0x000000 ` +
				`GroupingType="roll" ` +
				`DefaultLength=2 ` +
				`EntryCount=2 ` +
				`RollDistances=[4369, 8738] ` +
				`Unsupported=[]`,
		},
		{
			name: "sgpd: version 1 prol",
			src: &Sgpd{
				FullBox: FullBox{
					Version: 1,
					Flags:   [3]byte{0x00, 0x00, 0x00},
				},
				GroupingType:  [4]byte{'p', 'r', 'o', 'l'},
				DefaultLength: 2,
				EntryCount:    2,
				RollDistances: []int16{0x1111, 0x2222},
				Unsupported:   []byte{},
			},
			dst: &Sgpd{},
			bin: []byte{
				1,                // version
				0x00, 0x00, 0x00, // flags
				'p', 'r', 'o', 'l', // grouping type
				0x00, 0x00, 0x00, 0x02, // default length
				0x00, 0x00, 0x00, 0x02, // entry count
				0x11, 0x11, 0x22, 0x22, // unsupported
			},
			str: `Version=1 Flags=0x000000 ` +
				`GroupingType="prol" ` +
				`DefaultLength=2 ` +
				`EntryCount=2 ` +
				`RollDistances=[4369, 8738] ` +
				`Unsupported=[]`,
		},
		{
			name: "sgpd: version 1 alst",
			src: &Sgpd{
				FullBox: FullBox{
					Version: 1,
					Flags:   [3]byte{0x00, 0x00, 0x00},
				},
				GroupingType:  [4]byte{'a', 'l', 's', 't'},
				DefaultLength: 2,
				EntryCount:    2,
				Unsupported:   []byte{0x11, 0x22, 0x33, 0x44},
			},
			dst: &Sgpd{},
			bin: []byte{
				1,                // version
				0x00, 0x00, 0x00, // flags
				'a', 'l', 's', 't', // grouping type
				0x00, 0x00, 0x00, 0x02, // default length
				0x00, 0x00, 0x00, 0x02, // entry count
				0x11, 0x22, 0x33, 0x44, // unsupported
			},
			str: `Version=1 Flags=0x000000 ` +
				`GroupingType="alst" ` +
				`DefaultLength=2 ` +
				`EntryCount=2 ` +
				`Unsupported=[0x11, 0x22, 0x33, 0x44]`,
		},
		{
			name: "sgpd: version 1 rap",
			src: &Sgpd{
				FullBox: FullBox{
					Version: 1,
					Flags:   [3]byte{0x00, 0x00, 0x00},
				},
				GroupingType:  [4]byte{'r', 'a', 'p', ' '},
				DefaultLength: 1,
				EntryCount:    2,
				VisualRandomAccessEntries: []VisualRandomAccessEntry{
					{NumLeadingSamplesKnown: true, NumLeadingSamples: 0x27},
					{NumLeadingSamplesKnown: false, NumLeadingSamples: 0x1a},
				},
				Unsupported: []byte{},
			},
			dst: &Sgpd{},
			bin: []byte{
				1,                // version
				0x00, 0x00, 0x00, // flags
				'r', 'a', 'p', ' ', // grouping type
				0x00, 0x00, 0x00, 0x01, // default length
				0x00, 0x00, 0x00, 0x02, // entry count
				0xa7, 0x1a, // visual random access entry
			},
			str: `Version=1 Flags=0x000000 ` +
				`GroupingType="rap " ` +
				`DefaultLength=1 ` +
				`EntryCount=2 ` +
				`VisualRandomAccessEntries=[` +
				`{NumLeadingSamplesKnown=true NumLeadingSamples=0x27}, ` +
				`{NumLeadingSamplesKnown=false NumLeadingSamples=0x1a}] ` +
				`Unsupported=[]`,
		},
		{
			name: "sgpd: version 1 tele",
			src: &Sgpd{
				FullBox: FullBox{
					Version: 1,
					Flags:   [3]byte{0x00, 0x00, 0x00},
				},
				GroupingType:  [4]byte{'t', 'e', 'l', 'e'},
				DefaultLength: 1,
				EntryCount:    2,
				TemporalLevelEntries: []TemporalLevelEntry{
					{LevelUndependentlyUecodable: true},
					{LevelUndependentlyUecodable: false},
				},
				Unsupported: []byte{},
			},
			dst: &Sgpd{},
			bin: []byte{
				1,                // version
				0x00, 0x00, 0x00, // flags
				't', 'e', 'l', 'e', // grouping type
				0x00, 0x00, 0x00, 0x01, // default length
				0x00, 0x00, 0x00, 0x02, // entry count
				0x80, 0x00, // temporal level entry
			},
			str: `Version=1 Flags=0x000000 ` +
				`GroupingType="tele" ` +
				`DefaultLength=1 ` +
				`EntryCount=2 ` +
				`TemporalLevelEntries=[` +
				`{LevelUndependentlyUecodable=true}, ` +
				`{LevelUndependentlyUecodable=false}] ` +
				`Unsupported=[]`,
		},
		{
			name: "sgpd: version 2 roll",
			src: &Sgpd{
				FullBox: FullBox{
					Version: 2,
					Flags:   [3]byte{0x00, 0x00, 0x00},
				},
				GroupingType:                  [4]byte{'r', 'o', 'l', 'l'},
				DefaultSampleDescriptionIndex: 5,
				EntryCount:                    2,
				Unsupported:                   []byte{0x11, 0x22, 0x33, 0x44},
			},
			dst: &Sgpd{},
			bin: []byte{
				2,                // version
				0x00, 0x00, 0x00, // flags
				'r', 'o', 'l', 'l', // grouping type
				0x00, 0x00, 0x00, 0x05, // default sample description index
				0x00, 0x00, 0x00, 0x02, // entry count
				0x11, 0x22, 0x33, 0x44, // unsupported
			},
			str: `Version=2 Flags=0x000000 ` +
				`GroupingType="roll" ` +
				`DefaultSampleDescriptionIndex=5 ` +
				`EntryCount=2 ` +
				`Unsupported=[0x11, 0x22, 0x33, 0x44]`,
		},
		{
			name: "tfra: version 0",
			src: &Tfra{
				FullBox: FullBox{
					Version: 0,
					Flags:   [3]byte{0x00, 0x00, 0x00},
				},
				TrackID:               0x11111111,
				LengthSizeOfTrafNum:   0x1,
				LengthSizeOfTrunNum:   0x2,
				LengthSizeOfSampleNum: 0x3,
				NumberOfEntry:         2,
				Entries: []TfraEntry{
					{
						TimeV0:       0x22222222,
						MoofOffsetV0: 0x33333333,
						TrafNumber:   0x4444,
						TrunNumber:   0x555555,
						SampleNumber: 0x66666666,
					},
					{
						TimeV0:       0x77777777,
						MoofOffsetV0: 0x88888888,
						TrafNumber:   0x9999,
						TrunNumber:   0xaaaaaa,
						SampleNumber: 0xbbbbbbbb,
					},
				},
			},
			dst: &Tfra{},
			bin: []byte{
				0,                // version
				0x00, 0x00, 0x00, // flags
				0x11, 0x11, 0x11, 0x11, // trackID
				0x00, 0x00, 0x00, 0x1b, // rserved lengthSizeOfTrafNum lengthSizeOfTrunNum lengthSizeOfSampleNum
				0x00, 0x00, 0x00, 0x02, // numberOfEntry
				0x22, 0x22, 0x22, 0x22, // timeV0
				0x33, 0x33, 0x33, 0x33, // moofOffsetV0
				0x44, 0x44, // trafNumber
				0x55, 0x55, 0x55, // trunNumber
				0x66, 0x66, 0x66, 0x66, // sampleNumber
				0x77, 0x77, 0x77, 0x77, // timeV0
				0x88, 0x88, 0x88, 0x88, // moofOffsetV0
				0x99, 0x99, // trafNumber
				0xaa, 0xaa, 0xaa, // trunNumber
				0xbb, 0xbb, 0xbb, 0xbb, // sampleNumber
			},
			str: `Version=0 Flags=0x000000 ` +
				`TrackID=286331153 ` +
				`LengthSizeOfTrafNum=0x1 ` +
				`LengthSizeOfTrunNum=0x2 ` +
				`LengthSizeOfSampleNum=0x3 ` +
				`NumberOfEntry=2 ` +
				`Entries=[` +
				`{TimeV0=572662306 MoofOffsetV0=858993459 TrafNumber=17476 TrunNumber=5592405 SampleNumber=1717986918}, ` +
				`{TimeV0=2004318071 MoofOffsetV0=2290649224 TrafNumber=39321 TrunNumber=11184810 SampleNumber=3149642683}]`,
		},
		{
			name: "tfra: version 1",
			src: &Tfra{
				FullBox: FullBox{
					Version: 1,
					Flags:   [3]byte{0x00, 0x00, 0x00},
				},
				TrackID:               0x11111111,
				LengthSizeOfTrafNum:   0x1,
				LengthSizeOfTrunNum:   0x2,
				LengthSizeOfSampleNum: 0x3,
				NumberOfEntry:         2,
				Entries: []TfraEntry{
					{
						TimeV1:       0x2222222222222222,
						MoofOffsetV1: 0x3333333333333333,
						TrafNumber:   0x4444,
						TrunNumber:   0x555555,
						SampleNumber: 0x66666666,
					},
					{
						TimeV1:       0x7777777777777777,
						MoofOffsetV1: 0x8888888888888888,
						TrafNumber:   0x9999,
						TrunNumber:   0xaaaaaa,
						SampleNumber: 0xbbbbbbbb,
					},
				},
			},
			dst: &Tfra{},
			bin: []byte{
				1,                // version
				0x00, 0x00, 0x00, // flags
				0x11, 0x11, 0x11, 0x11, // trackID
				0x00, 0x00, 0x00, 0x1b, // rserved lengthSizeOfTrafNum lengthSizeOfTrunNum lengthSizeOfSampleNum
				0x00, 0x00, 0x00, 0x02, // numberOfEntry
				0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, // timeV1
				0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, // moofOffsetV1
				0x44, 0x44, // trafNumber
				0x55, 0x55, 0x55, // trunNumber
				0x66, 0x66, 0x66, 0x66, // sampleNumber
				0x77, 0x77, 0x77, 0x77, 0x77, 0x77, 0x77, 0x77, // timeV1
				0x88, 0x88, 0x88, 0x88, 0x88, 0x88, 0x88, 0x88, // moofOffsetV1
				0x99, 0x99, // trafNumber
				0xaa, 0xaa, 0xaa, // trunNumber
				0xbb, 0xbb, 0xbb, 0xbb, // sampleNumber
			},
			str: `Version=1 Flags=0x000000 ` +
				`TrackID=286331153 ` +
				`LengthSizeOfTrafNum=0x1 ` +
				`LengthSizeOfTrunNum=0x2 ` +
				`LengthSizeOfSampleNum=0x3 ` +
				`NumberOfEntry=2 ` +
				`Entries=[` +
				`{TimeV1=2459565876494606882 MoofOffsetV1=3689348814741910323 TrafNumber=17476 TrunNumber=5592405 SampleNumber=1717986918}, ` +
				`{TimeV1=8608480567731124087 MoofOffsetV1=9838263505978427528 TrafNumber=39321 TrunNumber=11184810 SampleNumber=3149642683}]`,
		},
		{
			name: "tfhd: no flags",
			src: &Tfhd{
				FullBox: FullBox{
					Version: 0,
					Flags:   [3]byte{0x00, 0x00, 0x00},
				},
				TrackID: 0x08404649,
			},
			dst: &Tfhd{},
			bin: []byte{
				0,                // version
				0x00, 0x00, 0x00, // flags
				0x08, 0x40, 0x46, 0x49, // track ID
			},
			str: `Version=0 Flags=0x000000 TrackID=138430025`,
		},
		{
			name: "tfhd: base data offset & default sample duration",
			src: &Tfhd{
				FullBox: FullBox{
					Version: 0,
					Flags:   [3]byte{0x00, 0x00, TfhdBaseDataOffsetPresent | TfhdDefaultSampleDurationPresent},
				},
				TrackID:               0x08404649,
				BaseDataOffset:        0x0123456789abcdef,
				DefaultSampleDuration: 0x23456789,
			},
			dst: &Tfhd{},
			bin: []byte{
				0,                // version
				0x00, 0x00, 0x09, // flags (0000 0000 1001)
				0x08, 0x40, 0x46, 0x49, // track ID
				0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
				0x23, 0x45, 0x67, 0x89,
			},
			str: `Version=0 Flags=0x000009 TrackID=138430025 BaseDataOffset=81985529216486895 DefaultSampleDuration=591751049`,
		},
		{
			name: "tkhd",
			src: &Tkhd{
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
			},
			dst: &Tkhd{},
			bin: []byte{
				0,                // version
				0x00, 0x00, 0x00, // flags
				0x01, 0x23, 0x45, 0x67, // creation time
				0x12, 0x34, 0x56, 0x78, // modification time
				0x23, 0x45, 0x67, 0x89, // track ID
				0x34, 0x56, 0x78, 0x9a, // reserved
				0x45, 0x67, 0x89, 0xab, // duration
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // reserved
				0x5b, 0xa0, // layer
				0xa4, 0x60, // alternate group
				0x01, 0x00, // volume
				0x00, 0x00, // reserved
				0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00, // matrix
				0x56, 0x78, 0x9a, 0xbc, // width
				0x67, 0x89, 0xab, 0xcd, // height
			},
			str: `Version=0 Flags=0x000000 ` +
				`CreationTimeV0=19088743 ` +
				`ModificationTimeV0=305419896 ` +
				`TrackIDV0=591751049 ` +
				`DurationV0=1164413355 ` +
				`Layer=23456 ` +
				`AlternateGroup=-23456 ` +
				`Volume=256 ` +
				`Matrix=[0x10000, 0x0, 0x0, 0x0, 0x10000, 0x0, 0x0, 0x0, 0x40000000] ` +
				`Width=1450744508 Height=1737075661`,
		},
		{
			name: "trun: version=0 flag=0x101",
			src: &Trun{
				FullBox: FullBox{
					Version: 0,
					Flags:   [3]byte{0x00, 0x01, 0x01},
				},
				SampleCount: 3,
				DataOffset:  50,
				Entries: []TrunEntry{
					{SampleDuration: 100},
					{SampleDuration: 101},
					{SampleDuration: 102},
				},
			},
			dst: &Trun{},
			bin: []byte{
				0,                // version
				0x00, 0x01, 0x01, // flags
				0x00, 0x00, 0x00, 0x03, // sample count
				0x00, 0x00, 0x00, 0x32, // data offset
				0x00, 0x00, 0x00, 0x64, // sample duration
				0x00, 0x00, 0x00, 0x65, // sample duration
				0x00, 0x00, 0x00, 0x66, // sample duration
			},
			str: `Version=0 Flags=0x000101 SampleCount=3 DataOffset=50 Entries=[{SampleDuration=100}, {SampleDuration=101}, {SampleDuration=102}]`,
		},
		{
			name: "trun: version=0 flag=0x204",
			src: &Trun{
				FullBox: FullBox{
					Version: 0,
					Flags:   [3]byte{0x00, 0x02, 0x04},
				},
				SampleCount:      3,
				FirstSampleFlags: 0x02468ace,
				Entries: []TrunEntry{
					{SampleSize: 100},
					{SampleSize: 101},
					{SampleSize: 102},
				},
			},
			dst: &Trun{},
			bin: []byte{
				0,                // version
				0x00, 0x02, 0x04, // flags
				0x00, 0x00, 0x00, 0x03, // sample count
				0x02, 0x46, 0x8a, 0xce, // first sample flags
				0x00, 0x00, 0x00, 0x64, // sample size
				0x00, 0x00, 0x00, 0x65, // sample size
				0x00, 0x00, 0x00, 0x66, // sample size
			},
			str: `Version=0 Flags=0x000204 SampleCount=3 FirstSampleFlags=0x2468ace Entries=[{SampleSize=100}, {SampleSize=101}, {SampleSize=102}]`,
		},
		{
			name: "trun: version=0 flag=0xc00",
			src: &Trun{
				FullBox: FullBox{
					Version: 0,
					Flags:   [3]byte{0x00, 0x0c, 0x00},
				},
				SampleCount: 3,
				Entries: []TrunEntry{
					{SampleFlags: 100, SampleCompositionTimeOffsetV0: 200},
					{SampleFlags: 101, SampleCompositionTimeOffsetV0: 201},
					{SampleFlags: 102, SampleCompositionTimeOffsetV0: 202},
				},
			},
			dst: &Trun{},
			bin: []byte{
				0,                // version
				0x00, 0x0c, 0x00, // flags
				0x00, 0x00, 0x00, 0x03, // sample count
				0x00, 0x00, 0x00, 0x64, // sample flags
				0x00, 0x00, 0x00, 0xc8, // sample composition time offset
				0x00, 0x00, 0x00, 0x65, // sample flags
				0x00, 0x00, 0x00, 0xc9, // sample composition time offset
				0x00, 0x00, 0x00, 0x66, // sample flags
				0x00, 0x00, 0x00, 0xca, // sample composition time offset
			},
			str: `Version=0 Flags=0x000c00 SampleCount=3 Entries=[` +
				`{SampleFlags=0x64 SampleCompositionTimeOffsetV0=200}, ` +
				`{SampleFlags=0x65 SampleCompositionTimeOffsetV0=201}, ` +
				`{SampleFlags=0x66 SampleCompositionTimeOffsetV0=202}]`,
		},
		{
			name: "trun: version=1 flag=0x800",
			src: &Trun{
				FullBox: FullBox{
					Version: 1,
					Flags:   [3]byte{0x00, 0x08, 0x00},
				},
				SampleCount: 3,
				Entries: []TrunEntry{
					{SampleCompositionTimeOffsetV1: 200},
					{SampleCompositionTimeOffsetV1: 201},
					{SampleCompositionTimeOffsetV1: -202},
				},
			},
			dst: &Trun{},
			bin: []byte{
				1,                // version
				0x00, 0x08, 0x00, // flags
				0x00, 0x00, 0x00, 0x03, // sample count
				0x00, 0x00, 0x00, 0xc8, // sample composition time offset
				0x00, 0x00, 0x00, 0xc9, // sample composition time offset
				0xff, 0xff, 0xff, 0x36, // sample composition time offset
			},
			str: `Version=1 Flags=0x000800 SampleCount=3 Entries=[` +
				`{SampleCompositionTimeOffsetV1=200}, ` +
				`{SampleCompositionTimeOffsetV1=201}, ` +
				`{SampleCompositionTimeOffsetV1=-202}]`,
		},
		{
			name: "VisualSampleEntry",
			src: &VisualSampleEntry{
				SampleEntry: SampleEntry{
					AnyTypeBox:         AnyTypeBox{Type: StrToBoxType("avc1")},
					DataReferenceIndex: 0x1234,
				},
				PreDefined:      0x0101,
				PreDefined2:     [3]uint32{0x01000001, 0x01000002, 0x01000003},
				Width:           0x0102,
				Height:          0x0103,
				Horizresolution: 0x01000004,
				Vertresolution:  0x01000005,
				Reserved2:       0x01000006,
				FrameCount:      0x0104,
				Compressorname:  [32]byte{5, 'a', 'b', 'e', 'm', 'a'},
				Depth:           0x0105,
				PreDefined3:     1001,
			},
			dst: &VisualSampleEntry{SampleEntry: SampleEntry{AnyTypeBox: AnyTypeBox{Type: StrToBoxType("avc1")}}},
			bin: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // reserved
				0x12, 0x34, // data reference index
				0x01, 0x01, // PreDefined
				0x00, 0x00, // Reserved
				0x01, 0x00, 0x00, 0x01,
				0x01, 0x00, 0x00, 0x02,
				0x01, 0x00, 0x00, 0x03, // PreDefined2
				0x01, 0x02, // Width
				0x01, 0x03, // Height
				0x01, 0x00, 0x00, 0x04, // Horizresolution
				0x01, 0x00, 0x00, 0x05, // Vertresolution
				0x01, 0x00, 0x00, 0x06, // Reserved2
				0x01, 0x04, // FrameCount
				5, 'a', 'b', 'e', 'm', 'a', 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Compressorname
				0x01, 0x05, // Depth
				0x03, 0xe9, // PreDefined3
			},
			str: `DataReferenceIndex=4660 ` +
				`PreDefined=257 ` +
				`PreDefined2=[16777217, 16777218, 16777219] ` +
				`Width=258 ` +
				`Height=259 ` +
				`Horizresolution=16777220 ` +
				`Vertresolution=16777221 ` +
				`FrameCount=260 ` +
				`Compressorname="abema" ` +
				`Depth=261 ` +
				`PreDefined3=1001`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Marshal
			buf := bytes.NewBuffer(nil)
			n, err := Marshal(buf, tc.src)
			require.NoError(t, err)
			assert.Equal(t, uint64(len(tc.bin)), n)
			assert.Equal(t, tc.bin, buf.Bytes())

			// Unmarshal
			n, err = Unmarshal(bytes.NewReader(tc.bin), uint64(len(tc.bin)), tc.dst)
			assert.NoError(t, err)
			assert.Equal(t, uint64(buf.Len()), n)
			assert.Equal(t, tc.src, tc.dst)

			// UnmarshalAny
			dst, n, err := UnmarshalAny(bytes.NewReader(tc.bin), tc.src.GetType(), uint64(len(tc.bin)))
			assert.NoError(t, err)
			assert.Equal(t, uint64(buf.Len()), n)
			assert.Equal(t, tc.src, dst)

			// Stringify
			str, err := Stringify(tc.src)
			require.NoError(t, err)
			assert.Equal(t, tc.str, str)
		})
	}
}

func TestHdlrUnmarshalHandlerName(t *testing.T) {
	testCases := []struct {
		name          string
		componentType []byte
		bytes         []byte
		want          string
		padding       int
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
			bytes:         []byte{0x00, 0x00},
			want:          "",
			padding:       1,
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
			assert.Len(t, dst.Padding, tc.padding)
		})
	}
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
	n, err := Unmarshal(r, uint64(len(bin)), &dst)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), n)
	s, _ := r.Seek(0, io.SeekCurrent)
	assert.Equal(t, int64(0), s)
	assert.Equal(t, uint8(0), dst.GetVersion())
	assert.Equal(t, uint32(0), dst.GetFlags())
}
