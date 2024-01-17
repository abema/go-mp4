package psshdump

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPsshdump(t *testing.T) {
	testCases := []struct {
		name    string
		file    string
		options []string
		wants   string
	}{
		{
			name: "sample_init.encv.mp4",
			file: "../../../../testdata/sample_init.encv.mp4",
			wants: "0:\n" +
				"  offset: 1307\n" +
				"  size: 52\n" +
				"  version: 1\n" +
				"  flags: 0x000000\n" +
				"  systemId: \n" +
				"  dataSize: 0\n" +
				"  base64: \"AAAANHBzc2gBAAAAEHfv7MCyTQKs4zweUuL7SwAAAAEBI0VniavN7wEjRWeJq83vAAAAAA==\"\n" +
				"\n",
		},
		{
			name: "sample_init.encv.mp4",
			file: "../../../../testdata/sample_init.enca.mp4",
			wants: "0:\n" +
				"  offset: 1307\n" +
				"  size: 52\n" +
				"  version: 1\n" +
				"  flags: 0x000000\n" +
				"  systemId: \n" +
				"  dataSize: 0\n" +
				"  base64: \"AAAANHBzc2gBAAAAEHfv7MCyTQKs4zweUuL7SwAAAAEBI0VniavN7wEjRWeJq83vAAAAAA==\"\n" +
				"\n",
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
			b, err := io.ReadAll(r)
			require.NoError(t, err)
			assert.Equal(t, tc.wants, string(b))
		})
	}
}
