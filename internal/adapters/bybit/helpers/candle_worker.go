package helpers

import (
	"context"
	"fmt"

	"github.com/hirokisan/bybit/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bybit/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
)

type CandleEventWorkerBybit struct {
	workers.CandleWorker
	WsClient *bybit.WebSocketClient
}

func (w *CandleEventWorkerBybit) SubscribeToCandle(
	pairSymbol string,
	eventCallback func(event workers.CandleEvent),
	errorHandler func(err error),
) error {
	wsSrv, err := w.WsClient.V5().Public(bybit.CategoryV5Spot)
	if err != nil {
		return fmt.Errorf("create trade events subscription service: %w", err)
	}

	eventHandler := CandleEventsHandler{
		pairSymbol: pairSymbol,
		callback:   eventCallback,
	}

	bybitInterval, isExists := mappers.CandleIntervalsToBybit[string(defaultCandleInterval)]
	if !isExists {
		return fmt.Errorf("interval %q not available", defaultCandleInterval)
	}

	if _, err := wsSrv.SubscribeKline(bybit.V5WebsocketPublicKlineParamKey{
		Interval: bybit.Interval(bybitInterval.Code),
		Symbol:   bybit.SymbolV5(pairSymbol),
	}, eventHandler.handle); err != nil {
		return fmt.Errorf("open candle events subscription: %w", err)
	}

	wsErrHandler := func(isWebsocketClosed bool, err error) {
		// TBD: handle reconnect: https://github.com/matrixbotio/exchange-gates-lib/issues/154
		errorHandler(err)
	}

	go func() {
		if err := wsSrv.Start(context.Background(), wsErrHandler); err != nil {
			errorHandler(fmt.Errorf(
				"start candle %q events subscription: %w",
				pairSymbol, err,
			))
		}
	}()
	return nil
}

func (w *CandleEventWorkerBybit) SubscribeToCandlesList(
	intervalsPerPair map[string]string,
	eventCallback func(event workers.CandleEvent),
	errorHandler func(err error),
) error {
	for pairSymbol := range intervalsPerPair {
		if err := w.SubscribeToCandle(pairSymbol, eventCallback, errorHandler); err != nil {
			return fmt.Errorf(
				"subscribe to %q candle events: %w",
				pairSymbol, err,
			)
		}
	}
	return nil
}
