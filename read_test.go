package mp4

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadBoxStructure(t *testing.T) {
	f, err := os.Open("./testdata/sample.mp4")
	require.NoError(t, err)
	defer f.Close()

	var n int
	_, err = ReadBoxStructure(f, func(h *ReadHandle) (interface{}, error) {
		n++
		switch n {
		case 57: // unsupported
			require.False(t, h.BoxInfo.IsSupportedType())
			buf := bytes.NewBuffer(nil)
			n, err := h.ReadData(buf)
			require.NoError(t, err)
			require.Equal(t, h.BoxInfo.Size-h.BoxInfo.HeaderSize, n)
			assert.Len(t, buf.Bytes(), int(n))
		case 41: // stbl
			require.True(t, h.BoxInfo.IsSupportedType())
			require.Equal(t, BoxTypeStbl(), h.BoxInfo.Type)
			infos, err := h.Expand()
			require.NoError(t, err)
			assert.Equal(t, []interface{}{"stsd", "stts", nil, nil, "stco", nil, nil}, infos)
		case 42: // stsd
			require.True(t, h.BoxInfo.IsSupportedType())
			require.Equal(t, BoxTypeStsd(), h.BoxInfo.Type)
			box, n, err := h.ReadPayload()
			require.NoError(t, err)
			require.Equal(t, uint64(8), n)
			assert.Equal(t, &Stsd{EntryCount: 1}, box)
			_, err = h.Expand()
			require.NoError(t, err)
			return "stsd", nil
		case 45: // stts
			require.True(t, h.BoxInfo.IsSupportedType())
			require.Equal(t, BoxTypeStts(), h.BoxInfo.Type)
			_, err = h.Expand()
			require.NoError(t, err)
			return "stts", nil
		case 48: // stco
			require.True(t, h.BoxInfo.IsSupportedType())
			require.Equal(t, BoxTypeStco(), h.BoxInfo.Type)
			_, err = h.Expand()
			require.NoError(t, err)
			return "stco", nil
		case 56: // data
			require.True(t, h.BoxInfo.IsSupportedType())
			require.Equal(t, BoxTypeData(), h.BoxInfo.Type)
			box, n, err := h.ReadPayload()
			require.NoError(t, err)
			require.Equal(t, uint64(21), n)
			assert.Equal(t, &Data{DataType: DataTypeStringUTF8, DataLang: 0, Data: []byte("Lavf58.29.100")}, box)
			_, err = h.Expand()
			require.NoError(t, err)
			return "stco", nil
		default: // otherwise
			require.True(t, h.BoxInfo.IsSupportedType())
			_, err = h.Expand()
			require.NoError(t, err)
		}
		return nil, nil
	})
	require.NoError(t, err)
	assert.Equal(t, 57, n)
}

