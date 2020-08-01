package mp4

import (
	"bytes"
	"fmt"
	"io"
)

/*************************** ctts ****************************/

func BoxTypeCtts() BoxType { return StrToBoxType("ctts") }

func init() {
	AddBoxDef(&Ctts{}, 0, 1)
}

type Ctts struct {
	FullBox    `mp4:"extend"`
	EntryCount uint32      `mp4:"size=32"`
	Entries    []CttsEntry `mp4:"len=dynamic,size=64"`
}

type CttsEntry struct {
	SampleCount    uint32 `mp4:"size=32"`
	SampleOffsetV0 uint32 `mp4:"size=32,ver=0"`
	SampleOffsetV1 int32  `mp4:"size=32,ver=1"`
}

// GetType returns the BoxType
func (*Ctts) GetType() BoxType {
	return BoxTypeCtts()
}

// GetFieldLength returns length of dynamic field
func (ctts *Ctts) GetFieldLength(name string) uint {
	switch name {
	case "Entries":
		return uint(ctts.EntryCount)
	}
	panic(fmt.Errorf("invalid name of dynamic-length field: boxType=ctts fieldName=%s", name))
}

/*************************** dinf ****************************/

func BoxTypeDinf() BoxType { return StrToBoxType("dinf") }

func init() {
	AddBoxDef(&Dinf{})
}

// Dinf is ISOBMFF dinf box type
type Dinf struct {
	Box
}

// GetType returns the BoxType
func (*Dinf) GetType() BoxType {
	return BoxTypeDinf()
}

/*************************** dref ****************************/

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

/*************************** edts ****************************/

func BoxTypeEdts() BoxType { return StrToBoxType("edts") }

func init() {
	AddBoxDef(&Edts{})
}

// Edts is ISOBMFF edts box type
type Edts struct {
	Box
}

// GetType returns the BoxType
func (*Edts) GetType() BoxType {
	return BoxTypeEdts()
}

/*************************** elst ****************************/

func BoxTypeElst() BoxType { return StrToBoxType("elst") }

func init() {
	AddBoxDef(&Elst{}, 0, 1)
}

// Elst is ISOBMFF elst box type
type Elst struct {
	FullBox    `mp4:"extend"`
	EntryCount uint32      `mp4:"size=32"`
	Entries    []ElstEntry `mp4:"len=dynamic,size=dynamic"`
}

type ElstEntry struct {
	SegmentDurationV0 uint32 `mp4:"size=32,ver=0"`
	MediaTimeV0       int32  `mp4:"size=32,ver=0"`
	SegmentDurationV1 uint64 `mp4:"size=64,ver=1"`
	MediaTimeV1       int64  `mp4:"size=64,ver=1"`
	MediaRateInteger  int16  `mp4:"size=16"`
	MediaRateFraction int16  `mp4:"size=16,const=0"`
}

// GetType returns the BoxType
func (*Elst) GetType() BoxType {
	return BoxTypeElst()
}

// GetFieldSize returns size of dynamic field
func (elst *Elst) GetFieldSize(name string) uint {
	switch name {
	case "Entries":
		switch elst.GetVersion() {
		case 0:
			return 0 +
				/* segmentDurationV0 */ 32 +
				/* mediaTimeV0       */ 32 +
				/* mediaRateInteger  */ 16 +
				/* mediaRateFraction */ 16
		case 1:
			return 0 +
				/* segmentDurationV1 */ 64 +
				/* mediaTimeV1       */ 64 +
				/* mediaRateInteger  */ 16 +
				/* mediaRateFraction */ 16
		}
	}
	panic(fmt.Errorf("invalid name of dynamic-size field: boxType=elst fieldName=%s", name))
}

// GetFieldLength returns length of dynamic field
func (elst *Elst) GetFieldLength(name string) uint {
	switch name {
	case "Entries":
		return uint(elst.EntryCount)
	}
	panic(fmt.Errorf("invalid name of dynamic-length field: boxType=elst fieldName=%s", name))
}

/*************************** emsg ****************************/

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

/*************************** esds ****************************/

// https://developer.apple.com/library/content/documentation/QuickTime/QTFF/QTFFChap3/qtff3.html

func BoxTypeEsds() BoxType { return StrToBoxType("esds") }

func init() {
	AddBoxDef(&Esds{}, 0)
}

const (
	ESDescrTag            = 0x03
	DecoderConfigDescrTag = 0x04
	DecSpecificInfoTag    = 0x05
	SLConfigDescrTag      = 0x06
)

// Esds is ES descripter box
type Esds struct {
	FullBox     `mp4:"extend"`
	Descriptors []Descriptor `mp4:"array"`
}

// GetType returns the BoxType
func (*Esds) GetType() BoxType {
	return BoxTypeEsds()
}

type Descriptor struct {
	BaseCustomFieldObject
	Tag                     int8                     `mp4:"size=8"` // must be 0x03
	Size                    uint32                   `mp4:"varint"`
	ESDescriptor            *ESDescriptor            `mp4:"extend,opt=dynamic"`
	DecoderConfigDescriptor *DecoderConfigDescriptor `mp4:"extend,opt=dynamic"`
	Data                    []byte                   `mp4:"size=8,opt=dynamic,len=dynamic"`
}

