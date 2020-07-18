package mp4

import (
	"bytes"
	"fmt"
)

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
