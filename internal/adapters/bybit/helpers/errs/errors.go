package errs

import "errors"

var ErrSpotBalanceNotFound = errors.New("spot balance not found")

const (
	ErrMsgOrderHasBeenFilled          = "Order has been filled"
	ErrMsgOrderHasBeenCancelled       = "Order has been canceled"
	ErrMsgOrderNotFound               = "Order does not exist"
	ErrMsgOrderCancellationInProgress = "Order cancellation in progress"
)
