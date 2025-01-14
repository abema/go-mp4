package mp4

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoxTypesISO14496_30(t *testing.T) {
	testCases := []struct {
		name string
		src  IImmutableBox
		dst  IBox
		bin  []byte
		str  string
		ctx  Context
	}{
		{
			name: "vttC",
			src: &WebVTTConfigurationBox{
				Config: "WEBVTT\n",
			},
			dst: &WebVTTConfigurationBox{},
			bin: []byte{'W', 'E', 'B', 'V', 'T', 'T', '\n'},
			str: `Config="WEBVTT."`,
		},
		{
			name: "vlab",
			src: &WebVTTSourceLabelBox{
				SourceLabel: "Source",
			},
			dst: &WebVTTSourceLabelBox{},
			bin: []byte{'S', 'o', 'u', 'r', 'c', 'e'},
			str: `SourceLabel="Source"`,
		},
		{
			name: "wvtt",
			src: &WVTTSampleEntry{
				SampleEntry: SampleEntry{
					AnyTypeBox:         AnyTypeBox{Type: StrToBoxType("wvtt")},
					DataReferenceIndex: 0x1234,
				},
			},
			dst: &WVTTSampleEntry{SampleEntry: SampleEntry{AnyTypeBox: AnyTypeBox{Type: StrToBoxType("wvtt")}}},
			bin: []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x12, 0x34},
			str: `DataReferenceIndex=4660`,
		},
		{
			name: "vttc",
			src:  &VTTCueBox{},
			dst:  &VTTCueBox{},
			bin:  []byte(nil),
			str:  ``,
		},
		{
			name: "vsid",
			src: &CueSourceIDBox{
				SourceId: 0,
			},
			dst: &CueSourceIDBox{},
			bin: []byte{0, 0, 0, 0},
			str: `SourceId=0`,
		},
		{
			name: "ctim",
			src: &CueTimeBox{
				CueCurrentTime: "00:00:00.000",
			},
			dst: &CueTimeBox{},
			bin: []byte{'0', '0', ':', '0', '0', ':', '0', '0', '.', '0', '0', '0'},
			str: `CueCurrentTime="00:00:00.000"`,
		},
		{
			name: "iden",
			src: &CueIDBox{
				CueId: "example_id",
			},
			dst: &CueIDBox{},
			bin: []byte{'e', 'x', 'a', 'm', 'p', 'l', 'e', '_', 'i', 'd'},
			str: `CueId="example_id"`,
		},
		{
			name: "sttg",
			src: &CueSettingsBox{
				Settings: "line=0",
			},
			dst: &CueSettingsBox{},
			bin: []byte{'l', 'i', 'n', 'e', '=', '0'},
			str: `Settings="line=0"`,
		},
		{
			name: "payl",
			src: &CuePayloadBox{
				CueText: "sample",
			},
			dst: &CuePayloadBox{},
			bin: []byte{'s', 'a', 'm', 'p', 'l', 'e'},
			str: `CueText="sample"`,
		},
		{
			name: "vtte",
			src:  &VTTEmptyCueBox{},
			dst:  &VTTEmptyCueBox{},
			bin:  []byte(nil),
			str:  ``,
		},
		{
			name: "vtta",
			src: &VTTAdditionalTextBox{
				CueAdditionalText: "test",
			},
			dst: &VTTAdditionalTextBox{},
			bin: []byte{'t', 'e', 's', 't'},
			str: `CueAdditionalText="test"`,
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
