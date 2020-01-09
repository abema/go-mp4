package mp4

import "fmt"

func BoxTypeStss() BoxType { return StrToBoxType("stss") }

func init() {
	AddBoxDef(&Stss{}, 0)
}

type Stss struct {
	FullBox      `mp4:"extend"`
	EntryCount   uint32   `mp4:"size=32"`
	SampleNumber []uint32 `mp4:"len=dynamic,size=32"`
}

// GetType returns the BoxType
func (*Stss) GetType() BoxType {
	return BoxTypeStss()
}

// GetFieldLength returns length of dynamic field
func (stss *Stss) GetFieldLength(name string) uint {
	switch name {
	case "SampleNumber":
		return uint(stss.EntryCount)
	}
	panic(fmt.Errorf("invalid name of dynamic-length field: boxType=stss fieldName=%s", name))
}
