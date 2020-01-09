package mp4

import "fmt"

func BoxTypeCtts() BoxType { return StrToBoxType("ctts") }

func init() {
	AddBoxDef(&Ctts{}, 0, 1)
}

type Ctts struct {
	FullBox    `mp4:"extend"`
	EntryCount uint32 `mp4:"size=32"`
	Entries    []struct {
		SampleCount    uint32 `mp4:"size=32"`
		SampleOffsetV0 uint32 `mp4:"size=32,ver=0"`
		SampleOffsetV1 int32  `mp4:"size=32,ver=1"`
	} `mp4:"len=dynamic,size=64"`
}

// GetType returns the BoxType
func (*Ctts) GetType() BoxType {
	return BoxTypeCtts()
}

// GetFieldLength returns length of dynamic field
func (ctts *Ctts) GetFieldLength(name string) uint {
	switch name {
	case "Entries":
		return uint(ctts.EntryCount)
	}
	panic(fmt.Errorf("invalid name of dynamic-length field: boxType=ctts fieldName=%s", name))
}
