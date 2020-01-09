package mp4

func BoxTypeUdta() BoxType { return StrToBoxType("udta") }

func init() {
	AddBoxDef(&Udta{}, noVersion)
}

// Udta is ISOBMFF udta box type
type Udta struct {
	Box
}

// GetType returns the BoxType
func (*Udta) GetType() BoxType {
	return BoxTypeUdta()
}
