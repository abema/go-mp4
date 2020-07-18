package mp4

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractBoxWithPayload(t *testing.T) {
	testCases := []struct {
		name     string
		path     BoxPath
		want     []*BoxInfoWithPayload
		hasError bool
	}{
		{
			name:     "empty box path",
			path:     BoxPath{},
			hasError: true,
		},
		{
			name: "invalid box path",
			path: BoxPath{BoxTypeUdta()},
			want: []*BoxInfoWithPayload{},
		},
		{
			name: "top level",
			path: BoxPath{BoxTypeMoov()},
			want: []*BoxInfoWithPayload{
				{
					Info:    BoxInfo{Offset: 6442, Size: 1836, HeaderSize: 8, Type: BoxTypeMoov()},
					Payload: &Moov{},
				},
			},
		},
		{
			name: "multi hit",
			path: BoxPath{BoxTypeMoov(), BoxTypeTrak(), BoxTypeMdia(), BoxTypeHdlr()},
			want: []*BoxInfoWithPayload{
				{
					Info: BoxInfo{Offset: 6734, Size: 44, HeaderSize: 8, Type: BoxTypeHdlr()},
					Payload: &Hdlr{
						HandlerType: [4]byte{'v', 'i', 'd', 'e'},
						Name:        "VideoHandle",
					},
				},
				{
					Info: BoxInfo{Offset: 7477, Size: 44, HeaderSize: 8, Type: BoxTypeHdlr()},
					Payload: &Hdlr{
						HandlerType: [4]byte{'s', 'o', 'u', 'n'},
						Name:        "SoundHandle",
					},
				},
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			f, err := os.Open("./_examples/sample.mp4")
			require.NoError(t, err)
			defer f.Close()

			bs, err := ExtractBoxWithPayload(f, nil, c.path)
			if c.hasError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, c.want, bs)
		})
	}
}

func TestExtractBox(t *testing.T) {
	testCases := []struct {
		name     string
		path     BoxPath
		want     []*BoxInfo
		hasError bool
	}{
		{
			name:     "empty box path",
			path:     BoxPath{},
			hasError: true,
		},
		{
			name: "invalid box path",
			path: BoxPath{BoxTypeUdta()},
			want: []*BoxInfo{},
		},
		{
			name: "top level",
			path: BoxPath{BoxTypeMoov()},
			want: []*BoxInfo{
				{Offset: 6442, Size: 1836, HeaderSize: 8, Type: BoxTypeMoov()},
			},
		},
		{
			name: "multi hit",
			path: BoxPath{BoxTypeMoov(), BoxTypeTrak(), BoxTypeTkhd()},
			want: []*BoxInfo{
				{Offset: 6566, Size: 92, HeaderSize: 8, Type: BoxTypeTkhd()},
				{Offset: 7309, Size: 92, HeaderSize: 8, Type: BoxTypeTkhd()},
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			f, err := os.Open("./_examples/sample.mp4")
			require.NoError(t, err)
			defer f.Close()

			boxes, err := ExtractBox(f, nil, c.path)
			if c.hasError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, c.want, boxes)
		})
	}
}

func TestExtractBoxes(t *testing.T) {
	testCases := []struct {
		name     string
		paths    []BoxPath
		want     []*BoxInfo
		hasError bool
	}{
		{
			name:  "empty path list",
			paths: []BoxPath{},
		},
		{
			name: "contains empty path",
			paths: []BoxPath{
				{BoxTypeMoov()},
				{},
			},
			hasError: true,
		},
		{
			name: "single path",
			paths: []BoxPath{
				{BoxTypeMoov(), BoxTypeUdta()},
			},
			want: []*BoxInfo{
				{Offset: 8145, Size: 133, HeaderSize: 8, Type: BoxTypeUdta()},
			},
		},
		{
			name: "multi path",
			paths: []BoxPath{
				{BoxTypeMoov()},
				{BoxTypeMoov(), BoxTypeUdta()},
			},
			want: []*BoxInfo{
				{Offset: 6442, Size: 1836, HeaderSize: 8, Type: BoxTypeMoov()},
				{Offset: 8145, Size: 133, HeaderSize: 8, Type: BoxTypeUdta()},
			},
		},
		{
			name: "multi hit",
			paths: []BoxPath{
				{BoxTypeMoov(), BoxTypeTrak(), BoxTypeTkhd()},
			},
			want: []*BoxInfo{
				{Offset: 6566, Size: 92, HeaderSize: 8, Type: BoxTypeTkhd()},
				{Offset: 7309, Size: 92, HeaderSize: 8, Type: BoxTypeTkhd()},
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			f, err := os.Open("./_examples/sample.mp4")
			require.NoError(t, err)
			defer f.Close()

			boxes, err := ExtractBoxes(f, nil, c.paths)
			if c.hasError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, c.want, boxes)
		})
	}
}

func TestExtractDescendantBox(t *testing.T) {
	f, err := os.Open("./_examples/sample.mp4")
	require.NoError(t, err)
	defer f.Close()

	boxes, err := ExtractBox(f, nil, BoxPath{BoxTypeMoov()})
	require.NoError(t, err)
	require.Equal(t, 1, len(boxes))

	descs, err := ExtractBox(f, boxes[0], BoxPath{BoxTypeTrak(), BoxTypeMdia()})
	require.NoError(t, err)
	require.Equal(t, 2, len(descs))
}
