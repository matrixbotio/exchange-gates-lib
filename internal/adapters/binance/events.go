package binance

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/adshao/go-binance/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
	"github.com/matrixbotio/exchange-gates-lib/pkg/utils"
)

// PriceWorkerBinance - MarketDataWorker for binance
type PriceWorkerBinance struct {
	workers.PriceWorker
}

// TradeEventWorkerBinance - TradeEventWorker for binance
type TradeEventWorkerBinance struct {
	workers.TradeEventWorker
}

func (a *adapter) GetPriceWorker(callback workers.PriceEventCallback) workers.IPriceWorker {
	w := PriceWorkerBinance{}
	w.PriceWorker.ExchangeTag = a.Tag
	w.PriceWorker.HandleEventCallback = callback
	return &w
}

func (w *PriceWorkerBinance) handlePriceEvent(event *binance.WsBookTickerEvent) {
	if event == nil {
		return
	}

	ask, bid, err := mappers.ConvertPriceEvent(*event)
	if err != nil {
		return // ignore broken price event
	}

	w.HandleEventCallback(workers.PriceEvent{
		ExchangeTag: w.ExchangeTag,
		Symbol:      event.Symbol,
		Ask:         ask,
		Bid:         bid,
	})
}

func (w *PriceWorkerBinance) SubscribeToPriceEvents(
	pairSymbols []string,
	errorHandler func(err error),
) (map[string]pkgStructs.WorkerChannels, error) {
	result := map[string]pkgStructs.WorkerChannels{}

	// event handler func
	w.WsChannels = new(pkgStructs.WorkerChannels)

	var openWsErr error
	for _, pairSymbol := range pairSymbols {
		newChannels := pkgStructs.WorkerChannels{}
		newChannels.WsDone, newChannels.WsStop, openWsErr = binance.WsBookTickerServe(
			pairSymbol,
			w.handlePriceEvent,
			errorHandler,
		)
		if openWsErr != nil {
			return result, fmt.Errorf("subscribe to %q price: %w", pairSymbol, openWsErr)
		}

		result[pairSymbol] = newChannels
	}

	return result, nil
}

func (a *adapter) GetTradeEventsWorker() workers.ITradeEventWorker {
	w := TradeEventWorkerBinance{}
	w.ExchangeTag = a.GetTag()
	return &w
}

// SubscribeToTradeEvents - websocket subscription to change trade candles on the exchange
func (w *TradeEventWorkerBinance) SubscribeToTradeEvents(
	symbol string,
	eventCallback func(event workers.TradeEvent),
	errorHandler func(err error),
) error {

	wsErrHandler := func(err error) {
		errorHandler(err)
	}

	wsTradeHandler := func(event *binance.WsTradeEvent) {
		if event != nil {
			// fix event.Time
			if strings.HasSuffix(strconv.FormatInt(event.Time, 10), "999") {
				event.Time++
			}
			wEvent := workers.TradeEvent{
				ID:            event.TradeID,
				Time:          event.Time,
				Symbol:        event.Symbol,
				ExchangeTag:   w.ExchangeTag,
				BuyerOrderID:  event.BuyerOrderID,
				SellerOrderID: event.SellerOrderID,
			}
			errs := make([]error, 2)
			wEvent.Price, errs[0] = strconv.ParseFloat(event.Price, 64)
			wEvent.Quantity, errs[0] = strconv.ParseFloat(event.Quantity, 64)
			if utils.LogNotNilError(errs) {
				return
			}
			eventCallback(wEvent)
		}
	}

	var err error
	w.WsChannels = new(pkgStructs.WorkerChannels)
	w.WsChannels.WsDone, w.WsChannels.WsStop, err = binance.WsTradeServe(
		symbol,
		wsTradeHandler,
		wsErrHandler,
	)
	if err != nil {
		return fmt.Errorf("subscribe to trade events: %w", err)
	}
	return nil
}
