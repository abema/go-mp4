package mp4

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"gopkg.in/src-d/go-billy.v4/memfs"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriter(t *testing.T) {
	output, err := memfs.New().Create("output.mp4")
	require.NoError(t, err)
	defer output.Close()
	w := NewWriter(output)

	// start ftyp
	bi, err := w.StartBox(&BoxInfo{Type: BoxTypeFtyp()})
	require.NoError(t, err)
	assert.Equal(t, uint64(0), bi.Offset)
	assert.Equal(t, uint64(8), bi.Size)

	_, err = Marshal(w, &Ftyp{
		MajorBrand:   [4]byte{'a', 'b', 'e', 'm'},
		MinorVersion: 0x12345678,
		CompatibleBrands: []CompatibleBrandElem{
			{CompatibleBrand: [4]byte{'a', 'b', 'c', 'd'}},
			{CompatibleBrand: [4]byte{'e', 'f', 'g', 'h'}},
		},
	}, Context{})
	require.NoError(t, err)

	// end ftyp
	bi, err = w.EndBox()
	require.NoError(t, err)
	assert.Equal(t, uint64(0), bi.Offset)
	assert.Equal(t, uint64(24), bi.Size)

	// start moov
	bi, err = w.StartBox(&BoxInfo{Type: BoxTypeMoov()})
	require.NoError(t, err)
	assert.Equal(t, uint64(24), bi.Offset)
	assert.Equal(t, uint64(8), bi.Size)

	// copy
	err = w.CopyBox(bytes.NewReader([]byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x0a,
		'u', 'd', 't', 'a',
		0x01, 0x02, 0x03, 0x04,
		0x05, 0x06, 0x07, 0x08,
	}), &BoxInfo{Offset: 6, Size: 15})
	require.NoError(t, err)

	// start trak
	bi, err = w.StartBox(&BoxInfo{Type: BoxTypeTrak()})
	require.NoError(t, err)
	assert.Equal(t, uint64(47), bi.Offset)
	assert.Equal(t, uint64(8), bi.Size)

	// start tkhd
	bi, err = w.StartBox(&BoxInfo{Type: BoxTypeTkhd()})
	require.NoError(t, err)
	assert.Equal(t, uint64(55), bi.Offset)
	assert.Equal(t, uint64(8), bi.Size)

	_, err = Marshal(w, &Tkhd{
		CreationTimeV0:     1,
		ModificationTimeV0: 2,
		TrackID:            3,
		DurationV0:         4,
		Layer:              5,
		AlternateGroup:     6,
		Volume:             7,
		Width:              8,
		Height:             9,
	}, Context{})
	require.NoError(t, err)

	// end tkhd
	bi, err = w.EndBox()
	require.NoError(t, err)
	assert.Equal(t, uint64(55), bi.Offset)
	assert.Equal(t, uint64(92), bi.Size)

	// end trak
	bi, err = w.EndBox()
	require.NoError(t, err)
	assert.Equal(t, uint64(47), bi.Offset)
	assert.Equal(t, uint64(100), bi.Size)

	// end moov
	bi, err = w.EndBox()
	require.NoError(t, err)
	assert.Equal(t, uint64(24), bi.Offset)
	assert.Equal(t, uint64(123), bi.Size)

	_, err = output.Seek(0, io.SeekStart)
	require.NoError(t, err)
	bin, err := ioutil.ReadAll(output)
	require.NoError(t, err)
	assert.Equal(t, []byte{
		// ftyp
		0x00, 0x00, 0x00, 0x18, // size
		'f', 't', 'y', 'p', // type
		'a', 'b', 'e', 'm', // major brand
		0x12, 0x34, 0x56, 0x78, // minor version
		'a', 'b', 'c', 'd', // compatible brand
		'e', 'f', 'g', 'h', // compatible brand
		// moov
		0x00, 0x00, 0x00, 0x7b, // size
		'm', 'o', 'o', 'v', // type
		// udta (copy)
		0x00, 0x00, 0x00, 0x0a,
		'u', 'd', 't', 'a',
		0x01, 0x02, 0x03, 0x04,
		0x05, 0x06, 0x07,
		// trak
		0x00, 0x00, 0x00, 0x64, // size
		't', 'r', 'a', 'k', // type
		// tkhd
		0x00, 0x00, 0x00, 0x5c, // size
		't', 'k', 'h', 'd', // type
		0,                // version
		0x00, 0x00, 0x00, // flags
		0x00, 0x00, 0x00, 0x01, // creation time
		0x00, 0x00, 0x00, 0x02, // modification time
		0x00, 0x00, 0x00, 0x03, // track ID
		0x00, 0x00, 0x00, 0x00, // reserved
		0x00, 0x00, 0x00, 0x04, // duration
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // reserved
		0x00, 0x05, // layer
		0x00, 0x06, // alternate group
		0x00, 0x07, // volume
		0x00, 0x00, // reserved
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // matrix
		0x00, 0x00, 0x00, 0x08, // width
		0x00, 0x00, 0x00, 0x09, // height
	}, bin)
}
