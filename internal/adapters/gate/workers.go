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

type GatePriceWorker struct {
	workers.PriceWorker
}

type GateCandleWorker struct {
	workers.CandleWorker

	creds      pkgStructs.APICredentials
	errHandler func(error)
}

type GateTradeWorker struct {
	workers.TradeEventWorker
}

func (a *adapter) GetPriceWorker(
	callback workers.PriceEventCallback,
) workers.IPriceWorker {
	w := &GatePriceWorker{}
	w.ExchangeTag = a.GetTag()
	return w
}

func (w *GatePriceWorker) SubscribeToPriceEvents(
	pairSymbols []string,
	errorHandler func(err error),
) (map[string]pkgStructs.WorkerChannels, error) {
	// not implemented yet
	return nil, nil
}

func (a *adapter) GetCandleWorker() workers.ICandleWorker {
	w := &GateCandleWorker{creds: a.creds}
	w.ExchangeTag = a.GetTag()
	return w
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

func (w *GateCandleWorker) Stop() {
	w.CandleWorker.Stop()
}

func (w *GateCandleWorker) SubscribeToCandle(
	pairSymbol string,
	interval consts.Interval,
	eventCallback func(event workers.CandleEvent),
	errorHandler func(err error),
) error {
	if errorHandler != nil {
		w.errHandler = errorHandler
	}

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

func (w *GateCandleWorker) SubscribeToCandlesList(
	intervalsPerPair map[string]consts.Interval,
	eventCallback func(event workers.CandleEvent),
	errorHandler func(err error),
) error {
	if errorHandler != nil {
		w.errHandler = errorHandler
	}

	for symbol, interval := range intervalsPerPair {
		if err := w.SubscribeToCandle(
			symbol, interval,
			eventCallback, errorHandler,
		); err != nil {
			return fmt.Errorf("subscribe to %q: %w", symbol, err)
		}
	}
	return nil
}

func (a *adapter) GetTradeEventsWorker() workers.ITradeEventWorker {
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
