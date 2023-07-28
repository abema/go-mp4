package dump

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDump(t *testing.T) {
	testCases := []struct {
		name    string
		file    string
		options []string
		wants   string
	}{
		{
			name:  "sample.mp4 no-options",
			file:  "../../_examples/sample.mp4",
			wants: sampleMP4Output,
		},
		{
			name:    "sample.mp4 with -full mvhd,loci option",
			file:    "../../_examples/sample.mp4",
			options: []string{"-full", "mvhd,loci"},
			wants:   sampleMP4OutputFullMvhdLoci,
		},
		{
			name:    "sample.mp4 with -offset option",
			file:    "../../_examples/sample.mp4",
			options: []string{"-offset"},
			wants:   sampleMP4OutputOffset,
		},
		{
			name:    "sample.mp4 with -hex option",
			file:    "../../_examples/sample.mp4",
			options: []string{"-hex"},
			wants:   sampleMP4OutputHex,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stdout := os.Stdout
			r, w, err := os.Pipe()
			require.NoError(t, err)
			defer func() {
				os.Stdout = stdout
			}()
			os.Stdout = w
			Main(append(tc.options, tc.file))
			w.Close()
			b, err := ioutil.ReadAll(r)
			require.NoError(t, err)
			assert.Equal(t, tc.wants, string(b))
		})
	}
}

