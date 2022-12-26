package helpers

import (
	"context"
	"fmt"

	"github.com/hirokisan/bybit/v2"
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

	if _, err := wsSrv.SubscribeKline(bybit.V5WebsocketPublicKlineParamKey{
		Interval: defaultCandleInterval,
		Symbol:   bybit.SymbolV5(pairSymbol),
	}, eventHandler.handle); err != nil {
		return fmt.Errorf("open candle events subscription: %w", err)
	}

	wsErrHandler := func(isWebsocketClosed bool, err error) {
		// TBD: handle reconnect: https://github.com/matrixbotio/exchange-gates-lib/issues/154
		errorHandler(err)
	}

	if err := wsSrv.Start(context.Background(), wsErrHandler); err != nil {
		return fmt.Errorf("start candle events subscription: %w", err)
	}
	return nil
}
