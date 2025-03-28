package binanceworkers

import (
	"fmt"

	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/wrapper"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	iWorkers "github.com/matrixbotio/exchange-gates-lib/internal/workers"
)

const TradeSubscriptionKey = "subsctiption"

// TradeEventWorkerBinance - TradeEventWorker for binance
type TradeEventWorkerBinance struct {
	iWorkers.TradeEventWorker
	binanceAPI wrapper.BinanceAPIWrapper
}

func NewTradeEventsWorker(binanceAPI wrapper.BinanceAPIWrapper) *TradeEventWorkerBinance {
	w := &TradeEventWorkerBinance{
		binanceAPI: binanceAPI,
	}
	w.ExchangeTag = consts.BinanceAdapterTag
	return w
}

func (w *TradeEventWorkerBinance) SubscribeToTradeEventsPrivate(
	eventCallback iWorkers.TradeEventPrivateCallback,
	errorHandler func(err error),
) error {
	if w.TradeEventWorker.IsSubscriptionExists(TradeSubscriptionKey) {
		return nil
	}

	wsDone, wsStop, err := w.binanceAPI.SubscribeToTradeEventsPrivate(
		w.ExchangeTag,
		eventCallback,
		errorHandler,
	)
	if err != nil {
		return fmt.Errorf("subscribe to trade events: %w", err)
	}

	w.TradeEventWorker.Save(
		iWorkers.CreateChannelsUnsubscriber(wsDone, wsStop),
		errorHandler,
		TradeSubscriptionKey,
	)
	return nil
}
