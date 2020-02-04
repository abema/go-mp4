package mp4

func BoxTypeTrak() BoxType { return StrToBoxType("trak") }

func init() {
	AddBoxDef(&Trak{}, noVersion)
}

// Trak is ISOBMFF trak box type
type Trak struct {
	Box
}

// GetType returns the BoxType
func (*Trak) GetType() BoxType {
	return BoxTypeTrak()
}
