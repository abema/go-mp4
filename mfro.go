package mp4

func BoxTypeMfro() BoxType { return StrToBoxType("mfro") }

func init() {
	AddBoxDef(&Mfro{}, 0)
}

// Mfro is ISOBMFF mfro box type
type Mfro struct {
	FullBox `mp4:"extend"`
	Size    uint32 `mp4:"size=32"`
}

// GetType returns the BoxType
func (*Mfro) GetType() BoxType {
	return BoxTypeMfro()
}
