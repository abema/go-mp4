package mp4

func BoxTypeDref() BoxType { return StrToBoxType("dref") }
func BoxTypeUrl() BoxType  { return StrToBoxType("url ") }
func BoxTypeUrn() BoxType  { return StrToBoxType("urn ") }

func init() {
	AddBoxDef(&Dref{}, 0)
	AddBoxDef(&Url{}, 0)
	AddBoxDef(&Urn{}, 0)
}

// Dref is ISOBMFF dref box type
type Dref struct {
	FullBox    `mp4:"extend"`
	EntryCount uint32 `mp4:"size=32"`
}

// GetType returns the BoxType
func (*Dref) GetType() BoxType {
	return BoxTypeDref()
}

type Url struct {
	FullBox  `mp4:"extend"`
	Location string `mp4:"string,nopt=0x000001"`
}

func (*Url) GetType() BoxType {
	return BoxTypeUrl()
}

const UrlSelfContained = 0x000001

type Urn struct {
	FullBox  `mp4:"extend"`
	Name     string `mp4:"string,nopt=0x000001"`
	Location string `mp4:"string,nopt=0x000001"`
}

func (*Urn) GetType() BoxType {
	return BoxTypeUrn()
}

const UrnSelfContained = 0x000001