// GetFieldLength returns length of dynamic field
func (ds *Descriptor) GetFieldLength(name string) uint {
	switch name {
	case "Data":
		return uint(ds.Size)
	}
	panic(fmt.Errorf("invalid name of dynamic-length field: boxType=esds fieldName=%s", name))
}

func (ds *Descriptor) IsOptFieldEnabled(name string) bool {
	switch ds.Tag {
	case ESDescrTag:
		return name == "ESDescriptor"
	case DecoderConfigDescrTag:
		return name == "DecoderConfigDescriptor"
	default:
		return name == "Data"
	}
}

// StringifyField returns field value as string
func (ds *Descriptor) StringifyField(name string, indent string, depth int) (string, bool) {
	switch name {
	case "Tag":
		switch ds.Tag {
		case ESDescrTag:
			return "ESDescr", true
		case DecoderConfigDescrTag:
			return "DecoderConfigDescr", true
		case DecSpecificInfoTag:
			return "DecSpecificInfo", true
		case SLConfigDescrTag:
			return "SLConfigDescr", true
		default:
			return "", false
		}
	default:
		return "", false
	}
}

type ESDescriptor struct {
	BaseCustomFieldObject
	ESID                 uint16 `mp4:"size=16"`
	StreamDependenceFlag bool   `mp4:"size=1"`
	UrlFlag              bool   `mp4:"size=1"`
	OcrStreamFlag        bool   `mp4:"size=1"`
	StreamPriority       int8   `mp4:"size=5"`
	DependsOnESID        uint16 `mp4:"size=16,opt=dynamic"`
	URLLength            uint8  `mp4:"size=8,opt=dynamic"`
	URLString            []byte `mp4:"size=8,len=dynamic,opt=dynamic,string"`
	OCRESID              uint16 `mp4:"size=16,opt=dynamic"`
}

func (esds *ESDescriptor) GetFieldLength(name string) uint {
	switch name {
	case "URLString":
		return uint(esds.URLLength)
	}
	panic(fmt.Errorf("invalid name of dynamic-length field: boxType=ESDescriptor fieldName=%s", name))
}

func (esds *ESDescriptor) IsOptFieldEnabled(name string) bool {
	switch name {
	case "DependsOnESID":
		return esds.StreamDependenceFlag
	case "URLLength", "URLString":
		return esds.UrlFlag
	case "OCRESID":
		return esds.OcrStreamFlag
	default:
		return false
	}
}

type DecoderConfigDescriptor struct {
	BaseCustomFieldObject
	ObjectTypeIndication byte   `mp4:"size=8"`
	StreamType           int8   `mp4:"size=6"`
	UpStream             bool   `mp4:"size=1"`
	Reserved             bool   `mp4:"size=1"`
	BufferSizeDB         uint32 `mp4:"size=24"`
	MaxBitrate           uint32 `mp4:"size=32"`
	AvgBitrate           uint32 `mp4:"size=32"`
}

/************************ free, skip *************************/

func BoxTypeFree() BoxType { return StrToBoxType("free") }
func BoxTypeSkip() BoxType { return StrToBoxType("skip") }

func init() {
	AddBoxDef(&Free{})
	AddBoxDef(&Skip{})
}

type FreeSpace struct {
	Box
	Data []uint8 `mp4:"size=8"`
}

type Free FreeSpace

func (*Free) GetType() BoxType {
	return BoxTypeFree()
}

type Skip FreeSpace

func (*Skip) GetType() BoxType {
	return BoxTypeSkip()
}

/*************************** ftyp ****************************/

func BoxTypeFtyp() BoxType { return StrToBoxType("ftyp") }

func init() {
	AddBoxDef(&Ftyp{})
}

// Ftyp is ISOBMFF ftyp box type
type Ftyp struct {
	Box
	MajorBrand       [4]byte               `mp4:"size=8,string"`
	MinorVersion     uint32                `mp4:"size=32"`
	CompatibleBrands []CompatibleBrandElem `mp4:"size=32"` // to end of the box
}

type CompatibleBrandElem struct {
	CompatibleBrand [4]byte `mp4:"size=8,string"`
}

// GetType returns the BoxType
func (*Ftyp) GetType() BoxType {
	return BoxTypeFtyp()
}

/*************************** hdlr ****************************/

func BoxTypeHdlr() BoxType { return StrToBoxType("hdlr") }

func init() {
	AddBoxDef(&Hdlr{}, 0)
}

// Hdlr is ISOBMFF hdlr box type
type Hdlr struct {
	FullBox `mp4:"extend"`
	// Predefined corresponds to component_type of QuickTime.
	// pre_defined of ISO-14496 has always zero,
	// however component_type has "mhlr" or "dhlr".
	PreDefined  uint32    `mp4:"size=32"`
	HandlerType [4]byte   `mp4:"size=8,string"`
	Reserved    [3]uint32 `mp4:"size=32,const=0"`
	Name        string    `mp4:"string=c_p"`
	Padding     []byte    `mp4:"size=8,const=0"`
}

