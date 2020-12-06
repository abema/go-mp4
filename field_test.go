package mp4

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildField(t *testing.T) {
	box := &struct {
		mockBox
		FullBox     `mp4:"0,extend"`
		Int32       int32    `mp4:"1,size=32"`
		Int17       int32    `mp4:"2,size=17"`
		Uint15      uint16   `mp4:"3,size=15"`
		Const       byte     `mp4:"4,size=8,const=0"`
		String      []byte   `mp4:"5,size=8,string"`
		PString     []byte   `mp4:"6,size=8,string=c_p"`
		Dec         byte     `mp4:"7,size=8,dec"`
		Hex         byte     `mp4:"8,size=8,hex"`
		ISO639_2    []byte   `mp4:"9,size=8,iso639-2"`
		UUID        [16]byte `mp4:"10,size=8,uuid"`
		Hidden      byte     `mp4:"11,size=8,hidden"`
		Opt         byte     `mp4:"12,size=8,opt=0x000010"`
		NOpt        byte     `mp4:"13,size=8,nopt=0x000010"`
		DynOpt      byte     `mp4:"14,size=8,opt=dynamic"`
		Varint      uint64   `mp4:"15,varint"`
		DynSize     uint64   `mp4:"16,size=dynamic"`
		FixedLen    []byte   `mp4:"17,size=8,len=5"`
		DynLen      []byte   `mp4:"18,size=8,len=dynamic"`
		Ver         byte     `mp4:"19,size=8,ver=1"`
		NVer        byte     `mp4:"20,size=8,nver=1"`
		NotSorted22 byte     `mp4:"22,size=8"`
		NotSorted23 byte     `mp4:"23,size=8"`
		NotSorted21 byte     `mp4:"21,size=8"`
	}{}

	fs := buildFields(box)
	require.Len(t, fs, 24)
	assert.Equal(t, &field{
		name:     "FullBox",
		order:    0,
		flags:    fieldExtend,
		version:  anyVersion,
		nVersion: anyVersion,
		length:   LengthUnlimited,
		children: []*field{
			{
				name:     "Version",
				order:    0,
				version:  anyVersion,
				nVersion: anyVersion,
				size:     8,
				length:   LengthUnlimited,
			}, {
				name:     "Flags",
				order:    1,
				version:  anyVersion,
				nVersion: anyVersion,
				size:     8,
				length:   LengthUnlimited,
			},
		},
	}, fs[0])
	assert.Equal(t, &field{
		name:     "Int32",
		order:    1,
		version:  anyVersion,
		nVersion: anyVersion,
		size:     32,
		length:   LengthUnlimited,
	}, fs[1])
	assert.Equal(t, &field{
		name:     "Int17",
		order:    2,
		version:  anyVersion,
		nVersion: anyVersion,
		size:     17,
		length:   LengthUnlimited,
	}, fs[2])
	assert.Equal(t, &field{
		name:     "Uint15",
		order:    3,
		version:  anyVersion,
		nVersion: anyVersion,
		size:     15,
		length:   LengthUnlimited,
	}, fs[3])
	assert.Equal(t, &field{
		name:     "Const",
		order:    4,
		version:  anyVersion,
		nVersion: anyVersion,
		size:     8,
		length:   LengthUnlimited,
		cnst:     "0",
	}, fs[4])
	assert.Equal(t, &field{
		name:     "String",
		order:    5,
		version:  anyVersion,
		nVersion: anyVersion,
		size:     8,
		length:   LengthUnlimited,
		flags:    fieldString,
		strType:  stringType_C,
	}, fs[5])
	assert.Equal(t, &field{
		name:     "PString",
		order:    6,
		version:  anyVersion,
		nVersion: anyVersion,
		size:     8,
		length:   LengthUnlimited,
		flags:    fieldString,
		strType:  stringType_C_P,
	}, fs[6])
	assert.Equal(t, &field{
		name:     "Dec",
		order:    7,
		version:  anyVersion,
		nVersion: anyVersion,
		size:     8,
		length:   LengthUnlimited,
		flags:    fieldDec,
	}, fs[7])
	assert.Equal(t, &field{
		name:     "Hex",
		order:    8,
		version:  anyVersion,
		nVersion: anyVersion,
		size:     8,
		length:   LengthUnlimited,
		flags:    fieldHex,
	}, fs[8])
	assert.Equal(t, &field{
		name:     "ISO639_2",
		order:    9,
		version:  anyVersion,
		nVersion: anyVersion,
		size:     8,
		length:   LengthUnlimited,
		flags:    fieldISO639_2,
	}, fs[9])
	assert.Equal(t, &field{
		name:     "UUID",
		order:    10,
		version:  anyVersion,
		nVersion: anyVersion,
		size:     8,
		length:   LengthUnlimited,
		flags:    fieldUUID,
	}, fs[10])
	assert.Equal(t, &field{
		name:     "Hidden",
		order:    11,
		version:  anyVersion,
		nVersion: anyVersion,
		size:     8,
		length:   LengthUnlimited,
		flags:    fieldHidden,
	}, fs[11])
	assert.Equal(t, &field{
		name:     "Opt",
		order:    12,
		version:  anyVersion,
		nVersion: anyVersion,
		size:     8,
		length:   LengthUnlimited,
		optFlag:  0x000010,
	}, fs[12])
	assert.Equal(t, &field{
		name:     "NOpt",
		order:    13,
		version:  anyVersion,
		nVersion: anyVersion,
		size:     8,
		length:   LengthUnlimited,
		nOptFlag: 0x000010,
	}, fs[13])
	assert.Equal(t, &field{
		name:     "DynOpt",
		order:    14,
		version:  anyVersion,
		nVersion: anyVersion,
		size:     8,
		length:   LengthUnlimited,
		flags:    fieldOptDynamic,
	}, fs[14])
	assert.Equal(t, &field{
		name:     "Varint",
		order:    15,
		version:  anyVersion,
		nVersion: anyVersion,
		length:   LengthUnlimited,
		flags:    fieldVarint,
	}, fs[15])
	assert.Equal(t, &field{
		name:     "DynSize",
		order:    16,
		version:  anyVersion,
		nVersion: anyVersion,
		length:   LengthUnlimited,
		flags:    fieldSizeDynamic,
	}, fs[16])
	assert.Equal(t, &field{
		name:     "FixedLen",
		order:    17,
		version:  anyVersion,
		nVersion: anyVersion,
		size:     8,
		length:   5,
	}, fs[17])
	assert.Equal(t, &field{
		name:     "DynLen",
		order:    18,
		version:  anyVersion,
		nVersion: anyVersion,
		size:     8,
		length:   LengthUnlimited,
		flags:    fieldLengthDynamic,
	}, fs[18])
	assert.Equal(t, &field{
		name:     "Ver",
		order:    19,
		version:  1,
		nVersion: anyVersion,
		size:     8,
		length:   LengthUnlimited,
	}, fs[19])
	assert.Equal(t, &field{
		name:     "NVer",
		order:    20,
		version:  anyVersion,
		nVersion: 1,
		size:     8,
		length:   LengthUnlimited,
	}, fs[20])
	assert.Equal(t, &field{
		name:     "NotSorted21",
		order:    21,
		version:  anyVersion,
		nVersion: anyVersion,
		size:     8,
		length:   LengthUnlimited,
	}, fs[21])
	assert.Equal(t, &field{
		name:     "NotSorted22",
		order:    22,
		version:  anyVersion,
		nVersion: anyVersion,
		size:     8,
		length:   LengthUnlimited,
	}, fs[22])
	assert.Equal(t, &field{
		name:     "NotSorted23",
		order:    23,
		version:  anyVersion,
		nVersion: anyVersion,
		size:     8,
		length:   LengthUnlimited,
	}, fs[23])
}

