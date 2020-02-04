package mp4

func BoxTypeTfdt() BoxType { return StrToBoxType("tfdt") }

func init() {
	AddBoxDef(&Tfdt{}, 0, 1)
}

// Tfdt is ISOBMFF tfdt box type
type Tfdt struct {
	FullBox               `mp4:"extend"`
	BaseMediaDecodeTimeV0 uint32 `mp4:"size=32,ver=0"`
	BaseMediaDecodeTimeV1 uint64 `mp4:"size=64,ver=1"`
}

// GetType returns the BoxType
func (*Tfdt) GetType() BoxType {
	return BoxTypeTfdt()
}