// > mp4tool dump testdata/sample.mp4 | cat -n
//  1	[ftyp] Size=32 MajorBrand="isom" MinorVersion=512 CompatibleBrands=[{CompatibleBrand="isom"}, {CompatibleBrand="iso2"}, {CompatibleBrand="avc1"}, {CompatibleBrand="mp41"}]
//  2	[free] Size=8 Data=[...] (use "-full free" to show all)
//  3	[mdat] Size=6402 Data=[...] (use "-full mdat" to show all)
//  4	[moov] Size=1836
//  5	  [mvhd] Size=108 ... (use "-full mvhd" to show all)
//  6	  [trak] Size=743
//  7	    [tkhd] Size=92 ... (use "-full tkhd" to show all)
//  8	    [edts] Size=36
//  9	      [elst] Size=28 Version=0 Flags=0x000000 EntryCount=1 Entries=[{SegmentDurationV0=1000 MediaTimeV0=2048 MediaRateInteger=1}]
// 10	    [mdia] Size=607
// 11	      [mdhd] Size=32 Version=0 Flags=0x000000 CreationTimeV0=0 ModificationTimeV0=0 Timescale=10240 DurationV0=10240 Pad=false Language="eng" PreDefined=0
// 12	      [hdlr] Size=44 Version=0 Flags=0x000000 PreDefined=0 HandlerType="vide" Name="VideoHandle"
// 13	      [minf] Size=523
// 14	        [vmhd] Size=20 Version=0 Flags=0x000001 Graphicsmode=0 Opcolor=[0, 0, 0]
// 15	        [dinf] Size=36
// 16	          [dref] Size=28 Version=0 Flags=0x000000 EntryCount=1
// 17	            [url ] Size=12 Version=0 Flags=0x000001
// 18	        [stbl] Size=459
// 19	          [stsd] Size=167 Version=0 Flags=0x000000 EntryCount=1
// 20	            [avc1] Size=151 ... (use "-full avc1" to show all)
// 21	              [avcC] Size=49 ConfigurationVersion=0x1 Profile=0x64 ProfileCompatibility=0x0 Level=0xc Data=[0xff, 0xe1, 0x0, 0x19, 0x67, 0x64, 0x0, 0xc, 0xac, 0xd9, 0x41, 0x41, 0x9f, 0x9f, 0x1, 0x6c, 0x80, 0x0, 0x0, 0x3, 0x0, 0x80, 0x0, 0x0, 0xa, 0x7, 0x8a, 0x14, 0xcb, 0x1, 0x0, 0x5, 0x68, 0xeb, 0xec, 0xb2, 0x2c]
// 22	              [pasp] Size=16 HSpacing=1 VSpacing=1
// 23	          [stts] Size=24 Version=0 Flags=0x000000 EntryCount=1 Entries=[{SampleCount=10 SampleDelta=1024}]
// 24	          [stss] Size=20 Version=0 Flags=0x000000 EntryCount=1 SampleNumber=[1]
// 25	          [ctts] Size=88 ... (use "-full ctts" to show all)
// 26	          [stsc] Size=40 Version=0 Flags=0x000000 EntryCount=2 Entries=[{FirstChunk=1 SamplesPerChunk=2 SampleDescriptionIndex=1}, {FirstChunk=2 SamplesPerChunk=1 SampleDescriptionIndex=1}]
// 27	          [stsz] Size=60 Version=0 Flags=0x000000 SampleSize=0 SampleCount=10 EntrySize=[3679, 86, 545, 180, 69, 60, 182, 22, 204, 15]
// 28	          [stco] Size=52 Version=0 Flags=0x000000 EntryCount=9 ChunkOffset=[48, 3836, 4527, 4864, 5043, 5227, 5560, 5702, 6038]
// 29	  [trak] Size=844
// 30	    [tkhd] Size=92 ... (use "-full tkhd" to show all)
// 31	    [edts] Size=36
// 32	      [elst] Size=28 Version=0 Flags=0x000000 EntryCount=1 Entries=[{SegmentDurationV0=1000 MediaTimeV0=1024 MediaRateInteger=1}]
// 33	    [mdia] Size=708
// 34	      [mdhd] Size=32 Version=0 Flags=0x000000 CreationTimeV0=0 ModificationTimeV0=0 Timescale=44100 DurationV0=45124 Pad=false Language="eng" PreDefined=0
// 35	      [hdlr] Size=44 Version=0 Flags=0x000000 PreDefined=0 HandlerType="soun" Name="SoundHandle"
// 36	      [minf] Size=624
// 37	        [smhd] Size=16 Version=0 Flags=0x000000 Balance=0
// 38	        [dinf] Size=36
// 39	          [dref] Size=28 Version=0 Flags=0x000000 EntryCount=1
// 40	            [url ] Size=12 Version=0 Flags=0x000001
// 41	        [stbl] Size=564
// 42	          [stsd] Size=106 Version=0 Flags=0x000000 EntryCount=1
// 43	            [mp4a] Size=90 DataReferenceIndex=1 EntryVersion=0 ChannelCount=2 SampleSize=16 PreDefined=0 SampleRate=2890137600
// 44	              [esds] Size=54 Version=0 Flags=0x000000 Descriptors=[{Tag=ESDescr Size=37 ESID=2 StreamDependenceFlag=false UrlFlag=false OcrStreamFlag=false StreamPriority=0}, {Tag=DecoderConfigDescr Size=23 ObjectTypeIndication=0x40 StreamType=5 UpStream=false Reserved=true BufferSizeDB=0 MaxBitrate=10570 AvgBitrate=10570}, {Tag=DecSpecificInfo Size=5 Data=[0x12, 0x10, 0x56, 0xe5, 0x0]}, {Tag=SLConfigDescr Size=1 Data=[0x2]}]
// 45	          [stts] Size=48 Version=0 Flags=0x000000 EntryCount=4 Entries=[{SampleCount=1 SampleDelta=1024}, {SampleCount=1 SampleDelta=1505}, {SampleCount=41 SampleDelta=1024}, {SampleCount=1 SampleDelta=611}]
// 46	          [stsc] Size=100 ... (use "-full stsc" to show all)
// 47	          [stsz] Size=196 ... (use "-full stsz" to show all)
// 48	          [stco] Size=52 Version=0 Flags=0x000000 EntryCount=9 ChunkOffset=[3813, 4381, 4707, 4933, 5103, 5409, 5582, 5906, 6053]
// 49	          [sgpd] Size=26 Version=1 Flags=0x000000 GroupingType="roll" DefaultLength=2 EntryCount=1 RollDistances=[-1] Unsupported=[]
// 50	          [sbgp] Size=28 Version=0 Flags=0x000000 GroupingType=1919904876 EntryCount=1 Entries=[{SampleCount=44 GroupDescriptionIndex=1}]
// 51	  [udta] Size=133
// 52	    [meta] Size=90 Version=0 Flags=0x000000
// 53	      [hdlr] Size=33 Version=0 Flags=0x000000 PreDefined=0 HandlerType="mdir" Name=""
// 54	      [ilst] Size=45
// 55	        [(c)too] Size=37
// 56	          [data] Size=29 DataType=UTF8 DataLang=0 Data="Lavf58.29.100"
// 57	    [loci] (unsupported box type) Size=35 Data=[...] (use "-full loci" to show all)

