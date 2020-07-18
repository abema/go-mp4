package mp4

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractBox(t *testing.T) {
	testCases := []struct {
		name     string
		path     BoxPath
		types    []BoxType
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
		},
		{
			name:  "top level",
			path:  BoxPath{BoxTypeMoov()},
			types: []BoxType{BoxTypeMoov()},
		},
		{
			name:  "multi hit",
			path:  BoxPath{BoxTypeMoov(), BoxTypeTrak(), BoxTypeTkhd()},
			types: []BoxType{BoxTypeTkhd(), BoxTypeTkhd()},
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
			assert.Equal(t, len(c.types), len(boxes))
			for bi := range boxes {
				assert.Equal(t, c.types[bi], boxes[bi].Type)
			}
		})
	}
}

func TestExtractBoxes(t *testing.T) {
	testCases := []struct {
		name     string
		paths    []BoxPath
		types    []BoxType
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
			types: []BoxType{BoxTypeUdta()},
		},
		{
			name: "multi path",
			paths: []BoxPath{
				{BoxTypeMoov()},
				{BoxTypeMoov(), BoxTypeUdta()},
			},
			types: []BoxType{BoxTypeMoov(), BoxTypeUdta()},
		},
		{
			name: "multi hit",
			paths: []BoxPath{
				{BoxTypeMoov(), BoxTypeTrak(), BoxTypeTkhd()},
			},
			types: []BoxType{BoxTypeTkhd(), BoxTypeTkhd()},
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
			assert.Equal(t, len(c.types), len(boxes))
			for bi := range boxes {
				assert.Equal(t, c.types[bi], boxes[bi].Type)
			}
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
