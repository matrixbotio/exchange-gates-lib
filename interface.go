package matrixgates

import (
	"./workers"
	sharederrs "github.com/matrixbotio/shared-errors"
)

//ExchangeInterface - universal exchange adapter interface
type ExchangeInterface interface {
	Connect(credentials APICredentials) *sharederrs.APIError
	GetOrderData(pairSymbol string, orderID int64) (*TradeEventData, *sharederrs.APIError)
	PlaceOrder(order BotOrder, pairLimits ExchangePairData) (*CreateOrderResponse, *sharederrs.APIError)
	GetAccountData() (*AccountData, *sharederrs.APIError)
	GetPairLastPrice(pairSymbol string) (float64, *sharederrs.APIError)
	CancelPairOrder(pairSymbol string, orderID int64) *sharederrs.APIError
	CancelPairOrders(pairSymbol string) *sharederrs.APIError
	GetPairOpenOrders(pairSymbol string) ([]*Order, *sharederrs.APIError)
	VerifyAPIKeys(keyPublic, keySecret string) *sharederrs.APIError
	GetPairs() ([]*ExchangePairData, *sharederrs.APIError)
	GetPriceWorker() workers.IPriceWorker
}

//ExchangeAdapters - map of all supported exchanges
var ExchangeAdapters map[int]ExchangeInterface = map[int]ExchangeInterface{
	1: NewBinanceSpotAdapter(),
}
