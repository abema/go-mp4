package mp4

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmsgMarshal(t *testing.T) {
	emsg := Emsg{
		FullBox: FullBox{
			Version: 0,
			Flags:   [3]byte{0x00, 0x00, 0x00},
		},
		SchemeIdUri:           "urn:test",
		Value:                 "foo",
		Timescale:             1000,
		PresentationTimeDelta: 123,
		EventDuration:         3000,
		Id:                    0xabcd,
		MessageData:           []byte("abema"),
	}
	expect := []byte{
		0,                // version
		0x00, 0x00, 0x00, // flags
		0x75, 0x72, 0x6e, 0x3a, 0x74, 0x65, 0x73, 0x74, 0x00, // scheme id uri
		0x66, 0x6f, 0x6f, 0x00, // value
		0x00, 0x00, 0x03, 0xe8, // timescale
		0x00, 0x00, 0x00, 0x7b, // presentation time delta
		0x00, 0x00, 0x0b, 0xb8, // event duration
		0x00, 0x00, 0xab, 0xcd, // id
		0x61, 0x62, 0x65, 0x6d, 0x61, // message data
	}

	buf := bytes.NewBuffer(nil)
	n, err := Marshal(buf, &emsg)
	require.NoError(t, err)
	assert.Equal(t, uint64(len(expect)), n)
	assert.Equal(t, expect, buf.Bytes())
}

func TestEmsgUnmarshal(t *testing.T) {
	data := []byte{
		0,                // version
		0x00, 0x00, 0x00, // flags
		0x75, 0x72, 0x6e, 0x3a, 0x74, 0x65, 0x73, 0x74, 0x00, // scheme id uri
		0x66, 0x6f, 0x6f, 0x00, // value
		0x00, 0x00, 0x03, 0xe8, // timescale
		0x00, 0x00, 0x00, 0x7b, // presentation time delta
		0x00, 0x00, 0x0b, 0xb8, // event duration
		0x00, 0x00, 0xab, 0xcd, // id
		0x61, 0x62, 0x65, 0x6d, 0x61, // message data
	}

	buf := bytes.NewReader(data)
	emsg := Emsg{}
	n, err := Unmarshal(buf, uint64(buf.Len()), &emsg)
	require.NoError(t, err)
	assert.Equal(t, uint64(len(data)), n)
	assert.Equal(t, uint8(0), emsg.Version)
	assert.Equal(t, uint32(0x000000), emsg.GetFlags())
	assert.Equal(t, "urn:test", emsg.SchemeIdUri)
	assert.Equal(t, "foo", emsg.Value)
	assert.Equal(t, uint32(1000), emsg.Timescale)
	assert.Equal(t, uint32(123), emsg.PresentationTimeDelta)
	assert.Equal(t, uint32(3000), emsg.EventDuration)
	assert.Equal(t, uint32(0xabcd), emsg.Id)
	assert.Equal(t, []byte("abema"), emsg.MessageData)

	buf = bytes.NewReader(data)
	result, n, err := UnmarshalAny(buf, BoxTypeEmsg(), uint64(len(data)))
	require.NoError(t, err)
	assert.Equal(t, uint64(len(data)), n)
	pemsg, ok := result.(*Emsg)
	require.True(t, ok)
	assert.Equal(t, uint8(0), pemsg.Version)
	assert.Equal(t, uint32(0x000000), pemsg.GetFlags())
	assert.Equal(t, "urn:test", pemsg.SchemeIdUri)
	assert.Equal(t, "foo", pemsg.Value)
	assert.Equal(t, uint32(1000), pemsg.Timescale)
	assert.Equal(t, uint32(123), pemsg.PresentationTimeDelta)
	assert.Equal(t, uint32(3000), pemsg.EventDuration)
	assert.Equal(t, uint32(0xabcd), pemsg.Id)
	assert.Equal(t, []byte("abema"), pemsg.MessageData)
}
