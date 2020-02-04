package mp4

func BoxTypeMehd() BoxType { return StrToBoxType("mehd") }

func init() {
	AddBoxDef(&Mehd{}, 0, 1)
}

// Mehd is ISOBMFF mehd box type
type Mehd struct {
	FullBox            `mp4:"extend"`
	FragmentDurationV0 uint32 `mp4:"size=32,ver=0"`
	FragmentDurationV1 uint64 `mp4:"size=64,ver=1"`
}

// GetType returns the BoxType
func (*Mehd) GetType() BoxType {
	return BoxTypeMehd()
}
