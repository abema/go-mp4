package mp4

func BoxTypeDinf() BoxType { return StrToBoxType("dinf") }

func init() {
	AddBoxDef(&Dinf{}, noVersion)
}

// Dinf is ISOBMFF dinf box type
type Dinf struct {
	Box
}

// GetType returns the BoxType
func (*Dinf) GetType() BoxType {
	return BoxTypeDinf()
}
