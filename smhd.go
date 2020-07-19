package mp4

func BoxTypeSmhd() BoxType { return StrToBoxType("smhd") }

func init() {
	AddBoxDef(&Smhd{}, 0)
}

type Smhd struct {
	FullBox  `mp4:"extend"`
	Balance  int16  `mp4:"size=16"` // template=0
	Reserved uint16 `mp4:"size=16,const=0"`
}

func (*Smhd) GetType() BoxType {
	return BoxTypeSmhd()
}
