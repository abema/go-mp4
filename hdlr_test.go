package mp4

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHdlrUnmarshalHandlerName(t *testing.T) {
	testCases := []struct {
		name          string
		componentType []byte
		bytes         []byte
		want          string
	}{
		{
			name:          "NormalString",
			componentType: []byte{0x00, 0x00, 0x00, 0x00},
			bytes:         []byte("abema"),
			want:          "abema",
		},
		{
			name:          "EmptyString",
			componentType: []byte{0x00, 0x00, 0x00, 0x00},
			bytes:         nil,
			want:          "",
		},
		{
			name:          "NormalLongString",
			componentType: []byte{0x00, 0x00, 0x00, 0x00},
			bytes:         []byte(" a 1st byte equals to this length"),
			want:          " a 1st byte equals to this length",
		},
		{
			name:          "AppleQuickTimePascalString",
			componentType: []byte("mhlr"),
			bytes:         []byte{5, 'a', 'b', 'e', 'm', 'a'},
			want:          "abema",
		},
		{
			name:          "AppleQuickTimePascalStringWithEmpty",
			componentType: []byte("mhlr"),
			bytes:         []byte{0x00},
			want:          "",
		},
		{
			name:          "AppleQuickTimePascalStringLong",
			componentType: []byte("mhlr"),
			bytes:         []byte(" a 1st byte equals to this length"),
			want:          "a 1st byte equals to this length",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bin := []byte{
				0,                // version
				0x00, 0x00, 0x00, // flags
			}
			bin = append(bin, tc.componentType...)
			bin = append(bin,
				'v', 'i', 'd', 'e', // handler type
				0x00, 0x00, 0x00, 0x00, // reserved
				0x00, 0x00, 0x00, 0x00, // reserved
				0x00, 0x00, 0x00, 0x00, // reserved
			)
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
