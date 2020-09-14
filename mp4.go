package mp4

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var ErrBoxInfoNotFound = errors.New("box info not found")

// BoxType is mpeg box type
type BoxType [4]byte

func StrToBoxType(code string) BoxType {
	if len(code) != 4 {
		panic(fmt.Errorf("invalid box type id length: [%s]", code))
	}
	return BoxType{code[0], code[1], code[2], code[3]}
}

func (boxType BoxType) String() string {
	if isPrintable(boxType[0]) && isPrintable(boxType[1]) && isPrintable(boxType[2]) && isPrintable(boxType[3]) {
		s := string([]byte{boxType[0], boxType[1], boxType[2], boxType[3]})
		s = strings.ReplaceAll(s, string([]byte{0xa9}), "(c)")
		return s
	}
	return fmt.Sprintf("0x%02x%02x%02x%02x", boxType[0], boxType[1], boxType[2], boxType[3])
}

func isASCII(c byte) bool {
	return c >= 0x20 && c <= 0x7e
}

func isPrintable(c byte) bool {
	return isASCII(c) || c == 0xa9
}

func (lhs BoxType) MatchWith(rhs BoxType) bool {
	if lhs == boxTypeAny || rhs == boxTypeAny {
		return true
	}
	return lhs == rhs
}

var boxTypeAny = BoxType{0x00, 0x00, 0x00, 0x00}

func BoxTypeAny() BoxType {
	return boxTypeAny
}

type boxDef struct {
	dataType reflect.Type
	versions []uint8
}

var boxMap = make(map[BoxType]boxDef, 64)

func AddBoxDef(payload IBox, versions ...uint8) {
	boxMap[payload.GetType()] = boxDef{
		dataType: reflect.TypeOf(payload).Elem(),
		versions: versions,
	}
}

func AddAnyTypeBoxDef(payload IAnyType, boxType BoxType, versions ...uint8) {
	boxMap[boxType] = boxDef{
		dataType: reflect.TypeOf(payload).Elem(),
		versions: versions,
	}
}

func (boxType BoxType) IsSupported() bool {
	_, ok := boxMap[boxType]
	return ok
}

func (boxType BoxType) New() (IBox, error) {
	boxDef, ok := boxMap[boxType]
	if !ok {
		return nil, ErrBoxInfoNotFound
	}

	box, ok := reflect.New(boxDef.dataType).Interface().(IBox)
	if !ok {
		return nil, fmt.Errorf("box type not implements IBox interface: %s", boxType.String())
	}

	anyTypeBox, ok := box.(IAnyType)
	if ok {
		anyTypeBox.SetType(boxType)
	}

	return box, nil
}

func (boxType BoxType) GetSupportedVersions() ([]uint8, error) {
	boxDef, ok := boxMap[boxType]
	if !ok {
		return nil, ErrBoxInfoNotFound
	}
	return boxDef.versions, nil
}

func (boxType BoxType) IsSupportedVersion(ver uint8) bool {
	boxDef, ok := boxMap[boxType]
	if !ok {
		return false
	}
	if len(boxDef.versions) == 0 {
		return true
	}
	for _, sver := range boxDef.versions {
		if ver == sver {
			return true
		}
	}
	return false
}
