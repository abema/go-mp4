package mp4

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProbe(t *testing.T) {
	f, err := os.Open("./_examples/sample.mp4")
	require.NoError(t, err)
	defer f.Close()

	info, err := Probe(f)
	require.NoError(t, err)

	assert.Equal(t, BrandISOM(), info.MajorBrand)
	assert.Equal(t, uint32(0x200), info.MinorVersion)
	require.Len(t, info.CompatibleBrands, 4)
	assert.Equal(t, BrandISOM(), info.CompatibleBrands[0])
	assert.Equal(t, BrandISO2(), info.CompatibleBrands[1])
	assert.Equal(t, BrandAVC1(), info.CompatibleBrands[2])
	assert.Equal(t, BrandMP41(), info.CompatibleBrands[3])
	assert.False(t, info.FastStart)
	assert.Equal(t, uint32(1000), info.Timescale)
	assert.Equal(t, uint64(1024), info.Duration)

	require.Len(t, info.Tracks, 2)

	assert.Equal(t, uint32(1), info.Tracks[0].TrackID)
	assert.Equal(t, uint32(10240), info.Tracks[0].Timescale)
	assert.Equal(t, uint64(10240), info.Tracks[0].Duration)
	assert.Equal(t, CodecAVC1, info.Tracks[0].Codec)
	assert.Equal(t, uint8(1), info.Tracks[0].AVC.ConfigurationVersion)
	assert.Equal(t, uint8(0x64), info.Tracks[0].AVC.Profile)
	assert.Equal(t, uint8(0), info.Tracks[0].AVC.ProfileCompatibility)
	assert.Equal(t, uint8(0xc), info.Tracks[0].AVC.Level)
	assert.Equal(t, uint16(0x04), info.Tracks[0].AVC.LengthSize)
	assert.Equal(t, uint16(320), info.Tracks[0].AVC.Width)
	assert.Equal(t, uint16(180), info.Tracks[0].AVC.Height)
	assert.False(t, info.Tracks[0].Encrypted)
	require.Len(t, info.Tracks[0].EditList, 1)
	assert.Equal(t, int64(2048), info.Tracks[0].EditList[0].MediaTime)
	assert.Equal(t, uint64(1000), info.Tracks[0].EditList[0].SegmentDuration)
	require.Len(t, info.Tracks[0].Samples, 10)
	assert.Equal(t, uint32(3679), info.Tracks[0].Samples[0].Size)
	assert.Equal(t, uint32(15), info.Tracks[0].Samples[9].Size)
	assert.Equal(t, uint32(1024), info.Tracks[0].Samples[0].TimeDelta)
	assert.Equal(t, uint32(1024), info.Tracks[0].Samples[9].TimeDelta)
	assert.Equal(t, int64(2048), info.Tracks[0].Samples[0].CompositionTimeOffset)
	assert.Equal(t, int64(1024), info.Tracks[0].Samples[9].CompositionTimeOffset)
	require.Len(t, info.Tracks[0].Chunks, 9)
	assert.Equal(t, uint32(48), info.Tracks[0].Chunks[0].DataOffset)
	assert.Equal(t, uint32(6038), info.Tracks[0].Chunks[8].DataOffset)
	assert.Equal(t, uint32(2), info.Tracks[0].Chunks[0].SamplesPerChunk)
	assert.Equal(t, uint32(1), info.Tracks[0].Chunks[8].SamplesPerChunk)

	assert.Equal(t, uint32(2), info.Tracks[1].TrackID)
	assert.Equal(t, uint32(44100), info.Tracks[1].Timescale)
	assert.Equal(t, uint64(45124), info.Tracks[1].Duration)
	assert.Equal(t, CodecMP4A, info.Tracks[1].Codec)
	assert.Equal(t, uint8(40), info.Tracks[1].MP4A.OTI)
	assert.Equal(t, uint8(2), info.Tracks[1].MP4A.AudOTI)
	assert.Equal(t, uint16(2), info.Tracks[1].MP4A.ChannelCount)
	assert.False(t, info.Tracks[1].Encrypted)

	require.Len(t, info.Segments, 0)

	idxs, err := FindIDRFrames(f, info.Tracks[0])
	require.NoError(t, err)
	require.Len(t, idxs, 1)
	assert.Equal(t, 0, idxs[0])
}

