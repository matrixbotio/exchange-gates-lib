package bybit

import (
	"context"
	"fmt"
	"time"

	"github.com/hirokisan/bybit/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bybit/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
)

const pingTimeout = time.Second * 20
const tradeSubscriptionKey = "subscription"

// TradeEventWorkerBybit :
type TradeEventWorkerBybit struct {
	workers.TradeEventWorker
	wsClient *bybit.WebSocketClient
}

func (w *TradeEventWorkerBybit) SubscribeToTradeEventsPrivate(
	eventCallback workers.TradeEventPrivateCallback,
	errorHandler func(err error),
) error {
	if w.TradeEventWorker.IsSubscriptionExists(tradeSubscriptionKey) {
		return nil
	}

	wsStop := make(chan struct{}, 1)
	wsDone := make(chan struct{}, 1)

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
		workers.CreateChannelsUnsubscriber(wsDone, wsStop),
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
		case <-wsStop:
			if err := unsubscribe(); err != nil {
				errorHandler(fmt.Errorf("unsubscribe from trade events: %w", err))
			}
		case <-wsDone:
		}

		pingActive = false
	}()

	wsErrHandler := func(isWebsocketClosed bool, wsErr error) {
		if !isWebsocketClosed {
			_ = service.Close()
		}

		wsDone <- struct{}{}

		errorHandler(fmt.Errorf("trade events subscription: %w", wsErr))
	}

	go func() {
		if err := service.Start(context.Background(), wsErrHandler); err != nil {
			wsErrHandler(false, fmt.Errorf("start trade events subscriber: %w", err))
		}
	}()

	return nil
}
