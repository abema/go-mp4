package mp4

import "fmt"

// https://www.etsi.org/deliver/etsi_ts/102300_102399/102366/01.04.01_60/ts_102366v010401p.pdf

/*************************** ac-3 ****************************/

func BoxTypeAC3() BoxType { return StrToBoxType("ac-3") }

func init() {
	AddAnyTypeBoxDef(&AudioSampleEntry{}, BoxTypeAC3())
}

/*************************** dac3 ****************************/

func BoxTypeDAC3() BoxType { return StrToBoxType("dac3") }

func init() {
	AddBoxDef(&Dac3{})
}

type Dac3 struct {
	Box
	Fscod       uint8 `mp4:"0,size=2"`
	Bsid        uint8 `mp4:"1,size=5"`
	Bsmod       uint8 `mp4:"2,size=3"`
	Acmod       uint8 `mp4:"3,size=3"`
	LfeOn       uint8 `mp4:"4,size=1"`
	BitRateCode uint8 `mp4:"5,size=5"`
	Reserved    uint8 `mp4:"6,size=5,const=0"`
}

func (Dac3) GetType() BoxType {
	return BoxTypeDAC3()
}

/*************************** ec-3 ****************************/

func BoxTypeEC3() BoxType { return StrToBoxType("ec-3") }

func init() {
	AddAnyTypeBoxDef(&AudioSampleEntry{}, BoxTypeEC3())
}

/*************************** dec3 ****************************/

func BoxTypeDec3() BoxType { return StrToBoxType("dec3") }

func init() {
	AddBoxDef(&Dec3{})
}

type Dec3 struct {
	Box
	DataRate  uint16       `mp4:"0,size=13"`
	NumIndSub uint8        `mp4:"1,size=3"`
	IndSubs   []Dec3IndSub `mp4:"2,len=dynamic"`
}

type Dec3IndSub struct {
	BaseCustomFieldObject
	Fscod     uint8  `mp4:"2,size=2"`
	Bsid      uint8  `mp4:"3,size=5"`
	Reserved1 uint8  `mp4:"4,size=1,const=0"`
	Asvc      uint8  `mp4:"5,size=1"`
	Bsmod     uint8  `mp4:"6,size=3"`
	Acmod     uint8  `mp4:"7,size=3"`
	LfeOn     uint8  `mp4:"8,size=1"`
	Reserved2 uint8  `mp4:"9,size=3,const=0"`
	NumDepSub uint8  `mp4:"10,size=4"`
	ChanLoc   uint16 `mp4:"11,size=9,opt=dynamic"`
	Reserved3 uint8  `mp4:"12,size=1,opt=dynamic,const=0"`
}

func (*Dec3) GetType() BoxType {
	return BoxTypeDec3()
}

func (dec3 *Dec3) GetFieldLength(name string, ctx Context) uint {
	switch name {
	case "IndSubs":
		return uint(dec3.NumIndSub + 1)
	}
	panic(fmt.Errorf("invalid name of dynamic-length field: boxType=dec3 fieldName=%s", name))
}

func (sub *Dec3IndSub) IsOptFieldEnabled(name string, ctx Context) bool {
	switch name {
	case "ChanLoc":
		return sub.NumDepSub > 0
	case "Reserved3":
		return sub.NumDepSub == 0
	default:
		return false
	}
}
