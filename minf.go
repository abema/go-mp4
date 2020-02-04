package mp4

func BoxTypeMinf() BoxType { return StrToBoxType("minf") }

func init() {
	AddBoxDef(&Minf{}, noVersion)
}

// Minf is ISOBMFF minf box type
type Minf struct {
	Box
}

// GetType returns the BoxType
func (*Minf) GetType() BoxType {
	return BoxTypeMinf()
}
