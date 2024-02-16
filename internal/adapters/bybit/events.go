package bybit

import (
	"context"
	"fmt"

	"github.com/hirokisan/bybit/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bybit/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
)

// PriceEventWorkerBybit :
type PriceEventWorkerBybit struct {
	workers.PriceWorker
	wsClient *bybit.WebSocketClient
}

// TradeEventWorkerBybit :
type TradeEventWorkerBybit struct {
	workers.TradeEventWorker
	wsClient *bybit.WebSocketClient
}

func (w *TradeEventWorkerBybit) SubscribeToTradeEvents(
	_ string,
	_ workers.TradeEventCallback,
	_ func(err error),
) error {
	// TBD: https://github.com/matrixbotio/exchange-gates-lib/issues/153
	return nil
}

func (w *TradeEventWorkerBybit) SubscribeToTradeEventsPrivate(
	eventCallback workers.TradeEventPrivateCallback,
	errorHandler func(err error),
) error {
	w.WsChannels = new(pkgStructs.WorkerChannels)
	w.WsChannels.WsStop = make(chan struct{}, 1)
	w.WsChannels.WsDone = make(chan struct{}, 1)

	service, err := w.wsClient.V5().Private()
	if err != nil {
		return fmt.Errorf("failed to get private subscription service: %w", err)
	}

	handler := func(e bybit.V5WebsocketPrivateOrderResponse) error {
		for _, datum := range e.Data {
			event, err := mappers.ParseTradeEvent(datum, e.CreationTime, w.ExchangeTag)
			if err != nil {
				return fmt.Errorf("parse trade event: %w", err)
			}

			eventCallback(event)
		}

		return nil
	}

	unsubscribe, err := service.SubscribeOrder(handler)
	if err != nil {
		return fmt.Errorf("subscribe to trade events: %w", err)
	}

	go func() {
		select {
		case <-w.WsChannels.WsStop:
			if err := unsubscribe(); err != nil {
				errorHandler(fmt.Errorf("unsubscribe from trade events: %w", err))
			}
		case <-w.WsChannels.WsDone:
		}
	}()

	wsErrHandler := func(isWebsocketClosed bool, wsErr error) {
		if !isWebsocketClosed {
			_ = service.Close()
		}

		w.WsChannels.WsDone <- struct{}{}

		errorHandler(fmt.Errorf("trade events subscription: %w", wsErr))
	}

	go func() {
		if err := service.Start(context.Background(), wsErrHandler); err != nil {
			w.WsChannels.WsDone <- struct{}{}
			errorHandler(fmt.Errorf("start trade events subscriber: %w", err))
		}
	}()

	return nil
}

func (w *PriceEventWorkerBybit) SubscribeToPriceEvents(
	pairSymbols []string,
	errorHandler func(err error),
) (map[string]pkgStructs.WorkerChannels, error) {
	result := map[string]pkgStructs.WorkerChannels{}

	w.WsChannels = new(pkgStructs.WorkerChannels)

	for _, pairSymbol := range pairSymbols {
		newChannels := pkgStructs.WorkerChannels{}
		newChannels.WsStop = make(chan struct{}, 1)
		newChannels.WsDone = make(chan struct{}, 1)

		wsSrv, err := w.wsClient.V5().Public(bybit.CategoryV5Spot)
		if err != nil {
			return nil, fmt.Errorf("create ticker service: %w", err)
		}

		eventHandler := func(r bybit.V5WebsocketPublicTickerResponse) error {
			event, err := mappers.ParsePriceEvent(r, w.ExchangeTag)
			if err != nil {
				return fmt.Errorf("parse price event: %w", err)
			}

			w.HandleEventCallback(event)
			return nil
		}

		unsubscribe, err := wsSrv.SubscribeTicker(bybit.V5WebsocketPublicTickerParamKey{
			Symbol: bybit.SymbolV5(pairSymbol),
		}, eventHandler)
		if err != nil {
			return nil, fmt.Errorf("create ticker subscriber: %w", err)
		}

		wsErrHandler := func(isWebsocketClosed bool, wsErr error) {
			if !isWebsocketClosed {
				_ = wsSrv.Close()
			}

			newChannels.WsDone <- struct{}{}

			errorHandler(fmt.Errorf("ticker subscription: %w", wsErr))
		}

		go func() {
			select {
			case <-newChannels.WsStop:
				if err := unsubscribe(); err != nil {
					errorHandler(fmt.Errorf("unsubscribe from ticker events: %w", err))
				}
			case <-newChannels.WsDone:
			}
		}()

		go func() {
			if err := wsSrv.Start(context.Background(), wsErrHandler); err != nil {
				newChannels.WsDone <- struct{}{}
				errorHandler(fmt.Errorf("start ticker subscriber: %w", err))
			}
		}()

		result[pairSymbol] = newChannels
	}

	return result, nil
}