func TestReadBoxStructureQT(t *testing.T) {
	f, err := os.Open("./testdata/sample_qt.mp4")
	require.NoError(t, err)
	defer f.Close()

	var n int
	_, err = ReadBoxStructure(f, func(h *ReadHandle) (interface{}, error) {
		n++
		switch n {
		case 5, 45: // unsupported
			require.False(t, h.BoxInfo.IsSupportedType())
			buf := bytes.NewBuffer(nil)
			n, err := h.ReadData(buf)
			require.NoError(t, err)
			require.Equal(t, h.BoxInfo.Size-h.BoxInfo.HeaderSize, n)
			assert.Len(t, buf.Bytes(), int(n))
		case 40: // mp4a
			require.True(t, h.BoxInfo.IsSupportedType())
			require.Equal(t, StrToBoxType("mp4a"), h.BoxInfo.Type)
			box, n, err := h.ReadPayload()
			require.NoError(t, err)
			require.Equal(t, uint64(44), n)
			assert.Equal(t, []byte{0x0, 0x0, 0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2}, box.(*AudioSampleEntry).QuickTimeData)
			_, err = h.Expand()
			require.NoError(t, err)
		case 43: // mp4a
			require.True(t, h.BoxInfo.IsSupportedType())
			require.Equal(t, StrToBoxType("mp4a"), h.BoxInfo.Type)
			box, n, err := h.ReadPayload()
			require.NoError(t, err)
			require.Equal(t, uint64(4), n)
			assert.Equal(t, []byte{0x0, 0x0, 0x0, 0x0}, box.(*AudioSampleEntry).QuickTimeData)
			_, err = h.Expand()
			require.NoError(t, err)
		default: // otherwise
			require.True(t, h.BoxInfo.IsSupportedType())
			_, err = h.Expand()
			require.NoError(t, err)
		}
		return nil, nil
	})
	require.NoError(t, err)
	assert.Equal(t, 49, n)
}

