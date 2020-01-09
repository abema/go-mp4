package mp4

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMdhdMarshal(t *testing.T) {
	// Version 0
	mdhd := Mdhd{
		FullBox: FullBox{
			Version: 0,
			Flags:   [3]byte{0x00, 0x00, 0x00},
		},
		CreationTimeV0:     0x12345678,
		ModificationTimeV0: 0x23456789,
		TimescaleV0:        0x01020304,
		DurationV0:         0x02030405,
		Pad:                true,
		Language:           [3]byte{'j' - 0x60, 'p' - 0x60, 'n' - 0x60}, // 0x0a, 0x10, 0x0e
		PreDefined:         0,
	}
	expect := []byte{
		0,                // version
		0x00, 0x00, 0x00, // flags
		0x12, 0x34, 0x56, 0x78, // creation time
		0x23, 0x45, 0x67, 0x89, // modification time
		0x01, 0x02, 0x03, 0x04, // timescale
		0x02, 0x03, 0x04, 0x05, // duration
		0xaa, 0x0e, // pad, language (1 01010 10000 01110)
		0x00, 0x00, // pre defined
	}

	buf := bytes.NewBuffer(nil)
	n, err := Marshal(buf, &mdhd)
	require.NoError(t, err)
	assert.Equal(t, uint64(len(expect)), n)
	assert.Equal(t, expect, buf.Bytes())

	// Version 1
	mdhd = Mdhd{
		FullBox: FullBox{
			Version: 1,
			Flags:   [3]byte{0x00, 0x00, 0x00},
		},
		CreationTimeV1:     0x123456789abcdef0,
		ModificationTimeV1: 0x23456789abcdef01,
		TimescaleV1:        0x01020304,
		DurationV1:         0x0203040506070809,
		Pad:                true,
		Language:           [3]byte{'j' - 0x60, 'p' - 0x60, 'n' - 0x60}, // 0x0a, 0x10, 0x0e
		PreDefined:         0,
	}
	expect = []byte{
		1,                // version
		0x00, 0x00, 0x00, // flags
		0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0, // creation time
		0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x01, // modification time
		0x01, 0x02, 0x03, 0x04, // timescale
		0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, // duration
		0xaa, 0x0e, // pad, language (1 01010 10000 01110)
		0x00, 0x00, // pre defined
	}

	buf = bytes.NewBuffer(nil)
	n, err = Marshal(buf, &mdhd)
	require.NoError(t, err)
	assert.Equal(t, uint64(len(expect)), n)
	assert.Equal(t, expect, buf.Bytes())
}

func TestMdhdUnmarshal(t *testing.T) {
	// Version 0
	data := []byte{
		0,                // version
		0x00, 0x00, 0x00, // flags
		0x12, 0x34, 0x56, 0x78, // creation time
		0x23, 0x45, 0x67, 0x89, // modification time
		0x34, 0x56, 0x78, 0x9a, // timescale
		0x45, 0x67, 0x89, 0xab, // duration
		0xaa, 0x67, // pad, language (1 01010 10011 00111)
		0x00, 0x00, // pre defined
	}

	buf := bytes.NewReader(data)
	mdhd := Mdhd{}
	n, err := Unmarshal(buf, uint64(buf.Len()), &mdhd)
	require.NoError(t, err)
	assert.Equal(t, uint64(len(data)), n)
	assert.Equal(t, uint8(0), mdhd.Version)
	assert.Equal(t, uint32(0x0), mdhd.GetFlags())
	assert.Equal(t, uint32(0x12345678), mdhd.CreationTimeV0)
	assert.Equal(t, uint32(0x23456789), mdhd.ModificationTimeV0)
	assert.Equal(t, uint32(0x3456789a), mdhd.TimescaleV0)
	assert.Equal(t, uint32(0x456789ab), mdhd.DurationV0)
	assert.Equal(t, uint64(0x0), mdhd.CreationTimeV1)
	assert.Equal(t, uint64(0x0), mdhd.ModificationTimeV1)
	assert.Equal(t, uint32(0x0), mdhd.TimescaleV1)
	assert.Equal(t, uint64(0x0), mdhd.DurationV1)
	assert.Equal(t, true, mdhd.Pad)
	assert.Equal(t, byte(0x0a), mdhd.Language[0])
	assert.Equal(t, byte(0x13), mdhd.Language[1])
	assert.Equal(t, byte(0x07), mdhd.Language[2])
	assert.Equal(t, uint16(0), mdhd.PreDefined)

	// Version 1
	data = []byte{
		1,                // version
		0x00, 0x00, 0x00, // flags
		0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0, // creation time
		0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x01, // modification time
		0x34, 0x56, 0x78, 0x9a, // timescale
		0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x01, 0x23, // duration
		0xaa, 0x67, // pad, language (1 01010 10011 00111)
		0x00, 0x00, // pre defined
	}

	buf = bytes.NewReader(data)
	mdhd = Mdhd{}
	n, err = Unmarshal(buf, uint64(buf.Len()), &mdhd)
	require.NoError(t, err)
	assert.Equal(t, uint64(len(data)), n)
	assert.Equal(t, uint8(1), mdhd.Version)
	assert.Equal(t, uint32(0x0), mdhd.GetFlags())
	assert.Equal(t, uint32(0x0), mdhd.CreationTimeV0)
	assert.Equal(t, uint32(0x0), mdhd.ModificationTimeV0)
	assert.Equal(t, uint32(0x0), mdhd.TimescaleV0)
	assert.Equal(t, uint32(0x0), mdhd.DurationV0)
	assert.Equal(t, uint64(0x123456789abcdef0), mdhd.CreationTimeV1)
	assert.Equal(t, uint64(0x23456789abcdef01), mdhd.ModificationTimeV1)
	assert.Equal(t, uint32(0x3456789a), mdhd.TimescaleV1)
	assert.Equal(t, uint64(0x456789abcdef0123), mdhd.DurationV1)
	assert.Equal(t, true, mdhd.Pad)
	assert.Equal(t, byte(0x0a), mdhd.Language[0])
	assert.Equal(t, byte(0x13), mdhd.Language[1])
	assert.Equal(t, byte(0x07), mdhd.Language[2])
	assert.Equal(t, uint16(0), mdhd.PreDefined)
}
