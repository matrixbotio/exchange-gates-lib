package gate

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
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
	srvs       sync.Map // symbol -> subscriptionData
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

	w.srvs.Range(func(iSymbol, iSub any) bool {
		if iSub == nil || iSymbol == nil {
			return true
		}

		pairSymbol, isConvertable := iSymbol.(string)
		if !isConvertable {
			if w.errHandler != nil {
				w.errHandler(fmt.Errorf(
					"unsubscribe: get symbol: unknown format: %s",
					reflect.ValueOf(iSymbol).String(),
				))
			}
			return true
		}

		// get subscription data
		subsData, isConvertable := iSub.(workers.SubscriptionData)
		if !isConvertable {
			if w.errHandler != nil {
				w.errHandler(fmt.Errorf(
					"unsubscribe: get subs data: unknown format: %s",
					reflect.ValueOf(iSub).String(),
				))
			}
			return true
		}

		// stop service
		if err := subsData.Service.Unsubscribe(); err != nil && w.errHandler != nil {
			w.errHandler(fmt.Errorf(
				"unsubscribe %q: %w",
				pairSymbol, err,
			))
		}

		// remove subscription data
		w.srvs.Delete(pairSymbol)
		return true
	})
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

	// TODO: check subscription exists?
	w.srvs.Load(pairSymbol) // TODO: use interval + symbol

	srv, err := gate.NewWsService(context.Background(), nil, nil)
	if err != nil {
		return fmt.Errorf("conn: %w", err)
	}

	gateInterval, err := mappers.ConvertIntervalToGate(interval)
	if err != nil {
		return fmt.Errorf("convert interval: %w", err)
	}

	reqPayload := getCandleSubsPayload(gateInterval, pairSymbol)

	w.srvs.Store(pairSymbol, workers.SubscriptionData{
		Service: getUnsubscriber(srv, reqPayload),
	})

	srv.SetCallBack(
		reqPayload.Channel,
		w.getCandleCallback(
			pairSymbol,
			interval,
			eventCallback,
			errorHandler,
		),
	)

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
