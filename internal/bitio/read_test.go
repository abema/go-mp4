package bitio

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	buf := bytes.NewReader([]byte{
		0xb5, 0x63, 0xd5, // 1011,0101,0110,0011,1101,0101
		0xa4, 0x6f,
		0xb4, 0xf1, 0xd7, // 1011,0100,1111,0001,1101,0111
	})
	r := NewReader(buf)
	var data []byte
	var err error
	var bit bool

	// 0101,1010
	//  ^^^ ^^^^
	data, err = r.ReadBits(7)
	require.NoError(t, err)
	require.Equal(t, []byte{0x5a}, data)

	// 0000,0001,0110,0011,1101,0101
	//         ^ ^^^^ ^^^^ ^^^^ ^^^^
	data, err = r.ReadBits(17)
	require.NoError(t, err)
	require.Equal(t, []byte{0x01, 0x63, 0xd5}, data)

	data = make([]byte, 2)
	n, err := r.Read(data)
	require.NoError(t, err)
	require.Equal(t, 2, n)
	assert.Equal(t, []byte{0xa4, 0x6f}, data)

	// 0000,0001,0110,1001,1110,0011
	//         ^ ^^^^ ^^^^ ^^^^ ^^^^
	data, err = r.ReadBits(17)
	require.NoError(t, err)
	require.Equal(t, []byte{0x01, 0x69, 0xe3}, data)

	bit, err = r.ReadBit()
	require.NoError(t, err)
	assert.True(t, bit)
	bit, err = r.ReadBit()
	require.NoError(t, err)
	assert.False(t, bit)

	// 0001,0111
	//    ^ ^^^^
	data, err = r.ReadBits(5)
	require.NoError(t, err)
	require.Equal(t, []byte{0x17}, data)
}

func TestReadBits(t *testing.T) {
	testCases := []struct {
		name         string
		octet        byte
		width        uint
		input        []byte
		size         uint
		err          bool
		expectedData []byte
	}{
		{
			name:         "no width",
			input:        []byte{0x6c, 0xa5},
			size:         10,
			expectedData: []byte{0x01, 0xb2},
		},
		{
			name:         "width 3",
			octet:        0x6c,
			width:        3,
			input:        []byte{0xa5},
			size:         10,
			expectedData: []byte{0x02, 0x52},
		},
		{
			name:         "reach to end of box",
			input:        []byte{0x6c, 0xa5},
			size:         16,
			expectedData: []byte{0x6c, 0xa5},
		},
		{
			name:  "overrun",
			input: []byte{0x6c, 0xa5},
			size:  17,
			err:   true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := NewReader(bytes.NewReader(tc.input))
			require.Zero(t, r.(*reader).octet)
			require.Zero(t, r.(*reader).width)
			r.(*reader).octet = tc.octet
			r.(*reader).width = tc.width
			data, err := r.ReadBits(tc.size)
			if tc.err {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedData, data)
		})
	}
}

func TestReadBit(t *testing.T) {
	r := NewReader(bytes.NewReader([]byte{0x6c, 0xa5})).(*reader)
	outputs := []struct {
		bit   bool
		octet byte
	}{
		{bit: false, octet: 0x6c},
		{bit: true, octet: 0x6c},
		{bit: true, octet: 0x6c},
		{bit: false, octet: 0x6c},
		{bit: true, octet: 0x6c},
		{bit: true, octet: 0x6c},
		{bit: false, octet: 0x6c},
		{bit: false, octet: 0x6c},
		{bit: true, octet: 0xa5},
		{bit: false, octet: 0xa5},
		{bit: true, octet: 0xa5},
		{bit: false, octet: 0xa5},
		{bit: false, octet: 0xa5},
		{bit: true, octet: 0xa5},
		{bit: false, octet: 0xa5},
		{bit: true, octet: 0xa5},
	}
	for _, o := range outputs {
		bit, err := r.ReadBit()
		require.NoError(t, err)
		assert.Equal(t, o.bit, bit)
		assert.Equal(t, o.octet, r.octet)
	}
	_, err := r.ReadBits(1)
	require.Error(t, err)
}

func TestReadInvalidAlignment(t *testing.T) {
	r := NewReader(bytes.NewReader([]byte{0x6c, 0x82, 0x41, 0x35, 0x71, 0xa4, 0xcd, 0x9f}))
	_, err := r.Read(make([]byte, 2))
	require.NoError(t, err)
	_, err = r.ReadBits(3)
	require.NoError(t, err)
	_, err = r.Read(make([]byte, 2))
	assert.Equal(t, ErrInvalidAlignment, err)
}

func TestSeekInvalidAlignment(t *testing.T) {
	r := NewReadSeeker(bytes.NewReader([]byte{0x6c, 0x82, 0x41, 0x35, 0x71, 0xa4, 0xcd, 0x9f}))

	_, err := r.Seek(2, io.SeekCurrent)
	require.NoError(t, err)

	data, err := r.ReadBits(3)
	require.NoError(t, err)
	require.Equal(t, []byte{0x02}, data)

	// When the head is not on 8 bits block border, SeekCurrent fails.
	_, err = r.Seek(2, io.SeekCurrent)
	assert.Equal(t, ErrInvalidAlignment, err)

	// SeekStart always succeeds.
	_, err = r.Seek(0, io.SeekStart)
	assert.NoError(t, err)

	data, err = r.ReadBits(3)
	require.NoError(t, err)
	require.Equal(t, []byte{0x03}, data)
}
