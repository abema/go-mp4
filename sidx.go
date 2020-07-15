package mp4

import "fmt"

func BoxTypeSidx() BoxType { return StrToBoxType("sidx") }

func init() {
	AddBoxDef(&Sidx{}, 0)
}

// Sidx is ISOBMFF Sidx box type
type Sidx struct {
	FullBox    `mp4:"extend"`
	ReferenceID uint32 `mp4:"size=32"`
	TimeScale uint32 `mp4:"size=32"`
	EarliestPresentationTime uint32 `mp4:"size=32"`
	FirstOffset uint32 `mp4:"size=32"`
	Reserved uint16 `mp4:"size=16"`
	ReferenceCount uint16 `mp4:"size=16"`
	References []struct {
		TypeSize uint32 `mp4:"size=32"`
		SubSegmentDuration uint32 `mp4:"size=32"`
		SapStartsTypeDeltaTime uint32 `mp4:"size=32"`
	} `mp4:"len=dynamic,size=96"`
}

// GetType returns the BoxType
func (*Sidx) GetType() BoxType {
	return BoxTypeSidx()
}

// GetFieldLength returns length of dynamic field
func (Sidx *Sidx) GetFieldLength(name string) uint {
	switch name {
	case "References":
		return uint(Sidx.ReferenceCount)
	}
	panic(fmt.Errorf("invalid name of dynamic-length field: boxType=Sidx fieldName=%s", name))
}
