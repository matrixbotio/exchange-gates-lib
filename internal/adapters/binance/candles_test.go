package binance

import (
	"testing"

	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/wrapper"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetCandlesSuccess(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)
	limit := 5
	pairSymbol := "LTCUSDT"
	interval := "5m"

	w.EXPECT().GetKlines(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(mappers.GetTestKlines(), nil)

	// when
	candles, err := a.GetCandles(limit, pairSymbol, interval)

	// then
	require.NoError(t, err)
	require.Len(t, candles, 2)
	assert.Equal(t, interval, candles[0].Interval)
	assert.Equal(t, 1000.0, candles[0].Open)
	assert.Equal(t, 2000.0, candles[0].Close)
	assert.Equal(t, 3000.0, candles[0].High)
	assert.Equal(t, 500.0, candles[0].Low)
	assert.Equal(t, 10000.0, candles[0].Volume)
}

func TestGetCandlesGetKlinesError(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)
	limit := 5
	pairSymbol := "LTCUSDT"
	interval := "5m"

	w.EXPECT().GetKlines(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errTestException)

	// when
	_, err := a.GetCandles(limit, pairSymbol, interval)

	// then
	require.ErrorIs(t, err, errTestException)
}

func TestGetCandlesConvertError(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)
	limit := 5
	pairSymbol := "LTCUSDT"
	interval := "5m"

	klines := mappers.GetTestKlines()
	klines[0].Low = "strange data"

	w.EXPECT().GetKlines(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(klines, nil)

	// when
	_, err := a.GetCandles(limit, pairSymbol, interval)

	// then
	require.ErrorContains(t, err, "convert candles")
}

func TestGetCandleWorkerSuccess(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)

	// when
	worker := a.GetCandleWorker()

	// then
	assert.NotEmpty(t, worker.GetExchangeTag())
	assert.Equal(t, a.GetTag(), worker.GetExchangeTag())
}

func TestSubscribeToCandleSuccess(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)
	worker := a.GetCandleWorker()

	w.EXPECT().SubscribeToCandle(
		mock.Anything, mock.Anything,
		mock.Anything, mock.Anything,
	).Return(
		make(chan struct{}),
		make(chan struct{}),
		nil,
	)

	// when
	err := worker.SubscribeToCandle(
		testPairSymbol,
		func(event workers.CandleEvent) {},
		func(err error) {},
	)

	// then
	require.NoError(t, err)
}

func TestSubscribeToCandleError(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)
	worker := a.GetCandleWorker()

	w.EXPECT().SubscribeToCandle(
		mock.Anything, mock.Anything,
		mock.Anything, mock.Anything,
	).Return(
		make(chan struct{}),
		make(chan struct{}),
		errTestException,
	)

	// when
	err := worker.SubscribeToCandle(
		testPairSymbol,
		func(event workers.CandleEvent) {},
		func(err error) {},
	)

	// then
	require.ErrorIs(t, err, errTestException)
}

func TestSubscribeToCandlesListSuccess(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)
	worker := a.GetCandleWorker()

	intervalsPerPair := map[string]string{
		"LTCUSDT": "1m",
	}

	w.EXPECT().SubscribeToCandlesList(
		mock.Anything, mock.Anything,
		mock.Anything,
	).Return(
		make(chan struct{}),
		make(chan struct{}),
		nil,
	)

	// when
	err := worker.SubscribeToCandlesList(
		intervalsPerPair,
		func(event workers.CandleEvent) {},
		func(err error) {},
	)

	// then
	require.NoError(t, err)
}

func TestSubscribeToCandlesListError(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)
	worker := a.GetCandleWorker()

	intervalsPerPair := map[string]string{
		"LTCUSDT": "1m",
	}

	w.EXPECT().SubscribeToCandlesList(
		mock.Anything, mock.Anything,
		mock.Anything,
	).Return(
		make(chan struct{}),
		make(chan struct{}),
		errTestException,
	)

	// when
	err := worker.SubscribeToCandlesList(
		intervalsPerPair,
		func(event workers.CandleEvent) {},
		func(err error) {},
	)

	// then
	require.ErrorIs(t, err, errTestException)
}