func TestResolveFieldInstance(t *testing.T) {
	fixedSize := uint(8)
	fixedLen := uint(1)
	dynSize1 := uint(16)
	dynLen1 := uint(2)
	dynSize2 := uint(32)
	dynLen2 := uint(4)
	cfo1 := struct {
		mockBox
		Box
	}{
		mockBox: mockBox{
			DynSizeMap: map[string]uint{"TestField": dynSize1},
			DynLenMap:  map[string]uint{"TestField": dynLen1},
		},
	}
	cfo2 := struct {
		mockBox
		Box
	}{
		mockBox: mockBox{
			DynSizeMap: map[string]uint{"TestField": dynSize2},
			DynLenMap:  map[string]uint{"TestField": dynLen2},
		},
	}
	nonCFO := struct{}{}

	testCases := []struct {
		name     string
		f        *field
		box      IImmutableBox
		parent   interface{}
		wantSize uint
		wantLen  uint
		wantCFO  ICustomFieldObject
	}{
		{
			name: "dynamic size with non CustomFieldObject",
			f: &field{
				name:   "TestField",
				flags:  fieldSizeDynamic,
				length: fixedLen,
			},
			box:      &cfo1,
			parent:   &nonCFO,
			wantSize: dynSize1,
			wantLen:  fixedLen,
			wantCFO:  &cfo1,
		},
		{
			name: "dynamic size with CustomFieldObject",
			f: &field{
				name:   "TestField",
				flags:  fieldSizeDynamic,
				length: fixedLen,
			},
			box:      &cfo1,
			parent:   &cfo2,
			wantSize: dynSize2,
			wantLen:  fixedLen,
			wantCFO:  &cfo2,
		},
		{
			name: "dynamic length with non CustomFieldObject",
			f: &field{
				name:  "TestField",
				flags: fieldLengthDynamic,
				size:  fixedSize,
			},
			box:      &cfo1,
			parent:   &nonCFO,
			wantSize: fixedSize,
			wantLen:  dynLen1,
			wantCFO:  &cfo1,
		},
		{
			name: "dynamic length with CustomFieldObject",
			f: &field{
				name:  "TestField",
				flags: fieldLengthDynamic,
				size:  fixedSize,
			},
			box:      &cfo1,
			parent:   &cfo2,
			wantSize: fixedSize,
			wantLen:  dynLen2,
			wantCFO:  &cfo2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fi := resolveFieldInstance(tc.f, tc.box, reflect.ValueOf(tc.parent).Elem(), Context{})
			assert.Equal(t, tc.wantSize, fi.size)
			assert.Equal(t, tc.wantLen, fi.length)
			assert.Same(t, tc.wantCFO, fi.cfo)
		})
	}
}

