package mp4

import "fmt"

func BoxTypeStsz() BoxType { return StrToBoxType("stsz") }

func init() {
	AddBoxDef(&Stsz{}, 0)
}

// Stsz is ISOBMFF stsz box type
type Stsz struct {
	FullBox     `mp4:"extend"`
	SampleSize  uint32   `mp4:"size=32"`
	SampleCount uint32   `mp4:"size=32"`
	EntrySize   []uint32 `mp4:"size=32,len=dynamic"`
}

// GetType returns the BoxType
func (*Stsz) GetType() BoxType {
	return BoxTypeStsz()
}

// GetFieldLength returns length of dynamic field
func (stsz *Stsz) GetFieldLength(name string) uint {
	switch name {
	case "EntrySize":
		if stsz.SampleSize == 0 {
			return uint(stsz.SampleCount)
		} else {
			return 0
		}
	}
	panic(fmt.Errorf("invalid name of dynamic-length field: boxType=stsz fieldName=%s", name))
}
