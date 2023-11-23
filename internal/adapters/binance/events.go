package binance

import (
	"fmt"

	"github.com/adshao/go-binance/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
)

// TBD: https://github.com/matrixbotio/exchange-gates-lib/issues/196
// PriceWorkerBinance - MarketDataWorker for binance
type PriceWorkerBinance struct {
	workers.PriceWorker
	binanceAPI BinanceAPIWrapper
}

// TradeEventWorkerBinance - TradeEventWorker for binance
type TradeEventWorkerBinance struct {
	workers.TradeEventWorker
	binanceAPI BinanceAPIWrapper
}

func (a *adapter) GetPriceWorker(callback workers.PriceEventCallback) workers.IPriceWorker {
	w := PriceWorkerBinance{
		binanceAPI: a.binanceAPI,
	}
	w.PriceWorker.ExchangeTag = a.Tag
	w.PriceWorker.HandleEventCallback = callback
	return &w
}

func (w *PriceWorkerBinance) handlePriceEvent(event *binance.WsBookTickerEvent) {
	if event == nil {
		return
	}

	ask, bid, err := mappers.ConvertPriceEvent(*event)
	if err != nil {
		return // ignore broken price event
	}

	w.HandleEventCallback(workers.PriceEvent{
		ExchangeTag: w.ExchangeTag,
		Symbol:      event.Symbol,
		Ask:         ask,
		Bid:         bid,
	})
}

func (w *PriceWorkerBinance) SubscribeToPriceEvents(
	pairSymbols []string,
	errorHandler func(err error),
) (map[string]pkgStructs.WorkerChannels, error) {
	result := map[string]pkgStructs.WorkerChannels{}
	w.WsChannels = new(pkgStructs.WorkerChannels)

	var err error
	for _, pairSymbol := range pairSymbols {
		newChannels := pkgStructs.WorkerChannels{}
		newChannels.WsDone, newChannels.WsStop, err = w.binanceAPI.SubscribeToPriceEvents(
			pairSymbol,
			w.handlePriceEvent,
			errorHandler,
		)
		if err != nil {
			return result, fmt.Errorf("subscribe to %q price: %w", pairSymbol, err)
		}

		result[pairSymbol] = newChannels
	}

	return result, nil
}

func (a *adapter) GetTradeEventsWorker() workers.ITradeEventWorker {
	w := TradeEventWorkerBinance{
		binanceAPI: a.binanceAPI,
	}
	w.ExchangeTag = a.GetTag()
	return &w
}

// SubscribeToTradeEvents - websocket subscription to change trade candles on the exchange
func (w *TradeEventWorkerBinance) SubscribeToTradeEvents(
	pairSymbol string,
	eventCallback func(event workers.TradeEvent),
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
