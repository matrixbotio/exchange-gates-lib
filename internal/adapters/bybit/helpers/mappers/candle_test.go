package mappers

import (
	"strconv"
	"testing"
	"time"

	"github.com/hirokisan/bybit/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCandleEndTimeData(t *testing.T) {
	// given
	expectedStartTimeMs := int64(1691831504753)
	event := bybit.V5GetKlineItem{
		StartTime: strconv.FormatInt(expectedStartTimeMs, 10),
	}
	intervalDuration := time.Hour
	expectedTimestampMs := int64(1691835104753)

	// when
	timeData, err := getCandleEndTimeData(event, intervalDuration)

	// then
	require.NoError(t, err)
	assert.Equal(t, expectedStartTimeMs, timeData.StartTimeMs)
	assert.Equal(t, expectedTimestampMs, timeData.EndTimeMs)
}

func TestConvertHistoricalCandle(t *testing.T) {
	// given
	pairSymbol := "BTCUSDT"
	eventData := bybit.V5GetKlineItem{
		StartTime: "1692119310600",
		Open:      "0.35",
		Close:     "0.45",
		Low:       "0.25",
		High:      "0.41",
		Volume:    "125061",
	}
	intervalDuration := time.Hour
	intervalCode := consts.Interval1hour

	// when
	candleData, err := ConvertHistoricalCandle(
		pairSymbol,
		eventData,
		intervalDuration,
		intervalCode,
	)

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(0.35), candleData.Open)
	assert.Equal(t, float64(0.45), candleData.Close)
	assert.Equal(t, float64(0.25), candleData.Low)
	assert.Equal(t, float64(0.41), candleData.High)
}

func TestConvertWsCandle(t *testing.T) {
	// given
	pairSymbol := "LTCUSDT"
	eventData := bybit.V5WebsocketPublicKlineData{
		Open:    "0.35",
		Close:   "0.45",
		Low:     "0.25",
		High:    "0.41",
		Volume:  "125061",
		Confirm: true,
	}
	interval := consts.Interval1hour

	// when
	event, err := ConvertWsCandle(pairSymbol, interval, eventData)

	// then
	require.NoError(t, err)
	assert.True(t, event.IsFinished)
	assert.Equal(t, pairSymbol, event.Symbol)
	assert.Equal(t, float64(0.35), event.Candle.Open)
	assert.Equal(t, float64(0.45), event.Candle.Close)
	assert.Equal(t, float64(0.25), event.Candle.Low)
	assert.Equal(t, float64(0.41), event.Candle.High)
	assert.Equal(t, float64(125061), event.Candle.Volume)
}
