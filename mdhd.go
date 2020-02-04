package mp4

func BoxTypeMdhd() BoxType { return StrToBoxType("mdhd") }

func init() {
	AddBoxDef(&Mdhd{}, 0, 1)
}

// Mdhd is ISOBMFF mdhd box type
type Mdhd struct {
	FullBox `mp4:"extend"`
	// Version 0
	CreationTimeV0     uint32 `mp4:"size=32,ver=0"`
	ModificationTimeV0 uint32 `mp4:"size=32,ver=0"`
	TimescaleV0        uint32 `mp4:"size=32,ver=0"`
	DurationV0         uint32 `mp4:"size=32,ver=0"`
	// Version 1
	CreationTimeV1     uint64 `mp4:"size=64,ver=1"`
	ModificationTimeV1 uint64 `mp4:"size=64,ver=1"`
	TimescaleV1        uint32 `mp4:"size=32,ver=1"`
	DurationV1         uint64 `mp4:"size=64,ver=1"`
	//
	Pad        bool    `mp4:"size=1"`
	Language   [3]byte `mp4:"size=5,iso639-2"` // ISO-639-2/T language code
	PreDefined uint16  `mp4:"size=16"`
}

// GetType returns the BoxType
func (*Mdhd) GetType() BoxType {
	return BoxTypeMdhd()
}
