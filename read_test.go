package mp4

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadBoxStructure(t *testing.T) {
	f, err := os.Open("./_examples/sample.mp4")
	require.NoError(t, err)
	defer f.Close()

	var n int
	_, err = ReadBoxStructure(f, func(h *ReadHandle) (interface{}, error) {
		n++
		switch n {
		case 2, 37, 49, 50, 54, 55: // free, smhd, sgpd, sbgp, ilst, loci
			require.False(t, h.BoxInfo.Type.IsSupported())
			buf := bytes.NewBuffer(nil)
			n, err := h.ReadData(buf)
			require.NoError(t, err)
			require.Equal(t, h.BoxInfo.Size-h.BoxInfo.HeaderSize, n)
			assert.Len(t, buf.Bytes(), int(n))
		case 41: // stbl
			require.True(t, h.BoxInfo.Type.IsSupported())
			require.Equal(t, BoxTypeStbl(), h.BoxInfo.Type)
			infos, err := h.Expand()
			require.NoError(t, err)
			assert.Equal(t, []interface{}{"stsd", "stts", nil, nil, "stco", nil, nil}, infos)
		case 42: // stsd
			require.True(t, h.BoxInfo.Type.IsSupported())
			require.Equal(t, BoxTypeStsd(), h.BoxInfo.Type)
			box, n, err := h.ReadPayload()
			require.NoError(t, err)
			require.Equal(t, uint64(8), n)
			assert.Equal(t, &Stsd{EntryCount: 1}, box)
			_, err = h.Expand()
			require.NoError(t, err)
			return "stsd", nil
		case 45: // stts
			require.True(t, h.BoxInfo.Type.IsSupported())
			require.Equal(t, BoxTypeStts(), h.BoxInfo.Type)
			_, err = h.Expand()
			require.NoError(t, err)
			return "stts", nil
		case 48: // stco
			require.True(t, h.BoxInfo.Type.IsSupported())
			require.Equal(t, BoxTypeStco(), h.BoxInfo.Type)
			_, err = h.Expand()
			require.NoError(t, err)
			return "stco", nil
		default: // otherwise
			require.True(t, h.BoxInfo.Type.IsSupported())
			_, err = h.Expand()
			require.NoError(t, err)
		}
		return nil, nil
	})
	require.NoError(t, err)
	assert.Equal(t, 55, n)
}
