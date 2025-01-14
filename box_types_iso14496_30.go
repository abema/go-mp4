package mp4

/*********************** WebVTT Sample Entry ****************************/

func BoxTypeVttC() BoxType { return StrToBoxType("vttC") }
func BoxTypeVlab() BoxType { return StrToBoxType("vlab") }
func BoxTypeWvtt() BoxType { return StrToBoxType("wvtt") }

func init() {
	AddBoxDef(&WebVTTConfigurationBox{})
	AddBoxDef(&WebVTTSourceLabelBox{})
	AddAnyTypeBoxDef(&WVTTSampleEntry{}, BoxTypeWvtt())
}

type WebVTTConfigurationBox struct {
	Box
	Config string `mp4:"0,boxstring"`
}

func (WebVTTConfigurationBox) GetType() BoxType {
	return BoxTypeVttC()
}

type WebVTTSourceLabelBox struct {
	Box
	SourceLabel string `mp4:"0,boxstring"`
}

func (WebVTTSourceLabelBox) GetType() BoxType {
	return BoxTypeVlab()
}

type WVTTSampleEntry struct {
	SampleEntry `mp4:"0,extend"`
}

/*********************** WebVTT Sample Format ****************************/

func BoxTypeVttc() BoxType { return StrToBoxType("vttc") }
func BoxTypeVsid() BoxType { return StrToBoxType("vsid") }
func BoxTypeCtim() BoxType { return StrToBoxType("ctim") }
func BoxTypeIden() BoxType { return StrToBoxType("iden") }
func BoxTypeSttg() BoxType { return StrToBoxType("sttg") }
func BoxTypePayl() BoxType { return StrToBoxType("payl") }
func BoxTypeVtte() BoxType { return StrToBoxType("vtte") }
func BoxTypeVtta() BoxType { return StrToBoxType("vtta") }

func init() {
	AddBoxDef(&VTTCueBox{})
	AddBoxDef(&CueSourceIDBox{})
	AddBoxDef(&CueTimeBox{})
	AddBoxDef(&CueIDBox{})
	AddBoxDef(&CueSettingsBox{})
	AddBoxDef(&CuePayloadBox{})
	AddBoxDef(&VTTEmptyCueBox{})
	AddBoxDef(&VTTAdditionalTextBox{})
}

type VTTCueBox struct {
	Box
}

func (VTTCueBox) GetType() BoxType {
	return BoxTypeVttc()
}

type CueSourceIDBox struct {
	Box
	SourceId uint32 `mp4:"0,size=32"`
}

func (CueSourceIDBox) GetType() BoxType {
	return BoxTypeVsid()
}

type CueTimeBox struct {
	Box
	CueCurrentTime string `mp4:"0,boxstring"`
}

func (CueTimeBox) GetType() BoxType {
	return BoxTypeCtim()
}

type CueIDBox struct {
	Box
	CueId string `mp4:"0,boxstring"`
}

func (CueIDBox) GetType() BoxType {
	return BoxTypeIden()
}

type CueSettingsBox struct {
	Box
	Settings string `mp4:"0,boxstring"`
}

func (CueSettingsBox) GetType() BoxType {
	return BoxTypeSttg()
}

type CuePayloadBox struct {
	Box
	CueText string `mp4:"0,boxstring"`
}

func (CuePayloadBox) GetType() BoxType {
	return BoxTypePayl()
}

type VTTEmptyCueBox struct {
	Box
}

func (VTTEmptyCueBox) GetType() BoxType {
	return BoxTypeVtte()
}

type VTTAdditionalTextBox struct {
	Box
	CueAdditionalText string `mp4:"0,boxstring"`
}

func (VTTAdditionalTextBox) GetType() BoxType {
	return BoxTypeVtta()
}
