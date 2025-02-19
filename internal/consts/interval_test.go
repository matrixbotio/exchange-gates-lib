package consts

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateIntervalNotSet(t *testing.T) {
	// given
	interval := ""

	// when
	err := ValidateInterval(interval)

	// then
	require.ErrorIs(t, err, errIntervalNotSet)
}

func TestValidateIntervalUnknown(t *testing.T) {
	// given
	interval := "wtf"

	// when
	err := ValidateInterval(interval)

	// then
	require.ErrorContains(t, err, "invalid")
}

func TestValidateIntervalSuccess(t *testing.T) {
	// given
	interval := "1h"

	// when
	err := ValidateInterval(interval)

	// then
	require.NoError(t, err)
}
