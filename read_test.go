package mp4

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadBoxStructure(t *testing.T) {
	f, err := os.Open("./_examples/sample.mp4")
	require.NoError(t, err)
	defer f.Close()

	var n int
	_, err = ReadBoxStructure(f, func(h *ReadHandle) (interface{}, error) {
		n++
		switch n {
		case 55, 56: // unsupported
			require.False(t, h.BoxInfo.Type.IsSupported())
			buf := bytes.NewBuffer(nil)
			n, err := h.ReadData(buf)
			require.NoError(t, err)
			require.Equal(t, h.BoxInfo.Size-h.BoxInfo.HeaderSize, n)
			assert.Len(t, buf.Bytes(), int(n))
		case 41: // stbl
			require.True(t, h.BoxInfo.Type.IsSupported())
			require.Equal(t, BoxTypeStbl(), h.BoxInfo.Type)
			infos, err := h.Expand()
			require.NoError(t, err)
			assert.Equal(t, []interface{}{"stsd", "stts", nil, nil, "stco", nil, nil}, infos)
		case 42: // stsd
			require.True(t, h.BoxInfo.Type.IsSupported())
			require.Equal(t, BoxTypeStsd(), h.BoxInfo.Type)
			box, n, err := h.ReadPayload()
			require.NoError(t, err)
			require.Equal(t, uint64(8), n)
			assert.Equal(t, &Stsd{EntryCount: 1}, box)
			_, err = h.Expand()
			require.NoError(t, err)
			return "stsd", nil
		case 45: // stts
			require.True(t, h.BoxInfo.Type.IsSupported())
			require.Equal(t, BoxTypeStts(), h.BoxInfo.Type)
			_, err = h.Expand()
			require.NoError(t, err)
			return "stts", nil
		case 48: // stco
			require.True(t, h.BoxInfo.Type.IsSupported())
			require.Equal(t, BoxTypeStco(), h.BoxInfo.Type)
			_, err = h.Expand()
			require.NoError(t, err)
			return "stco", nil
		default: // otherwise
			require.True(t, h.BoxInfo.Type.IsSupported())
			_, err = h.Expand()
			require.NoError(t, err)
		}
		return nil, nil
	})
	require.NoError(t, err)
	assert.Equal(t, 56, n)
}

// > mp4tool dump _examples/sample.mp4 | cat -n
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
// 55	        [0xa9746f6f] (unsupported box type) Size=37 Data=[...] (use "-full 0xa9746f6f" to show all)
// 56	    [loci] (unsupported box type) Size=35 Data=[...] (use "-full loci" to show all)
