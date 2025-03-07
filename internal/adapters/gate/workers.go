package gate

import (
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
	"github.com/matrixbotio/exchange-gates-lib/pkg/structs"
)

type GatePriceWorker struct {
	workers.PriceWorker
}

func (a *adapter) GetPriceWorker(
	callback workers.PriceEventCallback,
) workers.IPriceWorker {
	w := &GatePriceWorker{}
	w.ExchangeTag = a.GetTag()
	return w
}

func (w *GatePriceWorker) SubscribeToPriceEvents(
	pairSymbols []string,
	errorHandler func(err error),
) (map[string]structs.WorkerChannels, error) {
	// TODO
	return nil, nil
}

type GateCandleWorker struct {
	workers.CandleWorker
}

func (a *adapter) GetCandleWorker() workers.ICandleWorker {
	w := &GateCandleWorker{}
	w.ExchangeTag = a.GetTag()
	return w
}

func (w *GateCandleWorker) SubscribeToCandle(
	pairSymbol string,
	interval consts.Interval,
	eventCallback func(event workers.CandleEvent),
	errorHandler func(err error),
) error {
	// TODO
	return nil
}

func (w *GateCandleWorker) SubscribeToCandlesList(
	intervalsPerPair map[string]consts.Interval,
	eventCallback func(event workers.CandleEvent),
	errorHandler func(err error),
) error {
	// TODO
	return nil
}

func (a *adapter) GetTradeEventsWorker() workers.ITradeEventWorker {
	// TODO
	return nil
}
