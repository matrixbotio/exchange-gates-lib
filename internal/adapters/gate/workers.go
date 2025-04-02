package gate

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	gate "github.com/gateio/gatews/go"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/gate/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
	"github.com/matrixbotio/exchange-gates-lib/pkg/errs"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
)

const (
	wsConnTimeout        = time.Second * 15
	gateCandleChannel    = gate.ChannelSpotCandleStick
	gateTradeChannel     = "spot.usertrades_v2"
	tradeSubscriptionTag = "subscription"
)

type GateCandleWorker struct {
	workers.CandleWorker
}

type GateTradeWorker struct {
	workers.TradeEventWorker
	creds pkgStructs.APICredentials
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
	return a.tradeWorker.SubscribeToTradeEventsPrivate(eventCallback, errorHandler)
}

func getRawEventHandler[rawEventType any](
	eventCallback func(event rawEventType),
	errorHandler func(err error),
) gate.CallBack {
	return gate.NewCallBack(func(msg *gate.UpdateMsg) {
		if msg == nil {
			return
		}

		if msg.Error != nil {
			errorHandler(msg.Error)
			return
		}

		if msg.Data.Errs != nil {
			errorHandler(fmt.Errorf(
				"on candle: label: %s, message: %s",
				msg.Data.Errs.Label,
				msg.Data.Errs.Message,
			))
			return
		}

		if msg.Event != "update" {
			return
		}

		var event rawEventType
		if err := json.Unmarshal(msg.Result, &event); err != nil {
			errorHandler(fmt.Errorf("decode: %s", err.Error()))
			return
		}

		eventCallback(event)
	})
}

func (w *GateCandleWorker) SubscribeToCandle(
	pairSymbol string,
	interval consts.Interval,
	eventCallback func(event workers.CandleEvent),
	errorHandler func(err error),
) error {
	gateInterval, err := mappers.ConvertIntervalToGate(interval)
	if err != nil {
		return fmt.Errorf("convert interval: %w", err)
	}

	if w.CandleWorker.IsSubscriptionExists(pairSymbol, gateInterval) {
		return nil // already subscribed
	}

	// setup new ws connection
	srv, err := gate.NewWsService(context.Background(), nil, nil)
	if err != nil {
		return fmt.Errorf("conn: %w", err)
	}

	reqPayload := getCandleSubsPayload(gateInterval, pairSymbol)

	eventHandler := func(event gate.SpotCandleUpdateMsg) {
		eventParsed, err := mappers.ParseCandleEvent(
			event,
			pairSymbol,
			interval,
		)
		if err != nil {
			errorHandler(fmt.Errorf("parse candle: %s", err.Error()))
			return
		}

		eventCallback(eventParsed)
	}

	// set event handler
	srv.SetCallBack(
		reqPayload.Channel,
		getRawEventHandler(
			eventHandler,
			errorHandler,
		),
	)

	// subscribe
	go func() {
		if err := srv.Subscribe(
			reqPayload.Channel,
			reqPayload.Payload,
		); err != nil {
			errorHandler(fmt.Errorf("subscribe: %w", err))
		}
	}()

	// save subscription
	w.CandleWorker.Save(
		getUnsubscriber(srv, reqPayload),
		errorHandler,
		pairSymbol, gateInterval,
	)
	return nil
}

func (a *adapter) CreateTradeEventsWorker() *GateTradeWorker {
	w := &GateTradeWorker{}
	w.ExchangeTag = a.GetTag()
	return w
}

func (w *GateTradeWorker) SubscribeToTradeEventsPrivate(
	eventCallback workers.TradeEventPrivateCallback,
	errorHandler func(err error),
) error {
	if !w.creds.Keypair.IsSet() {
		return errs.ErrAPIKeyNotSet
	}

	if w.TradeEventWorker.IsSubscriptionExists(tradeSubscriptionTag) {
		return nil // already subscribed
	}

	cfg := gate.NewConnConfFromOption(&gate.ConfOptions{
		App:    "spot",
		Key:    w.creds.Keypair.Public,
		Secret: w.creds.Keypair.Secret,
	})

	srv, err := gate.NewWsService(context.Background(), nil, cfg)
	if err != nil {
		return fmt.Errorf("conn: %w", err)
	}

	reqPayload := getOrderSubsPayload()

	eventHandler := func(events []gate.SpotUserTradesMsg) {
		for _, event := range events {
			eventParsed, err := mappers.ParseOrderEvent(event)
			if err != nil {
				errorHandler(fmt.Errorf("parse order event: %s", err.Error()))
				return
			}

			eventCallback(eventParsed)
		}
	}

	// set event handler
	srv.SetCallBack(
		reqPayload.Channel,
		getRawEventHandler(
			eventHandler,
			errorHandler,
		),
	)

	// subscribe
	go func() {
		if err := srv.Subscribe(
			reqPayload.Channel,
			reqPayload.Payload,
		); err != nil {
			errorHandler(fmt.Errorf("subscribe: %w", err))
		}
	}()

	// save subscription
	w.TradeEventWorker.Save(
		getUnsubscriber(srv, reqPayload),
		errorHandler,
		tradeSubscriptionTag,
	)
	return nil
}

func (a *adapter) UnsubscribeCandle(
	pairSymbol string,
	interval consts.Interval,
) {
	gateInterval, err := mappers.ConvertIntervalToGate(interval)
	if err != nil {
		fmt.Printf("convert interval %q to gate\n", interval)
		return
	}

	a.candleWorker.Unsubscribe(pairSymbol, gateInterval)
}

func (a *adapter) UnsubscribeAccountTrades() {
	a.tradeWorker.UnsubscribeAll()
}
