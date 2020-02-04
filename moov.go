package mp4

func BoxTypeMoov() BoxType { return StrToBoxType("moov") }

func init() {
	AddBoxDef(&Moov{}, noVersion)
}

// Moov is ISOBMFF moov box type
type Moov struct {
	Box
}

// GetType returns the BoxType
func (*Moov) GetType() BoxType {
	return BoxTypeMoov()
}
