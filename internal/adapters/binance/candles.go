package binance

import (
	"context"
	"fmt"

	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/wrapper"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
)

// CandleWorkerBinance - MarketDataWorker for binance
type CandleWorkerBinance struct {
	workers.CandleWorker

	binanceAPI wrapper.BinanceAPIWrapper
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

func (a *adapter) GetCandleWorker() workers.ICandleWorker {
	w := CandleWorkerBinance{
		binanceAPI: a.binanceAPI,
	}
	w.ExchangeTag = a.GetTag()
	return &w
}

func (w *CandleWorkerBinance) SubscribeToCandle(
	pairSymbol string,
	eventCallback func(event workers.CandleEvent),
	errorHandler func(err error),
) error {
	var err error
	w.WsChannels = new(pkgStructs.WorkerChannels)
	w.WsChannels.WsDone, w.WsChannels.WsStop, err = w.binanceAPI.SubscribeToCandle(
		pairSymbol,
		consts.CandlesInterval,
		eventCallback,
		errorHandler,
	)
	return err
}

func (w *CandleWorkerBinance) SubscribeToCandlesList(
	intervalsPerPair map[string]string,
	eventCallback func(event workers.CandleEvent),
	errorHandler func(err error),
) error {
	var openWsErr error
	w.WsChannels = new(pkgStructs.WorkerChannels)
	w.WsChannels.WsDone, w.WsChannels.WsStop, openWsErr = w.binanceAPI.SubscribeToCandlesList(
		intervalsPerPair,
		eventCallback,
		errorHandler,
	)
	return openWsErr
}
