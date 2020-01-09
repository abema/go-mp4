package mp4

func BoxTypeMeta() BoxType { return StrToBoxType("meta") }

func init() {
	AddBoxDef(&Meta{}, 0)
}

// Meta is ISOBMFF meta box type
type Meta struct {
	FullBox `mp4:"extend"`
}

// GetType returns the BoxType
func (*Meta) GetType() BoxType {
	return BoxTypeMeta()
}
