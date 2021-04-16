package matrixgates

import "github.com/matrixbotio/errors"

//ExchangeInterface - universal exchange adapter interface
type ExchangeInterface interface {
	getOrderData() (*TradeEventData, *errors.APIError)
	placeOrder(order BotOrder) (*struct{}, *errors.APIError)
	getAccountData() (*struct{}, *errors.APIError)
	getPairLastPrice() (float64, *errors.APIError)
	cancelPairOrder() *errors.APIError
	cancelPairOrders() *errors.APIError
	getPairOpenOrders() ([]*struct{}, *errors.APIError)
	verifyAPIKeys() *errors.APIError
}
