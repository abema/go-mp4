package mp4

import "io"

func BoxTypeHdlr() BoxType { return StrToBoxType("hdlr") }

func init() {
	AddBoxDef(&Hdlr{}, 0)
}

// Hdlr is ISOBMFF hdlr box type
type Hdlr struct {
	FullBox     `mp4:"extend"`
	PreDefined  uint32    `mp4:"size=32"`
	HandlerType [4]byte   `mp4:"size=8,string"`
	Reserved    [3]uint32 `mp4:"size=32,const=0"`
	Name        string    `mp4:"string"`
}

// GetType returns the BoxType
func (*Hdlr) GetType() BoxType {
	return BoxTypeHdlr()
}

const (
	hdlrSizeWithoutName = 24 // FullBox header(4bytes) + PreDefined(4bytes) + HandlerType(4bytes) + Reserved(12bytes)
)

// handle a special case: the QuickTime files have a pascal
// string here, but ISO MP4 files have a C string.
// we try to detect a pascal encoding and correct it.
func (hdlr *Hdlr) unmarshalHandlerName(u *unmarshaller) error {
	nameSize := u.size - hdlrSizeWithoutName
	if nameSize <= 0 {
		return nil
	}
	if _, err := u.reader.Seek(-int64(nameSize), io.SeekCurrent); err != nil {
		return err
	}
	u.rbytes = hdlrSizeWithoutName
	nameb := make([]byte, nameSize)
	if _, err := io.ReadFull(u.reader, nameb); err != nil {
		return err
	}
	u.rbytes += nameSize
	if nameb[0] != 0x00 && nameb[0] == byte(nameSize-1) {
		hdlr.Name = string(nameb[1:])
	}
	return nil
}