func TestProbeWithFMP4(t *testing.T) {
	f, err := os.Open("./_examples/sample_fragmented.mp4")
	require.NoError(t, err)
	defer f.Close()

	info, err := Probe(f)
	require.NoError(t, err)

	require.Equal(t, 2, len(info.Tracks))
	require.Equal(t, 8, len(info.Segments))

	assert.Equal(t, uint32(1), info.Tracks[0].TrackID)
	assert.Equal(t, uint32(90000), info.Tracks[0].Timescale)

	assert.Equal(t, uint32(2), info.Tracks[1].TrackID)
	assert.Equal(t, uint32(44100), info.Tracks[1].Timescale)

	assert.Equal(t, uint32(1), info.Segments[0].TrackID)
	assert.Equal(t, uint64(1227), info.Segments[0].MoofOffset)
	assert.Equal(t, uint64(0), info.Segments[0].BaseMediaDecodeTime)
	assert.Equal(t, uint32(9000), info.Segments[0].DefaultSampleDuration)
	assert.Equal(t, uint32(3), info.Segments[0].SampleCount)
	assert.Equal(t, uint32(27000), info.Segments[0].Duration)
	assert.Equal(t, int32(18000), info.Segments[0].CompositionTimeOffset)
	assert.Equal(t, uint32(1054), info.Segments[0].Size)

	assert.Equal(t, uint32(2), info.Segments[1].TrackID)
	assert.Equal(t, uint64(2417), info.Segments[1].MoofOffset)
	assert.Equal(t, uint64(0), info.Segments[1].BaseMediaDecodeTime)
	assert.Equal(t, uint32(8830), info.Segments[1].DefaultSampleDuration)
	assert.Equal(t, uint32(5), info.Segments[1].SampleCount)
	assert.Equal(t, uint32(13407), info.Segments[1].Duration)
	assert.Equal(t, int32(0), info.Segments[1].CompositionTimeOffset)
	assert.Equal(t, uint32(177), info.Segments[1].Size)

	assert.Equal(t, uint32(1), info.Segments[2].TrackID)
	assert.Equal(t, uint64(2742), info.Segments[2].MoofOffset)
	assert.Equal(t, uint64(27000), info.Segments[2].BaseMediaDecodeTime)
	assert.Equal(t, uint32(9000), info.Segments[2].DefaultSampleDuration)
	assert.Equal(t, uint32(2), info.Segments[2].SampleCount)
	assert.Equal(t, uint32(18000), info.Segments[2].Duration)
	assert.Equal(t, int32(18000), info.Segments[2].CompositionTimeOffset)
	assert.Equal(t, uint32(282), info.Segments[2].Size)

	assert.Equal(t, uint32(2), info.Segments[3].TrackID)
	assert.Equal(t, uint64(3152), info.Segments[3].MoofOffset)
	assert.Equal(t, uint64(13407), info.Segments[3].BaseMediaDecodeTime)
	assert.Equal(t, uint32(1024), info.Segments[3].DefaultSampleDuration)
	assert.Equal(t, uint32(9), info.Segments[3].SampleCount)
	assert.Equal(t, uint32(9216), info.Segments[3].Duration)
	assert.Equal(t, int32(0), info.Segments[3].CompositionTimeOffset)
	assert.Equal(t, uint32(271), info.Segments[3].Size)
}

