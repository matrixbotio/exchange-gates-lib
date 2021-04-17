package matrixgates

import (
	sharederrs "github.com/matrixbotio/shared-errors"
)

//ExchangeInterface - universal exchange adapter interface
type ExchangeInterface interface {
	Connect(credentials APICredentials) *sharederrs.APIError
	GetOrderData() (*TradeEventData, *sharederrs.APIError)
	PlaceOrder(order BotOrder) (*CreateOrderResponse, *sharederrs.APIError)
	GetAccountData() (*AccountData, *sharederrs.APIError)
	GetPairLastPrice() (float64, *sharederrs.APIError)
	CancelPairOrder() *sharederrs.APIError
	CancelPairOrders() *sharederrs.APIError
	GetPairOpenOrders() ([]*Order, *sharederrs.APIError)
	VerifyAPIKeys(keyPublic, keySecret string) *sharederrs.APIError
	GetPairs() *sharederrs.APIError
}

//ExchangeAdapters - map of all supported exchanges
var ExchangeAdapters map[int]*ExchangeAdapter = map[int]*ExchangeAdapter{
	1: NewBinanceSpotAdapter(),
}
