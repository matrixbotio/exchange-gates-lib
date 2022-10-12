package adapters

import (
	"context"
	"errors"

	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/pkg/adapters/binance"
	"github.com/matrixbotio/exchange-gates-lib/pkg/structs"
	workers2 "github.com/matrixbotio/exchange-gates-lib/pkg/workers"
)

type Adapter interface {
	// Adapter
	GetName() string
	GetTag() string
	GetID() int

	// Methods
	Connect(credentials structs.APICredentials) error
	GetAccountData() (structs.AccountData, error) // account data with balances
	VerifyAPIKeys(keyPublic, keySecret string) error

	// Order
	GetPrices() ([]structs.SymbolPrice, error)
	GetOrderData(pairSymbol string, orderID int64) (structs.OrderData, error)
	GetOrderByClientOrderID(pairSymbol, clientOrderID string) (structs.OrderData, error)
	PlaceOrder(ctx context.Context, order structs.BotOrderAdjusted) (structs.CreateOrderResponse, error)

	// Pair
	GetPairData(pairSymbol string) (structs.ExchangePairData, error)
	GetPairLastPrice(pairSymbol string) (float64, error)
	CancelPairOrder(pairSymbol string, orderID int64, ctx context.Context) error
	CancelPairOrderByClientOrderID(pairSymbol string, clientOrderID string, ctx context.Context) error
	GetPairOpenOrders(pairSymbol string) ([]structs.OrderData, error)
	GetPairOrdersHistory(task structs.GetOrdersHistoryTask) ([]structs.OrderData, error)
	GetPairs() ([]structs.ExchangePairData, error)
	GetPairBalance(pair structs.PairSymbolData) (structs.PairBalance, error)

	// Workers
	GetPriceWorker(callback workers2.PriceEventCallback) workers2.IPriceWorker
	GetCandleWorker() workers2.ICandleWorker
	GetTradeEventsWorker() workers2.ITradeEventWorker
}

// GetExchangeAdapter - get supported exchange adapter with interface
func GetExchangeAdapter(exchangeID int) (Adapter, error) {
	switch exchangeID {
	default:
		return nil, errors.New("exchange not found")
	case consts.ExchangeIDbinanceSpot:
		return binance.NewBinanceSpotAdapter(), nil
	case consts.TestExchangeID:
		return GetTestExchangeAdapter(), nil
	}
}

// GetExchangeAdapters - get all supported exchange adapters
func GetExchangeAdapters() map[int]Adapter {
	return map[int]Adapter{
		consts.ExchangeIDbinanceSpot: binance.NewBinanceSpotAdapter(),
	}
}

func GetTestExchangeAdapter() Adapter {
	return &TestAdapter{
		ExchangeID: consts.TestExchangeID,
		Name:       "Test Exchange",
		Tag:        "test-exchange",
	}
}