// GetType returns the BoxType
func (*Hdlr) GetType() BoxType {
	return BoxTypeHdlr()
}

func (hdlr *Hdlr) IsPString(name string, bytes []byte, remainingSize uint64) bool {
	switch name {
	case "Name":
		return remainingSize == 0 && hdlr.PreDefined != 0
	default:
		panic(fmt.Errorf("invalid field name: name=%s", name))
	}
}

/*************************** mdat ****************************/

func BoxTypeMdat() BoxType { return StrToBoxType("mdat") }

func init() {
	AddBoxDef(&Mdat{})
}

// Mdat is ISOBMFF mdat box type
type Mdat struct {
	Box
	Data []byte `mp4:"size=8"`
}

// GetType returns the BoxType
func (*Mdat) GetType() BoxType {
	return BoxTypeMdat()
}

/*************************** mdhd ****************************/

func BoxTypeMdhd() BoxType { return StrToBoxType("mdhd") }

func init() {
	AddBoxDef(&Mdhd{}, 0, 1)
}

// Mdhd is ISOBMFF mdhd box type
type Mdhd struct {
	FullBox            `mp4:"extend"`
	CreationTimeV0     uint32 `mp4:"size=32,ver=0"`
	ModificationTimeV0 uint32 `mp4:"size=32,ver=0"`
	CreationTimeV1     uint64 `mp4:"size=64,ver=1"`
	ModificationTimeV1 uint64 `mp4:"size=64,ver=1"`
	Timescale          uint32 `mp4:"size=32"`
	DurationV0         uint32 `mp4:"size=32,ver=0"`
	DurationV1         uint64 `mp4:"size=64,ver=1"`
	//
	Pad        bool    `mp4:"size=1"`
	Language   [3]byte `mp4:"size=5,iso639-2"` // ISO-639-2/T language code
	PreDefined uint16  `mp4:"size=16"`
}

// GetType returns the BoxType
func (*Mdhd) GetType() BoxType {
	return BoxTypeMdhd()
}

/*************************** mdia ****************************/

func BoxTypeMdia() BoxType { return StrToBoxType("mdia") }

func init() {
	AddBoxDef(&Mdia{})
}

// Mdia is ISOBMFF mdia box type
type Mdia struct {
	Box
}

// GetType returns the BoxType
func (*Mdia) GetType() BoxType {
	return BoxTypeMdia()
}

/*************************** mehd ****************************/

func BoxTypeMehd() BoxType { return StrToBoxType("mehd") }

func init() {
	AddBoxDef(&Mehd{}, 0, 1)
}

// Mehd is ISOBMFF mehd box type
type Mehd struct {
	FullBox            `mp4:"extend"`
	FragmentDurationV0 uint32 `mp4:"size=32,ver=0"`
	FragmentDurationV1 uint64 `mp4:"size=64,ver=1"`
}

// GetType returns the BoxType
func (*Mehd) GetType() BoxType {
	return BoxTypeMehd()
}

/*************************** meta ****************************/

func BoxTypeMeta() BoxType { return StrToBoxType("meta") }

func init() {
	AddBoxDef(&Meta{}, 0)
}

// Meta is ISOBMFF meta box type
type Meta struct {
	FullBox `mp4:"extend"`
}

// GetType returns the BoxType
func (*Meta) GetType() BoxType {
	return BoxTypeMeta()
}

func (meta *Meta) BeforeUnmarshal(r io.ReadSeeker) (n uint64, override bool, err error) {
	// for Apple Quick Time
	buf := make([]byte, 4)
	if _, err := io.ReadFull(r, buf); err != nil {
		return 0, false, err
	}
	if _, err := r.Seek(-int64(len(buf)), io.SeekCurrent); err != nil {
		return 0, false, err
	}
	if buf[0]|buf[1]|buf[2]|buf[3] != 0x00 {
		meta.Version = 0
		meta.Flags = [3]byte{0, 0, 0}
		return 0, true, nil
	}
	return 0, false, nil
}

/*************************** mfhd ****************************/

func BoxTypeMfhd() BoxType { return StrToBoxType("mfhd") }

func init() {
	AddBoxDef(&Mfhd{}, 0)
}

// Mfhd is ISOBMFF mfhd box type
type Mfhd struct {
	FullBox        `mp4:"extend"`
	SequenceNumber uint32 `mp4:"size=32"`
}

// GetType returns the BoxType
func (*Mfhd) GetType() BoxType {
	return BoxTypeMfhd()
}

/*************************** mfra ****************************/

func BoxTypeMfra() BoxType { return StrToBoxType("mfra") }

func init() {
	AddBoxDef(&Mfra{})
}

// Mfra is ISOBMFF mfra box type
type Mfra struct {
	Box
}

// GetType returns the BoxType
func (*Mfra) GetType() BoxType {
	return BoxTypeMfra()
}

/*************************** mfro ****************************/

func BoxTypeMfro() BoxType { return StrToBoxType("mfro") }

