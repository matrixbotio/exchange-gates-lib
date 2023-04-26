package binance

import (
	"testing"

	"github.com/Sagleft/go-binance/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertBinanceCandleEvent(t *testing.T) {
	// given
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

	// when
	event, err := convertBinanceCandleEvent(&rawEvent)

	// then
	require.NoError(t, err)
	assert.Equal(t, rawEvent.Symbol, event.Symbol)
	assert.Equal(t, float64(100), event.Candle.Open)
	assert.Equal(t, float64(105), event.Candle.Close)
	assert.Equal(t, float64(120), event.Candle.High)
	assert.Equal(t, float64(98), event.Candle.Low)
	assert.Equal(t, float64(500), event.Candle.Volume)
	assert.Equal(t, int64(1682506268000), event.Candle.EndTime)
}

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
	candleHandler := getCandleEventsHandler(eventHandler, errorHandler)
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
	candleHandler := getCandleEventsHandler(eventHandler, errorHandler)
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
	candleHandler := getCandleEventsHandler(eventHandler, errorHandler)
	candleHandler(nil)

	// then
	require.NoError(t, testErr)
	assert.Empty(t, convertedEvent.Symbol)
}
