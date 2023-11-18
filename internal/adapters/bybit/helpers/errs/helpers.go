package errs

import (
	"fmt"
	"strings"

	"github.com/matrixbotio/exchange-gates-lib/pkg/errs"
)

func HandleCancelOrderError(
	orderIDFormatted string,
	pairSymbol string,
	err error,
) error {
	if err == nil {
		return nil
	}

	if strings.Contains(err.Error(), "Order has been canceled") {
		return errs.OrderNotFound
	}
	if strings.Contains(err.Error(), errs.ErrOrderFilled.Error()) {
		return errs.ErrOrderFilled
	}
	if strings.Contains(err.Error(), "Order does not exist") {
		return errs.OrderNotFound
	}

	return fmt.Errorf(
		"cancel order %s in %q: %w",
		orderIDFormatted, pairSymbol, err,
	)
}
