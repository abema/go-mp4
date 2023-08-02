package mp4

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoxTypesISO23001_7(t *testing.T) {
	testCases := []struct {
		name string
		src  IImmutableBox
		dst  IBox
		bin  []byte
		str  string
		ctx  Context
	}{
		{
			name: "pssh: version 0: no KIDs",
			src: &Pssh{
				FullBox: FullBox{
					Version: 0,
					Flags:   [3]byte{0x00, 0x00, 0x00},
				},
				SystemID: [16]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10},
				DataSize: 5,
				Data:     []byte{0x21, 0x22, 0x23, 0x24, 0x25},
			},
			dst: &Pssh{},
			bin: []byte{
				0,                // version
				0x00, 0x00, 0x00, // flags
				// system ID
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
				0x00, 0x00, 0x00, 0x05, // data size
				0x21, 0x22, 0x23, 0x24, 0x25, // data
			},
			str: `Version=0 Flags=0x000000 ` +
				`SystemID=01020304-0506-0708-090a-0b0c0d0e0f10 ` +
				`DataSize=5 ` +
				`Data=[0x21, 0x22, 0x23, 0x24, 0x25]`,
		},
		{
			name: "pssh: version 1: with KIDs",
			src: &Pssh{
				FullBox: FullBox{
					Version: 1,
					Flags:   [3]byte{0x00, 0x00, 0x00},
				},
				SystemID: [16]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10},
				KIDCount: 2,
				KIDs: []PsshKID{
					{KID: [16]byte{0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x10}},
					{KID: [16]byte{0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x20}},
				},
				DataSize: 5,
				Data:     []byte{0x21, 0x22, 0x23, 0x24, 0x25},
			},
			dst: &Pssh{},
			bin: []byte{
				1,                // version
				0x00, 0x00, 0x00, // flags
				// system ID
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
				0x00, 0x00, 0x00, 0x02, // KID count
				// KIDs
				0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x10,
				0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x20,
				0x00, 0x00, 0x00, 0x05, // data size
				0x21, 0x22, 0x23, 0x24, 0x25, // data
			},
			str: `Version=1 Flags=0x000000 ` +
				`SystemID=01020304-0506-0708-090a-0b0c0d0e0f10 ` +
				`KIDCount=2 ` +
				`KIDs=[11121314-1516-1718-191a-1b1c1d1e1f10, 21222324-2526-2728-292a-2b2c2d2e2f20] ` +
				`DataSize=5 ` +
				`Data=[0x21, 0x22, 0x23, 0x24, 0x25]`,
		},
		{
			name: "tenc: DefaultIsProtected=1 DefaultPerSampleIVSize=0",
			src: &Tenc{
				FullBox: FullBox{
					Version: 1,
					Flags:   [3]byte{0x00, 0x00, 0x00},
				},
				Reserved:               0x00,
				DefaultCryptByteBlock:  0x0a,
				DefaultSkipByteBlock:   0x0b,
				DefaultIsProtected:     1,
				DefaultPerSampleIVSize: 0,
				DefaultKID: [16]byte{
					0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
					0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
				},
				DefaultConstantIVSize: 4,
				DefaultConstantIV:     []byte{0x01, 0x23, 0x45, 0x67},
			},
			dst: &Tenc{},
			bin: []byte{
				1,                // version
				0x00, 0x00, 0x00, // flags
				0x00,       // reserved
				0xab,       // default crypt byte block / default skip byte block
				0x01, 0x00, // default is protected / default per sample iv size
				0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, // default kid
				0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
				0x04,                   // default constant iv size
				0x01, 0x23, 0x45, 0x67, // default constant iv
			},
			str: `Version=1 Flags=0x000000 ` +
				`Reserved=0 ` +
				`DefaultCryptByteBlock=10 ` +
				`DefaultSkipByteBlock=11 ` +
				`DefaultIsProtected=1 ` +
				`DefaultPerSampleIVSize=0 ` +
				`DefaultKID=01234567-89ab-cdef-0123-456789abcdef ` +
				`DefaultConstantIVSize=4 ` +
				`DefaultConstantIV=[0x1, 0x23, 0x45, 0x67]`,
		},
		{
			name: "tenc: DefaultIsProtected=0 DefaultPerSampleIVSize=0",
			src: &Tenc{
				FullBox: FullBox{
					Version: 1,
					Flags:   [3]byte{0x00, 0x00, 0x00},
				},
				Reserved:               0x00,
				DefaultCryptByteBlock:  0x0a,
				DefaultSkipByteBlock:   0x0b,
				DefaultIsProtected:     0,
				DefaultPerSampleIVSize: 0,
				DefaultKID: [16]byte{
					0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
					0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
				},
			},
			dst: &Tenc{},
			bin: []byte{
				1,                // version
				0x00, 0x00, 0x00, // flags
				0x00,       // reserved
				0xab,       // default crypt byte block / default skip byte block
				0x00, 0x00, // default is protected / default per sample iv size
				0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, // default kid
				0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
			},
			str: `Version=1 Flags=0x000000 ` +
				`Reserved=0 ` +
				`DefaultCryptByteBlock=10 ` +
				`DefaultSkipByteBlock=11 ` +
				`DefaultIsProtected=0 ` +
				`DefaultPerSampleIVSize=0 ` +
				`DefaultKID=01234567-89ab-cdef-0123-456789abcdef`,
		},
		{
			name: "tenc: DefaultIsProtected=1 DefaultPerSampleIVSize=1",
			src: &Tenc{
				FullBox: FullBox{
					Version: 1,
					Flags:   [3]byte{0x00, 0x00, 0x00},
				},
				Reserved:               0x00,
				DefaultCryptByteBlock:  0x0a,
				DefaultSkipByteBlock:   0x0b,
				DefaultIsProtected:     1,
				DefaultPerSampleIVSize: 1,
				DefaultKID: [16]byte{
					0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
					0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
				},
			},
			dst: &Tenc{},
			bin: []byte{
				1,                // version
				0x00, 0x00, 0x00, // flags
				0x00,       // reserved
				0xab,       // default crypt byte block / default skip byte block
				0x01, 0x01, // default is protected / default per sample iv size
				0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, // default kid
				0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
			},
			str: `Version=1 Flags=0x000000 ` +
				`Reserved=0 ` +
				`DefaultCryptByteBlock=10 ` +
				`DefaultSkipByteBlock=11 ` +
				`DefaultIsProtected=1 ` +
				`DefaultPerSampleIVSize=1 ` +
				`DefaultKID=01234567-89ab-cdef-0123-456789abcdef`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Marshal
			buf := bytes.NewBuffer(nil)
			n, err := Marshal(buf, tc.src, tc.ctx)
			require.NoError(t, err)
			assert.Equal(t, uint64(len(tc.bin)), n)
			assert.Equal(t, tc.bin, buf.Bytes())

			// Unmarshal
			r := bytes.NewReader(tc.bin)
			n, err = Unmarshal(r, uint64(len(tc.bin)), tc.dst, tc.ctx)
			require.NoError(t, err)
			assert.Equal(t, uint64(buf.Len()), n)
			assert.Equal(t, tc.src, tc.dst)
			s, err := r.Seek(0, io.SeekCurrent)
			require.NoError(t, err)
			assert.Equal(t, int64(buf.Len()), s)

			// UnmarshalAny
			dst, n, err := UnmarshalAny(bytes.NewReader(tc.bin), tc.src.GetType(), uint64(len(tc.bin)), tc.ctx)
			require.NoError(t, err)
			assert.Equal(t, uint64(buf.Len()), n)
			assert.Equal(t, tc.src, dst)
			s, err = r.Seek(0, io.SeekCurrent)
			require.NoError(t, err)
			assert.Equal(t, int64(buf.Len()), s)

			// Stringify
			str, err := Stringify(tc.src, tc.ctx)
			require.NoError(t, err)
			assert.Equal(t, tc.str, str)
		})
	}
}
