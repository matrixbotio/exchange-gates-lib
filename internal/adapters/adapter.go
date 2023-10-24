package adapters

import (
	"context"

	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
)

type Adapter interface {
	// Adapter
	GetName() string
	GetTag() string
	GetID() int

	// Methods
	// TBD: https://github.com/matrixbotio/exchange-gates-lib/issues/156
	Connect(credentials pkgStructs.APICredentials) error
	CanTrade() (bool, error)
	VerifyAPIKeys(keyPublic, keySecret string) error
	GetAccountBalance() ([]structs.Balance, error)

	// Order
	GetOrderData(pairSymbol string, orderID int64) (structs.OrderData, error)
	GetOrderByClientOrderID(pairSymbol, clientOrderID string) (structs.OrderData, error)
	PlaceOrder(
		ctx context.Context,
		order structs.BotOrderAdjusted,
	) (structs.CreateOrderResponse, error)
	GetOrderExecFee(
		pairSymbol string,
		orderSide string,
		orderID int64,
	) (structs.OrderFees, error)

	// Pair
	GetPairData(pairSymbol string) (structs.ExchangePairData, error)
	GetPairLastPrice(pairSymbol string) (float64, error)
	CancelPairOrder(pairSymbol string, orderID int64, ctx context.Context) error
	CancelPairOrderByClientOrderID(
		pairSymbol string,
		clientOrderID string,
		ctx context.Context,
	) error
	GetPairOpenOrders(pairSymbol string) ([]structs.OrderData, error)
	GetPairs() ([]structs.ExchangePairData, error)
	GetPairBalance(pair structs.PairSymbolData) (structs.PairBalance, error)

	// Workers
	GetPriceWorker(callback workers.PriceEventCallback) workers.IPriceWorker
	GetCandleWorker() workers.ICandleWorker
	GetTradeEventsWorker() workers.ITradeEventWorker

	// Candle
	GetCandles(limit int, symbol string, interval string) ([]workers.CandleData, error)

	// TBD: remove: https://github.com/matrixbotio/exchange-gates-lib/issues/149
	GetPrices() ([]structs.SymbolPrice, error)
	GetPairOrdersHistory(task structs.GetOrdersHistoryTask) ([]structs.OrderData, error)
}
