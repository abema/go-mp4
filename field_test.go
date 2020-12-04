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
		FullBox    `mp4:"0,extend"`
		ByteArray  []byte `mp4:"1,size=8"`
		Int        int32  `mp4:"2,size=32"`
		NotSorted4 byte   `mp4:"4,size=8"`
		NotSorted5 byte   `mp4:"5,size=8"`
		NotSorted3 byte   `mp4:"3,size=8"`
	}{}

	fs := buildFields(box)
	assert.Equal(t, []*field{
		{name: "FullBox", order: 0, children: []*field{
			{name: "Version", order: 0},
			{name: "Flags", order: 1},
		}},
		{name: "ByteArray", order: 1},
		{name: "Int", order: 2},
		{name: "NotSorted3", order: 3},
		{name: "NotSorted4", order: 4},
		{name: "NotSorted5", order: 5},
	}, fs)
}

func TestReadFieldConfig(t *testing.T) {
	box := &struct {
		mockBox
		FullBox
		ByteArray []byte
		String    string
		Int       int32
	}{
		mockBox: mockBox{
			DynSizeMap: map[string]uint{
				"ByteArray": 3,
			},
			DynLenMap: map[string]uint{
				"ByteArray": 7,
			},
		},
	}

	testCases := []struct {
		name      string
		box       IImmutableBox
		fieldName string
		fieldTag  fieldTag
		err       bool
		expected  fieldConfig
	}{
		{
			name:      "static size",
			box:       box,
			fieldName: "ByteArray",
			fieldTag:  fieldTag{"size": "8"},
			expected: fieldConfig{
				name:     "ByteArray",
				cfo:      box,
				size:     8,
				length:   LengthUnlimited,
				version:  anyVersion,
				nVersion: anyVersion,
			},
		},
		{
			name:      "invalid size",
			box:       box,
			fieldName: "ByteArray",
			fieldTag:  fieldTag{"size": "invalid"},
			err:       true,
		},
		{
			name:      "dynamic size",
			box:       box,
			fieldName: "ByteArray",
			fieldTag:  fieldTag{"size": "dynamic"},
			expected: fieldConfig{
				name:     "ByteArray",
				cfo:      box,
				size:     3,
				length:   LengthUnlimited,
				version:  anyVersion,
				nVersion: anyVersion,
			},
		},
		{
			name:      "static length",
			box:       box,
			fieldName: "ByteArray",
			fieldTag:  fieldTag{"len": "16", "size": "8"},
			expected: fieldConfig{
				name:     "ByteArray",
				cfo:      box,
				size:     8,
				length:   16,
				version:  anyVersion,
				nVersion: anyVersion,
			},
		},
		{
			name:      "invalid length",
			box:       box,
			fieldName: "ByteArray",
			fieldTag:  fieldTag{"len": "foo", "size": "8"},
			err:       true,
		},
		{
			name:      "dynamic length",
			box:       box,
			fieldName: "ByteArray",
			fieldTag:  fieldTag{"len": "dynamic", "size": "8"},
			expected: fieldConfig{
				name:     "ByteArray",
				cfo:      box,
				size:     8,
				length:   7,
				version:  anyVersion,
				nVersion: anyVersion,
			},
		},
		{
			name:      "varint",
			box:       box,
			fieldName: "Int",
			fieldTag:  fieldTag{"varint": "", "size": "13"},
			expected: fieldConfig{
				name:     "Int",
				cfo:      box,
				size:     13,
				length:   LengthUnlimited,
				version:  anyVersion,
				nVersion: anyVersion,
				varint:   true,
			},
		},
		{
			name:      "ver 0",
			box:       box,
			fieldName: "Int",
			fieldTag:  fieldTag{"ver": "0", "size": "32"},
			expected: fieldConfig{
				name:     "Int",
				cfo:      box,
				size:     32,
				length:   LengthUnlimited,
				version:  0,
				nVersion: anyVersion,
			},
		},
		{
			name:      "ver 1",
			box:       box,
			fieldName: "Int",
			fieldTag:  fieldTag{"ver": "1", "size": "32"},
			expected: fieldConfig{
				name:     "Int",
				cfo:      box,
				size:     32,
				length:   LengthUnlimited,
				version:  1,
				nVersion: anyVersion,
			},
		},
		{
			name:      "invalid ver",
			box:       box,
			fieldName: "Int",
			fieldTag:  fieldTag{"ver": "foo", "size": "32"},
			err:       true,
		},
		{
			name:      "nver 0",
			box:       box,
			fieldName: "Int",
			fieldTag:  fieldTag{"nver": "0", "size": "32"},
			expected: fieldConfig{
				name:     "Int",
				cfo:      box,
				size:     32,
				length:   LengthUnlimited,
				version:  anyVersion,
				nVersion: 0,
			},
		},
		{
			name:      "nver 1",
			box:       box,
			fieldName: "Int",
			fieldTag:  fieldTag{"nver": "1", "size": "32"},
			expected: fieldConfig{
				name:     "Int",
				cfo:      box,
				size:     32,
				length:   LengthUnlimited,
				version:  anyVersion,
				nVersion: 1,
			},
		},
		{
			name:      "invalid nver",
			box:       box,
			fieldName: "Int",
			fieldTag:  fieldTag{"nver": "foo", "size": "32"},
			err:       true,
		},
		{
			name:      "opt dynamic",
			box:       box,
			fieldName: "String",
			fieldTag:  fieldTag{"opt": "dynamic"},
			expected: fieldConfig{
				name:       "String",
				cfo:        box,
				length:     LengthUnlimited,
				version:    anyVersion,
				nVersion:   anyVersion,
				optDynamic: true,
			},
		},
		{
			name:      "opt hex",
			box:       box,
			fieldName: "String",
			fieldTag:  fieldTag{"opt": "0x0100"},
			expected: fieldConfig{
				name:     "String",
				cfo:      box,
				length:   LengthUnlimited,
				version:  anyVersion,
				nVersion: anyVersion,
				optFlag:  0x0100,
			},
		},
		{
			name:      "opt dec",
			box:       box,
			fieldName: "String",
			fieldTag:  fieldTag{"opt": "32"},
			expected: fieldConfig{
				name:     "String",
				cfo:      box,
				length:   LengthUnlimited,
				version:  anyVersion,
				nVersion: anyVersion,
				optFlag:  0x0020,
			},
		},
		{
			name:      "invalid opt",
			box:       box,
			fieldName: "String",
			fieldTag:  fieldTag{"opt": "foo"},
			err:       true,
		},
		{
			name:      "nopt hex",
			box:       box,
			fieldName: "String",
			fieldTag:  fieldTag{"nopt": "0x0100"},
			expected: fieldConfig{
				name:     "String",
				cfo:      box,
				length:   LengthUnlimited,
				version:  anyVersion,
				nVersion: anyVersion,
				nOptFlag: 0x0100,
			},
		},
		{
			name:      "nopt dec",
			box:       box,
			fieldName: "String",
			fieldTag:  fieldTag{"nopt": "32"},
			expected: fieldConfig{
				name:     "String",
				cfo:      box,
				length:   LengthUnlimited,
				version:  anyVersion,
				nVersion: anyVersion,
				nOptFlag: 0x0020,
			},
		},
		{
			name:      "invalid nopt",
			box:       box,
			fieldName: "String",
			fieldTag:  fieldTag{"nopt": "foo"},
			err:       true,
		},
		{
			name:      "const",
			box:       box,
			fieldName: "Int",
			fieldTag:  fieldTag{"const": "0", "size": "32"},
			expected: fieldConfig{
				name:     "Int",
				cfo:      box,
				size:     32,
				length:   LengthUnlimited,
				version:  anyVersion,
				nVersion: anyVersion,
				cnst:     "0",
			},
		},
		{
			name:      "extend",
			box:       box,
			fieldName: "FullBox",
			fieldTag:  fieldTag{"extend": ""},
			expected: fieldConfig{
				name:     "FullBox",
				cfo:      box,
				length:   LengthUnlimited,
				version:  anyVersion,
				nVersion: anyVersion,
				extend:   true,
			},
		},
		{
			name:      "hex",
			box:       box,
			fieldName: "Int",
			fieldTag:  fieldTag{"hex": "", "size": "32"},
			expected: fieldConfig{
				name:     "Int",
				cfo:      box,
				size:     32,
				length:   LengthUnlimited,
				version:  anyVersion,
				nVersion: anyVersion,
				hex:      true,
			},
		},
		{
			name:      "string - c style",
			box:       box,
			fieldName: "String",
			fieldTag:  fieldTag{"string": ""},
			expected: fieldConfig{
				name:     "String",
				cfo:      box,
				length:   LengthUnlimited,
				version:  anyVersion,
				nVersion: anyVersion,
				str:      true,
				strType:  StringType_C,
			},
		},
		{
			name:      "string - c style or pascal style",
			box:       box,
			fieldName: "String",
			fieldTag:  fieldTag{"string": "c_p"},
			expected: fieldConfig{
				name:     "String",
				cfo:      box,
				length:   LengthUnlimited,
				version:  anyVersion,
				nVersion: anyVersion,
				str:      true,
				strType:  StringType_C_P,
			},
		},
		{
			name:      "iso639-2",
			box:       box,
			fieldName: "String",
			fieldTag:  fieldTag{"iso639-2": ""},
			expected: fieldConfig{
				name:     "String",
				cfo:      box,
				length:   LengthUnlimited,
				version:  anyVersion,
				nVersion: anyVersion,
				iso639_2: true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(tc.box).Elem()
			config, err := readFieldConfig(tc.box, v, tc.fieldName, tc.fieldTag, Context{})
			if tc.err {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expected, config)
		})
	}
}
