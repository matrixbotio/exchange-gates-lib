package helpers

import (
	"testing"

	"github.com/adshao/go-binance/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCandleEventsHandlerSuccess(t *testing.T) {
	// given
	var testErr error
	var rawEvent = binance.WsKlineEvent{
		Symbol: "BTCBUSD",
		Kline: binance.WsKline{
			EndTime: 1682506327999,
			Open:    "100",
			Close:   "105",
			High:    "120",
			Low:     "98",
			Volume:  "500",
		},
	}
	var convertedEvent workers.CandleEvent
	var eventHandler = func(event workers.CandleEvent) {
		convertedEvent = event
	}
	var errorHandler = func(err error) {
		testErr = err
	}

	// when
	candleHandler := GetCandleEventsHandler(eventHandler, errorHandler)
	candleHandler(&rawEvent)

	// then
	require.NoError(t, testErr)
	assert.Equal(t, rawEvent.Symbol, convertedEvent.Symbol)
}

func TestGetCandleEventsHandlerError(t *testing.T) {
	// given
	var testErr error
	var rawEvent = binance.WsKlineEvent{}
	var eventHandler = func(event workers.CandleEvent) {}
	var errorHandler = func(err error) {
		testErr = err
	}

	// when
	candleHandler := GetCandleEventsHandler(eventHandler, errorHandler)
	candleHandler(&rawEvent)

	// then
	require.Error(t, testErr)
}

func TestGetCandleEventsHandlerEmptyEvent(t *testing.T) {
	// given
	var testErr error
	var convertedEvent workers.CandleEvent
	var eventHandler = func(event workers.CandleEvent) {
		convertedEvent = event
	}
	var errorHandler = func(err error) {
		testErr = err
	}

	// when
	candleHandler := GetCandleEventsHandler(eventHandler, errorHandler)
	candleHandler(nil)

	// then
	require.NoError(t, testErr)
	assert.Empty(t, convertedEvent.Symbol)
}
