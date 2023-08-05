package mp4

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoxTypesOpus(t *testing.T) {
	testCases := []struct {
		name string
		src  IImmutableBox
		dst  IBox
		bin  []byte
		str  string
		ctx  Context
	}{
		{
			name: "dOps",
			src: &DOps{
				OutputChannelCount:   2,
				PreSkip:              312,
				InputSampleRate:      48000,
				OutputGain:           0,
				ChannelMappingFamily: 2,
				StreamCount:          1,
				CoupledCount:         1,
				ChannelMapping:       []uint8{1, 2},
			},
			dst: &DOps{},
			bin: []byte{
				0x00, 0x02, 0x01, 0x38, 0x00, 0x00, 0xbb, 0x80,
				0x00, 0x00, 0x02, 0x01, 0x01, 0x01, 0x02,
			},
			str: `Version=0 OutputChannelCount=0x2 PreSkip=312 InputSampleRate=48000 OutputGain=0 ChannelMappingFamily=0x2 StreamCount=0x1 CoupledCount=0x1 ChannelMapping=[0x1, 0x2]`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Marshal
			buf := bytes.NewBuffer(nil)
			n, err := Marshal(buf, tc.src, tc.ctx)
			require.NoError(t, err)
			assert.Equal(t, uint64(len(tc.bin)), n)
			assert.Equal(t, tc.bin, buf.Bytes())

			// Unmarshal
			r := bytes.NewReader(tc.bin)
			n, err = Unmarshal(r, uint64(len(tc.bin)), tc.dst, tc.ctx)
			require.NoError(t, err)
			assert.Equal(t, uint64(buf.Len()), n)
			assert.Equal(t, tc.src, tc.dst)
			s, err := r.Seek(0, io.SeekCurrent)
			require.NoError(t, err)
			assert.Equal(t, int64(buf.Len()), s)

			// UnmarshalAny
			dst, n, err := UnmarshalAny(bytes.NewReader(tc.bin), tc.src.GetType(), uint64(len(tc.bin)), tc.ctx)
			require.NoError(t, err)
			assert.Equal(t, uint64(buf.Len()), n)
			assert.Equal(t, tc.src, dst)
			s, err = r.Seek(0, io.SeekCurrent)
			require.NoError(t, err)
			assert.Equal(t, int64(buf.Len()), s)

			// Stringify
			str, err := Stringify(tc.src, tc.ctx)
			require.NoError(t, err)
			assert.Equal(t, tc.str, str)
		})
	}
}
