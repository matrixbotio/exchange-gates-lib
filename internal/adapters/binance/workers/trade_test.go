package workers

import (
	"errors"
	"testing"

	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/wrapper"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var errTestException = errors.New("test error")

func TestGetTradeEventsWorker(t *testing.T) {
	// given
	var w = wrapper.NewMockBinanceAPIWrapper(t)

	// when
	var worker = NewTradeEventsWorker(testExchangeTag, w)

	// then
	assert.NotEmpty(t, worker.GetExchangeTag())
	assert.Equal(t, testExchangeTag, worker.GetExchangeTag())
}

func TestSubscribeToTradeEventsSuccess(t *testing.T) {
	// given
	var w = wrapper.NewMockBinanceAPIWrapper(t)
	var cb = func(event workers.TradeEvent) {}
	var worker = NewTradeEventsWorker(testExchangeTag, w)
	var pairSymbol = "LTCUSDT"

	var lastErr error
	var errHandler = func(err error) {
		lastErr = err
	}

	w.EXPECT().SubscribeToTradeEvents(
		pairSymbol,
		worker.GetExchangeTag(),
		mock.Anything,
		mock.Anything,
	).Return(
		make(chan struct{}),
		make(chan struct{}),
		nil,
	)

	// when
	err := worker.SubscribeToTradeEvents(pairSymbol, cb, errHandler)

	// then
	require.NoError(t, err)
	require.NoError(t, lastErr)
}

func TestSubscribeToTradeEventsError(t *testing.T) {
	// given
	var w = wrapper.NewMockBinanceAPIWrapper(t)
	var cb = func(event workers.TradeEvent) {}
	var worker = NewTradeEventsWorker(testExchangeTag, w)
	var pairSymbol = "LTCUSDT"

	w.EXPECT().SubscribeToTradeEvents(
		pairSymbol,
		worker.GetExchangeTag(),
		mock.Anything,
		mock.Anything,
	).Return(
		make(chan struct{}),
		make(chan struct{}),
		errTestException,
	)

	// when
	err := worker.SubscribeToTradeEvents(
		pairSymbol,
		cb,
		func(err error) {},
	)

	// then
	require.ErrorIs(t, err, errTestException)
}
