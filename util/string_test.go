package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatSignedFixedFloat1616(t *testing.T) {
	assert.Equal(t, "1", FormatSignedFixedFloat1616(0x00010000))
	assert.Equal(t, "1234", FormatSignedFixedFloat1616(0x04d20000))
	assert.Equal(t, "-1234", FormatSignedFixedFloat1616(-0x04d20000))
	assert.Equal(t, "1.50000", FormatSignedFixedFloat1616(0x00018000))
	assert.Equal(t, "1234.56789", FormatSignedFixedFloat1616(0x04d29161))
	assert.Equal(t, "-1234.56789", FormatSignedFixedFloat1616(-0x04d29161))
}

func TestFormatUnsignedFixedFloat1616(t *testing.T) {
	assert.Equal(t, "1", FormatUnsignedFixedFloat1616(0x00010000))
	assert.Equal(t, "1234", FormatUnsignedFixedFloat1616(0x04d20000))
	assert.Equal(t, "1.50000", FormatUnsignedFixedFloat1616(0x00018000))
	assert.Equal(t, "1234.56789", FormatUnsignedFixedFloat1616(0x04d29161))
	assert.Equal(t, "65535.99998", FormatUnsignedFixedFloat1616(0xffffffff))
}

func TestFormatSignedFixedFloat88(t *testing.T) {
	assert.Equal(t, "1", FormatSignedFixedFloat88(0x0100))
	assert.Equal(t, "123", FormatSignedFixedFloat88(0x7b00))
	assert.Equal(t, "-123", FormatSignedFixedFloat88(-0x7b00))
	assert.Equal(t, "1.500", FormatSignedFixedFloat88(0x0180))
	assert.Equal(t, "123.457", FormatSignedFixedFloat88(0x7b75))
	assert.Equal(t, "-123.457", FormatSignedFixedFloat88(-0x7b75))
}
