package util

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadString(t *testing.T) {
	r := bytes.NewReader([]byte{
		'f', 'i', 'r', 's', 't', 0,
		's', 'e', 'c', 'o', 'n', 'd', 0,
		't', 'h', 'i', 'r', 'd', 0,
	})
	s, err := ReadString(r)
	require.NoError(t, err)
	assert.Equal(t, "first", s)
	s, err = ReadString(r)
	require.NoError(t, err)
	assert.Equal(t, "second", s)
	s, err = ReadString(r)
	require.NoError(t, err)
	assert.Equal(t, "third", s)
	_, err = ReadString(r)
	assert.Equal(t, io.EOF, err)
}

func TestWriteString(t *testing.T) {
	w := bytes.NewBuffer(nil)
	require.NoError(t, WriteString(w, "first"))
	require.NoError(t, WriteString(w, "second"))
	require.NoError(t, WriteString(w, "third"))
	assert.Equal(t, []byte{
		'f', 'i', 'r', 's', 't', 0,
		's', 'e', 'c', 'o', 'n', 'd', 0,
		't', 'h', 'i', 'r', 'd', 0,
	}, w.Bytes())
}
