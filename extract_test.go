package mp4

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractBox(t *testing.T) {
	patterns := []struct {
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

	for _, p := range patterns {
		func() {
			f, err := os.Open("./_examples/sample.mp4")
			require.NoError(t, err)
			defer f.Close()

			boxes, err := ExtractBox(f, nil, p.path)
			if p.hasError {
				require.Error(t, err, p.name)
				return
			}
			require.NoError(t, err, p.name)
			assert.Equal(t, len(p.types), len(boxes), p.name)
			for bi := range boxes {
				assert.Equal(t, p.types[bi], boxes[bi].Type, p.name)
			}
		}()
	}
}

func TestExtractBoxes(t *testing.T) {
	patterns := []struct {
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

	for _, p := range patterns {
		func() {
			f, err := os.Open("./_examples/sample.mp4")
			require.NoError(t, err)
			defer f.Close()

			boxes, err := ExtractBoxes(f, nil, p.paths)
			if p.hasError {
				require.Error(t, err, p.name)
				return
			}

			require.NoError(t, err, p.name)
			assert.Equal(t, len(p.types), len(boxes), p.name)
			for bi := range boxes {
				assert.Equal(t, p.types[bi], boxes[bi].Type, p.name)
			}
		}()
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