var sampleMP4Output = "" +
	`[ftyp] Size=32 MajorBrand="isom" MinorVersion=512 CompatibleBrands=[{CompatibleBrand="isom"}, {CompatibleBrand="iso2"}, {CompatibleBrand="avc1"}, {CompatibleBrand="mp41"}]` + "\n" +
	`[free] Size=8 Data=[...] (use "-full free" to show all)` + "\n" +
	`[mdat] Size=6402 Data=[...] (use "-full mdat" to show all)` + "\n" +
	`[moov] Size=1836` + "\n" +
	`  [mvhd] Size=108 ... (use "-full mvhd" to show all)` + "\n" +
	`  [trak] Size=743` + "\n" +
	`    [tkhd] Size=92 ... (use "-full tkhd" to show all)` + "\n" +
	`    [edts] Size=36` + "\n" +
	`      [elst] Size=28 Version=0 Flags=0x000000 EntryCount=1 Entries=[{SegmentDurationV0=1000 MediaTimeV0=2048 MediaRateInteger=1}]` + "\n" +
	`    [mdia] Size=607` + "\n" +
	`      [mdhd] Size=32 Version=0 Flags=0x000000 CreationTimeV0=0 ModificationTimeV0=0 Timescale=10240 DurationV0=10240 Language="eng" PreDefined=0` + "\n" +
	`      [hdlr] Size=44 Version=0 Flags=0x000000 PreDefined=0 HandlerType="vide" Name="VideoHandle"` + "\n" +
	`      [minf] Size=523` + "\n" +
	`        [vmhd] Size=20 Version=0 Flags=0x000001 Graphicsmode=0 Opcolor=[0, 0, 0]` + "\n" +
	`        [dinf] Size=36` + "\n" +
	`          [dref] Size=28 Version=0 Flags=0x000000 EntryCount=1` + "\n" +
	`            [url ] Size=12 Version=0 Flags=0x000001` + "\n" +
	`        [stbl] Size=459` + "\n" +
	`          [stsd] Size=167 Version=0 Flags=0x000000 EntryCount=1` + "\n" +
	`            [avc1] Size=151 ... (use "-full avc1" to show all)` + "\n" +
	`              [avcC] Size=49 ... (use "-full avcC" to show all)` + "\n" +
	`              [pasp] Size=16 HSpacing=1 VSpacing=1` + "\n" +
	`          [stts] Size=24 Version=0 Flags=0x000000 EntryCount=1 Entries=[{SampleCount=10 SampleDelta=1024}]` + "\n" +
	`          [stss] Size=20 Version=0 Flags=0x000000 EntryCount=1 SampleNumber=[1]` + "\n" +
	`          [ctts] Size=88 ... (use "-full ctts" to show all)` + "\n" +
	`          [stsc] Size=40 ... (use "-full stsc" to show all)` + "\n" +
	`          [stsz] Size=60 Version=0 Flags=0x000000 SampleSize=0 SampleCount=10 EntrySize=[3679, 86, 545, 180, 69, 60, 182, 22, 204, 15]` + "\n" +
	`          [stco] Size=52 Version=0 Flags=0x000000 EntryCount=9 ChunkOffset=[48, 3836, 4527, 4864, 5043, 5227, 5560, 5702, 6038]` + "\n" +
	`  [trak] Size=844` + "\n" +
	`    [tkhd] Size=92 ... (use "-full tkhd" to show all)` + "\n" +
	`    [edts] Size=36` + "\n" +
	`      [elst] Size=28 Version=0 Flags=0x000000 EntryCount=1 Entries=[{SegmentDurationV0=1000 MediaTimeV0=1024 MediaRateInteger=1}]` + "\n" +
	`    [mdia] Size=708` + "\n" +
	`      [mdhd] Size=32 Version=0 Flags=0x000000 CreationTimeV0=0 ModificationTimeV0=0 Timescale=44100 DurationV0=45124 Language="eng" PreDefined=0` + "\n" +
	`      [hdlr] Size=44 Version=0 Flags=0x000000 PreDefined=0 HandlerType="soun" Name="SoundHandle"` + "\n" +
	`      [minf] Size=624` + "\n" +
	`        [smhd] Size=16 Version=0 Flags=0x000000 Balance=0` + "\n" +
	`        [dinf] Size=36` + "\n" +
	`          [dref] Size=28 Version=0 Flags=0x000000 EntryCount=1` + "\n" +
	`            [url ] Size=12 Version=0 Flags=0x000001` + "\n" +
	`        [stbl] Size=564` + "\n" +
	`          [stsd] Size=106 Version=0 Flags=0x000000 EntryCount=1` + "\n" +
	`            [mp4a] Size=90 DataReferenceIndex=1 EntryVersion=0 ChannelCount=2 SampleSize=16 PreDefined=0 SampleRate=44100` + "\n" +
	`              [esds] Size=54 ... (use "-full esds" to show all)` + "\n" +
	`          [stts] Size=48 ... (use "-full stts" to show all)` + "\n" +
	`          [stsc] Size=100 ... (use "-full stsc" to show all)` + "\n" +
	`          [stsz] Size=196 ... (use "-full stsz" to show all)` + "\n" +
	`          [stco] Size=52 Version=0 Flags=0x000000 EntryCount=9 ChunkOffset=[3813, 4381, 4707, 4933, 5103, 5409, 5582, 5906, 6053]` + "\n" +
	`          [sgpd] Size=26 Version=1 Flags=0x000000 GroupingType="roll" DefaultLength=2 EntryCount=1 RollDistances=[-1]` + "\n" +
	`          [sbgp] Size=28 Version=0 Flags=0x000000 GroupingType=1919904876 EntryCount=1 Entries=[{SampleCount=44 GroupDescriptionIndex=1}]` + "\n" +
	`  [udta] Size=133` + "\n" +
	`    [meta] Size=90 Version=0 Flags=0x000000` + "\n" +
	`      [hdlr] Size=33 Version=0 Flags=0x000000 PreDefined=0 HandlerType="mdir" Name=""` + "\n" +
	`      [ilst] Size=45` + "\n" +
	`        [(c)too] Size=37` + "\n" +
	`          [data] Size=29 DataType=UTF8 DataLang=0 Data="Lavf58.29.100"` + "\n" +
	`    [loci] (unsupported box type) Size=35 Data=[...] (use "-full loci" to show all)` + "\n"

