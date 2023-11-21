package errs

import (
	"errors"
	"strings"

	pkgErrs "github.com/matrixbotio/exchange-gates-lib/pkg/errs"
)

const (
	UnknownOrderMsg   = "Unknown order sent"
	orderNotExistsMsg = "Order does not exist"
	orderFilledMsg    = "Order has been filled"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials to connect to Binance")
	ErrOrderIDNotSet      = errors.New("orderID is not set")
	ErrAccountDataEmpty   = errors.New("account data response is empty")
)

func IsErrorAboutUnknownOrder(err error) bool {
	return strings.Contains(err.Error(), UnknownOrderMsg) ||
		strings.Contains(err.Error(), orderNotExistsMsg)
}

func IsErrorAboutOrderFilled(err error) bool {
	return strings.Contains(err.Error(), orderFilledMsg)
}

func HandleCancelOrderError(err error) error {
	if err == nil {
		return nil
	}

	if IsErrorAboutUnknownOrder(err) {
		return pkgErrs.OrderNotFound
	}
	if IsErrorAboutOrderFilled(err) {
		return pkgErrs.ErrOrderFilled
	}
	return err
}
