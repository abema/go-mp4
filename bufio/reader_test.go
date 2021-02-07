package bufio

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadSeeker(t *testing.T) {
	r := NewReadSeeker(bytes.NewReader([]byte(""+
		"ABCDEFGH"+
		"IJKLMNOP"+
		"QRSTUVWX"+
		"YZabcdef"+
		"ghijklmn"+
		"opqrstuv"+
		"wxyz",
	)), 8, 4)

	testRead(t, r, 5, "ABCDE", nil)
	testRead(t, r, 5, "FGHIJ", nil)
	testSeekRead(t, r, -2, io.SeekCurrent, 8, 3, "IJK")
	testSeekRead(t, r, 17, io.SeekStart, 17, 3, "RST")
	testSeekRead(t, r, -8, io.SeekEnd, 44, 5, "stuvw")
	testRead(t, r, 5, "xyz", io.EOF)
	testRead(t, r, 5, "", io.EOF)
}

func testRead(t *testing.T, r io.Reader, n int64, expected string, expectedErr error) {
	w := bytes.NewBuffer(nil)
	read, err := io.CopyN(w, r, n)
	require.Equal(t, expectedErr, err)
	require.Equal(t, int64(len(expected)), read)
	require.Equal(t, expected, w.String())
}

func testSeekRead(t *testing.T, r io.ReadSeeker, offset int64, whence int, newOffset int64, n int, expected string) {
	o, err := r.Seek(offset, whence)
	require.NoError(t, err)
	require.Equal(t, newOffset, o)
	buf := make([]byte, n)
	read, err := io.ReadFull(r, buf)
	require.NoError(t, err)
	require.Equal(t, n, read)
	require.Equal(t, expected, string(buf))
}
