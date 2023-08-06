package mp4

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoxTypesVp(t *testing.T) {
	testCases := []struct {
		name string
		src  IImmutableBox
		dst  IBox
		bin  []byte
		str  string
		ctx  Context
	}{
		{
			name: "vpcC",
			src: &VpcC{
				FullBox: FullBox{
					Version: 1,
				},
				Profile:                     1,
				Level:                       50,
				BitDepth:                    10,
				ChromaSubsampling:           3,
				VideoFullRangeFlag:          1,
				ColourPrimaries:             0,
				TransferCharacteristics:     1,
				MatrixCoefficients:          10,
				CodecInitializationDataSize: 3,
				CodecInitializationData:     []byte{5, 4, 3},
			},
			dst: &VpcC{},
			bin: []byte{
				0x01, 0x00, 0x00, 0x00, 0x01, 0x32, 0xa7, 0x00,
				0x01, 0x0a, 0x00, 0x03, 0x05, 0x04, 0x03,
			},
			str: `Version=1 Flags=0x000000 Profile=0x1 Level=0x32 BitDepth=0xa ChromaSubsampling=0x3 VideoFullRangeFlag=0x1 ColourPrimaries=0x0 TransferCharacteristics=0x1 MatrixCoefficients=0xa CodecInitializationDataSize=3 CodecInitializationData=[0x5, 0x4, 0x3]`,
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
