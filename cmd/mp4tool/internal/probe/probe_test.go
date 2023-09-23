package probe

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProbe(t *testing.T) {
	testCases := []struct {
		name    string
		file    string
		options []string
		wants   string
	}{
		{
			name:  "sample.mp4 no-options",
			file:  "../../../../testdata/sample.mp4",
			wants: sampleMP4JSONOutput,
		},
		{
			name:    "sample.mp4 format-json",
			file:    "../../../../testdata/sample.mp4",
			options: []string{"-format", "json"},
			wants:   sampleMP4JSONOutput,
		},
		{
			name:    "sample.mp4 format-json",
			file:    "../../../../testdata/sample.mp4",
			options: []string{"-format", "yaml"},
			wants:   sampleMP4YamlOutput,
		},
	}
	for _, tc := range testCases {
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
	}
}

var sampleMP4JSONOutput = "" +
	`{` + "\n" +
	`  "MajorBrand": "isom",` + "\n" +
	`  "MinorVersion": 512,` + "\n" +
	`  "CompatibleBrands": [` + "\n" +
	`    "isom",` + "\n" +
	`    "iso2",` + "\n" +
	`    "avc1",` + "\n" +
	`    "mp41"` + "\n" +
	`  ],` + "\n" +
	`  "FastStart": false,` + "\n" +
	`  "Timescale": 1000,` + "\n" +
	`  "Duration": 1024,` + "\n" +
	`  "DurationSeconds": 1.024,` + "\n" +
	`  "Tracks": [` + "\n" +
	`    {` + "\n" +
	`      "TrackID": 1,` + "\n" +
	`      "Timescale": 10240,` + "\n" +
	`      "Duration": 10240,` + "\n" +
	`      "DurationSeconds": 1,` + "\n" +
	`      "Codec": "avc1.64000C",` + "\n" +
	`      "Encrypted": false,` + "\n" +
	`      "Width": 320,` + "\n" +
	`      "Height": 180,` + "\n" +
	`      "SampleNum": 10,` + "\n" +
	`      "ChunkNum": 9,` + "\n" +
	`      "IDRFrameNum": 1,` + "\n" +
	`      "Bitrate": 40336,` + "\n" +
	`      "MaxBitrate": 40336` + "\n" +
	`    },` + "\n" +
	`    {` + "\n" +
	`      "TrackID": 2,` + "\n" +
	`      "Timescale": 44100,` + "\n" +
	`      "Duration": 45124,` + "\n" +
	`      "DurationSeconds": 1.02322,` + "\n" +
	`      "Codec": "mp4a.40.2",` + "\n" +
	`      "Encrypted": false,` + "\n" +
	`      "SampleNum": 44,` + "\n" +
	`      "ChunkNum": 9,` + "\n" +
	`      "Bitrate": 10570,` + "\n" +
	`      "MaxBitrate": 10632` + "\n" +
	`    }` + "\n" +
	`  ]` + "\n" +
	`}` + "\n"

var sampleMP4YamlOutput = "" +
	`major_brand: isom` + "\n" +
	`minor_version: 512` + "\n" +
	`compatible_brands:` + "\n" +
	`- isom` + "\n" +
	`- iso2` + "\n" +
	`- avc1` + "\n" +
	`- mp41` + "\n" +
	`fast_start: false` + "\n" +
	`timescale: 1000` + "\n" +
	`duration: 1024` + "\n" +
	`duration_seconds: 1.024` + "\n" +
	`tracks:` + "\n" +
	`- track_id: 1` + "\n" +
	`  timescale: 10240` + "\n" +
	`  duration: 10240` + "\n" +
	`  duration_seconds: 1` + "\n" +
	`  codec: avc1.64000C` + "\n" +
	`  encrypted: false` + "\n" +
	`  width: 320` + "\n" +
	`  height: 180` + "\n" +
	`  sample_num: 10` + "\n" +
	`  chunk_num: 9` + "\n" +
	`  idr_frame_num: 1` + "\n" +
	`  bitrate: 40336` + "\n" +
	`  max_bitrate: 40336` + "\n" +
	`- track_id: 2` + "\n" +
	`  timescale: 44100` + "\n" +
	`  duration: 45124` + "\n" +
	`  duration_seconds: 1.02322` + "\n" +
	`  codec: mp4a.40.2` + "\n" +
	`  encrypted: false` + "\n" +
	`  sample_num: 44` + "\n" +
	`  chunk_num: 9` + "\n" +
	`  bitrate: 10570` + "\n" +
	`  max_bitrate: 10632` + "\n"
