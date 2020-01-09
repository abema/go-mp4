package mp4

func BoxTypeMvex() BoxType { return StrToBoxType("mvex") }

func init() {
	AddBoxDef(&Mvex{}, noVersion)
}

// Mvex is ISOBMFF mvex box type
type Mvex struct {
	Box
}

// GetType returns the BoxType
func (*Mvex) GetType() BoxType {
	return BoxTypeMvex()
}
