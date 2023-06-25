package binance

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/adshao/go-binance/v2"
)

func TestConvertCandles(t *testing.T) {
	// given
	klines := []*binance.Kline{
		{
			OpenTime:  time.Date(2023, 6, 25, 0, 0, 0, 0, time.UTC).UnixMilli(),
			CloseTime: time.Date(2023, 6, 25, 0, 59, 59, 0, time.UTC).UnixMilli(),
			Open:      "1000",
			Close:     "2000",
			High:      "3000",
			Low:       "500",
			Volume:    "10000",
		},
		{
			OpenTime:  time.Date(2023, 6, 25, 1, 0, 0, 0, time.UTC).UnixMilli(),
			CloseTime: time.Date(2023, 6, 25, 1, 59, 59, 0, time.UTC).UnixMilli(),
			Open:      "2000",
			Close:     "3000",
			High:      "4000",
			Low:       "1000",
			Volume:    "20000",
		},
	}
	interval := "1h"

	// when
	candles, err := ConvertCandles(klines, interval)

	// then
	assert.NoError(t, err)
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
	interval := "1h"

	// when
	_, err := ConvertCandles(klines, interval)

	// then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "convert candles `open` value is empty")
}
