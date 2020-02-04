package mp4

func BoxTypeMvhd() BoxType { return StrToBoxType("mvhd") }

func init() {
	AddBoxDef(&Mvhd{}, 0, 1)
}

// Mvhd is ISOBMFF mvhd box type
type Mvhd struct {
	FullBox            `mp4:"extend"`
	CreationTimeV0     uint32    `mp4:"size=32,ver=0"`
	ModificationTimeV0 uint32    `mp4:"size=32,ver=0"`
	TimescaleV0        uint32    `mp4:"size=32,ver=0"`
	DurationV0         uint32    `mp4:"size=32,ver=0"`
	CreationTimeV1     uint64    `mp4:"size=64,ver=1"`
	ModificationTimeV1 uint64    `mp4:"size=64,ver=1"`
	TimescaleV1        uint32    `mp4:"size=32,ver=1"`
	DurationV1         uint64    `mp4:"size=64,ver=1"`
	Rate               int32     `mp4:"size=32"` // template=0x00010000
	Volume             int16     `mp4:"size=16"` // template=0x0100
	Reserved           int16     `mp4:"size=16,const=0"`
	Reserved2          [2]uint32 `mp4:"size=32,const=0"`
	Matrix             [9]int32  `mp4:"size=32,hex"` // template={ 0x00010000,0,0,0,0x00010000,0,0,0,0x40000000 }
	PreDefined         [6]int32  `mp4:"size=32"`
	NextTrackID        uint32    `mp4:"size=32"`
}

// GetType returns the BoxType
func (*Mvhd) GetType() BoxType {
	return BoxTypeMvhd()
}