func init() {
	AddBoxDef(&Mfro{}, 0)
}

// Mfro is ISOBMFF mfro box type
type Mfro struct {
	FullBox `mp4:"extend"`
	Size    uint32 `mp4:"size=32"`
}

// GetType returns the BoxType
func (*Mfro) GetType() BoxType {
	return BoxTypeMfro()
}

/*************************** minf ****************************/

func BoxTypeMinf() BoxType { return StrToBoxType("minf") }

func init() {
	AddBoxDef(&Minf{})
}

// Minf is ISOBMFF minf box type
type Minf struct {
	Box
}

// GetType returns the BoxType
func (*Minf) GetType() BoxType {
	return BoxTypeMinf()
}

/*************************** moof ****************************/

func BoxTypeMoof() BoxType { return StrToBoxType("moof") }

func init() {
	AddBoxDef(&Moof{})
}

// Moof is ISOBMFF moof box type
type Moof struct {
	Box
}

// GetType returns the BoxType
func (*Moof) GetType() BoxType {
	return BoxTypeMoof()
}

/*************************** moov ****************************/

func BoxTypeMoov() BoxType { return StrToBoxType("moov") }

func init() {
	AddBoxDef(&Moov{})
}

// Moov is ISOBMFF moov box type
type Moov struct {
	Box
}

// GetType returns the BoxType
func (*Moov) GetType() BoxType {
	return BoxTypeMoov()
}

/*************************** mvex ****************************/

func BoxTypeMvex() BoxType { return StrToBoxType("mvex") }

func init() {
	AddBoxDef(&Mvex{})
}

// Mvex is ISOBMFF mvex box type
type Mvex struct {
	Box
}

// GetType returns the BoxType
func (*Mvex) GetType() BoxType {
	return BoxTypeMvex()
}

/*************************** mvhd ****************************/

func BoxTypeMvhd() BoxType { return StrToBoxType("mvhd") }

func init() {
	AddBoxDef(&Mvhd{}, 0, 1)
}

// Mvhd is ISOBMFF mvhd box type
type Mvhd struct {
	FullBox            `mp4:"extend"`
	CreationTimeV0     uint32    `mp4:"size=32,ver=0"`
	ModificationTimeV0 uint32    `mp4:"size=32,ver=0"`
	CreationTimeV1     uint64    `mp4:"size=64,ver=1"`
	ModificationTimeV1 uint64    `mp4:"size=64,ver=1"`
	Timescale          uint32    `mp4:"size=32"`
	DurationV0         uint32    `mp4:"size=32,ver=0"`
	DurationV1         uint64    `mp4:"size=64,ver=1"`
	Rate               int32     `mp4:"size=32"` // template=0x00010000
	Volume             int16     `mp4:"size=16"` // template=0x0100
	Reserved           int16     `mp4:"size=16,const=0"`
	Reserved2          [2]uint32 `mp4:"size=32,const=0"`
	Matrix             [9]int32  `mp4:"size=32,hex"` // template={ 0x00010000,0,0,0,0x00010000,0,0,0,0x40000000 }
	PreDefined         [6]int32  `mp4:"size=32"`
	NextTrackID        uint32    `mp4:"size=32"`
}

// GetType returns the BoxType
func (*Mvhd) GetType() BoxType {
	return BoxTypeMvhd()
}

/*************************** pssh ****************************/

func BoxTypePssh() BoxType { return StrToBoxType("pssh") }

func init() {
	AddBoxDef(&Pssh{}, 0, 1)
}

// Pssh is ISOBMFF pssh box type
type Pssh struct {
	FullBox  `mp4:"extend"`
	SystemID [16]byte  `mp4:"size=8"`
	KIDCount uint32    `mp4:"size=32,nver=0"`
	KIDs     []PsshKID `mp4:"nver=0,len=dynamic,size=128"`
	DataSize int32     `mp4:"size=32"`
	Data     []byte    `mp4:"size=8,len=dynamic"`
}

type PsshKID struct {
	KID [16]byte `mp4:"size=8"`
}

// GetFieldLength returns length of dynamic field
func (pssh *Pssh) GetFieldLength(name string) uint {
	switch name {
	case "KIDs":
		return uint(pssh.KIDCount)
	case "Data":
		return uint(pssh.DataSize)
	}
	panic(fmt.Errorf("invalid name of dynamic-length field: boxType=pssh fieldName=%s", name))
}

// StringifyField returns field value as string
func (pssh *Pssh) StringifyField(name string, indent string, depth int) (string, bool) {
	switch name {
	case "SystemID":
		buf := bytes.NewBuffer(nil)
		buf.WriteString("\"")
		for _, b := range pssh.SystemID {
			buf.WriteString(fmt.Sprintf("%02x", b))
		}
		buf.WriteString("\"")
		return buf.String(), true

	case "KIDs":
		buf := bytes.NewBuffer(nil)
		buf.WriteString("[")
		for i, e := range pssh.KIDs {
			if i != 0 {
				buf.WriteString(" \"")
			} else {
				buf.WriteString("\"")
			}
			for _, b := range e.KID {
				buf.WriteString(fmt.Sprintf("%02x", b))
			}
			buf.WriteString("\"")
		}
		buf.WriteString("]")
		return buf.String(), true

	default:
		return "", false
	}
}

