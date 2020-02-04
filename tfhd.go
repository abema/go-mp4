package mp4

func BoxTypeTfhd() BoxType { return StrToBoxType("tfhd") }

func init() {
	AddBoxDef(&Tfhd{}, 0)
}

// Tfhd is ISOBMFF tfhd box type
type Tfhd struct {
	FullBox `mp4:"extend"`
	TrackID uint32 `mp4:"size=32"`

	// optional
	BaseDataOffset         uint64 `mp4:"size=64,opt=0x000001"`
	SampleDescriptionIndex uint32 `mp4:"size=32,opt=0x000002"`
	DefaultSampleDuration  uint32 `mp4:"size=32,opt=0x000008"`
	DefaultSampleSize      uint32 `mp4:"size=32,opt=0x000010"`
	DefaultSampleFlags     uint32 `mp4:"size=32,opt=0x000020,hex"`
}

const (
	TfhdBaseDataOffsetPresent         = 0x000001
	TfhdSampleDescriptionIndexPresent = 0x000002
	TfhdDefaultSampleDurationPresent  = 0x000008
	TfhdDefaultSampleSizePresent      = 0x000010
	TfhdDefaultSampleFlagsPresent     = 0x000020
	TfhdDurationIsEmpty               = 0x010000
	TfhdDefaultBaseIsMoof             = 0x020000
)

// GetType returns the BoxType
func (*Tfhd) GetType() BoxType {
	return BoxTypeTfhd()
}
