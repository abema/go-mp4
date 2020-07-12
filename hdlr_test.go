package mp4

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHdlrUnmarshalHandlerName(t *testing.T) {
	testCases := []struct {
		name  string
		bytes []byte
		want  string
	}{
		{name: "NormalString", bytes: []byte("abema"), want: "abema"},
		{name: "EmptyString", bytes: nil, want: ""},
		{name: "AppleQuickTimePascalString", bytes: []byte{5, 'a', 'b', 'e', 'm', 'a'}, want: "abema"},
		{name: "AppleQuickTimePascalStringWithEmpty", bytes: []byte{0x00, 0x00}, want: ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bin := []byte{
				0,                // version
				0x00, 0x00, 0x00, // flags
				0x00, 0x00, 0x00, 0x00, // predefined
				'v', 'i', 'd', 'e', // handler type
				0x00, 0x00, 0x00, 0x00, // reserved
				0x00, 0x00, 0x00, 0x00, // reserved
				0x00, 0x00, 0x00, 0x00, // reserved
			}
			bin = append(bin, tc.bytes...)

			// unmarshal
			dst := Hdlr{}
			r := bytes.NewReader(bin)
			n, err := Unmarshal(r, uint64(len(bin)), &dst)
			assert.NoError(t, err)
			assert.Equal(t, uint64(len(bin)), n)
			assert.Equal(t, [4]byte{'v', 'i', 'd', 'e'}, dst.HandlerType)
			assert.Equal(t, tc.want, dst.Name)
		})
	}
}
