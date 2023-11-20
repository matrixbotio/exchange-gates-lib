package errs

import (
	"errors"
	"strings"

	pkgErrs "github.com/matrixbotio/exchange-gates-lib/pkg/errs"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials to connect to Binance")
)

func IsErrorAboutUnknownOrder(err error) bool {
	return strings.Contains(err.Error(), "Unknown order sent") ||
		strings.Contains(err.Error(), "Order does not exist")
}

func IsErrorAboutOrderFilled(err error) bool {
	return strings.Contains(err.Error(), "Order has been filled")
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
