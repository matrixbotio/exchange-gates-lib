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

const pingTimeout = time.Second * 20
const tradeSubscriptionKey = "subscription"

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

// DEPRECATED
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

	pingActive := true

	service, err := w.wsClient.V5().Private()
	if err != nil {
		return fmt.Errorf("failed to get private subscription service: %w", err)
	}

	if err := service.Subscribe(); err != nil {
		return fmt.Errorf("init service: %w", err)
	}

	handler := func(e bybit.V5WebsocketPrivateExecutionResponse) error {
		for _, eventRaw := range e.Data {
			event, err := mappers.ParseTradeEventPrivate(eventRaw, e.CreationTime, w.ExchangeTag)
			if err != nil {
				return fmt.Errorf("parse trade event: %w", err)
			}

			eventCallback(event)
		}

		return nil
	}

	unsubscribe, err := service.SubscribeExecution(handler)
	if err != nil {
		return fmt.Errorf("subscribe to trade events: %w", err)
	}

	// set unsibscriber & save subscription data
	w.TradeEventWorker.Save(
		nil,
		errorHandler,
		tradeSubscriptionKey,
	)

	go func() {
		for pingActive {
			if err := service.Ping(); err != nil {
				errorHandler(fmt.Errorf("ping: %w", err))
			}

			time.Sleep(pingTimeout)
		}
	}()

	go func() {
		select {
		case <-w.WsChannels.WsStop:
			if err := unsubscribe(); err != nil {
				errorHandler(fmt.Errorf("unsubscribe from trade events: %w", err))
			}
		case <-w.WsChannels.WsDone:
		}

		pingActive = false
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
			wsErrHandler(false, fmt.Errorf("start trade events subscriber: %w", err))
		}
	}()

	return nil
}

// DEPRECATED
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
