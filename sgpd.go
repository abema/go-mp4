package mp4

func BoxTypeSgpd() BoxType { return StrToBoxType("sgpd") }

func init() {
	AddBoxDef(&Sgpd{}, 0, 1, 2)
}

type Sgpd struct {
	FullBox                         `mp4:"extend"`
	GroupingType                    [4]byte `mp4:"size=8"`
	DefaultLength                   uint32  `mp4:"size=32,ver=1"`
	DefaultSample_description_index uint32  `mp4:"size=32,ver=2"`
	EntryCount                      uint32  `mp4:"size=32"`
	RollDistances                   []int16 `mp4:"size=16,opt=dynamic"`
	//AlternativeStartupEntry         []AlternativeStartupEntry `mp4:"size=dynamic,opt=dynamic"`
	VisualRandomAccessEntry []VisualRandomAccessEntry `mp4:"size=dynamic,opt=dynamic"`
	TemporalLevelEntry      []TemporalLevelEntry      `mp4:"size=dynamic,opt=dynamic"`
	Unsupported             []byte                    `mp4:"size=8,opt=dynamic"`
}

/* go-mp4 has not supported nested dynamic field.
type AlternativeStartupEntry struct {
	RollCount         uint16                       `mp4:"size=16`
	FirstOutputSample uint16                       `mp4:"size=16`
	SampleOffset      []uint32                     `mp4:"size=32,len=dynamic`
	Opts              []AlternativeStartupEntryOpt `mp4:"size=32`
}

type AlternativeStartupEntryOpt struct {
	NumOutputSamples uint16 `mp4:"size=16`
	NumTotalSamples  uint16 `mp4:"size=16`
}
*/

type VisualRandomAccessEntry struct {
	NumLeadingSamplesKnown uint8 `mp4:"size=1"`
	NumLeadingSamples      uint8 `mp4:"size=7"`
}

type TemporalLevelEntry struct {
	LevelUndependentlyUecodable bool  `mp4:"size=1"`
	Reserved                    uint8 `mp4:"size=7,const=0"`
}

func (sgpd *Sgpd) GetFieldSize(name string) uint {
	switch name {
	case "RollDistances",
		//"AlternativeStartupEntry",
		"VisualRandomAccessEntry",
		"TemporalLevelEntry":
		return uint(sgpd.DefaultLength * 8)
	default:
		return 0
	}
}

func (sgpd *Sgpd) IsOptFieldEnabled(name string) bool {
	switch name {
	case "RollDistances":
		return sgpd.Version == 1 &&
			(sgpd.GroupingType == [4]byte{'r', 'o', 'l', 'l'} ||
				sgpd.GroupingType == [4]byte{'p', 'r', 'o', 'l'}) &&
			sgpd.DefaultLength == 2
		/*
			case "AlternativeStartupEntry":
				return sgpd.Version == 1 &&
					sgpd.GroupingType == [4]byte{'a', 'l', 's', 't'} &&
					sgpd.DefaultLength != 0
		*/
	case "VisualRandomAccessEntry":
		return sgpd.Version == 1 &&
			sgpd.GroupingType == [4]byte{'r', 'a', 'p', ' '} &&
			sgpd.DefaultLength == 1
	case "TemporalLevelEntry":
		return sgpd.Version == 1 &&
			sgpd.GroupingType == [4]byte{'t', 'e', 'l', 'e'} &&
			sgpd.DefaultLength == 1
	case "Unsupported":
		return sgpd.Version == 0 || (sgpd.Version == 1 && sgpd.DefaultLength == 0)
	default:
		return false
	}
}

func (sgpd *Sgpd) StringifyField(name string, indent string, depth int) (string, bool) {
	switch name {
	case "GroupingType":
		return string([]byte{sgpd.GroupingType[0], sgpd.GroupingType[1], sgpd.GroupingType[2], sgpd.GroupingType[3]}), true
	default:
		return "", false
	}
}

func (*Sgpd) GetType() BoxType {
	return BoxTypeSgpd()
}
