package mp4

func BoxTypeMdat() BoxType { return StrToBoxType("mdat") }

func init() {
	AddBoxDef(&Mdat{}, noVersion)
}

// Mdat is ISOBMFF mdat box type
type Mdat struct {
	Box
	Data []byte `mp4:"size=8"`
}

// GetType returns the BoxType
func (*Mdat) GetType() BoxType {
	return BoxTypeMdat()
}
