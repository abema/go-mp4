package mp4

func BoxTypeMdia() BoxType { return StrToBoxType("mdia") }

func init() {
	AddBoxDef(&Mdia{}, noVersion)
}

// Mdia is ISOBMFF mdia box type
type Mdia struct {
	Box
}

// GetType returns the BoxType
func (*Mdia) GetType() BoxType {
	return BoxTypeMdia()
}
