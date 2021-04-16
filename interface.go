package matrixgates

import (
	sharederrs "github.com/matrixbotio/shared-errors"
)

//ExchangeInterface - universal exchange adapter interface
type ExchangeInterface interface {
	getOrderData() (*TradeEventData, *sharederrs.APIError)
	placeOrder(order BotOrder) (*struct{}, *sharederrs.APIError)
	getAccountData() (*struct{}, *sharederrs.APIError)
	getPairLastPrice() (float64, *sharederrs.APIError)
	cancelPairOrder() *sharederrs.APIError
	cancelPairOrders() *sharederrs.APIError
	getPairOpenOrders() ([]*struct{}, *sharederrs.APIError)
	verifyAPIKeys() *sharederrs.APIError
}
