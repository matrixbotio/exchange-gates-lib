package errs

import "errors"

var (
	ErrOrderFilled                 = errors.New("order has been filled")
	ErrOrderNotFound               = errors.New("order not found")
	ErrOrderCancellationInProgress = errors.New("order cancellation in progress")
	ErrOrderDuplicate              = errors.New("order has already been placed")
	ErrMinimumTP                   = errors.New("minimum TP order not passed")

	// ErrOrderDataNotActual returned when it is necessary to search for order data in history
	ErrOrderDataNotActual = errors.New("order data not actual")
)
