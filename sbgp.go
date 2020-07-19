package mp4

import "fmt"

func BoxTypeSbgp() BoxType { return StrToBoxType("sbgp") }

func init() {
	AddBoxDef(&Sbgp{}, 0, 1)
}

type Sbgp struct {
	FullBox                 `mp4:"extend"`
	GroupingType            uint32      `mp4:"size=32"`
	grouping_type_parameter uint32      `mp4:"size=32,ver=1"`
	EntryCount              uint32      `mp4:"size=32"`
	Entries                 []SbgpEntry `mp4:"len=dynamic,size=64"`
}

type SbgpEntry struct {
	SampleCount           uint32 `mp4:"size=32"`
	GroupDescriptionIndex uint32 `mp4:"size=32"`
}

func (sbgp *Sbgp) GetFieldLength(name string) uint {
	switch name {
	case "Entries":
		return uint(sbgp.EntryCount)
	}
	panic(fmt.Errorf("invalid name of dynamic-length field: boxType=sbgp fieldName=%s", name))
}

func (*Sbgp) GetType() BoxType {
	return BoxTypeSbgp()
}