func TestIsTargetField(t *testing.T) {
	box := &struct {
		AnyTypeBox
		FullBox
	}{
		FullBox: FullBox{
			Version: 1,
			Flags:   [3]byte{0x00, 0x00, 0x06},
		},
	}

	cfo := struct {
		mockBox
		Box
	}{
		mockBox: mockBox{
			DynOptMap: map[string]bool{
				"DynEnabledField":  true,
				"DynDisabledField": false,
			},
		},
	}

	testCases := []struct {
		name  string
		fi    *fieldInstance
		wants bool
	}{
		{
			name: "normal",
			fi: &fieldInstance{
				field: field{
					name:     "TestField",
					version:  anyVersion,
					nVersion: anyVersion,
				},
				cfo: &cfo,
			},
			wants: true,
		},
		{
			name: "ver=0",
			fi: &fieldInstance{
				field: field{
					name:     "TestField",
					version:  0,
					nVersion: anyVersion,
				},
				cfo: &cfo,
			},
			wants: false,
		},
		{
			name: "ver=1",
			fi: &fieldInstance{
				field: field{
					name:     "TestField",
					version:  1,
					nVersion: anyVersion,
				},
				cfo: &cfo,
			},
			wants: true,
		},
		{
			name: "nver=0",
			fi: &fieldInstance{
				field: field{
					name:     "TestField",
					version:  anyVersion,
					nVersion: 0,
				},
				cfo: &cfo,
			},
			wants: true,
		},
		{
			name: "nver=1",
			fi: &fieldInstance{
				field: field{
					name:     "TestField",
					version:  anyVersion,
					nVersion: 1,
				},
				cfo: &cfo,
			},
			wants: false,
		},
		{
			name: "opt=0x000001",
			fi: &fieldInstance{
				field: field{
					name:     "TestField",
					version:  anyVersion,
					nVersion: anyVersion,
					optFlag:  0x000001,
				},
				cfo: &cfo,
			},
			wants: false,
		},
		{
			name: "opt=0x000002",
			fi: &fieldInstance{
				field: field{
					name:     "TestField",
					version:  anyVersion,
					nVersion: anyVersion,
					optFlag:  0x000002,
				},
				cfo: &cfo,
			},
			wants: true,
		},
		{
			name: "opt=0x000004",
			fi: &fieldInstance{
				field: field{
					name:     "TestField",
					version:  anyVersion,
					nVersion: anyVersion,
					optFlag:  0x000004,
				},
				cfo: &cfo,
			},
			wants: true,
		},
		{
			name: "opt=0x000008",
			fi: &fieldInstance{
				field: field{
					name:     "TestField",
					version:  anyVersion,
					nVersion: anyVersion,
					optFlag:  0x000008,
				},
				cfo: &cfo,
			},
			wants: false,
		},
		{
			name: "nopt=0x000001",
			fi: &fieldInstance{
				field: field{
					name:     "TestField",
					version:  anyVersion,
					nVersion: anyVersion,
					nOptFlag: 0x000001,
				},
				cfo: &cfo,
			},
			wants: true,
		},
		{
			name: "nopt=0x000002",
			fi: &fieldInstance{
				field: field{
					name:     "TestField",
					version:  anyVersion,
					nVersion: anyVersion,
					nOptFlag: 0x000002,
				},
				cfo: &cfo,
			},
			wants: false,
		},
		{
			name: "opt=dynamic enabled",
			fi: &fieldInstance{
				field: field{
					name:     "DynEnabledField",
					version:  anyVersion,
					nVersion: anyVersion,
					flags:    fieldOptDynamic,
				},
				cfo: &cfo,
			},
			wants: true,
		},
		{
			name: "opt=dynamic disabled",
			fi: &fieldInstance{
				field: field{
					name:     "DynDisabledField",
					version:  anyVersion,
					nVersion: anyVersion,
					flags:    fieldOptDynamic,
				},
				cfo: &cfo,
			},
			wants: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.wants, isTargetField(box, tc.fi, Context{}))
		})
	}
}
