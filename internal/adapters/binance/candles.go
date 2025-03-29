package binance

import (
	"context"
	"fmt"

	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/wrapper"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
)

// CandleWorkerBinance - MarketDataWorker for binance
type CandleWorkerBinance struct {
	workers.CandleWorker

	binanceAPI wrapper.BinanceAPIWrapper
}

func NewCandleWorker(binanceAPI wrapper.BinanceAPIWrapper) *CandleWorkerBinance {
	w := &CandleWorkerBinance{
		binanceAPI: binanceAPI,
	}
	w.ExchangeTag = consts.BinanceAdapterTag
	return w
}

func (a *adapter) GetCandles(limit int, pairSymbol string, interval consts.Interval) (
	[]workers.CandleData,
	error,
) {
	ctx, cancel := context.WithTimeout(context.Background(), consts.ReadTimeout)
	defer cancel()

	klines, err := a.binanceAPI.GetKlines(
		ctx, pairSymbol,
		convertInterval(interval),
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("get klines: %w", err)
	}

	candles, err := mappers.ConvertCandles(klines, interval)
	if err != nil {
		return nil, fmt.Errorf("convert candles: %w", err)
	}

	return candles, nil
}

func convertInterval(ourFormat consts.Interval) string {
	return string(ourFormat)
}

func (w *CandleWorkerBinance) SubscribeToCandle(
	pairSymbol string,
	interval consts.Interval,
	eventCallback func(event workers.CandleEvent),
	errorHandler func(err error),
) error {
	if w.CandleWorker.IsSubscriptionExists(pairSymbol, convertInterval(interval)) {
		return nil
	}

	wsDone, wsStop, err := w.binanceAPI.SubscribeToCandle(
		pairSymbol,
		convertInterval(interval),
		eventCallback,
		errorHandler,
	)
	if err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}

	// save subscription
	w.CandleWorker.Save(
		workers.CreateChannelsUnsubscriber(wsDone, wsStop),
		errorHandler,
		pairSymbol, convertInterval(interval),
	)
	return nil
}
