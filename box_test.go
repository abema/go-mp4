package mp4

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFullBoxFlags(t *testing.T) {
	box := FullBox{}
	box.SetFlags(0x35ac68)
	assert.Equal(t, byte(0x35), box.Flags[0])
	assert.Equal(t, byte(0xac), box.Flags[1])
	assert.Equal(t, byte(0x68), box.Flags[2])
	assert.Equal(t, uint32(0x35ac68), box.GetFlags())

	box.AddFlag(0x030000)
	assert.Equal(t, uint32(0x37ac68), box.GetFlags())

	box.RemoveFlag(0x000900)
	assert.Equal(t, uint32(0x37a468), box.GetFlags())
}
