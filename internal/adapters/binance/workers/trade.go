package workers

import (
	"fmt"

	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/wrapper"
	iWorkers "github.com/matrixbotio/exchange-gates-lib/internal/workers"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
)

// TradeEventWorkerBinance - TradeEventWorker for binance
type TradeEventWorkerBinance struct {
	iWorkers.TradeEventWorker
	binanceAPI wrapper.BinanceAPIWrapper
}

func NewTradeEventsWorker(
	exchangeTag string,
	binanceAPI wrapper.BinanceAPIWrapper,
) iWorkers.ITradeEventWorker {
	w := TradeEventWorkerBinance{
		binanceAPI: binanceAPI,
	}
	w.ExchangeTag = exchangeTag
	return &w
}

// SubscribeToTradeEvents - websocket subscription to change trade candles on the exchange
func (w *TradeEventWorkerBinance) SubscribeToTradeEvents(
	pairSymbol string,
	eventCallback iWorkers.TradeEventCallback,
	errorHandler func(err error),
) error {

	var err error
	w.WsChannels = new(pkgStructs.WorkerChannels)
	w.WsChannels.WsDone, w.WsChannels.WsStop, err = w.binanceAPI.SubscribeToTradeEvents(
		pairSymbol,
		w.GetExchangeTag(),
		eventCallback,
		errorHandler,
	)
	if err != nil {
		return fmt.Errorf("subscribe to trade events: %w", err)
	}
	return nil
}

func (w *TradeEventWorkerBinance) SubscribeToTradeEventsPrivate(
	eventCallback iWorkers.TradeEventPrivateCallback,
	errorHandler func(err error),
) error {
	var err error
	w.WsChannels = new(pkgStructs.WorkerChannels)

	w.WsChannels.WsDone, w.WsChannels.WsStop, err = w.binanceAPI.SubscribeToTradeEventsPrivate(
		w.ExchangeTag,
		eventCallback,
		errorHandler,
	)
	if err != nil {
		return fmt.Errorf("subscribe to trade events: %w", err)
	}

	return nil
}