// > mp4tool dump -full mp4a sample_qt.mp4 | cat -n
//  1	[ftyp] Size=20 MajorBrand="qt  " MinorVersion=512 CompatibleBrands=[{CompatibleBrand="qt  "}]
//  2	[free] Size=42 Data=[...] (use "-full free" to show all)
//  3	[moov] Size=340232
//  4	  [udta] Size=31
//  5	    [(c)enc] (unsupported box type) Size=23 Data=[...] (use "-full (c)enc" to show all)
//  6	  [mvhd] Size=108 ... (use "-full mvhd" to show all)
//  7	  [trak] Size=115889
//  8	    [tkhd] Size=92 ... (use "-full tkhd" to show all)
//  9	    [mdia] Size=115789
// 10	      [mdhd] Size=32 Version=0 Flags=0x000000 CreationTimeV0=2082844800 ModificationTimeV0=2082844800 Timescale=24 DurationV0=14315 Language="```" PreDefined=0
// 11	      [hdlr] Size=45 Version=0 Flags=0x000000 PreDefined=1835560050 HandlerType="vide" Name="VideoHandler"
// 12	      [minf] Size=115704
// 13	        [hdlr] Size=44 Version=0 Flags=0x000000 PreDefined=1684565106 HandlerType="url " Name="DataHandler"
// 14	        [vmhd] Size=20 Version=0 Flags=0x000001 Graphicsmode=0 Opcolor=[0, 0, 0]
// 15	        [dinf] Size=36
// 16	          [dref] Size=28 Version=0 Flags=0x000000 EntryCount=1
// 17	            [url ] Size=12 Version=0 Flags=0x000001
// 18	        [stbl] Size=115596
// 19	          [stsd] Size=148 Version=0 Flags=0x000000 EntryCount=1
// 20	            [avc1] Size=132 ... (use "-full avc1" to show all)
// 21	              [avcC] Size=46 ... (use "-full avcC" to show all)
// 22	          [stts] Size=24 Version=0 Flags=0x000000 EntryCount=1 Entries=[{SampleCount=14315 SampleDelta=1}]
// 23	          [stss] Size=832 ... (use "-full stss" to show all)
// 24	          [stsc] Size=28 Version=0 Flags=0x000000 EntryCount=1 Entries=[{FirstChunk=1 SamplesPerChunk=1 SampleDescriptionIndex=1}]
// 25	          [stsz] Size=57280 ... (use "-full stsz" to show all)
// 26	          [stco] Size=57276 ... (use "-full stco" to show all)
// 27	  [trak] Size=224196
// 28	    [tkhd] Size=92 ... (use "-full tkhd" to show all)
// 29	    [mdia] Size=224096
// 30	      [mdhd] Size=32 Version=0 Flags=0x000000 CreationTimeV0=2082844800 ModificationTimeV0=2082844800 Timescale=48000 DurationV0=28628992 Language="```" PreDefined=0
// 31	      [hdlr] Size=45 Version=0 Flags=0x000000 PreDefined=1835560050 HandlerType="soun" Name="SoundHandler"
// 32	      [minf] Size=224011
// 33	        [hdlr] Size=44 Version=0 Flags=0x000000 PreDefined=1684565106 HandlerType="url " Name="DataHandler"
// 34	        [smhd] Size=16 Version=0 Flags=0x000000 Balance=0
// 35	        [dinf] Size=36
// 36	          [dref] Size=28 Version=0 Flags=0x000000 EntryCount=1
// 37	            [url ] Size=12 Version=0 Flags=0x000001
// 38	        [stbl] Size=223907
// 39	          [stsd] Size=147 Version=0 Flags=0x000000 EntryCount=1
// 40	            [mp4a] Size=131 DataReferenceIndex=1 EntryVersion=1 ChannelCount=2 SampleSize=16 PreDefined=65534 SampleRate=3145728000 QuickTimeData=[0x0, 0x0, 0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2]
// 41	              [wave] Size=79
// 42	                [frma] Size=12 DataFormat="mp4a"
// 43	                [mp4a] Size=12 QuickTimeData=[0x0, 0x0, 0x0, 0x0]
// 44	                [esds] Size=39 ... (use "-full esds" to show all)
// 45	                [0x00000000] (unsupported box type) Size=8 Data=[...] (use "-full 0x00000000" to show all)
// 46	          [stts] Size=24 Version=0 Flags=0x000000 EntryCount=1 Entries=[{SampleCount=27958 SampleDelta=1024}]
// 47	          [stsc] Size=28 Version=0 Flags=0x000000 EntryCount=1 Entries=[{FirstChunk=1 SamplesPerChunk=1 SampleDescriptionIndex=1}]
// 48	          [stsz] Size=111852 ... (use "-full stsz" to show all)
// 49	          [stco] Size=111848 ... (use "-full stco" to show all)

