package mp4

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoxTypeString(t *testing.T) {
	assert.Equal(t, "1234", BoxType{'1', '2', '3', '4'}.String())
	assert.Equal(t, "abcd", BoxType{'a', 'b', 'c', 'd'}.String())
	assert.Equal(t, "xx x", BoxType{'x', 'x', ' ', 'x'}.String())
	assert.Equal(t, "xx~x", BoxType{'x', 'x', '~', 'x'}.String())
	assert.Equal(t, "xx(c)x", BoxType{'x', 'x', 0xa9, 'x'}.String())
	assert.Equal(t, "0x7878ab78", BoxType{'x', 'x', 0xab, 'x'}.String())
}

func TestIsSupported(t *testing.T) {
	assert.True(t, StrToBoxType("pssh").IsSupported(Context{}))
	assert.False(t, StrToBoxType("1234").IsSupported(Context{}))
}

func TestGetSupportedVersions(t *testing.T) {
	vers, err := BoxTypePssh().GetSupportedVersions(Context{})
	require.NoError(t, err)
	assert.Equal(t, []uint8{0, 1}, vers)
}

func TestIsSupportedVersion(t *testing.T) {
	assert.True(t, BoxTypePssh().IsSupportedVersion(0, Context{}))
	assert.True(t, BoxTypePssh().IsSupportedVersion(1, Context{}))
	assert.False(t, BoxTypePssh().IsSupportedVersion(2, Context{}))
}
