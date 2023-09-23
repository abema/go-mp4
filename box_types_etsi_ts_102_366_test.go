package mp4

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoxTypesETSI_TS_102_366(t *testing.T) {
	testCases := []struct {
		name string
		src  IImmutableBox
		dst  IBox
		bin  []byte
		str  string
		ctx  Context
	}{
		{
			name: "dac3",
			src: &Dac3{
				Fscod:       0,
				Bsid:        8,
				Acmod:       7,
				LfeOn:       1,
				BitRateCode: 0x7,
			},
			dst: &Dac3{},
			bin: []byte{0x10, 0x3c, 0xe0},
			str: `Fscod=0x0 Bsid=0x8 Bsmod=0x0 Acmod=0x7 LfeOn=0x1 BitRateCode=0x7`,
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
