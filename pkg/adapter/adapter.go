package adapter

import (
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
)

type Adapter = adapters.Adapter
type MockAdapter = adapters.MockAdapter

var NewMockAdapter = adapters.NewMockAdapter

type (
	AccountData          = structs.AccountData
	Balance              = structs.Balance
	SymbolPrice          = structs.SymbolPrice
	OrderData            = structs.OrderData
	BotOrderAdjusted     = structs.BotOrderAdjusted
	CreateOrderResponse  = structs.CreateOrderResponse
	ExchangePairData     = structs.ExchangePairData
	GetOrdersHistoryTask = structs.GetOrdersHistoryTask
	PairBalance          = structs.PairBalance
	PairSymbolData       = structs.PairSymbolData
	AssetBalance         = structs.AssetBalance
	OrderFees            = structs.OrderFees
)

type Interval = consts.Interval

const (
	Interval1min   = consts.Interval1min
	Interval5min   = consts.Interval5min
	Interval15min  = consts.Interval15min
	Interval30min  = consts.Interval30min
	Interval1hour  = consts.Interval1hour
	Interval4hour  = consts.Interval4hour
	Interval6hour  = consts.Interval6hour
	Interval12hour = consts.Interval12hour
	Interval1day   = consts.Interval1day
)

var GetIntervals = consts.GetIntervals

// events
type (
	TradeEventPrivate = workers.TradeEventPrivate
	OrderEvent        = workers.OrderEvent
	CandleEvent       = workers.CandleEvent
)

type CandleData = workers.CandleData

const PairStatusTrading = consts.PairDefaultStatus
