package mp4

func BoxTypeMfra() BoxType { return StrToBoxType("mfra") }

func init() {
	AddBoxDef(&Mfra{}, noVersion)
}

// Mfra is ISOBMFF mfra box type
type Mfra struct {
	Box
}

// GetType returns the BoxType
func (*Mfra) GetType() BoxType {
	return BoxTypeMfra()
}
