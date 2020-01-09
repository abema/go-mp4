package mp4

func BoxTypeTkhd() BoxType { return StrToBoxType("tkhd") }

func init() {
	AddBoxDef(&Tkhd{}, 0, 1)
}

// Tkhd is ISOBMFF tkhd box type
type Tkhd struct {
	FullBox `mp4:"extend"`
	// Version 0
	CreationTimeV0     uint32 `mp4:"size=32,ver=0"`
	ModificationTimeV0 uint32 `mp4:"size=32,ver=0"`
	TrackIDV0          uint32 `mp4:"size=32,ver=0"`
	ReservedV0         uint32 `mp4:"size=32,ver=0,const=0"`
	DurationV0         uint32 `mp4:"size=32,ver=0"`
	// Version 1
	CreationTimeV1     uint64 `mp4:"size=64,ver=1"`
	ModificationTimeV1 uint64 `mp4:"size=64,ver=1"`
	TrackIDV1          uint32 `mp4:"size=32,ver=1"`
	ReservedV1         uint32 `mp4:"size=32,ver=1,const=0"`
	DurationV1         uint64 `mp4:"size=64,ver=1"`
	//
	Reserved       [2]uint32 `mp4:"size=32,const=0"`
	Layer          int16     `mp4:"size=16"` // template=0
	AlternateGroup int16     `mp4:"size=16"` // template=0
	Volume         int16     `mp4:"size=16"` // template={if track_is_audio 0x0100 else 0}
	Reserved2      uint16    `mp4:"size=16,const=0"`
	Matrix         [9]int32  `mp4:"size=32,hex"` // template={ 0x00010000,0,0,0,0x00010000,0,0,0,0x40000000 };
	Width          uint32    `mp4:"size=32"`
	Height         uint32    `mp4:"size=32"`
}

// GetType returns the BoxType
func (*Tkhd) GetType() BoxType {
	return BoxTypeTkhd()
}
