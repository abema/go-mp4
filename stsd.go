package mp4

func BoxTypeStsd() BoxType { return StrToBoxType("stsd") }

func init() {
	AddBoxDef(&Stsd{}, 0)
}

// Stsd is ISOBMFF stsd box type
type Stsd struct {
	FullBox    `mp4:"extend"`
	EntryCount uint32 `mp4:"size=32"`
}

// GetType returns the BoxType
func (*Stsd) GetType() BoxType {
	return BoxTypeStsd()
}
