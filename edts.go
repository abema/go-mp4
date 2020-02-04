package mp4

func BoxTypeEdts() BoxType { return StrToBoxType("edts") }

func init() {
	AddBoxDef(&Edts{}, noVersion)
}

// Edts is ISOBMFF edts box type
type Edts struct {
	Box
}

// GetType returns the BoxType
func (*Edts) GetType() BoxType {
	return BoxTypeEdts()
}
