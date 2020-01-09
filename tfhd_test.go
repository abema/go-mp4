package mp4

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTfhdMarshal(t *testing.T) {
	tfhd := Tfhd{
		FullBox: FullBox{
			Version: 0,
			Flags:   [3]byte{0x00, 0x00, 0x00},
		},
		TrackID:                0x08404649,
		BaseDataOffset:         0x0123456789abcdef,
		SampleDescriptionIndex: 0x12345678,
		DefaultSampleDuration:  0x23456789,
		DefaultSampleSize:      0x3456789a,
		DefaultSampleFlags:     0x456789ab,
	}
	expect := []byte{
		0,                // version
		0x00, 0x00, 0x00, // flags
		0x08, 0x40, 0x46, 0x49, // track ID
	}

	buf := bytes.NewBuffer(nil)
	n, err := Marshal(buf, &tfhd)
	require.NoError(t, err)
	assert.Equal(t, uint64(len(expect)), n)
	assert.Equal(t, expect, buf.Bytes())

	tfhd.SetFlags(TfhdBaseDataOffsetPresent | TfhdDefaultSampleDurationPresent)
	expect = []byte{
		0,                // version
		0x00, 0x00, 0x09, // flags (0000 0000 1001)
		0x08, 0x40, 0x46, 0x49, // track ID
		0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
		0x23, 0x45, 0x67, 0x89,
	}

	buf = bytes.NewBuffer(nil)
	n, err = Marshal(buf, &tfhd)
	require.NoError(t, err)
	assert.Equal(t, uint64(len(expect)), n)
	assert.Equal(t, expect, buf.Bytes())
}

func TestTfhdUnmarshal(t *testing.T) {
	data := []byte{
		0,                // version
		0x00, 0x00, 0x00, // flags (0000 0000 1001)
		0x08, 0x40, 0x46, 0x49, // track ID
	}

	buf := bytes.NewReader(data)
	tfhd := Tfhd{}
	n, err := Unmarshal(buf, uint64(buf.Len()), &tfhd)
	require.NoError(t, err)
	assert.Equal(t, uint64(len(data)), n)
	assert.Equal(t, uint8(0), tfhd.Version)
	assert.Equal(t, uint32(0x00), tfhd.GetFlags())
	assert.Equal(t, uint32(0x08404649), tfhd.TrackID)

	data = []byte{
		0,                // version
		0x00, 0x00, 0x09, // flags (0000 0000 1001)
		0x08, 0x40, 0x46, 0x49, // track ID
		0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
		0x23, 0x45, 0x67, 0x89,
	}

	buf = bytes.NewReader(data)
	tfhd = Tfhd{}
	n, err = Unmarshal(buf, uint64(buf.Len()), &tfhd)
	require.NoError(t, err)
	assert.Equal(t, uint64(len(data)), n)
	assert.Equal(t, uint8(0), tfhd.Version)
	assert.Equal(t, uint32(0x09), tfhd.GetFlags())
	assert.Equal(t, uint32(0x08404649), tfhd.TrackID)
	assert.Equal(t, uint64(0x0123456789abcdef), tfhd.BaseDataOffset)
	assert.Equal(t, uint32(0x23456789), tfhd.DefaultSampleDuration)
}
