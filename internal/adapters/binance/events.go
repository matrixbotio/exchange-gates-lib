package binance

import (
	binanceworkers "github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/workers"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
)

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
	return a.tradeWorker.SubscribeToTradeEventsPrivate(eventCallback, errorHandler)
}

func (a *adapter) UnsubscribeCandle(
	pairSymbol string,
	interval consts.Interval,
) {
	a.candleWorker.Unsubscribe(
		pairSymbol,
		convertInterval(interval),
	)
}

func (a *adapter) UnsubscribeAccountTrades() {
	a.tradeWorker.Unsubscribe(binanceworkers.TradeSubscriptionKey)
}
