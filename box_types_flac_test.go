package mp4

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoxTypesFLAC(t *testing.T) {
	testCases := []struct {
		name string
		src  IImmutableBox
		dst  IBox
		bin  []byte
		str  string
		ctx  Context
	}{
		{
			name: "fLaC",
			src: &AudioSampleEntry{
				SampleEntry: SampleEntry{
					AnyTypeBox:         AnyTypeBox{Type: BoxTypeFLaC()},
					DataReferenceIndex: 0x1234,
				},
				EntryVersion: 0x0123,
				ChannelCount: 0x2345,
				SampleSize:   0x4567,
				PreDefined:   0x6789,
				SampleRate:   0x01234567,
			},
			dst: &AudioSampleEntry{SampleEntry: SampleEntry{AnyTypeBox: AnyTypeBox{Type: BoxTypeFLaC()}}},
			bin: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // reserved
				0x12, 0x34, // data reference index
				0x01, 0x23, // entry version
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // reserved
				0x23, 0x45, // channel count
				0x45, 0x67, // sample size
				0x67, 0x89, // pre-defined
				0x00, 0x00, // reserved
				0x01, 0x23, 0x45, 0x67, // sample rate
			},
			str: `DataReferenceIndex=4660 ` +
				`EntryVersion=291 ` +
				`ChannelCount=9029 ` +
				`SampleSize=17767 ` +
				`PreDefined=26505 ` +
				`SampleRate=291.27110`,
		},
		{
			name: "dfLa",
			src: &DfLa{
				FullBox: FullBox{
					Version: 0,
					Flags:   [3]byte{0x00, 0x00, 0x00},
				},
				Blocks: []byte{0x80, 0x00, 0x00, 0x03, 0x11, 0x22, 0x33},
			},
			dst: &DfLa{},
			bin: []byte{
				0x00, 0x00, 0x00, 0x00, // version + flags
				0x80, 0x00, 0x00, 0x03, 0x11, 0x22, 0x33,
			},
			str: `Version=0 Flags=0x000000 Blocks=[0x80, 0x0, 0x0, 0x3, 0x11, 0x22, 0x33]`,
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
