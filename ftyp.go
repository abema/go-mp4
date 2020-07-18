package mp4

func BoxTypeFtyp() BoxType { return StrToBoxType("ftyp") }

func init() {
	AddBoxDef(&Ftyp{}, noVersion)
}

// Ftyp is ISOBMFF ftyp box type
type Ftyp struct {
	Box
	MajorBrand       [4]byte               `mp4:"size=8,string"`
	MinorVersion     uint32                `mp4:"size=32"`
	CompatibleBrands []CompatibleBrandElem `mp4:"size=32"` // to end of the box
}

type CompatibleBrandElem struct {
	CompatibleBrand [4]byte `mp4:"size=8,string"`
}

// GetType returns the BoxType
func (*Ftyp) GetType() BoxType {
	return BoxTypeFtyp()
}
