package adapter

import (
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
)

type Adapter = adapters.Adapter

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

type (
	ICandleWorker      = workers.ICandleWorker
	IPriceWorker       = workers.IPriceWorker
	ITradeEventWorker  = workers.ITradeEventWorker
	PriceEventCallback = workers.PriceEventCallback
)

// events
type (
	PriceEvent  = workers.PriceEvent
	TradeEvent  = workers.TradeEvent
	CandleEvent = workers.CandleEvent
)

const PairStatusTrading = consts.PairDefaultStatus
