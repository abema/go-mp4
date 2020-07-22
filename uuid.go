package mp4

func BoxTypeUuid() BoxType { return StrToBoxType("uuid") }

func init() {
	AddBoxDef(&Uuid{}, noVersion)
}

// Uuid is ISOBMFF uuid box type
type Uuid struct {
	Box
	UserType [16]uint8 `mp4:"size=8,len=16"`
	Data     []byte    `mp4:"size=8"`
}

// GetType returns the BoxType
func (*Uuid) GetType() BoxType {
	return BoxTypeUuid()
}
