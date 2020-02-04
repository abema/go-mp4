package mp4

func BoxTypeMfhd() BoxType { return StrToBoxType("mfhd") }

func init() {
	AddBoxDef(&Mfhd{}, 0)
}

// Mfhd is ISOBMFF mfhd box type
type Mfhd struct {
	FullBox        `mp4:"extend"`
	SequenceNumber uint32 `mp4:"size=32"`
}

// GetType returns the BoxType
func (*Mfhd) GetType() BoxType {
	return BoxTypeMfhd()
}
