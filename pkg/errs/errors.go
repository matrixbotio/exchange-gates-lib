package errs

import "errors"

var (
	ErrOrderFilled        = errors.New("Order has been filled")
	ErrOrderPendingCancel = errors.New("order cancellation in progress")

	// TBD: https://github.com/matrixbotio/exchange-gates-lib/issues/191
	OrderNotFound = errors.New("order not found")
)