var sampleMP4OutputFullMvhdLoci = "" +
	`[ftyp] Size=32 MajorBrand="isom" MinorVersion=512 CompatibleBrands=[{CompatibleBrand="isom"}, {CompatibleBrand="iso2"}, {CompatibleBrand="avc1"}, {CompatibleBrand="mp41"}]` + "\n" +
	`[free] Size=8 Data=[...] (use "-full free" to show all)` + "\n" +
	`[mdat] Size=6402 Data=[...] (use "-full mdat" to show all)` + "\n" +
	`[moov] Size=1836` + "\n" +
	`  [mvhd] Size=108 Version=0 Flags=0x000000 CreationTimeV0=0 ModificationTimeV0=0 Timescale=1000 DurationV0=1024 Rate=1 Volume=256 Matrix=[0x10000, 0x0, 0x0, 0x0, 0x10000, 0x0, 0x0, 0x0, 0x40000000] PreDefined=[0, 0, 0, 0, 0, 0] NextTrackID=3` + "\n" +
	`  [trak] Size=743` + "\n" +
	`    [tkhd] Size=92 ... (use "-full tkhd" to show all)` + "\n" +
	`    [edts] Size=36` + "\n" +
	`      [elst] Size=28 Version=0 Flags=0x000000 EntryCount=1 Entries=[{SegmentDurationV0=1000 MediaTimeV0=2048 MediaRateInteger=1}]` + "\n" +
	`    [mdia] Size=607` + "\n" +
	`      [mdhd] Size=32 Version=0 Flags=0x000000 CreationTimeV0=0 ModificationTimeV0=0 Timescale=10240 DurationV0=10240 Language="eng" PreDefined=0` + "\n" +
	`      [hdlr] Size=44 Version=0 Flags=0x000000 PreDefined=0 HandlerType="vide" Name="VideoHandle"` + "\n" +
	`      [minf] Size=523` + "\n" +
	`        [vmhd] Size=20 Version=0 Flags=0x000001 Graphicsmode=0 Opcolor=[0, 0, 0]` + "\n" +
	`        [dinf] Size=36` + "\n" +
	`          [dref] Size=28 Version=0 Flags=0x000000 EntryCount=1` + "\n" +
	`            [url ] Size=12 Version=0 Flags=0x000001` + "\n" +
	`        [stbl] Size=459` + "\n" +
	`          [stsd] Size=167 Version=0 Flags=0x000000 EntryCount=1` + "\n" +
	`            [avc1] Size=151 ... (use "-full avc1" to show all)` + "\n" +
	`              [avcC] Size=49 ... (use "-full avcC" to show all)` + "\n" +
	`              [pasp] Size=16 HSpacing=1 VSpacing=1` + "\n" +
	`          [stts] Size=24 Version=0 Flags=0x000000 EntryCount=1 Entries=[{SampleCount=10 SampleDelta=1024}]` + "\n" +
	`          [stss] Size=20 Version=0 Flags=0x000000 EntryCount=1 SampleNumber=[1]` + "\n" +
	`          [ctts] Size=88 ... (use "-full ctts" to show all)` + "\n" +
	`          [stsc] Size=40 ... (use "-full stsc" to show all)` + "\n" +
	`          [stsz] Size=60 Version=0 Flags=0x000000 SampleSize=0 SampleCount=10 EntrySize=[3679, 86, 545, 180, 69, 60, 182, 22, 204, 15]` + "\n" +
	`          [stco] Size=52 Version=0 Flags=0x000000 EntryCount=9 ChunkOffset=[48, 3836, 4527, 4864, 5043, 5227, 5560, 5702, 6038]` + "\n" +
	`  [trak] Size=844` + "\n" +
	`    [tkhd] Size=92 ... (use "-full tkhd" to show all)` + "\n" +
	`    [edts] Size=36` + "\n" +
	`      [elst] Size=28 Version=0 Flags=0x000000 EntryCount=1 Entries=[{SegmentDurationV0=1000 MediaTimeV0=1024 MediaRateInteger=1}]` + "\n" +
	`    [mdia] Size=708` + "\n" +
	`      [mdhd] Size=32 Version=0 Flags=0x000000 CreationTimeV0=0 ModificationTimeV0=0 Timescale=44100 DurationV0=45124 Language="eng" PreDefined=0` + "\n" +
	`      [hdlr] Size=44 Version=0 Flags=0x000000 PreDefined=0 HandlerType="soun" Name="SoundHandle"` + "\n" +
	`      [minf] Size=624` + "\n" +
	`        [smhd] Size=16 Version=0 Flags=0x000000 Balance=0` + "\n" +
	`        [dinf] Size=36` + "\n" +
	`          [dref] Size=28 Version=0 Flags=0x000000 EntryCount=1` + "\n" +
	`            [url ] Size=12 Version=0 Flags=0x000001` + "\n" +
	`        [stbl] Size=564` + "\n" +
	`          [stsd] Size=106 Version=0 Flags=0x000000 EntryCount=1` + "\n" +
	`            [mp4a] Size=90 DataReferenceIndex=1 EntryVersion=0 ChannelCount=2 SampleSize=16 PreDefined=0 SampleRate=44100` + "\n" +
	`              [esds] Size=54 ... (use "-full esds" to show all)` + "\n" +
	`          [stts] Size=48 ... (use "-full stts" to show all)` + "\n" +
	`          [stsc] Size=100 ... (use "-full stsc" to show all)` + "\n" +
	`          [stsz] Size=196 ... (use "-full stsz" to show all)` + "\n" +
	`          [stco] Size=52 Version=0 Flags=0x000000 EntryCount=9 ChunkOffset=[3813, 4381, 4707, 4933, 5103, 5409, 5582, 5906, 6053]` + "\n" +
	`          [sgpd] Size=26 Version=1 Flags=0x000000 GroupingType="roll" DefaultLength=2 EntryCount=1 RollDistances=[-1]` + "\n" +
	`          [sbgp] Size=28 Version=0 Flags=0x000000 GroupingType=1919904876 EntryCount=1 Entries=[{SampleCount=44 GroupDescriptionIndex=1}]` + "\n" +
	`  [udta] Size=133` + "\n" +
	`    [meta] Size=90 Version=0 Flags=0x000000` + "\n" +
	`      [hdlr] Size=33 Version=0 Flags=0x000000 PreDefined=0 HandlerType="mdir" Name=""` + "\n" +
	`      [ilst] Size=45` + "\n" +
	`        [(c)too] Size=37` + "\n" +
	`          [data] Size=29 DataType=UTF8 DataLang=0 Data="Lavf58.29.100"` + "\n" +
	`    [loci] (unsupported box type) Size=35 Data=[0x00 0x00 0x00 0x00 0x15 0xc7 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x65 0x61 0x72 0x74 0x68 0x00 0x00]` + "\n"

