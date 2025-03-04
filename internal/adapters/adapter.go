//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package adapters

import (
	"context"

	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
)

type Adapter interface {
	// ADAPTER
	GetName() string
	GetTag() string
	GetID() int
	GetPairSymbol(baseTicker string, quoteTicker string) string

	// BASIC
	// TBD: call Connect on adapter init:
	//	https://github.com/matrixbotio/exchange-gates-lib/issues/156
	// Connect to exchange
	Connect(credentials pkgStructs.APICredentials) error
	// CanTrade - check the permission of the API key for trading
	CanTrade() (bool, error)
	// VerifyAPIKeys - Check if the API key has expired
	VerifyAPIKeys(keyPublic, keySecret string) error
	// GetAccountBalance - get account balances for individual tickers
	GetAccountBalance() ([]structs.Balance, error)
	GetLimits() pkgStructs.ExchangeLimits

	// ORDER
	// GetOrderData - get order data
	GetOrderData(pairSymbol string, orderID int64) (structs.OrderData, error)
	// GetClientOrderData - get order data by client order ID
	GetOrderByClientOrderID(pairSymbol, clientOrderID string) (structs.OrderData, error)
	// PlaceOrder - place order on exchange
	PlaceOrder(
		ctx context.Context,
		order structs.BotOrderAdjusted,
	) (structs.CreateOrderResponse, error)
	// Get the amount of fees for order execution
	GetOrderExecFee(
		baseAssetTicker string,
		quoteAssetTicker string,
		orderSide consts.OrderSide,
		orderID int64,
	) (structs.OrderFees, error)
	GenClientOrderID() string

	/*
		GetOrdersHistory - get orders history.

		NOTE: orderID is optional.
		Time in unix timestamp ms.
	*/
	GetHistoryOrder(
		baseAssetTicker string,
		quoteAssetTicker string,
		orderID int64,
	) (structs.OrderHistory, error)

	// PAIR
	// GetPairData - get pair data & limits
	GetPairData(pairSymbol string) (structs.ExchangePairData, error)
	// GetPairLastPrice - get pair last price ^ↀᴥↀ^
	GetPairLastPrice(pairSymbol string) (float64, error)
	// CancelPairOrder - cancel one exchange pair order by ID
	CancelPairOrder(pairSymbol string, orderID int64, ctx context.Context) error
	// CancelPairOrder - cancel one exchange pair order by client order ID
	CancelPairOrderByClientOrderID(
		pairSymbol string,
		clientOrderID string,
		ctx context.Context,
	) error
	// GetPairs get all Binance pairs
	GetPairs() ([]structs.ExchangePairData, error)
	// GetPairBalance - get pair balance: ticker, quote asset balance for pair symbol
	GetPairBalance(pair structs.PairSymbolData) (structs.PairBalance, error)

	// WORKERS
	// GetPriceWorker - create new market data worker
	GetPriceWorker(callback workers.PriceEventCallback) workers.IPriceWorker
	// GetCandleWorker - create new market candle worker
	GetCandleWorker() workers.ICandleWorker
	// GetTradeEventsWorker - create new market candle worker
	GetTradeEventsWorker() workers.ITradeEventWorker

	// CANDLE
	GetCandles(limit int, symbol string, interval consts.Interval) ([]workers.CandleData, error)
}
