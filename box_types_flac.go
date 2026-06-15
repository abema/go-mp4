package mp4

// https://github.com/xiph/flac/blob/master/doc/isoflac.txt

/*************************** fLaC ****************************/

func BoxTypeFLaC() BoxType { return StrToBoxType("fLaC") }

func init() {
	AddAnyTypeBoxDef(&AudioSampleEntry{}, BoxTypeFLaC())
}

/*************************** dfLa ****************************/

func BoxTypeDfLa() BoxType { return StrToBoxType("dfLa") }

func init() {
	AddBoxDef(&DfLa{}, 0)
}

type DfLa struct {
	FullBox `mp4:"0,extend"`
	Blocks  []byte `mp4:"1,size=8"`
}

func (*DfLa) GetType() BoxType {
	return BoxTypeDfLa()
}
