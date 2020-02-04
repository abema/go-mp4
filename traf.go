package mp4

func BoxTypeTraf() BoxType { return StrToBoxType("traf") }

func init() {
	AddBoxDef(&Traf{}, noVersion)
}

// Traf is ISOBMFF traf box type
type Traf struct {
	Box
}

// GetType returns the BoxType
func (*Traf) GetType() BoxType {
	return BoxTypeTraf()
}
