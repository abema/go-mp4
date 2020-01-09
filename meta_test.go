package mp4

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetaMarshal(t *testing.T) {
	src := Meta{
		FullBox: FullBox{
			Version: 0,
			Flags:   [3]byte{0x00, 0x00, 0x00},
		},
	}
	bin := []byte{
		0,                // version
		0x00, 0x00, 0x00, // flags
	}

	// marshal
	buf := bytes.NewBuffer(nil)
	n, err := Marshal(buf, &src)
	require.NoError(t, err)
	assert.Equal(t, uint64(len(bin)), n)
	assert.Equal(t, bin, buf.Bytes())

	// unmarshal
	dst := Meta{}
	n, err = Unmarshal(bytes.NewReader(bin), uint64(len(bin)+8), &dst)
	assert.NoError(t, err)
	assert.Equal(t, uint64(buf.Len()), n)
	assert.Equal(t, src, dst)
}

func TestMetaMarshalAppleQuickTime(t *testing.T) {
	bin := []byte{
		0x00, 0x00, 0x01, 0x00, // size
		'h', 'd', 'l', 'r', // type
		0,                // version
		0x00, 0x00, 0x00, // flags
	}

	// unmarshal
	dst := Meta{}
	r := bytes.NewReader(bin)
	n, err := Unmarshal(r, uint64(len(bin)+8), &dst)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), n)
	s, _ := r.Seek(0, io.SeekCurrent)
	assert.Equal(t, int64(0), s)
}
