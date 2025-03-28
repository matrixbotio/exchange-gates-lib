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
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
)

const (
	wsConnTimeout     = time.Second * 15
	gateCandleChannel = gate.ChannelSpotCandleStick
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

func (w *GateCandleWorker) getCandleCallback(
	pairSymbol string,
	interval consts.Interval,
	eventCallback func(event workers.CandleEvent),
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

		var event gate.SpotCandleUpdateMsg
		if err := json.Unmarshal(msg.Result, &event); err != nil {
			errorHandler(fmt.Errorf("decode candle: %s", err.Error()))
			return
		}

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

	srv, err := gate.NewWsService(context.Background(), nil, nil)
	if err != nil {
		return fmt.Errorf("conn: %w", err)
	}

	reqPayload := getCandleSubsPayload(gateInterval, pairSymbol)

	// save subscription
	w.CandleWorker.Save(
		getUnsubscriber(srv, reqPayload),
		errorHandler,
		pairSymbol, gateInterval,
	)

	// set event handler
	srv.SetCallBack(
		reqPayload.Channel,
		w.getCandleCallback(
			pairSymbol,
			interval,
			eventCallback,
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
	// TODO
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
