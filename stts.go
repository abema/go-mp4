package mp4

import "fmt"

func BoxTypeStts() BoxType { return StrToBoxType("stts") }

func init() {
	AddBoxDef(&Stts{}, 0)
}

// Stts is ISOBMFF stts box type
type Stts struct {
	FullBox    `mp4:"extend"`
	EntryCount uint32      `mp4:"size=32"`
	Entries    []SttsEntry `mp4:"len=dynamic,size=64"`
}

type SttsEntry struct {
	SampleCount uint32 `mp4:"size=32"`
	SampleDelta uint32 `mp4:"size=32"`
}

// GetType returns the BoxType
func (*Stts) GetType() BoxType {
	return BoxTypeStts()
}

// GetFieldLength returns length of dynamic field
func (stts *Stts) GetFieldLength(name string) uint {
	switch name {
	case "Entries":
		return uint(stts.EntryCount)
	}
	panic(fmt.Errorf("invalid name of dynamic-length field: boxType=stts fieldName=%s", name))
}
