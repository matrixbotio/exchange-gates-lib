package errs

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"

	pkgErrs "github.com/matrixbotio/exchange-gates-lib/pkg/errs"
)

func TestHandleCancelOrderErrorEmpty(t *testing.T) {
	// given
	var orderIDFormatted = "123"
	var pairSymbol = "LTCBUSD"
	var testErr error

	// when
	err := MapCancelOrderError(orderIDFormatted, pairSymbol, testErr)

	// then
	assert.NoError(t, err)
}

func TestHandleCancelOrderErrorAlreadyCancelled(t *testing.T) {
	// given
	var orderIDFormatted = "123"
	var pairSymbol = "LTCBUSD"
	var testErr = errors.New("error: Order has been canceled")

	// when
	err := MapCancelOrderError(orderIDFormatted, pairSymbol, testErr)

	// then
	assert.Error(t, err)
	assert.Equal(t, pkgErrs.ErrOrderNotFound, err)
}

func TestHandleCancelOrderErrorFilled(t *testing.T) {
	// given
	var orderIDFormatted = "123"
	var pairSymbol = "LTCBUSD"
	var testErr = pkgErrs.ErrOrderFilled

	// when
	err := MapCancelOrderError(orderIDFormatted, pairSymbol, testErr)

	// then
	assert.Error(t, err)
}

func TestHandleCancelOrderError(t *testing.T) {
	// given
	var orderIDFormatted = "123"
	var pairSymbol = "LTCBUSD"
	var testErr error = errors.New("unknown error")

	// when
	err := MapCancelOrderError(orderIDFormatted, pairSymbol, testErr)

	// then
	assert.Error(t, err)
}
