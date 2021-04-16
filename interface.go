package matrixgates

import (
	sharederrs "github.com/matrixbotio/shared-errors"
)

//ExchangeInterface - universal exchange adapter interface
type ExchangeInterface interface {
	GetOrderData() (*TradeEventData, *sharederrs.APIError)
	PlaceOrder(order BotOrder) (*CreateOrderResponse, *sharederrs.APIError)
	GetAccountData() (*AccountData, *sharederrs.APIError)
	GetPairLastPrice() (float64, *sharederrs.APIError)
	CancelPairOrder() *sharederrs.APIError
	CancelPairOrders() *sharederrs.APIError
	GetPairOpenOrders() ([]*Order, *sharederrs.APIError)
	VerifyAPIKeys() *sharederrs.APIError
}
