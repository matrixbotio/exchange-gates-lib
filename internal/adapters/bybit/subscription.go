package bybit

import (
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bybit/helpers"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
	"github.com/matrixbotio/exchange-gates-lib/pkg/structs"
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
	w := &TradeEventWorkerBybit{
		wsClient: a.wsClient,
	}
	w.TradeEventWorker.WsChannels = new(structs.WorkerChannels)
	return w
}
