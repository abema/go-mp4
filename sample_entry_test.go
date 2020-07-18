package mp4

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVisualSampleEntryStringify(t *testing.T) {
	vse := VisualSampleEntry{
		SampleEntry:     SampleEntry{},
		PreDefined:      0x0101,
		PreDefined2:     [3]uint32{0x01000001, 0x01000002, 0x01000003},
		Width:           0x0102,
		Height:          0x0103,
		Horizresolution: 0x01000004,
		Vertresolution:  0x01000005,
		Reserved2:       0x01000006,
		FrameCount:      0x0104,
		Compressorname:  [32]byte{'a', 'b', 'e', 'm', 'a'},
		Depth:           0x0105,
		PreDefined3:     1001,
	}
	str, err := Stringify(&vse)
	require.NoError(t, err)
	assert.Equal(t, "DataReferenceIndex=0"+
		" PreDefined=257"+
		" PreDefined2=[16777217,"+
		" 16777218,"+
		" 16777219]"+
		" Width=258"+
		" Height=259"+
		" Horizresolution=16777220"+
		" Vertresolution=16777221"+
		" FrameCount=260"+
		" Compressorname=\"abema\""+
		" Depth=261"+
		" PreDefined3=1001", str)
}
