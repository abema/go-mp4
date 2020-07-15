package mp4

func BoxTypeStyp() BoxType { return StrToBoxType("styp") }

func init() {
	AddBoxDef(&Styp{}, noVersion)
}

// Styp is as the same as Ftype
type Styp struct {
	Box
	MajorBrand       [4]byte `mp4:"size=8,string"`
	MinorVersion     uint32  `mp4:"size=32"`
	CompatibleBrands []struct {
		CompatibleBrand [4]byte `mp4:"size=8,string"`
	} `mp4:"size=32"` // to end of the box
}

// GetType returns the BoxType
func (*Styp) GetType() BoxType {
	return BoxTypeStyp()
}
