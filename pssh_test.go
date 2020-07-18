package mp4

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPsshStringify(t *testing.T) {
	flags := [3]byte{0x00, 0x00, 0x00}
	systemID := [16]byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
	}
	kid1 := PsshKID{KID: [16]byte{
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x10,
	}}
	kid2 := PsshKID{KID: [16]byte{
		0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
		0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x20,
	}}
	data := []byte{0x21, 0x22, 0x23, 0x24, 0x25}

	testCases := []struct {
		name string
		pssh Pssh
		want string
	}{
		{
			name: "version 0: no KIDs",
			pssh: Pssh{
				FullBox: FullBox{
					Version: 0,
					Flags:   flags,
				},
				SystemID: systemID,
				DataSize: int32(len(data)),
				Data:     data,
			},
			want: `Version=0 ` +
				`Flags=0x000000 ` +
				`SystemID="0102030405060708090a0b0c0d0e0f10" ` +
				`DataSize=5 ` +
				`Data=[0x21, 0x22, 0x23, 0x24, 0x25]`,
		},
		{
			name: "version 1: with KIDs",
			pssh: Pssh{
				FullBox: FullBox{
					Version: 1,
					Flags:   flags,
				},
				SystemID: systemID,
				KIDCount: 2,
				KIDs:     []PsshKID{kid1, kid2},
				DataSize: int32(len(data)),
				Data:     data,
			},
			want: `Version=1 ` +
				`Flags=0x000000 ` +
				`SystemID="0102030405060708090a0b0c0d0e0f10" ` +
				`KIDCount=2 ` +
				`KIDs=["1112131415161718191a1b1c1d1e1f10" "2122232425262728292a2b2c2d2e2f20"] ` +
				`DataSize=5 ` +
				`Data=[0x21, 0x22, 0x23, 0x24, 0x25]`,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			str, err := Stringify(&c.pssh)
			require.NoError(t, err)
			assert.Equal(t, c.want, str)
		})
	}
}
