package errs

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	pkgErrs "github.com/matrixbotio/exchange-gates-lib/pkg/errs"
)

func TestHandleCancelOrderErrorEmpty(t *testing.T) {
	// given
	var orderIDFormatted = "123"
	var pairSymbol = "LTCBUSD"
	var testErr error

	// when
	err := HandleCancelOrderError(orderIDFormatted, pairSymbol, testErr)

	// then
	assert.NoError(t, err)
}

func TestHandleCancelOrderErrorAlreadyCancelled(t *testing.T) {
	// given
	var orderIDFormatted = "123"
	var pairSymbol = "LTCBUSD"
	var testErr = errors.New("error: Order has been canceled")

	// when
	err := HandleCancelOrderError(orderIDFormatted, pairSymbol, testErr)

	// then
	assert.NoError(t, err)
}

func TestHandleCancelOrderErrorFilled(t *testing.T) {
	// given
	var orderIDFormatted = "123"
	var pairSymbol = "LTCBUSD"
	var testErr = pkgErrs.ErrOrderFilled

	// when
	err := HandleCancelOrderError(orderIDFormatted, pairSymbol, testErr)

	// then
	assert.Error(t, err)
}

func TestHandleCancelOrderError(t *testing.T) {
	// given
	var orderIDFormatted = "123"
	var pairSymbol = "LTCBUSD"
	var testErr error = errors.New("unknown error")

	// when
	err := HandleCancelOrderError(orderIDFormatted, pairSymbol, testErr)

	// then
	assert.Error(t, err)
}
