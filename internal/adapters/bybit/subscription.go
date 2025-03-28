package bybit

import (
	"fmt"

	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bybit/helpers"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bybit/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
)

func (a *adapter) CreateCandleWorker() *helpers.CandleEventWorkerBybit {
	w := &helpers.CandleEventWorkerBybit{
		WsClient: a.wsClient,
	}
	w.CandleWorker.ExchangeTag = a.GetTag()
	return w
}

func (a *adapter) CreateTradeEventsWorker() *TradeEventWorkerBybit {
	w := &TradeEventWorkerBybit{
		wsClient: a.wsClient,
	}
	w.TradeEventWorker.ExchangeTag = a.GetTag()
	return w
}

func (a *adapter) SubscribeCandle(
	pairSymbol string,
	interval consts.Interval,
	eventCallback func(event workers.CandleEvent),
	errorHandler func(err error),
) error {
	return a.candleWorker.SubscribeToCandle(
		pairSymbol,
		interval,
		eventCallback,
		errorHandler,
	)
}

func (a *adapter) SubscribeAccountTrades(
	eventCallback workers.TradeEventPrivateCallback,
	errorHandler func(err error),
) error {
	return a.tradeWorker.SubscribeToTradeEventsPrivate(
		eventCallback, errorHandler,
	)
}

func (a *adapter) UnsubscribeCandle(
	pairSymbol string,
	interval consts.Interval,
) {
	bybitInterval, isExists := mappers.CandleIntervalsToBybit[interval]
	if !isExists {
		fmt.Printf("convert interval %q to bingx\n", interval)
		return
	}

	a.candleWorker.Unsubscribe(pairSymbol, bybitInterval.Code)
}

func (a *adapter) UnsubscribeAccountTrades() {
	a.tradeWorker.UnsubscribeAll()
}