var sampleMP4OutputOffset = "" +
	`[ftyp] Offset=0 Size=32 ... (use "-full ftyp" to show all)` + "\n" +
	`[free] Offset=32 Size=8 Data=[...] (use "-full free" to show all)` + "\n" +
	`[mdat] Offset=40 Size=6402 Data=[...] (use "-full mdat" to show all)` + "\n" +
	`[moov] Offset=6442 Size=1836` + "\n" +
	`  [mvhd] Offset=6450 Size=108 ... (use "-full mvhd" to show all)` + "\n" +
	`  [trak] Offset=6558 Size=743` + "\n" +
	`    [tkhd] Offset=6566 Size=92 ... (use "-full tkhd" to show all)` + "\n" +
	`    [edts] Offset=6658 Size=36` + "\n" +
	`      [elst] Offset=6666 Size=28 Version=0 Flags=0x000000 EntryCount=1 Entries=[{SegmentDurationV0=1000 MediaTimeV0=2048 MediaRateInteger=1}]` + "\n" +
	`    [mdia] Offset=6694 Size=607` + "\n" +
	`      [mdhd] Offset=6702 Size=32 Version=0 Flags=0x000000 CreationTimeV0=0 ModificationTimeV0=0 Timescale=10240 DurationV0=10240 Language="eng" PreDefined=0` + "\n" +
	`      [hdlr] Offset=6734 Size=44 Version=0 Flags=0x000000 PreDefined=0 HandlerType="vide" Name="VideoHandle"` + "\n" +
	`      [minf] Offset=6778 Size=523` + "\n" +
	`        [vmhd] Offset=6786 Size=20 Version=0 Flags=0x000001 Graphicsmode=0 Opcolor=[0, 0, 0]` + "\n" +
	`        [dinf] Offset=6806 Size=36` + "\n" +
	`          [dref] Offset=6814 Size=28 Version=0 Flags=0x000000 EntryCount=1` + "\n" +
	`            [url ] Offset=6830 Size=12 Version=0 Flags=0x000001` + "\n" +
	`        [stbl] Offset=6842 Size=459` + "\n" +
	`          [stsd] Offset=6850 Size=167 Version=0 Flags=0x000000 EntryCount=1` + "\n" +
	`            [avc1] Offset=6866 Size=151 ... (use "-full avc1" to show all)` + "\n" +
	`              [avcC] Offset=6952 Size=49 ... (use "-full avcC" to show all)` + "\n" +
	`              [pasp] Offset=7001 Size=16 HSpacing=1 VSpacing=1` + "\n" +
	`          [stts] Offset=7017 Size=24 Version=0 Flags=0x000000 EntryCount=1 Entries=[{SampleCount=10 SampleDelta=1024}]` + "\n" +
	`          [stss] Offset=7041 Size=20 Version=0 Flags=0x000000 EntryCount=1 SampleNumber=[1]` + "\n" +
	`          [ctts] Offset=7061 Size=88 ... (use "-full ctts" to show all)` + "\n" +
	`          [stsc] Offset=7149 Size=40 ... (use "-full stsc" to show all)` + "\n" +
	`          [stsz] Offset=7189 Size=60 Version=0 Flags=0x000000 SampleSize=0 SampleCount=10 EntrySize=[3679, 86, 545, 180, 69, 60, 182, 22, 204, 15]` + "\n" +
	`          [stco] Offset=7249 Size=52 Version=0 Flags=0x000000 EntryCount=9 ChunkOffset=[48, 3836, 4527, 4864, 5043, 5227, 5560, 5702, 6038]` + "\n" +
	`  [trak] Offset=7301 Size=844` + "\n" +
	`    [tkhd] Offset=7309 Size=92 ... (use "-full tkhd" to show all)` + "\n" +
	`    [edts] Offset=7401 Size=36` + "\n" +
	`      [elst] Offset=7409 Size=28 Version=0 Flags=0x000000 EntryCount=1 Entries=[{SegmentDurationV0=1000 MediaTimeV0=1024 MediaRateInteger=1}]` + "\n" +
	`    [mdia] Offset=7437 Size=708` + "\n" +
	`      [mdhd] Offset=7445 Size=32 Version=0 Flags=0x000000 CreationTimeV0=0 ModificationTimeV0=0 Timescale=44100 DurationV0=45124 Language="eng" PreDefined=0` + "\n" +
	`      [hdlr] Offset=7477 Size=44 Version=0 Flags=0x000000 PreDefined=0 HandlerType="soun" Name="SoundHandle"` + "\n" +
	`      [minf] Offset=7521 Size=624` + "\n" +
	`        [smhd] Offset=7529 Size=16 Version=0 Flags=0x000000 Balance=0` + "\n" +
	`        [dinf] Offset=7545 Size=36` + "\n" +
	`          [dref] Offset=7553 Size=28 Version=0 Flags=0x000000 EntryCount=1` + "\n" +
	`            [url ] Offset=7569 Size=12 Version=0 Flags=0x000001` + "\n" +
	`        [stbl] Offset=7581 Size=564` + "\n" +
	`          [stsd] Offset=7589 Size=106 Version=0 Flags=0x000000 EntryCount=1` + "\n" +
	`            [mp4a] Offset=7605 Size=90 DataReferenceIndex=1 EntryVersion=0 ChannelCount=2 SampleSize=16 PreDefined=0 SampleRate=44100` + "\n" +
	`              [esds] Offset=7641 Size=54 ... (use "-full esds" to show all)` + "\n" +
	`          [stts] Offset=7695 Size=48 ... (use "-full stts" to show all)` + "\n" +
	`          [stsc] Offset=7743 Size=100 ... (use "-full stsc" to show all)` + "\n" +
	`          [stsz] Offset=7843 Size=196 ... (use "-full stsz" to show all)` + "\n" +
	`          [stco] Offset=8039 Size=52 Version=0 Flags=0x000000 EntryCount=9 ChunkOffset=[3813, 4381, 4707, 4933, 5103, 5409, 5582, 5906, 6053]` + "\n" +
	`          [sgpd] Offset=8091 Size=26 Version=1 Flags=0x000000 GroupingType="roll" DefaultLength=2 EntryCount=1 RollDistances=[-1]` + "\n" +
	`          [sbgp] Offset=8117 Size=28 Version=0 Flags=0x000000 GroupingType=1919904876 EntryCount=1 Entries=[{SampleCount=44 GroupDescriptionIndex=1}]` + "\n" +
	`  [udta] Offset=8145 Size=133` + "\n" +
	`    [meta] Offset=8153 Size=90 Version=0 Flags=0x000000` + "\n" +
	`      [hdlr] Offset=8165 Size=33 Version=0 Flags=0x000000 PreDefined=0 HandlerType="mdir" Name=""` + "\n" +
	`      [ilst] Offset=8198 Size=45` + "\n" +
	`        [(c)too] Offset=8206 Size=37` + "\n" +
	`          [data] Offset=8214 Size=29 DataType=UTF8 DataLang=0 Data="Lavf58.29.100"` + "\n" +
	`    [loci] (unsupported box type) Offset=8243 Size=35 Data=[...] (use "-full loci" to show all)` + "\n"

