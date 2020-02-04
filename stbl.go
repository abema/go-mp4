package mp4

func BoxTypeStbl() BoxType { return StrToBoxType("stbl") }

func init() {
	AddBoxDef(&Stbl{}, noVersion)
}

// Stbl is ISOBMFF stbl box type
type Stbl struct {
	Box
}

// GetType returns the BoxType
func (*Stbl) GetType() BoxType {
	return BoxTypeStbl()
}
