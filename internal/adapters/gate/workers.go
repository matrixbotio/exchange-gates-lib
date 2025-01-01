package gate

import "github.com/matrixbotio/exchange-gates-lib/internal/workers"

func (a *adapter) GetPriceWorker(
	callback workers.PriceEventCallback,
) workers.IPriceWorker {
	// TODO
	return nil
}

func (a *adapter) GetCandleWorker() workers.ICandleWorker {
	// TODO
	return nil
}

func (a *adapter) GetTradeEventsWorker() workers.ITradeEventWorker {
	// TODO
	return nil
}
