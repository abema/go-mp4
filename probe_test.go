package mp4

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
