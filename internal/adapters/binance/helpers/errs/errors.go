package errs

import (
	"errors"
	"strings"
)

const (
	UnknownOrderMsg      = "Unknown order sent"
	ErrMsgOrderDuplicate = "Duplicate order sent"
	orderNotExistsMsg    = "Order does not exist"
	orderFilledMsg       = "Order has been filled"
)

var (
	ErrInvalidCredentials  = errors.New("invalid credentials to connect to Binance")
	ErrOrderIDNotSet       = errors.New("orderID is not set")
	ErrAccountDataEmpty    = errors.New("account data response is empty")
	ErrClientOrderIDNotSet = errors.New("client order ID is not set")
	ErrOrderResponseEmpty  = errors.New("order response is empty")
	ErrPairResponseEmpty   = errors.New("pairs response is empty")
	ErrTradingNotAllowed   = errors.New("your API key does not have permission to trade," +
		" change its restrictions")
)

func IsErrorAboutUnknownOrder(err error) bool {
	return strings.Contains(err.Error(), UnknownOrderMsg) ||
		strings.Contains(err.Error(), orderNotExistsMsg)
}

func IsErrorAboutOrderFilled(err error) bool {
	return strings.Contains(err.Error(), orderFilledMsg)
}
