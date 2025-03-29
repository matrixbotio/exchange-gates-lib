package bingx

import (
	"fmt"

	bingxgo "github.com/matrixbotio/go-bingx"

	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bingx/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
	"github.com/matrixbotio/exchange-gates-lib/pkg/structs"
)

const tradeSubscriptionKey = "subscription"

type CandleEventWorkerBingX struct {
	workers.CandleWorker
}

func (a *adapter) CreateCandleWorker() *CandleEventWorkerBingX {
	w := &CandleEventWorkerBingX{}
	w.CandleWorker.ExchangeTag = a.GetTag()
	return w
}

type TradeEventWorkerBingX struct {
	workers.TradeEventWorker
	client *bingxgo.SpotClient
	creds  structs.APICredentials
}

func (a *adapter) CreateTradeEventsWorker() *TradeEventWorkerBingX {
	w := &TradeEventWorkerBingX{client: &a.client, creds: a.creds}
	w.TradeEventWorker.ExchangeTag = a.GetTag()
	return w
}

func (w *TradeEventWorkerBingX) SubscribeToTradeEventsPrivate(
	eventCallback workers.TradeEventPrivateCallback,
	errorHandler func(err error),
) error {
	if w.TradeEventWorker.IsSubscriptionExists(tradeSubscriptionKey) {
		return nil
	}

	wsDone, wsStop, err := bingxgo.WsOrderUpdateServe(
		w.creds.Keypair.Public,
		w.creds.Keypair.Secret,
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

	w.TradeEventWorker.Save(
		workers.CreateChannelsUnsubscriber(wsDone, wsStop),
		errorHandler,
		tradeSubscriptionKey,
	)
	return nil
}

func (w *CandleEventWorkerBingX) SubscribeToCandle(
	pairSymbol string,
	interval consts.Interval,
	eventCallback func(event workers.CandleEvent),
	errorHandler func(err error),
) error {
	bingxInterval, err := ConvertIntervalToBingXWs(interval)
	if err != nil {
		return fmt.Errorf("convert interval: %w", err)
	}

	if w.CandleWorker.IsSubscriptionExists(pairSymbol, string(bingxInterval)) {
		return nil
	}

	wsDone, wsStop, err := bingxgo.WsKlineServe(
		pairSymbol,
		bingxInterval,
		GetBingXCandleEventsHandler(
			eventCallback,
			errorHandler,
		),
		errorHandler,
	)
	if err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}

	w.CandleWorker.Save(
		workers.CreateChannelsUnsubscriber(wsDone, wsStop),
		errorHandler,
		pairSymbol, string(bingxInterval),
	)
	return nil
}

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
	return a.tradeWorker.SubscribeToTradeEventsPrivate(
		eventCallback, errorHandler,
	)
}

func (a *adapter) UnsubscribeCandle(
	pairSymbol string,
	interval consts.Interval,
) {
	bingxInterval, err := ConvertIntervalToBingXWs(interval)
	if err != nil {
		fmt.Printf(
			"convert interval %q to bingx: %s\n",
			interval, err.Error(),
		)
		return
	}

	a.candleWorker.Unsubscribe(pairSymbol, string(bingxInterval))
}

func (a *adapter) UnsubscribeAccountTrades() {
	a.tradeWorker.UnsubscribeAll()
}
