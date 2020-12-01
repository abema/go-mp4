package mp4

import (
	"reflect"
	"strconv"
	"strings"
)

type StringType int

const (
	StringType_C StringType = iota
	StringType_C_P
)

type fieldConfig struct {
	name       string
	cfo        ICustomFieldObject
	size       uint
	length     uint
	cnst       string
	strType    StringType
	varint     bool
	version    uint8
	nVersion   uint8
	optDynamic bool
	optFlag    uint32
	nOptFlag   uint32
	extend     bool
	dec        bool
	hex        bool
	str        bool
	iso639_2   bool
	uuid       bool
	hidden     bool
}

func readFieldConfig(box IImmutableBox, parent reflect.Value, fieldName string, tag fieldTag, ctx Context) (config fieldConfig, err error) {
	config.name = fieldName
	cfo, ok := parent.Addr().Interface().(ICustomFieldObject)
	if ok {
		config.cfo = cfo
	} else {
		config.cfo = box
	}

	if val, contained := tag["size"]; contained {
		if val == "dynamic" {
			config.size = config.cfo.GetFieldSize(fieldName, ctx)
		} else {
			var size uint64
			size, err = strconv.ParseUint(val, 10, 32)
			if err != nil {
				return
			}
			config.size = uint(size)
		}
	}

	config.length = LengthUnlimited
	if val, contained := tag["len"]; contained {
		if val == "dynamic" {
			config.length = config.cfo.GetFieldLength(fieldName, ctx)
		} else {
			var l uint64
			l, err = strconv.ParseUint(val, 10, 32)
			if err != nil {
				return
			}
			config.length = uint(l)
		}
	}

	if _, contained := tag["varint"]; contained {
		config.varint = true
	}

	config.version = anyVersion
	if val, contained := tag["ver"]; contained {
		var ver int
		ver, err = strconv.Atoi(val)
		if err != nil {
			return
		}
		config.version = uint8(ver)
	}

	config.nVersion = anyVersion
	if val, contained := tag["nver"]; contained {
		var ver int
		ver, err = strconv.Atoi(val)
		if err != nil {
			return
		}
		config.nVersion = uint8(ver)
	}

	if val, contained := tag["opt"]; contained {
		if val == "dynamic" {
			config.optDynamic = true
		} else {
			var opt uint64
			if strings.HasPrefix(val, "0x") {
				opt, err = strconv.ParseUint(val[2:], 16, 32)
			} else {
				opt, err = strconv.ParseUint(val, 10, 32)
			}
			if err != nil {
				return
			}
			config.optFlag = uint32(opt)
		}
	}

	if val, contained := tag["nopt"]; contained {
		var nopt uint64
		if strings.HasPrefix(val, "0x") {
			nopt, err = strconv.ParseUint(val[2:], 16, 32)
		} else {
			nopt, err = strconv.ParseUint(val, 10, 32)
		}
		if err != nil {
			return
		}
		config.nOptFlag = uint32(nopt)
	}

	if val, contained := tag["const"]; contained {
		config.cnst = val
	}

	if _, contained := tag["extend"]; contained {
		config.extend = true
	}

	if _, contained := tag["dec"]; contained {
		config.dec = true
	}

	if _, contained := tag["hex"]; contained {
		config.hex = true
	}

	if val, contained := tag["string"]; contained {
		config.str = true
		if val == "c_p" {
			config.strType = StringType_C_P
		}
	}

	if _, contained := tag["iso639-2"]; contained {
		config.iso639_2 = true
	}

	if _, contained := tag["uuid"]; contained {
		config.uuid = true
	}

	if _, contained := tag["hidden"]; contained {
		config.hidden = true
	}

	return
}

func isTargetField(box IImmutableBox, config fieldConfig, ctx Context) bool {
	if box.GetVersion() != anyVersion {
		if config.version != anyVersion && box.GetVersion() != config.version {
			return false
		}

		if config.nVersion != anyVersion && box.GetVersion() == config.nVersion {
			return false
		}
	}

	if config.optFlag != 0 && box.GetFlags()&config.optFlag == 0 {
		return false
	}

	if config.nOptFlag != 0 && box.GetFlags()&config.nOptFlag != 0 {
		return false
	}

	if config.optDynamic && !config.cfo.IsOptFieldEnabled(config.name, ctx) {
		return false
	}

	return true
}

type fieldTag map[string]string

func parseFieldTag(str string) fieldTag {
	tag := make(map[string]string, 8)

	list := strings.Split(str, ",")
	for _, e := range list {
		kv := strings.SplitN(e, "=", 2)
		if len(kv) == 2 {
			tag[strings.Trim(kv[0], " ")] = strings.Trim(kv[1], " ")
		} else {
			tag[strings.Trim(kv[0], " ")] = ""
		}
	}

	return tag
}
