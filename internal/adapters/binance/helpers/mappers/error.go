package mappers

import (
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/errs"
	pkgErrs "github.com/matrixbotio/exchange-gates-lib/pkg/errs"
)

func MapCancelOrderError(err error) error {
	if err == nil {
		return nil
	}

	if errs.IsErrorAboutUnknownOrder(err) {
		return pkgErrs.ErrOrderNotFound
	}
	if errs.IsErrorAboutOrderFilled(err) {
		return pkgErrs.ErrOrderFilled
	}
	return err
}
