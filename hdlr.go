package mp4

import "fmt"

func BoxTypeHdlr() BoxType { return StrToBoxType("hdlr") }

func init() {
	AddBoxDef(&Hdlr{}, 0)
}

// Hdlr is ISOBMFF hdlr box type
type Hdlr struct {
	FullBox `mp4:"extend"`
	// Predefined corresponds to component_type of QuickTime.
	// pre_defined of ISO-14496 has always zero,
	// however component_type has "mhlr" or "dhlr".
	PreDefined  uint32    `mp4:"size=32"`
	HandlerType [4]byte   `mp4:"size=8,string"`
	Reserved    [3]uint32 `mp4:"size=32,const=0"`
	Name        string    `mp4:"string=c_p"`
}

// GetType returns the BoxType
func (*Hdlr) GetType() BoxType {
	return BoxTypeHdlr()
}

func (hdlr *Hdlr) IsPString(name string, bytes []byte, remainingSize uint64) bool {
	switch name {
	case "Name":
		return remainingSize == 0 && hdlr.PreDefined != 0
	default:
		panic(fmt.Errorf("invalid field name: name=%s", name))
	}
}
