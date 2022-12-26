package binance

import (
	"context"
	"fmt"

	"github.com/adshao/go-binance/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
)

// CandleWorkerBinance - MarketDataWorker for binance
type CandleWorkerBinance struct {
	workers.CandleWorker
}

func (a *adapter) GetCandles(limit int, symbol string, interval string) ([]workers.CandleData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), consts.ReadTimeout)
	defer cancel()

	klines, err := a.binanceAPI.NewKlinesService().Symbol(symbol).Interval(interval).Limit(limit).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("binance adapter get candles rq: %w", err)
	}

	candles, err := ConvertCandles(klines, interval)
	if err != nil {
		return nil, fmt.Errorf("binance adapter get candles convert: %w", err)
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
		getCandleEventsHandler(eventCallback, errorHandler),
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
		getCandleEventsHandler(eventCallback, errorHandler),
		errorHandler,
	)
	return openWsErr
}