// GetType returns the BoxType
func (*Pssh) GetType() BoxType {
	return BoxTypePssh()
}

/*********************** SampleEntry *************************/

func init() {
	AddAnyTypeBoxDef(&VisualSampleEntry{}, StrToBoxType("avc1"))
	AddAnyTypeBoxDef(&VisualSampleEntry{}, StrToBoxType("encv"))
	AddAnyTypeBoxDef(&AudioSampleEntry{}, StrToBoxType("mp4a"))
	AddAnyTypeBoxDef(&AudioSampleEntry{}, StrToBoxType("enca"))
	AddAnyTypeBoxDef(&AVCDecoderConfiguration{}, StrToBoxType("avcC"))
	AddAnyTypeBoxDef(&PixelAspectRatioBox{}, StrToBoxType("pasp"))
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
		if vse.Compressorname[0] <= 31 {
			return `"` + string(vse.Compressorname[1:vse.Compressorname[0]+1]) + `"`, true
		}
		return "", false
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

/*************************** sbgp ****************************/

func BoxTypeSbgp() BoxType { return StrToBoxType("sbgp") }

func init() {
	AddBoxDef(&Sbgp{}, 0, 1)
}

type Sbgp struct {
	FullBox                 `mp4:"extend"`
	GroupingType            uint32      `mp4:"size=32"`
	grouping_type_parameter uint32      `mp4:"size=32,ver=1"`
	EntryCount              uint32      `mp4:"size=32"`
	Entries                 []SbgpEntry `mp4:"len=dynamic,size=64"`
}

type SbgpEntry struct {
	SampleCount           uint32 `mp4:"size=32"`
	GroupDescriptionIndex uint32 `mp4:"size=32"`
}

func (sbgp *Sbgp) GetFieldLength(name string) uint {
	switch name {
	case "Entries":
		return uint(sbgp.EntryCount)
	}
	panic(fmt.Errorf("invalid name of dynamic-length field: boxType=sbgp fieldName=%s", name))
}

func (*Sbgp) GetType() BoxType {
	return BoxTypeSbgp()
}

/*************************** sgpd ****************************/

func BoxTypeSgpd() BoxType { return StrToBoxType("sgpd") }

func init() {
	AddBoxDef(&Sgpd{}, 0, 1, 2)
}

type Sgpd struct {
	FullBox                       `mp4:"extend"`
	GroupingType                  [4]byte `mp4:"size=8,string"`
	DefaultLength                 uint32  `mp4:"size=32,ver=1"`
	DefaultSampleDescriptionIndex uint32  `mp4:"size=32,ver=2"`
	EntryCount                    uint32  `mp4:"size=32"`
	RollDistances                 []int16 `mp4:"size=16,opt=dynamic"`
	//AlternativeStartupEntries     []AlternativeStartupEntry `mp4:"size=dynamic,opt=dynamic"`
	VisualRandomAccessEntries []VisualRandomAccessEntry `mp4:"size=dynamic,opt=dynamic"`
	TemporalLevelEntries      []TemporalLevelEntry      `mp4:"size=dynamic,opt=dynamic"`
	Unsupported               []byte                    `mp4:"size=8"`
}

/* go-mp4 has never supported nested dynamic field.
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
	NumLeadingSamplesKnown bool  `mp4:"size=1"`
	NumLeadingSamples      uint8 `mp4:"size=7"`
}

type TemporalLevelEntry struct {
	LevelUndependentlyUecodable bool  `mp4:"size=1"`
	Reserved                    uint8 `mp4:"size=7,const=0"`
}

func (sgpd *Sgpd) GetFieldSize(name string) uint {
	switch name {
	case "RollDistances",
		//"AlternativeStartupEntries",
		"VisualRandomAccessEntries",
		"TemporalLevelEntries":
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
			case "AlternativeStartupEntries":
				return sgpd.Version == 1 &&
					sgpd.GroupingType == [4]byte{'a', 'l', 's', 't'} &&
					sgpd.DefaultLength != 0
		*/
	case "VisualRandomAccessEntries":
		return sgpd.Version == 1 &&
			sgpd.GroupingType == [4]byte{'r', 'a', 'p', ' '} &&
			sgpd.DefaultLength == 1
	case "TemporalLevelEntries":
		return sgpd.Version == 1 &&
			sgpd.GroupingType == [4]byte{'t', 'e', 'l', 'e'} &&
			sgpd.DefaultLength == 1
	default:
		return false
	}
}

func (*Sgpd) GetType() BoxType {
	return BoxTypeSgpd()
}

/*************************** smhd ****************************/

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

/*************************** stbl ****************************/

func BoxTypeStbl() BoxType { return StrToBoxType("stbl") }

func init() {
	AddBoxDef(&Stbl{})
}

// Stbl is ISOBMFF stbl box type
type Stbl struct {
	Box
}

// GetType returns the BoxType
func (*Stbl) GetType() BoxType {
	return BoxTypeStbl()
}

