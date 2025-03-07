package mappers

import (
	"errors"
	"fmt"
	"strconv"

	gate "github.com/gateio/gatews/go"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
	"github.com/shopspring/decimal"
)

const gateCandleDataLen = 8

func ConvertCandles(
	rawData [][]string,
	interval consts.Interval,
) ([]workers.CandleData, error) {
	var result []workers.CandleData

	for _, candleRaw := range rawData {
		candle, err := ConvertCandle(candleRaw, interval)
		if err != nil {
			return nil, err
		}

		result = append(result, candle)
	}
	return result, nil
}

/*
- 0 Unix timestamp with second precision
- 1 Trading volume in quote currency
- 2 Closing price
- 3 Highest price
- 4 Lowest price
- 5 Opening price
- 6 Trading volume in base currency
- 7 Whether the window is closed; true indicates the end of this
*/
func ConvertCandle(rawData []string, interval consts.Interval) (workers.CandleData, error) {
	if len(rawData) < gateCandleDataLen {
		return workers.CandleData{}, errors.New("invalid candle data len")
	}

	timestampSeconds, err := strconv.ParseInt(rawData[0], 10, 64)
	if err != nil {
		return workers.CandleData{}, fmt.Errorf("parse timestamp: %w", err)
	}

	baseVolume, err := decimal.NewFromString(rawData[6])
	if err != nil {
		return workers.CandleData{}, fmt.Errorf("parse volume: %w", err)
	}

	priceOpen, err := decimal.NewFromString(rawData[5])
	if err != nil {
		return workers.CandleData{}, fmt.Errorf("parse open price: %w", err)
	}

	priceClose, err := decimal.NewFromString(rawData[2])
	if err != nil {
		return workers.CandleData{}, fmt.Errorf("parse close price: %w", err)
	}

	priceHigh, err := decimal.NewFromString(rawData[3])
	if err != nil {
		return workers.CandleData{}, fmt.Errorf("parse high price: %w", err)
	}

	priceLow, err := decimal.NewFromString(rawData[4])
	if err != nil {
		return workers.CandleData{}, fmt.Errorf("parse low price: %w", err)
	}

	return workers.CandleData{
		EndTime:  timestampSeconds * 1000,
		Interval: interval,
		Open:     priceOpen.InexactFloat64(),
		Close:    priceClose.InexactFloat64(),
		High:     priceHigh.InexactFloat64(),
		Low:      priceLow.InexactFloat64(),
		Volume:   baseVolume.InexactFloat64(),
	}, nil
}

func ParseCandleEvent(
	event gate.SpotCandleUpdateMsg,
	pairSymbol string,
	interval consts.Interval,
) (workers.CandleEvent, error) {
	timestampSeconds, err := strconv.ParseInt(event.Time, 10, 64)
	if err != nil {
		return workers.CandleEvent{}, fmt.Errorf("parse timestamp: %w", err)
	}

	priceOpen, err := decimal.NewFromString(event.Open)
	if err != nil {
		return workers.CandleEvent{}, fmt.Errorf("parse open price: %w", err)
	}

	priceClose, err := decimal.NewFromString(event.Close)
	if err != nil {
		return workers.CandleEvent{}, fmt.Errorf("parse close price: %w", err)
	}

	priceHigh, err := decimal.NewFromString(event.High)
	if err != nil {
		return workers.CandleEvent{}, fmt.Errorf("parse high price: %w", err)
	}

	priceLow, err := decimal.NewFromString(event.Low)
	if err != nil {
		return workers.CandleEvent{}, fmt.Errorf("parse low price: %w", err)
	}

	baseVolume, err := decimal.NewFromString(event.Volume)
	if err != nil {
		return workers.CandleEvent{}, fmt.Errorf("parse volume: %w", err)
	}

	return workers.CandleEvent{
		Symbol: pairSymbol,
		Candle: workers.CandleData{
			EndTime:  timestampSeconds,
			Interval: interval,
			Open:     priceOpen.InexactFloat64(),
			Close:    priceClose.InexactFloat64(),
			High:     priceHigh.InexactFloat64(),
			Low:      priceLow.InexactFloat64(),
			Volume:   baseVolume.InexactFloat64(),
		},
		Time:       timestampSeconds,
		IsFinished: event.WindowClose,
	}, nil
}
