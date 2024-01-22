package errs

import "errors"

var (
	ErrOrderFilled                 = errors.New("order has been filled")
	ErrOrderNotFound               = errors.New("order not found")
	ErrOrderCancellationInProgress = errors.New("order cancellation in progress")
	ErrOrderDuplicate              = errors.New("order has already been placed")
)
