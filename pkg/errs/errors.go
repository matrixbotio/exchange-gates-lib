package errs

import "errors"

var (
	ErrOrderFilled = errors.New("Order has been filled")
	OrderNotFound  = errors.New("order not found")
)