var sampleMP4OutputHex = "" +
	`[ftyp] Size=0x20 MajorBrand="isom" MinorVersion=512 CompatibleBrands=[{CompatibleBrand="isom"}, {CompatibleBrand="iso2"}, {CompatibleBrand="avc1"}, {CompatibleBrand="mp41"}]` + "\n" +
	`[free] Size=0x8 Data=[...] (use "-full free" to show all)` + "\n" +
	`[mdat] Size=0x1902 Data=[...] (use "-full mdat" to show all)` + "\n" +
	`[moov] Size=0x72c` + "\n" +
	`  [mvhd] Size=0x6c ... (use "-full mvhd" to show all)` + "\n" +
	`  [trak] Size=0x2e7` + "\n" +
	`    [tkhd] Size=0x5c ... (use "-full tkhd" to show all)` + "\n" +
	`    [edts] Size=0x24` + "\n" +
	`      [elst] Size=0x1c Version=0 Flags=0x000000 EntryCount=1 Entries=[{SegmentDurationV0=1000 MediaTimeV0=2048 MediaRateInteger=1}]` + "\n" +
	`    [mdia] Size=0x25f` + "\n" +
	`      [mdhd] Size=0x20 Version=0 Flags=0x000000 CreationTimeV0=0 ModificationTimeV0=0 Timescale=10240 DurationV0=10240 Language="eng" PreDefined=0` + "\n" +
	`      [hdlr] Size=0x2c Version=0 Flags=0x000000 PreDefined=0 HandlerType="vide" Name="VideoHandle"` + "\n" +
	`      [minf] Size=0x20b` + "\n" +
	`        [vmhd] Size=0x14 Version=0 Flags=0x000001 Graphicsmode=0 Opcolor=[0, 0, 0]` + "\n" +
	`        [dinf] Size=0x24` + "\n" +
	`          [dref] Size=0x1c Version=0 Flags=0x000000 EntryCount=1` + "\n" +
	`            [url ] Size=0xc Version=0 Flags=0x000001` + "\n" +
	`        [stbl] Size=0x1cb` + "\n" +
	`          [stsd] Size=0xa7 Version=0 Flags=0x000000 EntryCount=1` + "\n" +
	`            [avc1] Size=0x97 ... (use "-full avc1" to show all)` + "\n" +
	`              [avcC] Size=0x31 ... (use "-full avcC" to show all)` + "\n" +
	`              [pasp] Size=0x10 HSpacing=1 VSpacing=1` + "\n" +
	`          [stts] Size=0x18 Version=0 Flags=0x000000 EntryCount=1 Entries=[{SampleCount=10 SampleDelta=1024}]` + "\n" +
	`          [stss] Size=0x14 Version=0 Flags=0x000000 EntryCount=1 SampleNumber=[1]` + "\n" +
	`          [ctts] Size=0x58 ... (use "-full ctts" to show all)` + "\n" +
	`          [stsc] Size=0x28 ... (use "-full stsc" to show all)` + "\n" +
	`          [stsz] Size=0x3c Version=0 Flags=0x000000 SampleSize=0 SampleCount=10 EntrySize=[3679, 86, 545, 180, 69, 60, 182, 22, 204, 15]` + "\n" +
	`          [stco] Size=0x34 Version=0 Flags=0x000000 EntryCount=9 ChunkOffset=[48, 3836, 4527, 4864, 5043, 5227, 5560, 5702, 6038]` + "\n" +
	`  [trak] Size=0x34c` + "\n" +
	`    [tkhd] Size=0x5c ... (use "-full tkhd" to show all)` + "\n" +
	`    [edts] Size=0x24` + "\n" +
	`      [elst] Size=0x1c Version=0 Flags=0x000000 EntryCount=1 Entries=[{SegmentDurationV0=1000 MediaTimeV0=1024 MediaRateInteger=1}]` + "\n" +
	`    [mdia] Size=0x2c4` + "\n" +
	`      [mdhd] Size=0x20 Version=0 Flags=0x000000 CreationTimeV0=0 ModificationTimeV0=0 Timescale=44100 DurationV0=45124 Language="eng" PreDefined=0` + "\n" +
	`      [hdlr] Size=0x2c Version=0 Flags=0x000000 PreDefined=0 HandlerType="soun" Name="SoundHandle"` + "\n" +
	`      [minf] Size=0x270` + "\n" +
	`        [smhd] Size=0x10 Version=0 Flags=0x000000 Balance=0` + "\n" +
	`        [dinf] Size=0x24` + "\n" +
	`          [dref] Size=0x1c Version=0 Flags=0x000000 EntryCount=1` + "\n" +
	`            [url ] Size=0xc Version=0 Flags=0x000001` + "\n" +
	`        [stbl] Size=0x234` + "\n" +
	`          [stsd] Size=0x6a Version=0 Flags=0x000000 EntryCount=1` + "\n" +
	`            [mp4a] Size=0x5a DataReferenceIndex=1 EntryVersion=0 ChannelCount=2 SampleSize=16 PreDefined=0 SampleRate=44100` + "\n" +
	`              [esds] Size=0x36 ... (use "-full esds" to show all)` + "\n" +
	`          [stts] Size=0x30 ... (use "-full stts" to show all)` + "\n" +
	`          [stsc] Size=0x64 ... (use "-full stsc" to show all)` + "\n" +
	`          [stsz] Size=0xc4 ... (use "-full stsz" to show all)` + "\n" +
	`          [stco] Size=0x34 Version=0 Flags=0x000000 EntryCount=9 ChunkOffset=[3813, 4381, 4707, 4933, 5103, 5409, 5582, 5906, 6053]` + "\n" +
	`          [sgpd] Size=0x1a Version=1 Flags=0x000000 GroupingType="roll" DefaultLength=2 EntryCount=1 RollDistances=[-1]` + "\n" +
	`          [sbgp] Size=0x1c Version=0 Flags=0x000000 GroupingType=1919904876 EntryCount=1 Entries=[{SampleCount=44 GroupDescriptionIndex=1}]` + "\n" +
	`  [udta] Size=0x85` + "\n" +
	`    [meta] Size=0x5a Version=0 Flags=0x000000` + "\n" +
	`      [hdlr] Size=0x21 Version=0 Flags=0x000000 PreDefined=0 HandlerType="mdir" Name=""` + "\n" +
	`      [ilst] Size=0x2d` + "\n" +
	`        [(c)too] Size=0x25` + "\n" +
	`          [data] Size=0x1d DataType=UTF8 DataLang=0 Data="Lavf58.29.100"` + "\n" +
	`    [loci] (unsupported box type) Size=0x23 Data=[...] (use "-full loci" to show all)` + "\n"
