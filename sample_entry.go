package mp4

func init() {
	AddAnyTypeBoxDef(&VisualSampleEntry{},
		StrToBoxType("avc1"),
		StrToBoxType("encv"))
	AddAnyTypeBoxDef(&AudioSampleEntry{},
		StrToBoxType("mp4a"),
		StrToBoxType("enca"))
	AddAnyTypeBoxDef(&AVCDecoderConfiguration{},
		StrToBoxType("avcC"))
	AddAnyTypeBoxDef(&PixelAspectRatioBox{},
		StrToBoxType("pasp"))
}

type SampleEntry struct {
	AnyTypeBox
	Reserved           [6]uint8 `mp4:"size=8,const=0"`
	DataReferenceIndex uint16   `mp4:"size=16"`
}

type VisualSampleEntry struct {
	SampleEntry     `mp4:"extend"`
	PreDefined      uint16    `mp4:"size=16"`
	Reserved        uint16    `mp4:"size=16,const=0"`
	PreDefined2     [3]uint32 `mp4:"size=32"`
	Width           uint16    `mp4:"size=16"`
	Height          uint16    `mp4:"size=16"`
	Horizresolution uint32    `mp4:"size=32"`
	Vertresolution  uint32    `mp4:"size=32"`
	Reserved2       uint32    `mp4:"size=32,const=0"`
	FrameCount      uint16    `mp4:"size=16"`
	Compressorname  [32]byte  `mp4:"size=8"`
	Depth           uint16    `mp4:"size=16"`
	PreDefined3     int16     `mp4:"size=16"`
}

// StringifyField returns field value as string
func (vse *VisualSampleEntry) StringifyField(name string, indent string, depth int) (string, bool) {
	switch name {
	case "Compressorname":
		end := 0
		for ; end < len(vse.Compressorname); end++ {
			if vse.Compressorname[end] == 0 {
				break
			}
		}
		return `"` + string(vse.Compressorname[:end]) + `"`, true
	default:
		return "", false
	}
}

type AudioSampleEntry struct {
	SampleEntry  `mp4:"extend"`
	EntryVersion uint16    `mp4:"size=16"`
	Reserved     [3]uint16 `mp4:"size=16,const=0,hidden"`
	ChannelCount uint16    `mp4:"size=16"`
	SampleSize   uint16    `mp4:"size=16"`
	PreDefined   uint16    `mp4:"size=16"`
	Reserved2    uint16    `mp4:"size=16,const=0,hidden"`
	SampleRate   uint32    `mp4:"size=32"`
}

type AVCDecoderConfiguration struct {
	AnyTypeBox
	ConfigurationVersion uint8 `mp4:"size=8"`
	Profile              uint8 `mp4:"size=8"`
	ProfileCompatibility uint8 `mp4:"size=8"`
	Level                uint8 `mp4:"size=8"`
	// TODO: Refer to ISO/IEC 14496-15
	Data []byte `mp4:"size=8"`
}

type PixelAspectRatioBox struct {
	AnyTypeBox
	HSpacing uint32 `mp4:"size=32"`
	VSpacing uint32 `mp4:"size=32"`
}