/*************************** stco ****************************/

func BoxTypeStco() BoxType { return StrToBoxType("stco") }

func init() {
	AddBoxDef(&Stco{}, 0)
}

// Stco is ISOBMFF stco box type
type Stco struct {
	FullBox     `mp4:"extend"`
	EntryCount  uint32   `mp4:"size=32"`
	ChunkOffset []uint32 `mp4:"size=32,len=dynamic"`
}

// GetType returns the BoxType
func (*Stco) GetType() BoxType {
	return BoxTypeStco()
}

// GetFieldLength returns length of dynamic field
func (stco *Stco) GetFieldLength(name string) uint {
	switch name {
	case "ChunkOffset":
		return uint(stco.EntryCount)
	}
	panic(fmt.Errorf("invalid name of dynamic-length field: boxType=stco fieldName=%s", name))
}

/*************************** stsc ****************************/

func BoxTypeStsc() BoxType { return StrToBoxType("stsc") }

func init() {
	AddBoxDef(&Stsc{}, 0)
}

// Stsc is ISOBMFF stsc box type
type Stsc struct {
	FullBox    `mp4:"extend"`
	EntryCount uint32      `mp4:"size=32"`
	Entries    []StscEntry `mp4:"len=dynamic,size=96"`
}

type StscEntry struct {
	FirstChunk             uint32 `mp4:"size=32"`
	SamplesPerChunk        uint32 `mp4:"size=32"`
	SampleDescriptionIndex uint32 `mp4:"size=32"`
}

// GetType returns the BoxType
func (*Stsc) GetType() BoxType {
	return BoxTypeStsc()
}

// GetFieldLength returns length of dynamic field
func (stsc *Stsc) GetFieldLength(name string) uint {
	switch name {
	case "Entries":
		return uint(stsc.EntryCount)
	}
	panic(fmt.Errorf("invalid name of dynamic-length field: boxType=stsc fieldName=%s", name))
}

/*************************** stsd ****************************/

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

/*************************** stss ****************************/

func BoxTypeStss() BoxType { return StrToBoxType("stss") }

func init() {
	AddBoxDef(&Stss{}, 0)
}

type Stss struct {
	FullBox      `mp4:"extend"`
	EntryCount   uint32   `mp4:"size=32"`
	SampleNumber []uint32 `mp4:"len=dynamic,size=32"`
}

// GetType returns the BoxType
func (*Stss) GetType() BoxType {
	return BoxTypeStss()
}

// GetFieldLength returns length of dynamic field
func (stss *Stss) GetFieldLength(name string) uint {
	switch name {
	case "SampleNumber":
		return uint(stss.EntryCount)
	}
	panic(fmt.Errorf("invalid name of dynamic-length field: boxType=stss fieldName=%s", name))
}

/*************************** stsz ****************************/

func BoxTypeStsz() BoxType { return StrToBoxType("stsz") }

func init() {
	AddBoxDef(&Stsz{}, 0)
}

// Stsz is ISOBMFF stsz box type
type Stsz struct {
	FullBox     `mp4:"extend"`
	SampleSize  uint32   `mp4:"size=32"`
	SampleCount uint32   `mp4:"size=32"`
	EntrySize   []uint32 `mp4:"size=32,len=dynamic"`
}

// GetType returns the BoxType
func (*Stsz) GetType() BoxType {
	return BoxTypeStsz()
}

// GetFieldLength returns length of dynamic field
func (stsz *Stsz) GetFieldLength(name string) uint {
	switch name {
	case "EntrySize":
		if stsz.SampleSize == 0 {
			return uint(stsz.SampleCount)
		} else {
			return 0
		}
	}
	panic(fmt.Errorf("invalid name of dynamic-length field: boxType=stsz fieldName=%s", name))
}

/*************************** stts ****************************/

func BoxTypeStts() BoxType { return StrToBoxType("stts") }

func init() {
	AddBoxDef(&Stts{}, 0)
}

// Stts is ISOBMFF stts box type
type Stts struct {
	FullBox    `mp4:"extend"`
	EntryCount uint32      `mp4:"size=32"`
	Entries    []SttsEntry `mp4:"len=dynamic,size=64"`
}

type SttsEntry struct {
	SampleCount uint32 `mp4:"size=32"`
	SampleDelta uint32 `mp4:"size=32"`
}

// GetType returns the BoxType
func (*Stts) GetType() BoxType {
	return BoxTypeStts()
}

// GetFieldLength returns length of dynamic field
func (stts *Stts) GetFieldLength(name string) uint {
	switch name {
	case "Entries":
		return uint(stts.EntryCount)
	}
	panic(fmt.Errorf("invalid name of dynamic-length field: boxType=stts fieldName=%s", name))
}

/*************************** tfdt ****************************/

func BoxTypeTfdt() BoxType { return StrToBoxType("tfdt") }

func init() {
	AddBoxDef(&Tfdt{}, 0, 1)
}

