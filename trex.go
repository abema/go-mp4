package mp4

func BoxTypeTrex() BoxType { return StrToBoxType("trex") }

func init() {
	AddBoxDef(&Trex{}, 0)
}

// Trex is ISOBMFF trex box type
type Trex struct {
	FullBox                       `mp4:"extend"`
	TrackID                       uint32 `mp4:"size=32"`
	DefaultSampleDescriptionIndex uint32 `mp4:"size=32"`
	DefaultSampleDuration         uint32 `mp4:"size=32"`
	DefaultSampleSize             uint32 `mp4:"size=32"`
	DefaultSampleFlags            uint32 `mp4:"size=32,hex"`
}

// GetType returns the BoxType
func (*Trex) GetType() BoxType {
	return BoxTypeTrex()
}
