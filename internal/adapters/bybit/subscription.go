package bybit

import (
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bybit/helpers"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
)

func (a *adapter) GetPriceWorker(callback workers.PriceEventCallback) workers.IPriceWorker {
	w := &PriceEventWorkerBybit{
		wsClient: a.wsClient,
	}
	w.PriceWorker.ExchangeTag = a.Tag
	w.PriceWorker.HandleEventCallback = callback
	return w
}

func (a *adapter) GetCandleWorker() workers.ICandleWorker {
	return &helpers.CandleEventWorkerBybit{
		WsClient: a.wsClient,
	}
}

func (a *adapter) GetTradeEventsWorker() workers.ITradeEventWorker {
	return &TradeEventWorkerBybit{
		wsClient: a.wsClient,
	}
}
