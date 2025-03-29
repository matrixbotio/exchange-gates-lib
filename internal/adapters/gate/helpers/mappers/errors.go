package mappers

import (
	"strings"

	"github.com/matrixbotio/exchange-gates-lib/pkg/errs"
)

const ErrOrderNotActualMessage = "Order not found"

func MapCancelOrderErr(err error) error {
	if err == nil {
		return nil
	}

	if strings.Contains(err.Error(), ErrOrderNotActualMessage) {
		return errs.ErrOrderNotFound
	}
	return err
}
