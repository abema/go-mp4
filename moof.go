package mp4

func BoxTypeMoof() BoxType { return StrToBoxType("moof") }

func init() {
	AddBoxDef(&Moof{}, noVersion)
}

// Moof is ISOBMFF moof box type
type Moof struct {
	Box
}

// GetType returns the BoxType
func (*Moof) GetType() BoxType {
	return BoxTypeMoof()
}
