package adapter

import (
	"context"

	"github.com/matrixbotio/exchange-gates-lib/pkg/structs"
	"github.com/matrixbotio/exchange-gates-lib/pkg/workers"
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
	GetPriceWorker(callback workers.PriceEventCallback) workers.IPriceWorker
	GetCandleWorker() workers.ICandleWorker
	GetTradeEventsWorker() workers.ITradeEventWorker
}
