package bingx

import (
	"fmt"

	bingxgo "github.com/Sagleft/go-bingx"

	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bingx/helpers"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bingx/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
	"github.com/matrixbotio/exchange-gates-lib/pkg/structs"
)

type PriceEventWorkerBingX struct {
	workers.PriceWorker
}

func (a *adapter) GetPriceWorker(callback workers.PriceEventCallback) workers.IPriceWorker {
	w := &PriceEventWorkerBingX{}
	w.PriceWorker.ExchangeTag = a.Tag
	w.PriceWorker.HandleEventCallback = callback
	return w
}

type CandleEventWorkerBingX struct {
	workers.CandleWorker
}

func (a *adapter) GetCandleWorker() workers.ICandleWorker {
	w := &CandleEventWorkerBingX{}
	w.CandleWorker.ExchangeTag = a.GetTag()
	return w
}

type TradeEventWorkerBingX struct {
	workers.TradeEventWorker
	client *bingxgo.SpotClient
}

func (a *adapter) GetTradeEventsWorker() workers.ITradeEventWorker {
	w := &TradeEventWorkerBingX{client: &a.client}
	w.TradeEventWorker.ExchangeTag = a.GetTag()
	w.TradeEventWorker.WsChannels = new(structs.WorkerChannels)
	return w
}

func (w *TradeEventWorkerBingX) SubscribeToTradeEventsPrivate(
	eventCallback workers.TradeEventPrivateCallback,
	errorHandler func(err error),
) error {
	var err error
	w.WsChannels = new(structs.WorkerChannels)
	w.WsChannels.WsDone, w.WsChannels.WsStop, err = w.client.WsOrderUpdateServe(
		func(o *bingxgo.WsOrder) {
			event, err := mappers.ConvertOrderEvent(o)
			if err != nil {
				errorHandler(fmt.Errorf("convert: %w", err))
				return
			}

			eventCallback(event)
		},
		errorHandler,
	)
	if err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}
	return nil
}

func (w *PriceEventWorkerBingX) SubscribeToPriceEvents(
	pairSymbols []string,
	errorHandler func(err error),
) (map[string]structs.WorkerChannels, error) {
	// not implemented
	return nil, nil
}

func (w *CandleEventWorkerBingX) SubscribeToCandle(
	pairSymbol string,
	eventCallback func(event workers.CandleEvent),
	errorHandler func(err error),
) error {
	var err error
	w.WsChannels = new(structs.WorkerChannels)
	w.WsChannels.WsDone, w.WsChannels.WsStop, err = bingxgo.WsKlineServe(
		pairSymbol,
		bingxgo.Interval1,
		helpers.GetBingXCandleEventsHandler(
			eventCallback,
			errorHandler,
		),
		errorHandler,
	)
	if err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}
	return nil
}
