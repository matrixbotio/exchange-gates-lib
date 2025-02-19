package helpers

import (
	"context"
	"fmt"

	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"

	"github.com/hirokisan/bybit/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bybit/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
)

type CandleEventWorkerBybit struct {
	workers.CandleWorker
	WsClient *bybit.WebSocketClient
}

func (w *CandleEventWorkerBybit) subscribe(
	intervalsPerPair map[string]consts.Interval,
	eventCallback func(event workers.CandleEvent),
	errorHandler func(err error),
) error {
	w.WsChannels = new(pkgStructs.WorkerChannels)
	w.WsChannels.WsStop = make(chan struct{}, 1)
	w.WsChannels.WsDone = make(chan struct{}, 1)

	wsSrv, err := w.WsClient.V5().Public(bybit.CategoryV5Spot)
	if err != nil {
		return fmt.Errorf("create candle events subscription service: %w", err)
	}

	eventHandler := CandleEventsHandler{
		symbols:  make(symbolPerTopic),
		callback: eventCallback,
	}

	var keys []bybit.V5WebsocketPublicKlineParamKey
	for pairSymbol, interval := range intervalsPerPair {
		bybitInterval, isExists := mappers.CandleIntervalsToBybit[interval]
		if !isExists {
			return fmt.Errorf("interval %q not available", interval)
		}

		key := bybit.V5WebsocketPublicKlineParamKey{
			Interval: bybit.Interval(bybitInterval.Code),
			Symbol:   bybit.SymbolV5(pairSymbol),
		}

		eventHandler.symbols[key.Topic()] = symbolData{
			Symbol:   pairSymbol,
			Interval: interval,
		}
		keys = append(keys, key)
	}

	unsubscribe, err := wsSrv.SubscribeKlines(keys, eventHandler.handle)
	if err != nil {
		return fmt.Errorf("open candle events subscription: %w", err)
	}

	go func() {
		select {
		case <-w.WsChannels.WsStop:
			if err := unsubscribe(); err != nil {
				errorHandler(fmt.Errorf("unsubscribe from ticker events: %w", err))
			}
		case <-w.WsChannels.WsDone:
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
			w.WsChannels.WsDone <- struct{}{}

			errorHandler(fmt.Errorf(
				"start candle events subscription: %w",
				err,
			))
		}
	}()

	return nil
}

func (w *CandleEventWorkerBybit) SubscribeToCandle(
	pairSymbol string,
	interval consts.Interval,
	eventCallback func(event workers.CandleEvent),
	errorHandler func(err error),
) error {
	return w.subscribe(map[string]consts.Interval{
		pairSymbol: interval,
	}, eventCallback, errorHandler)
}

func (w *CandleEventWorkerBybit) SubscribeToCandlesList(
	intervalsPerPair map[string]consts.Interval,
	eventCallback func(event workers.CandleEvent),
	errorHandler func(err error),
) error {
	return w.subscribe(intervalsPerPair, eventCallback, errorHandler)
}
