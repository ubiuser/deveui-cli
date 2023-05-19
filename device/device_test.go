package device

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewDevice(t *testing.T) {
	t.Parallel()

	got, err := NewDevice()
	require.NoError(t, err)

	// regex can be tested and explained here: https://regex101.com/

	reId := regexp.MustCompile("^[0-9A-F]{16}$")
	assert.True(t, reId.MatchString(got.GetIdentifier()))

	reCode := regexp.MustCompile("^[0-9A-F]{5}$")
	assert.True(t, reCode.MatchString(got.GetCode()))
}
