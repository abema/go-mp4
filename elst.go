package mp4

import "fmt"

func BoxTypeElst() BoxType { return StrToBoxType("elst") }

func init() {
	AddBoxDef(&Elst{}, 0, 1)
}

// Elst is ISOBMFF elst box type
type Elst struct {
	FullBox    `mp4:"extend"`
	EntryCount uint32 `mp4:"size=32"`
	Entries    []struct {
		SegmentDurationV0 uint32 `mp4:"size=32,ver=0"`
		MediaTimeV0       int32  `mp4:"size=32,ver=0"`
		SegmentDurationV1 uint64 `mp4:"size=64,ver=1"`
		MediaTimeV1       int64  `mp4:"size=64,ver=1"`
		MediaRateInteger  int16  `mp4:"size=16"`
		MediaRateFraction int16  `mp4:"size=16,const=0"`
	} `mp4:"len=dynamic,size=dynamic"`
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
