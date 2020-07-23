package mp4

import "io"

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
