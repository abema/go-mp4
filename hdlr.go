package mp4

func BoxTypeHdlr() BoxType { return StrToBoxType("hdlr") }

func init() {
	AddBoxDef(&Hdlr{}, 0)
}

// Hdlr is ISOBMFF hdlr box type
type Hdlr struct {
	FullBox     `mp4:"extend"`
	PreDefined  uint32    `mp4:"size=32"`
	HandlerType [4]byte   `mp4:"size=8,string"`
	Reserved    [3]uint32 `mp4:"size=32,const=0"`
	Name        string    `mp4:"string"`
}

// GetType returns the BoxType
func (*Hdlr) GetType() BoxType {
	return BoxTypeHdlr()
}
