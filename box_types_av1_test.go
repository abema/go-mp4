package mp4

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoxTypesAV1(t *testing.T) {
	testCases := []struct {
		name string
		src  IImmutableBox
		dst  IBox
		bin  []byte
		str  string
		ctx  Context
	}{
		{
			name: "Av1C",
			src: &Av1C{
				Marker:               1,
				Version:              1,
				SeqProfile:           2,
				SeqLevelIdx0:         1,
				SeqTier0:             1,
				HighBitdepth:         1,
				TwelveBit:            0,
				Monochrome:           0,
				ChromaSubsamplingX:   1,
				ChromaSubsamplingY:   1,
				ChromaSamplePosition: 0,
				ConfigOBUs: []byte{
					0x08, 0x00, 0x00, 0x00, 0x42, 0xa7, 0xbf, 0xe4,
					0x60, 0x0d, 0x00, 0x40,
				},
			},
			dst: &Av1C{},
			bin: []byte{
				0x81, 0x41, 0xcc, 0x00, 0x08, 0x00, 0x00, 0x00,
				0x42, 0xa7, 0xbf, 0xe4, 0x60, 0x0d, 0x00, 0x40,
			},
			str: `SeqProfile=0x2 SeqLevelIdx0=0x1 SeqTier0=0x1 HighBitdepth=0x1 TwelveBit=0x0 Monochrome=0x0 ChromaSubsamplingX=0x1 ChromaSubsamplingY=0x1 ChromaSamplePosition=0x0 InitialPresentationDelayPresent=0x0 InitialPresentationDelayMinusOne=0x0 ConfigOBUs=[0x8, 0x0, 0x0, 0x0, 0x42, 0xa7, 0xbf, 0xe4, 0x60, 0xd, 0x0, 0x40]`,
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
