package binance

import (
	"context"
	"fmt"

	"github.com/adshao/go-binance/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
)

// CandleWorkerBinance - MarketDataWorker for binance
type CandleWorkerBinance struct {
	workers.CandleWorker
}

func (a *adapter) GetCandles(limit int, pairSymbol string, interval string) (
	[]workers.CandleData,
	error,
) {
	ctx, cancel := context.WithTimeout(context.Background(), consts.ReadTimeout)
	defer cancel()

	klines, err := a.binanceAPI.GetKlines(ctx, pairSymbol, interval, limit)
	if err != nil {
		return nil, fmt.Errorf("get klines: %w", err)
	}

	candles, err := mappers.ConvertCandles(klines, interval)
	if err != nil {
		return nil, fmt.Errorf("convert candles: %w", err)
	}

	return candles, nil
}

// GetCandleWorker - create new market candle worker
func (a *adapter) GetCandleWorker() workers.ICandleWorker {
	w := CandleWorkerBinance{}
	w.ExchangeTag = a.GetTag()
	return &w
}

func (w *CandleWorkerBinance) SubscribeToCandle(
	pairSymbol string,
	eventCallback func(event workers.CandleEvent),
	errorHandler func(err error),
) error {
	var openWsErr error
	w.WsChannels = new(pkgStructs.WorkerChannels)
	w.WsChannels.WsDone, w.WsChannels.WsStop, openWsErr = binance.WsKlineServe(
		pairSymbol,
		consts.CandlesInterval,
		helpers.GetCandleEventsHandler(eventCallback, errorHandler),
		errorHandler,
	)
	return openWsErr
}

func (w *CandleWorkerBinance) SubscribeToCandlesList(
	intervalsPerPair map[string]string,
	eventCallback func(event workers.CandleEvent),
	errorHandler func(err error),
) error {
	var openWsErr error
	w.WsChannels = new(pkgStructs.WorkerChannels)
	w.WsChannels.WsDone, w.WsChannels.WsStop, openWsErr = binance.WsCombinedKlineServe(
		intervalsPerPair,
		helpers.GetCandleEventsHandler(eventCallback, errorHandler),
		errorHandler,
	)
	return openWsErr
}