func TestProbeFra(t *testing.T) {
	f, err := os.Open("./_examples/sample_fragmented.mp4")
	require.NoError(t, err)
	defer f.Close()

	info, err := ProbeFra(f)
	require.NoError(t, err)

	require.Equal(t, 2, len(info.Tracks))
	require.Equal(t, 8, len(info.Segments))

	assert.Equal(t, uint32(1), info.Tracks[0].TrackID)
	assert.Equal(t, uint32(90000), info.Tracks[0].Timescale)

	assert.Equal(t, uint32(2), info.Tracks[1].TrackID)
	assert.Equal(t, uint32(44100), info.Tracks[1].Timescale)

	assert.Equal(t, uint32(1), info.Segments[0].TrackID)
	assert.Equal(t, uint64(1227), info.Segments[0].MoofOffset)
	assert.Equal(t, uint64(0), info.Segments[0].BaseMediaDecodeTime)
	assert.Equal(t, uint32(9000), info.Segments[0].DefaultSampleDuration)
	assert.Equal(t, uint32(3), info.Segments[0].SampleCount)
	assert.Equal(t, uint32(27000), info.Segments[0].Duration)
	assert.Equal(t, int32(18000), info.Segments[0].CompositionTimeOffset)

	assert.Equal(t, uint32(2), info.Segments[1].TrackID)
	assert.Equal(t, uint64(2417), info.Segments[1].MoofOffset)
	assert.Equal(t, uint64(0), info.Segments[1].BaseMediaDecodeTime)
	assert.Equal(t, uint32(8830), info.Segments[1].DefaultSampleDuration)
	assert.Equal(t, uint32(5), info.Segments[1].SampleCount)
	assert.Equal(t, uint32(13407), info.Segments[1].Duration)
	assert.Equal(t, int32(0), info.Segments[1].CompositionTimeOffset)

	assert.Equal(t, uint32(1), info.Segments[2].TrackID)
	assert.Equal(t, uint64(2742), info.Segments[2].MoofOffset)
	assert.Equal(t, uint64(27000), info.Segments[2].BaseMediaDecodeTime)
	assert.Equal(t, uint32(9000), info.Segments[2].DefaultSampleDuration)
	assert.Equal(t, uint32(2), info.Segments[2].SampleCount)
	assert.Equal(t, uint32(18000), info.Segments[2].Duration)
	assert.Equal(t, int32(18000), info.Segments[2].CompositionTimeOffset)

	assert.Equal(t, uint32(2), info.Segments[3].TrackID)
	assert.Equal(t, uint64(3152), info.Segments[3].MoofOffset)
	assert.Equal(t, uint64(13407), info.Segments[3].BaseMediaDecodeTime)
	assert.Equal(t, uint32(1024), info.Segments[3].DefaultSampleDuration)
	assert.Equal(t, uint32(9), info.Segments[3].SampleCount)
	assert.Equal(t, uint32(9216), info.Segments[3].Duration)
	assert.Equal(t, int32(0), info.Segments[3].CompositionTimeOffset)
}

func TestSamplesGetBitrate(t *testing.T) {
	assert.Equal(t, uint64(0), Samples{}.GetBitrate(100))

	assert.Equal(t, uint64(14400), // = 900 / (50 / 100) * 8
		Samples{
			{TimeDelta: 10, Size: 100},
			{TimeDelta: 10, Size: 200},
			{TimeDelta: 10, Size: 300},
			{TimeDelta: 10, Size: 100},
			{TimeDelta: 10, Size: 200},
		}.GetBitrate(100))
}

func TestSamplesGetMaxBitrate(t *testing.T) {
	assert.Equal(t, uint64(0), Samples{}.GetMaxBitrate(100, 20))

	assert.Equal(t, uint64(20000), // = 500 / (20 / 100) * 8
		Samples{
			{TimeDelta: 10, Size: 100},
			{TimeDelta: 10, Size: 200},
			{TimeDelta: 10, Size: 300},
			{TimeDelta: 10, Size: 100},
			{TimeDelta: 10, Size: 200},
		}.GetMaxBitrate(100, 20))
}

func TestSegmentsGetBitrate(t *testing.T) {
	assert.Equal(t, uint64(0), Segments{}.GetBitrate(2, 100))

	assert.Equal(t, uint64(14400), // = 900 / (50 / 100) * 8
		Segments{
			{TrackID: 1, Duration: 10, Size: 300},
			{TrackID: 2, Duration: 10, Size: 100},
			{TrackID: 2, Duration: 10, Size: 200},
			{TrackID: 1, Duration: 10, Size: 200},
			{TrackID: 2, Duration: 10, Size: 300},
			{TrackID: 3, Duration: 10, Size: 700},
			{TrackID: 2, Duration: 10, Size: 100},
			{TrackID: 1, Duration: 10, Size: 800},
			{TrackID: 2, Duration: 10, Size: 200},
		}.GetBitrate(2, 100))
}

func TestSegmentsGetMaxBitrate(t *testing.T) {
	assert.Equal(t, uint64(0), Segments{}.GetMaxBitrate(2, 100))

	assert.Equal(t, uint64(24000), // = 300 / (10 / 100) * 8
		Segments{
			{TrackID: 1, Duration: 10, Size: 300},
			{TrackID: 2, Duration: 10, Size: 100},
			{TrackID: 2, Duration: 10, Size: 200},
			{TrackID: 1, Duration: 10, Size: 200},
			{TrackID: 2, Duration: 10, Size: 300},
			{TrackID: 3, Duration: 10, Size: 700},
			{TrackID: 2, Duration: 10, Size: 100},
			{TrackID: 1, Duration: 10, Size: 800},
			{TrackID: 2, Duration: 10, Size: 200},
		}.GetMaxBitrate(2, 100))
}