// Tfdt is ISOBMFF tfdt box type
type Tfdt struct {
	FullBox               `mp4:"extend"`
	BaseMediaDecodeTimeV0 uint32 `mp4:"size=32,ver=0"`
	BaseMediaDecodeTimeV1 uint64 `mp4:"size=64,ver=1"`
}

// GetType returns the BoxType
func (*Tfdt) GetType() BoxType {
	return BoxTypeTfdt()
}

/*************************** tfhd ****************************/

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

/*************************** tfra ****************************/

func BoxTypeTfra() BoxType { return StrToBoxType("tfra") }

func init() {
	AddBoxDef(&Tfra{}, 0, 1)
}

// Tfra is ISOBMFF tfra box type
type Tfra struct {
	FullBox               `mp4:"extend"`
	TrackID               uint32      `mp4:"size=32"`
	Reserved              uint32      `mp4:"size=26,const=0"`
	LengthSizeOfTrafNum   byte        `mp4:"size=2"`
	LengthSizeOfTrunNum   byte        `mp4:"size=2"`
	LengthSizeOfSampleNum byte        `mp4:"size=2"`
	NumberOfEntry         uint32      `mp4:"size=32"`
	Entries               []TfraEntry `mp4:"len=dynamic,size=dynamic"`
}

type TfraEntry struct {
	TimeV0       uint32 `mp4:"size=32,ver=0"`
	MoofOffsetV0 uint32 `mp4:"size=32,ver=0"`
	TimeV1       uint64 `mp4:"size=64,ver=1"`
	MoofOffsetV1 uint64 `mp4:"size=64,ver=1"`
	TrafNumber   uint32 `mp4:"size=dynamic"`
	TrunNumber   uint32 `mp4:"size=dynamic"`
	SampleNumber uint32 `mp4:"size=dynamic"`
}

// GetType returns the BoxType
func (*Tfra) GetType() BoxType {
	return BoxTypeTfra()
}

// GetFieldSize returns size of dynamic field
func (tfra *Tfra) GetFieldSize(name string) uint {
	switch name {
	case "TrafNumber":
		return (uint(tfra.LengthSizeOfTrafNum) + 1) * 8
	case "TrunNumber":
		return (uint(tfra.LengthSizeOfTrunNum) + 1) * 8
	case "SampleNumber":
		return (uint(tfra.LengthSizeOfSampleNum) + 1) * 8
	case "Entries":
		switch tfra.GetVersion() {
		case 0:
			return 0 +
				/* TimeV0       */ 32 +
				/* MoofOffsetV0 */ 32 +
				/* TrafNumber   */ (uint(tfra.LengthSizeOfTrafNum)+1)*8 +
				/* TrunNumber   */ (uint(tfra.LengthSizeOfTrunNum)+1)*8 +
				/* SampleNumber */ (uint(tfra.LengthSizeOfSampleNum)+1)*8
		case 1:
			return 0 +
				/* TimeV1       */ 64 +
				/* MoofOffsetV1 */ 64 +
				/* TrafNumber   */ (uint(tfra.LengthSizeOfTrafNum)+1)*8 +
				/* TrunNumber   */ (uint(tfra.LengthSizeOfTrunNum)+1)*8 +
				/* SampleNumber */ (uint(tfra.LengthSizeOfSampleNum)+1)*8
		}
	}
	panic(fmt.Errorf("invalid name of dynamic-size field: boxType=tfra fieldName=%s", name))
}

// GetFieldLength returns length of dynamic field
func (tfra *Tfra) GetFieldLength(name string) uint {
	switch name {
	case "Entries":
		return uint(tfra.NumberOfEntry)
	}
	panic(fmt.Errorf("invalid name of dynamic-length field: boxType=tfra fieldName=%s", name))
}

/*************************** tkhd ****************************/

func BoxTypeTkhd() BoxType { return StrToBoxType("tkhd") }

func init() {
	AddBoxDef(&Tkhd{}, 0, 1)
}

// Tkhd is ISOBMFF tkhd box type
type Tkhd struct {
	FullBox `mp4:"extend"`
	// Version 0
	CreationTimeV0     uint32 `mp4:"size=32,ver=0"`
	ModificationTimeV0 uint32 `mp4:"size=32,ver=0"`
	TrackIDV0          uint32 `mp4:"size=32,ver=0"`
	ReservedV0         uint32 `mp4:"size=32,ver=0,const=0"`
	DurationV0         uint32 `mp4:"size=32,ver=0"`
	// Version 1
	CreationTimeV1     uint64 `mp4:"size=64,ver=1"`
	ModificationTimeV1 uint64 `mp4:"size=64,ver=1"`
	TrackIDV1          uint32 `mp4:"size=32,ver=1"`
	ReservedV1         uint32 `mp4:"size=32,ver=1,const=0"`
	DurationV1         uint64 `mp4:"size=64,ver=1"`
	//
	Reserved       [2]uint32 `mp4:"size=32,const=0"`
	Layer          int16     `mp4:"size=16"` // template=0
	AlternateGroup int16     `mp4:"size=16"` // template=0
	Volume         int16     `mp4:"size=16"` // template={if track_is_audio 0x0100 else 0}
	Reserved2      uint16    `mp4:"size=16,const=0"`
	Matrix         [9]int32  `mp4:"size=32,hex"` // template={ 0x00010000,0,0,0,0x00010000,0,0,0,0x40000000 };
	Width          uint32    `mp4:"size=32"`
	Height         uint32    `mp4:"size=32"`
}

