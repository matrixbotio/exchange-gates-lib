package workers

import (
	"fmt"

	"github.com/adshao/go-binance/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/wrapper"
	iWorkers "github.com/matrixbotio/exchange-gates-lib/internal/workers"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
)

// PriceWorkerBinance - MarketDataWorker for binance
type PriceWorkerBinance struct {
	iWorkers.PriceWorker
	binanceAPI wrapper.BinanceAPIWrapper
}

func NewPriceWorker(
	exchangeTag string,
	binanceAPI wrapper.BinanceAPIWrapper,
	callback iWorkers.PriceEventCallback,
) iWorkers.IPriceWorker {
	w := PriceWorkerBinance{
		binanceAPI: binanceAPI,
	}
	w.PriceWorker.ExchangeTag = exchangeTag
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

	w.HandleEventCallback(iWorkers.PriceEvent{
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
			// TBD: add reconnect
			// https://github.com/matrixbotio/exchange-gates-lib/issues/199
			return result, fmt.Errorf("subscribe to %q price: %w", pairSymbol, err)
		}

		result[pairSymbol] = newChannels
	}

	return result, nil
}
