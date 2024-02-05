package bybit

import (
	"context"
	"fmt"
	"time"

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
		return fmt.Errorf("subscribe to order events: %w", err)
	}

	w.WsChannels.WsStop = make(chan struct{}, 1)
	go func() {
		for {
			select {
			case <-time.After(20 * time.Second):
				if err := service.Ping(); err != nil {
					// TBD: handle reconnect: https://github.com/matrixbotio/exchange-gates-lib/issues/154
					errorHandler(fmt.Errorf("trade events ping: %w", err))
				}
			case <-w.WsChannels.WsStop:
				if err := unsubscribe(); err != nil {
					errorHandler(fmt.Errorf("unsubscribe from ticker events: %w", err))
				}
			}
		}
	}()

	wsErrHandler := func(isWebsocketClosed bool, wsErr error) {
		// TBD: handle reconnect: https://github.com/matrixbotio/exchange-gates-lib/issues/154
		errorHandler(fmt.Errorf("trade events subscription: %w", wsErr))
	}

	if err := service.Start(context.Background(), wsErrHandler); err != nil {
		return fmt.Errorf("start trade events subscriber: %w", err)
	}

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
			// TBD: handle reconnect: https://github.com/matrixbotio/exchange-gates-lib/issues/154
			errorHandler(fmt.Errorf("ticker subscription: %w", wsErr))
		}

		newChannels.WsStop = make(chan struct{}, 1)
		go func() {
			for {
				select {
				case <-time.After(20 * time.Second):
					if err := wsSrv.Ping(); err != nil {
						// TBD: handle reconnect: https://github.com/matrixbotio/exchange-gates-lib/issues/154
						errorHandler(fmt.Errorf("ticker ping: %w", err))
					}
				case <-newChannels.WsStop:
					if err := unsubscribe(); err != nil {
						errorHandler(fmt.Errorf("unsubscribe from ticker events: %w", err))
					}
				}
			}
		}()

		if err := wsSrv.Start(context.Background(), wsErrHandler); err != nil {
			return nil, fmt.Errorf("start ticker subscriber: %w", err)
		}

		result[pairSymbol] = newChannels
	}

	return result, nil
}
