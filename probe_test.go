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

	assert.Equal(t, 2, len(info.Tracks))
	assert.Equal(t, 8, len(info.Segments))
}
