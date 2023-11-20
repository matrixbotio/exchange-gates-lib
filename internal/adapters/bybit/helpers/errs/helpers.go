package errs

import (
	"fmt"
	"strings"

	"github.com/matrixbotio/exchange-gates-lib/pkg/errs"
)

func MapCancelOrderError(orderIDFormatted string, pairSymbol string, err error) error {
	if err == nil {
		return nil
	}

	switch true {
	case strings.Contains(err.Error(), ErrMsgOrderHasBeenCancelled):
		return errs.ErrOrderNotFound
	case strings.Contains(err.Error(), ErrMsgOrderHasBeenFilled):
		return errs.ErrOrderFilled
	case strings.Contains(err.Error(), ErrMsgOrderNotFound):
		return errs.ErrOrderNotFound
	case strings.Contains(err.Error(), ErrMsgOrderCancellationInProgress):
		return errs.ErrOrderCancellationInProgress
	default:
		return fmt.Errorf(
			"cancel order %s in %q: %w",
			orderIDFormatted, pairSymbol, err,
		)
	}
}
