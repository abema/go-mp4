package mp4

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteBoxInfo(t *testing.T) {
	type pattern struct {
		name       string
		input      io.Writer
		bi         *BoxInfo
		hasError   bool
		expectedBI *BoxInfo
		assert     func(*pattern)
	}

	patterns := []pattern{
		{
			name:  "small-size",
			input: &bytes.Buffer{},
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
			assert: func(p *pattern) {
				assert.Equal(t, []byte{
					0x00, 0x01, 0x23, 0x45,
					't', 'e', 's', 't',
				}, p.input.(*bytes.Buffer).Bytes(), "%s", p.name)
			},
		},
		{
			name:  "large-size",
			input: &bytes.Buffer{},
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
			assert: func(p *pattern) {
				assert.Equal(t, []byte{
					0x00, 0x00, 0x00, 0x01,
					't', 'e', 's', 't',
					0x00, 0x00, 0x12, 0x34,
					0x56, 0x78, 0x9a, 0xbc,
				}, p.input.(*bytes.Buffer).Bytes(), "%s", p.name)
			},
		},
		{
			name:  "extend to eof",
			input: &bytes.Buffer{},
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
			assert: func(p *pattern) {
				assert.Equal(t, []byte{
					0x00, 0x00, 0x00, 0x00,
					't', 'e', 's', 't',
				}, p.input.(*bytes.Buffer).Bytes(), "%s", p.name)
			},
		},
	}

	for _, p := range patterns {
		bi, err := WriteBoxInfo(p.input, p.bi)
		if !p.hasError {
			require.NoError(t, err, "%s", p.name)
			assert.Equal(t, p.expectedBI, bi, "%s", p.name)
			p.assert(&p)
		} else {
			assert.Error(t, err, "%s", p.name)
		}
	}
}

func TestReadBoxInfo(t *testing.T) {
	patterns := []struct {
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

	for _, p := range patterns {
		buf := bytes.NewReader(p.buf)
		buf.Seek(p.seek, io.SeekStart)
		bi, err := ReadBoxInfo(buf)
		if !p.hasError {
			require.NoError(t, err)
			assert.Equal(t, p.expected, bi)
		} else {
			assert.Error(t, err)
		}
	}
}
