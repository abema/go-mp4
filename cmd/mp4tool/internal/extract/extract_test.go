package extract

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtract(t *testing.T) {
	testCases := []struct {
		name         string
		file         string
		boxType      string
		expectedSize int
	}{
		{
			name:         "sample.mp4/ftyp",
			file:         "../../../../testdata/sample.mp4",
			boxType:      "ftyp",
			expectedSize: 32,
		},
		{
			name:         "sample.mp4/mdhd",
			file:         "../../../../testdata/sample.mp4",
			boxType:      "mdhd",
			expectedSize: 64, // = 32 (1st trak) + 32 (2nd trak)
		},
		{
			name:         "sample_fragmented.mp4/trun",
			file:         "../../../../testdata/sample_fragmented.mp4",
			boxType:      "trun",
			expectedSize: 452,
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
			require.Zero(t, Main([]string{tc.boxType, tc.file}))
			w.Close()
			b, err := io.ReadAll(r)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedSize, len(b))
			assert.Equal(t, tc.boxType, string(b[4:8]))
		})
	}
}

func TestValidation(t *testing.T) {
	// valid
	require.Zero(t, Main([]string{"xxxx", "../../../../testdata/sample.mp4"}))

	// invalid
	require.NotZero(t, Main([]string{}))
	require.NotZero(t, Main([]string{"xxxx"}))
	require.NotZero(t, Main([]string{"xxxxx", "../../../../testdata/sample.mp4"}))
	require.NotZero(t, Main([]string{"xxxx", "not_found.mp4"}))
}
