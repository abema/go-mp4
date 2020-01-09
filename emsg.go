package mp4

func BoxTypeEmsg() BoxType { return StrToBoxType("emsg") }

func init() {
	AddBoxDef(&Emsg{}, 0)
}

// Emsg is ISOBMFF emsg box type
type Emsg struct {
	FullBox               `mp4:"extend"`
	SchemeIdUri           string `mp4:"string"`
	Value                 string `mp4:"string"`
	Timescale             uint32 `mp4:"size=32"`
	PresentationTimeDelta uint32 `mp4:"size=32"`
	EventDuration         uint32 `mp4:"size=32"`
	Id                    uint32 `mp4:"size=32"`
	MessageData           []byte `mp4:"size=8,string"`
}

// GetType returns the BoxType
func (*Emsg) GetType() BoxType {
	return BoxTypeEmsg()
}