// this used to cause an infinite loop.
func TestReadBoxStructureZeroSize(t *testing.T) {
	b := []byte("\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01")

	_, err := ReadBoxStructure(bytes.NewReader(b), func(h *ReadHandle) (interface{}, error) {
		return nil, nil
	})
	require.Error(t, err)
}

func FuzzReadBoxStructure(f *testing.F) {
	// AC-3 track from Apple HLS
	f.Add([]byte{
		0x00, 0x00, 0x00, 0x20, 0x66, 0x74, 0x79, 0x70,
		0x6d, 0x70, 0x34, 0x32, 0x00, 0x00, 0x00, 0x01,
		0x6d, 0x70, 0x34, 0x31, 0x6d, 0x70, 0x34, 0x32,
		0x69, 0x73, 0x6f, 0x6d, 0x68, 0x6c, 0x73, 0x66,
		0x00, 0x00, 0x02, 0x20, 0x6d, 0x6f, 0x6f, 0x76,
		0x00, 0x00, 0x00, 0x6c, 0x6d, 0x76, 0x68, 0x64,
		0x00, 0x00, 0x00, 0x00, 0xd5, 0x5b, 0xc6, 0x5d,
		0xd5, 0x5b, 0xc6, 0x5d, 0x00, 0x00, 0x02, 0x58,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00,
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x01, 0x84,
		0x74, 0x72, 0x61, 0x6b, 0x00, 0x00, 0x00, 0x5c,
		0x74, 0x6b, 0x68, 0x64, 0x00, 0x00, 0x00, 0x01,
		0xd5, 0x5b, 0xc6, 0x5d, 0xd5, 0x5b, 0xc6, 0x5d,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x01, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x01, 0x20, 0x6d, 0x64, 0x69, 0x61,
		0x00, 0x00, 0x00, 0x20, 0x6d, 0x64, 0x68, 0x64,
		0x00, 0x00, 0x00, 0x00, 0xd5, 0x5b, 0xc6, 0x5d,
		0xd5, 0x5b, 0xc6, 0x5d, 0x00, 0x00, 0xbb, 0x80,
		0x00, 0x00, 0x00, 0x00, 0x55, 0xc4, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x31, 0x68, 0x64, 0x6c, 0x72,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x73, 0x6f, 0x75, 0x6e, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x43, 0x6f, 0x72, 0x65, 0x20, 0x4d, 0x65, 0x64,
		0x69, 0x61, 0x20, 0x41, 0x75, 0x64, 0x69, 0x6f,
		0x00, 0x00, 0x00, 0x00, 0xc7, 0x6d, 0x69, 0x6e,
		0x66, 0x00, 0x00, 0x00, 0x10, 0x73, 0x6d, 0x68,
		0x64, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x24, 0x64, 0x69, 0x6e,
		0x66, 0x00, 0x00, 0x00, 0x1c, 0x64, 0x72, 0x65,
		0x66, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x01, 0x00, 0x00, 0x00, 0x0c, 0x75, 0x72, 0x6c,
		0x20, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
		0x8b, 0x73, 0x74, 0x62, 0x6c, 0x00, 0x00, 0x00,
		0x3f, 0x73, 0x74, 0x73, 0x64, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
		0x2f, 0x61, 0x63, 0x2d, 0x33, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00,
		0x10, 0x00, 0x00, 0x00, 0x00, 0xbb, 0x80, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x0b, 0x64, 0x61, 0x63,
		0x33, 0x0c, 0x3d, 0x40, 0x00, 0x00, 0x00, 0x10,
		0x73, 0x74, 0x74, 0x73, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10,
		0x73, 0x74, 0x73, 0x63, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x14,
		0x73, 0x74, 0x73, 0x7a, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x10, 0x73, 0x74, 0x63, 0x6f,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x28, 0x6d, 0x76, 0x65, 0x78,
		0x00, 0x00, 0x00, 0x20, 0x74, 0x72, 0x65, 0x78,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	})

	f.Fuzz(func(t *testing.T, b []byte) {
		ReadBoxStructure(bytes.NewReader(b), func(h *ReadHandle) (interface{}, error) {
			if h.BoxInfo.IsSupportedType() {
				_, _, err := h.ReadPayload()
				if err != nil {
					return nil, err
				}

				return h.Expand()
			}

			return nil, nil
		})
	})
}
