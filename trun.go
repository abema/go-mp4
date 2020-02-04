package mp4

import "fmt"

func BoxTypeTrun() BoxType { return StrToBoxType("trun") }

func init() {
	AddBoxDef(&Trun{}, 0, 1)
}

// Trun is ISOBMFF trun box type
type Trun struct {
	FullBox     `mp4:"extend"`
	SampleCount uint32 `mp4:"size=32"`

	// optional fields
	DataOffset       int32  `mp4:"size=32,opt=0x000001"`
	FirstSampleFlags uint32 `mp4:"size=32,opt=0x000004,hex"`
	Entries          []struct {
		SampleDuration                uint32 `mp4:"size=32,opt=0x000100"`
		SampleSize                    uint32 `mp4:"size=32,opt=0x000200"`
		SampleFlags                   uint32 `mp4:"size=32,opt=0x000400,hex"`
		SampleCompositionTimeOffsetV0 uint32 `mp4:"size=32,opt=0x000800,ver=0"`
		SampleCompositionTimeOffsetV1 int32  `mp4:"size=32,opt=0x000800,nver=0"`
	} `mp4:"len=dynamic,size=dynamic"`
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
