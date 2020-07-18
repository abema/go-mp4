package mp4

import "fmt"

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
