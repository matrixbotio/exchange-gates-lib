package helpers

import (
	"context"
	"fmt"

	"github.com/hirokisan/bybit/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bybit/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
)

type CandleEventWorkerBybit struct {
	workers.CandleWorker
	WsClient *bybit.WebSocketClient
}

func (w *CandleEventWorkerBybit) SubscribeToCandle(
	pairSymbol string,
	interval consts.Interval,
	eventCallback func(event workers.CandleEvent),
	errorHandler func(err error),
) error {
	eventHandler := CandleEventsHandler{
		symbols:  make(symbolPerTopic),
		callback: eventCallback,
	}

	wsStop := make(chan struct{}, 1)
	wsDone := make(chan struct{}, 1)

	bybitInterval, isExists := mappers.CandleIntervalsToBybit[interval]
	if !isExists {
		return fmt.Errorf("interval %q not available", interval)
	}

	if w.CandleWorker.IsSubscriptionExists(pairSymbol, bybitInterval.Code) {
		return nil // already subscribed
	}

	key := bybit.V5WebsocketPublicKlineParamKey{
		Interval: bybit.Interval(bybitInterval.Code),
		Symbol:   bybit.SymbolV5(pairSymbol),
	}

	eventHandler.symbols[key.Topic()] = symbolData{
		Symbol:   pairSymbol,
		Interval: interval,
	}

	w.CandleWorker.Save(
		workers.CreateChannelsUnsubscriber(wsDone, wsStop),
		errorHandler,
		pairSymbol, bybitInterval.Code,
	)

	wsSrv, err := w.WsClient.V5().Public(bybit.CategoryV5Spot)
	if err != nil {
		return fmt.Errorf("create candle events subscription service: %w", err)
	}

	unsubscribe, err := wsSrv.SubscribeKline(key, eventHandler.handle)
	if err != nil {
		return fmt.Errorf("open candle events subscription: %w", err)
	}

	go func() {
		select {
		case <-wsStop:
			if err := unsubscribe(); err != nil {
				errorHandler(fmt.Errorf("unsubscribe from ticker events: %w", err))
			}
		case <-wsDone:
		}
	}()

	wsErrHandler := func(isWebsocketClosed bool, err error) {
		if !isWebsocketClosed {
			_ = wsSrv.Close()
		}

		errorHandler(fmt.Errorf("bybit candles subscription: %w", err))
	}

	go func() {
		if err := wsSrv.Start(context.Background(), wsErrHandler); err != nil {
			wsDone <- struct{}{}

			errorHandler(fmt.Errorf(
				"start candle events subscription: %w",
				err,
			))
		}
	}()

	return nil
}
