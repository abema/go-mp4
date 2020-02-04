package mp4

import "fmt"

func BoxTypeStsc() BoxType { return StrToBoxType("stsc") }

func init() {
	AddBoxDef(&Stsc{}, 0)
}

// Stsc is ISOBMFF stsc box type
type Stsc struct {
	FullBox    `mp4:"extend"`
	EntryCount uint32 `mp4:"size=32"`
	Entries    []struct {
		FirstChunk             uint32 `mp4:"size=32"`
		SamplesPerChunk        uint32 `mp4:"size=32"`
		SampleDescriptionIndex uint32 `mp4:"size=32"`
	} `mp4:"len=dynamic,size=96"`
}

// GetType returns the BoxType
func (*Stsc) GetType() BoxType {
	return BoxTypeStsc()
}

// GetFieldLength returns length of dynamic field
func (stsc *Stsc) GetFieldLength(name string) uint {
	switch name {
	case "Entries":
		return uint(stsc.EntryCount)
	}
	panic(fmt.Errorf("invalid name of dynamic-length field: boxType=stsc fieldName=%s", name))
}
