package pkg

import (
	"context"
	"errors"

	workers2 "github.com/matrixbotio/exchange-gates-lib/pkg/workers"
)

// ExchangeInterface - universal exchange adapter interface
type ExchangeInterface interface {
	// Adapter
	GetName() string
	GetTag() string
	GetID() int

	// Methods
	Connect(credentials APICredentials) error
	GetAccountData() (AccountData, error) // account data with balances
	VerifyAPIKeys(keyPublic, keySecret string) error

	// Order
	GetPrices() ([]SymbolPrice, error)
	GetOrderData(pairSymbol string, orderID int64) (OrderData, error)
	GetOrderByClientOrderID(pairSymbol, clientOrderID string) (OrderData, error)
	PlaceOrder(ctx context.Context, order BotOrderAdjusted) (CreateOrderResponse, error)

	// Pair
	GetPairData(pairSymbol string) (ExchangePairData, error)
	GetPairLastPrice(pairSymbol string) (float64, error)
	CancelPairOrder(pairSymbol string, orderID int64, ctx context.Context) error
	CancelPairOrderByClientOrderID(pairSymbol string, clientOrderID string, ctx context.Context) error
	GetPairOpenOrders(pairSymbol string) ([]OrderData, error)
	GetPairOrdersHistory(task GetOrdersHistoryTask) ([]OrderData, error)
	GetPairs() ([]ExchangePairData, error)
	GetPairBalance(pair PairSymbolData) (PairBalance, error)

	// Workers
	GetPriceWorker(callback workers2.PriceEventCallback) workers2.IPriceWorker
	GetCandleWorker() workers2.ICandleWorker
	GetTradeEventsWorker() workers2.ITradeEventWorker
}

// GetExchangeAdapter - get supported exchange adapter with interface
func GetExchangeAdapter(exchangeID int) (ExchangeInterface, error) {
	switch exchangeID {
	default:
		return nil, errors.New("exchange not found")
	case exchangeIDbinanceSpot:
		return NewBinanceSpotAdapter(), nil
	case TestExchangeID:
		return GetTestExchangeAdapter(), nil
	}
}

// GetExchangeAdapters - get all supported exchange adapters
func GetExchangeAdapters() map[int]ExchangeInterface {
	return map[int]ExchangeInterface{
		exchangeIDbinanceSpot: NewBinanceSpotAdapter(),
	}
}

func GetTestExchangeAdapter() ExchangeInterface {
	return &ExchangeAdapter{
		ExchangeID: TestExchangeID,
		Name:       "Test Exchange",
		Tag:        "test-exchange",
	}
}
