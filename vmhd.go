package mp4

func BoxTypeVmhd() BoxType { return StrToBoxType("vmhd") }

func init() {
	AddBoxDef(&Vmhd{}, 0)
}

// Vmhd is ISOBMFF vmhd box type
type Vmhd struct {
	FullBox      `mp4:"extend"`
	Graphicsmode uint16    `mp4:"size=16"` // template=0
	Opcolor      [3]uint16 `mp4:"size=16"` // template={0, 0, 0}
}

// GetType returns the BoxType
func (*Vmhd) GetType() BoxType {
	return BoxTypeVmhd()
}
