package mp4

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockBox struct {
	Type         BoxType
	DynSizeMap   map[string]uint
	DynLenMap    map[string]uint
	DynOptMap    map[string]bool
	IsPStringMap map[string]bool
}

func (m *mockBox) GetType() BoxType {
	return m.Type
}

func (m *mockBox) GetFieldSize(n string, ctx Context) uint {
	if s, ok := m.DynSizeMap[n]; !ok {
		panic(fmt.Errorf("invalid name of dynamic-size field: %s", n))
	} else {
		return s
	}
}

func (m *mockBox) GetFieldLength(n string, ctx Context) uint {
	if l, ok := m.DynLenMap[n]; !ok {
		panic(fmt.Errorf("invalid name of dynamic-length field: %s", n))
	} else {
		return l
	}
}

func (m *mockBox) IsOptFieldEnabled(n string, ctx Context) bool {
	if enabled, ok := m.DynOptMap[n]; !ok {
		panic(fmt.Errorf("invalid name of dynamic-opt field: %s", n))
	} else {
		return enabled
	}
}

func (m *mockBox) IsPString(name string, bytes []byte, remainingSize uint64, ctx Context) bool {
	if b, ok := m.IsPStringMap[name]; ok {
		return b
	}
	return true
}

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
