package mp4

import (
	"bytes"
	"io"
	"testing"

	"github.com/orcaman/writerseeker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteBoxInfo(t *testing.T) {
	type testCase struct {
		name          string
		pre           []byte
		bi            *BoxInfo
		hasError      bool
		expectedBI    *BoxInfo
		expectedBytes []byte
	}

	testCases := []testCase{
		{
			name: "small-size",
			bi: &BoxInfo{
				Size:       0x12345,
				HeaderSize: 8,
				Type:       StrToBoxType("test"),
			},
			expectedBI: &BoxInfo{
				Size:       0x12345,
				HeaderSize: 8,
				Type:       StrToBoxType("test"),
			},
			expectedBytes: []byte{
				0x00, 0x01, 0x23, 0x45,
				't', 'e', 's', 't',
			},
		},
		{
			name: "large-size",
			bi: &BoxInfo{
				Size:       0x123456789abc,
				HeaderSize: 8,
				Type:       StrToBoxType("test"),
			},
			expectedBI: &BoxInfo{
				Size:       0x123456789abc + 8,
				HeaderSize: 16,
				Type:       StrToBoxType("test"),
			},
			expectedBytes: []byte{
				0x00, 0x00, 0x00, 0x01,
				't', 'e', 's', 't',
				0x00, 0x00, 0x12, 0x34,
				0x56, 0x78, 0x9a, 0xbc,
			},
		},
		{
			name: "extend to eof",
			bi: &BoxInfo{
				Size:        0x123,
				HeaderSize:  8,
				Type:        StrToBoxType("test"),
				ExtendToEOF: true,
			},
			expectedBI: &BoxInfo{
				Size:        0x123,
				HeaderSize:  8,
				Type:        StrToBoxType("test"),
				ExtendToEOF: true,
			},
			expectedBytes: []byte{
				0x00, 0x00, 0x00, 0x00,
				't', 'e', 's', 't',
			},
		},
		{
			name: "with offset",
			pre:  []byte{0x00, 0x00, 0x00},
			bi: &BoxInfo{
				Size:       0x12345,
				HeaderSize: 8,
				Type:       StrToBoxType("test"),
			},
			expectedBI: &BoxInfo{
				Offset:     3,
				Size:       0x12345,
				HeaderSize: 8,
				Type:       StrToBoxType("test"),
			},
			expectedBytes: []byte{
				0x00, 0x00, 0x00, // pre inserted
				0x00, 0x01, 0x23, 0x45,
				't', 'e', 's', 't',
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			w := &writerseeker.WriterSeeker{}
			_, err := w.Write(c.pre)
			require.NoError(t, err)
			bi, err := WriteBoxInfo(w, c.bi)
			if !c.hasError {
				require.NoError(t, err)
				assert.Equal(t, c.expectedBI, bi)
				b, err := io.ReadAll(w.Reader())
				require.NoError(t, err)
				assert.Equal(t, c.expectedBytes, b)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestReadBoxInfo(t *testing.T) {
	testCases := []struct {
		name     string
		buf      []byte
		seek     int64
		hasError bool
		expected *BoxInfo
	}{
		{
			name: "no offset",
			buf: []byte{
				0x00, 0x01, 0x23, 0x45,
				't', 'e', 's', 't',
			},
			expected: &BoxInfo{
				Size:       0x12345,
				HeaderSize: 8,
				Type:       StrToBoxType("test"),
			},
		},
		{
			name: "has offset",
			buf: []byte{
				0x00, 0x00,
				0x00, 0x01, 0x23, 0x45,
				't', 'e', 's', 't',
			},
			seek: 2,
			expected: &BoxInfo{
				Offset:     2,
				Size:       0x12345,
				HeaderSize: 8,
				Type:       StrToBoxType("test"),
			},
		},
		{
			name: "large-size",
			buf: []byte{
				0x00, 0x00, 0x00, 0x01,
				't', 'e', 's', 't',
				0x01, 0x23, 0x45, 0x67,
				0x89, 0xab, 0xcd, 0xef,
			},
			expected: &BoxInfo{
				Size:       0x123456789abcdef,
				HeaderSize: 16,
				Type:       StrToBoxType("test"),
			},
		},
		{
			name: "extend to eof",
			buf: []byte{
				0x00, 0x00, 0x00, 0x00,
				't', 'e', 's', 't',
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
			},
			expected: &BoxInfo{
				Size:        20,
				HeaderSize:  8,
				Type:        StrToBoxType("test"),
				ExtendToEOF: true,
			},
		},
		{
			name: "end-of-file",
			buf: []byte{
				0x00, 0x01, 0x23, 0x45,
				't', 'e', 's',
			},
			hasError: true,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			buf := bytes.NewReader(c.buf)
			buf.Seek(c.seek, io.SeekStart)
			bi, err := ReadBoxInfo(buf)
			if !c.hasError {
				require.NoError(t, err)
				assert.Equal(t, c.expected, bi)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
