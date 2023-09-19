package errs

import "errors"

var OrderNotFound = errors.New("order not found")

var (
	ErrOrderFilled = errors.New("Order has been filled")
)
