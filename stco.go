package mp4

import "fmt"

func BoxTypeStco() BoxType { return StrToBoxType("stco") }

func init() {
	AddBoxDef(&Stco{}, 0)
}

// Stco is ISOBMFF stco box type
type Stco struct {
	FullBox     `mp4:"extend"`
	EntryCount  uint32   `mp4:"size=32"`
	ChunkOffset []uint32 `mp4:"size=32,len=dynamic"`
}

// GetType returns the BoxType
func (*Stco) GetType() BoxType {
	return BoxTypeStco()
}

// GetFieldLength returns length of dynamic field
func (stco *Stco) GetFieldLength(name string) uint {
	switch name {
	case "ChunkOffset":
		return uint(stco.EntryCount)
	}
	panic(fmt.Errorf("invalid name of dynamic-length field: boxType=stco fieldName=%s", name))
}
