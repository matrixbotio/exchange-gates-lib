package matrixgates

import (
	"github.com/matrixbotio/exchange-gates/workers"
)

//ExchangeInterface - universal exchange adapter interface
type ExchangeInterface interface {
	// Methods
	Connect(credentials APICredentials) error
	GetOrderData(pairSymbol string, orderID int64) (*TradeEventData, error)
	PlaceOrder(order BotOrder, pairLimits ExchangePairData) (*CreateOrderResponse, error)
	GetAccountData() (*AccountData, error)
	GetPairLastPrice(pairSymbol string) (float64, error)
	CancelPairOrder(pairSymbol string, orderID int64) error
	CancelPairOrders(pairSymbol string) error
	GetPairOpenOrders(pairSymbol string) ([]*Order, error)
	VerifyAPIKeys(keyPublic, keySecret string) error
	GetPairs() ([]*ExchangePairData, error)
	// Workers
	GetPriceWorker() workers.IPriceWorker
	GetCandleWorker() workers.ICandleWorker
}

//ExchangeAdapters - map of all supported exchanges
var ExchangeAdapters map[int]ExchangeInterface = map[int]ExchangeInterface{
	1: NewBinanceSpotAdapter(),
}
