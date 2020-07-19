package mp4

func BoxTypeFree() BoxType { return StrToBoxType("free") }
func BoxTypeSkip() BoxType { return StrToBoxType("skip") }

func init() {
	AddBoxDef(&Free{}, noVersion)
	AddBoxDef(&Skip{}, noVersion)
}

type FreeSpace struct {
	Box
	Data []uint8 `mp4:"size=8"`
}

type Free FreeSpace

func (*Free) GetType() BoxType {
	return BoxTypeFree()
}

type Skip FreeSpace

func (*Skip) GetType() BoxType {
	return BoxTypeSkip()
}
