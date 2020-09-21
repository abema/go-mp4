package bitio

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrite(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	w := NewWriter(buf)

	// 1101,1010
	//  ^^^ ^^^^
	require.NoError(t, w.WriteBits([]byte{0xda}, 7))

	// 0000,0111,0110,0011,1101,0101
	//         ^ ^^^^ ^^^^ ^^^^ ^^^^
	require.NoError(t, w.WriteBits([]byte{0x07, 0x63, 0xd5}, 17))

	_, err := w.Write([]byte{0xa4, 0x6f})
	require.NoError(t, err)

	// 0000,0111,0110,1001,1110,0011
	//         ^ ^^^^ ^^^^ ^^^^ ^^^^
	require.NoError(t, w.WriteBits([]byte{0x07, 0x69, 0xe3}, 17))

	require.NoError(t, w.WriteBit(true))
	require.NoError(t, w.WriteBit(false))

	// 1111,0111
	//    ^ ^^^^
	require.NoError(t, w.WriteBits([]byte{0xf7}, 5))

	assert.Equal(t, []byte{
		0xb5, 0x63, 0xd5, // 1011,0101,0110,0011,1101,0101
		0xa4, 0x6f,
		0xb4, 0xf1, 0xd7, // 1011,0100,1111,0001,1101,0111
	}, buf.Bytes())
}

func TestWriteInvalidAlignment(t *testing.T) {
	w := NewWriter(bytes.NewBuffer(nil))
	_, err := w.Write([]byte{0xa4, 0x6f})
	require.NoError(t, err)
	require.NoError(t, w.WriteBits([]byte{0xda}, 7))
	_, err = w.Write([]byte{0xa4, 0x6f})
	require.Equal(t, ErrInvalidAlignment, err)
}