// GetType returns the BoxType
func (*Tkhd) GetType() BoxType {
	return BoxTypeTkhd()
}

/*************************** traf ****************************/

func BoxTypeTraf() BoxType { return StrToBoxType("traf") }

func init() {
	AddBoxDef(&Traf{})
}

// Traf is ISOBMFF traf box type
type Traf struct {
	Box
}

// GetType returns the BoxType
func (*Traf) GetType() BoxType {
	return BoxTypeTraf()
}

/*************************** trak ****************************/

func BoxTypeTrak() BoxType { return StrToBoxType("trak") }

func init() {
	AddBoxDef(&Trak{})
}

// Trak is ISOBMFF trak box type
type Trak struct {
	Box
}

// GetType returns the BoxType
func (*Trak) GetType() BoxType {
	return BoxTypeTrak()
}

/*************************** trex ****************************/

func BoxTypeTrex() BoxType { return StrToBoxType("trex") }

func init() {
	AddBoxDef(&Trex{}, 0)
}

// Trex is ISOBMFF trex box type
type Trex struct {
	FullBox                       `mp4:"extend"`
	TrackID                       uint32 `mp4:"size=32"`
	DefaultSampleDescriptionIndex uint32 `mp4:"size=32"`
	DefaultSampleDuration         uint32 `mp4:"size=32"`
	DefaultSampleSize             uint32 `mp4:"size=32"`
	DefaultSampleFlags            uint32 `mp4:"size=32,hex"`
}

// GetType returns the BoxType
func (*Trex) GetType() BoxType {
	return BoxTypeTrex()
}

/*************************** trun ****************************/

func BoxTypeTrun() BoxType { return StrToBoxType("trun") }

func init() {
	AddBoxDef(&Trun{}, 0, 1)
}

// Trun is ISOBMFF trun box type
type Trun struct {
	FullBox     `mp4:"extend"`
	SampleCount uint32 `mp4:"size=32"`

	// optional fields
	DataOffset       int32       `mp4:"size=32,opt=0x000001"`
	FirstSampleFlags uint32      `mp4:"size=32,opt=0x000004,hex"`
	Entries          []TrunEntry `mp4:"len=dynamic,size=dynamic"`
}

type TrunEntry struct {
	SampleDuration                uint32 `mp4:"size=32,opt=0x000100"`
	SampleSize                    uint32 `mp4:"size=32,opt=0x000200"`
	SampleFlags                   uint32 `mp4:"size=32,opt=0x000400,hex"`
	SampleCompositionTimeOffsetV0 uint32 `mp4:"size=32,opt=0x000800,ver=0"`
	SampleCompositionTimeOffsetV1 int32  `mp4:"size=32,opt=0x000800,nver=0"`
}

// GetType returns the BoxType
func (*Trun) GetType() BoxType {
	return BoxTypeTrun()
}

// GetFieldSize returns size of dynamic field
func (trun *Trun) GetFieldSize(name string) uint {
	switch name {
	case "Entries":
		var size uint
		flags := trun.GetFlags()
		if flags&0x100 != 0 {
			size += 32 // SampleDuration
		}
		if flags&0x200 != 0 {
			size += 32 // SampleSize
		}
		if flags&0x400 != 0 {
			size += 32 // SampleFlags
		}
		if flags&0x800 != 0 {
			size += 32 // SampleCompositionTimeOffsetV0 or V1
		}
		return size
	}
	panic(fmt.Errorf("invalid name of dynamic-size field: boxType=trun fieldName=%s", name))
}

// GetFieldLength returns length of dynamic field
func (trun *Trun) GetFieldLength(name string) uint {
	switch name {
	case "Entries":
		return uint(trun.SampleCount)
	}
	panic(fmt.Errorf("invalid name of dynamic-length field: boxType=trun fieldName=%s", name))
}

/*************************** udta ****************************/

func BoxTypeUdta() BoxType { return StrToBoxType("udta") }

func init() {
	AddBoxDef(&Udta{})
}

// Udta is ISOBMFF udta box type
type Udta struct {
	Box
}

// GetType returns the BoxType
func (*Udta) GetType() BoxType {
	return BoxTypeUdta()
}

/*************************** vmhd ****************************/

func BoxTypeVmhd() BoxType { return StrToBoxType("vmhd") }

func init() {
	AddBoxDef(&Vmhd{}, 0)
}

// Vmhd is ISOBMFF vmhd box type
type Vmhd struct {
	FullBox      `mp4:"extend"`
	Graphicsmode uint16    `mp4:"size=16"` // template=0
	Opcolor      [3]uint16 `mp4:"size=16"` // template={0, 0, 0}
}

// GetType returns the BoxType
func (*Vmhd) GetType() BoxType {
	return BoxTypeVmhd()
}
