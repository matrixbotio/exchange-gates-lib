package binance

import (
	"testing"

	"github.com/adshao/go-binance/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetPriceWorkerSuccess(t *testing.T) {
	// given
	w := NewMockBinanceAPIWrapper(t)
	a := createAdapter(w)
	cb := func(event workers.PriceEvent) {}

	// when
	worker := a.GetPriceWorker(cb)

	// then
	assert.NotEmpty(t, worker.GetExchangeTag())
	assert.Equal(t, a.GetTag(), worker.GetExchangeTag())
}

func TestHandlePriceEventSuccess(t *testing.T) {
	// given
	var w = NewMockBinanceAPIWrapper(t)
	var a = createAdapter(w)

	var lastEvent workers.PriceEvent
	var cb = func(event workers.PriceEvent) {
		lastEvent = event
	}
	var worker = a.GetPriceWorker(cb).(*PriceWorkerBinance)

	var testEvent = &binance.WsBookTickerEvent{
		Symbol:       "LTCUSDT",
		BestAskPrice: "65.081",
		BestBidPrice: "65.019",
	}

	// when
	worker.handlePriceEvent(testEvent)

	// then
	assert.Equal(t, testEvent.Symbol, lastEvent.Symbol)
	assert.Equal(t, float64(65.081), lastEvent.Ask)
	assert.Equal(t, float64(65.019), lastEvent.Bid)
}

func TestHandlePriceEventEmpty(t *testing.T) {
	// given
	var w = NewMockBinanceAPIWrapper(t)
	var a = createAdapter(w)

	var lastEvent workers.PriceEvent
	var cb = func(event workers.PriceEvent) {
		lastEvent = event
	}
	var worker = a.GetPriceWorker(cb).(*PriceWorkerBinance)

	var testEvent *binance.WsBookTickerEvent

	// when
	worker.handlePriceEvent(testEvent)

	// then
	assert.Empty(t, lastEvent.Symbol)
	assert.Empty(t, lastEvent.Ask)
	assert.Empty(t, lastEvent.Bid)
}

func TestHandlePriceEventBroken(t *testing.T) {
	// given
	var w = NewMockBinanceAPIWrapper(t)
	var a = createAdapter(w)

	var lastEvent workers.PriceEvent
	var cb = func(event workers.PriceEvent) {
		lastEvent = event
	}
	var worker = a.GetPriceWorker(cb).(*PriceWorkerBinance)

	var testEvent = &binance.WsBookTickerEvent{
		BestAskPrice: "strange data",
	}

	// when
	worker.handlePriceEvent(testEvent)

	// then
	assert.Empty(t, lastEvent.Symbol)
	assert.Empty(t, lastEvent.Ask)
	assert.Empty(t, lastEvent.Bid)
}

func TestSubscribeToPriceEventsSuccess(t *testing.T) {
	// given
	var w = NewMockBinanceAPIWrapper(t)
	var a = createAdapter(w)
	var cb = func(event workers.PriceEvent) {}
	var worker = a.GetPriceWorker(cb)
	var pairSymbols = []string{"LTCUSDT", "MTXBBTC"}

	var lastErr error
	var errHandler = func(err error) {
		lastErr = err
	}

	for _, pairSymbol := range pairSymbols {
		w.EXPECT().SubscribeToPriceEvents(
			pairSymbol,
			mock.Anything,
			mock.Anything,
		).Return(
			make(chan struct{}),
			make(chan struct{}),
			nil,
		)
	}

	// when
	_, subscribeErr := worker.SubscribeToPriceEvents(pairSymbols, errHandler)

	// then
	require.NoError(t, subscribeErr)
	require.NoError(t, lastErr)
}

func TestSubscribeToPriceEventsError(t *testing.T) {
	// given
	var w = NewMockBinanceAPIWrapper(t)
	var a = createAdapter(w)
	var cb = func(event workers.PriceEvent) {}
	var worker = a.GetPriceWorker(cb)
	var pairSymbols = []string{"LTCUSDT"}

	w.EXPECT().SubscribeToPriceEvents(
		pairSymbols[0],
		mock.Anything,
		mock.Anything,
	).Return(
		make(chan struct{}),
		make(chan struct{}),
		errTestException,
	)

	// when
	_, subscribeErr := worker.SubscribeToPriceEvents(
		pairSymbols,
		func(err error) {},
	)

	// then
	require.ErrorIs(t, subscribeErr, errTestException)
}
