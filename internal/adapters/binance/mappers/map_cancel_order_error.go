package mappers

import (
	"github.com/matrixbotio/exchange-gates-lib/pkg/errs"
	"strings"
)

func MapCancelOrderError(err error) error {
	if err == nil {
		return nil
	}

	switch true {
	case strings.Contains(err.Error(), "Unknown order sent"):
		return errs.ErrOrderNotFound
	case strings.Contains(err.Error(), "Order has been filled"):
		return errs.ErrOrderFilled
	default:
		return err
	}
}
