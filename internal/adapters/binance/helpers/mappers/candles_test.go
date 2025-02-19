package mappers

import (
	"testing"
	"time"

	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/adshao/go-binance/v2"
)

const testKlineInterval = consts.Interval1hour

func TestConvertCandles(t *testing.T) {
	// given
	klines := GetTestKlines()
	interval := "1h"

	// when
	candles, err := ConvertCandles(klines, testKlineInterval)

	// then
	require.NoError(t, err)
	assert.Len(t, candles, 2)

	// validate the contents of the first candle
	assert.Equal(t, klines[0].OpenTime, candles[0].StartTime)
	assert.Equal(t, fixCandleEndTime(klines[0].CloseTime), candles[0].EndTime)
	assert.Equal(t, interval, candles[0].Interval)
	assert.Equal(t, 1000.0, candles[0].Open)
	assert.Equal(t, 2000.0, candles[0].Close)
	assert.Equal(t, 3000.0, candles[0].High)
	assert.Equal(t, 500.0, candles[0].Low)
	assert.Equal(t, 10000.0, candles[0].Volume)
}

func TestConvertCandlesWithError(t *testing.T) {
	// given
	klines := []*binance.Kline{
		{
			OpenTime:  time.Date(2023, 6, 25, 0, 0, 0, 0, time.UTC).UnixMilli(),
			CloseTime: time.Date(2023, 6, 25, 0, 59, 59, 0, time.UTC).UnixMilli(),
			Open:      "", // Open is empty string
			Close:     "2000",
			High:      "3000",
			Low:       "500",
			Volume:    "10000",
		},
	}

	// when
	_, err := ConvertCandles(klines, testKlineInterval)

	// then
	require.ErrorContains(t, err, "`open` value is empty")
}

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
	event, err := ConvertBinanceCandleEvent(&rawEvent)

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
