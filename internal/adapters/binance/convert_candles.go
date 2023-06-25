package binance

import (
	"errors"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
	"strconv"
)

func ConvertCandles(klines []*binance.Kline, interval string) ([]workers.CandleData, error) {
	var candles []workers.CandleData

	for _, kline := range klines {
		candle := workers.CandleData{
			StartTime: kline.OpenTime,
			EndTime:   fixCandleEndTime(kline.CloseTime),
			Interval:  interval,
		}

		var err error
		if kline.Open == "" {
			return nil, errors.New("convert candles `open` value is empty")
		}
		if candle.Open, err = strconv.ParseFloat(kline.Open, 64); err != nil {
			return nil, fmt.Errorf("convert candles `open` value: %w", err)
		}

		if kline.Close == "" {
			return nil, errors.New("convert candles `close` value is empty")
		}
		if candle.Close, err = strconv.ParseFloat(kline.Close, 64); err != nil {
			return nil, fmt.Errorf("convert candles `close` value: %w", err)
		}

		if kline.High == "" {
			return nil, errors.New("convert candles `high` value is empty")
		}
		if candle.High, err = strconv.ParseFloat(kline.High, 64); err != nil {
			return nil, fmt.Errorf("convert candles `high` value: %w", err)
		}

		if kline.Low == "" {
			return nil, errors.New("convert candles `low` value is empty")
		}
		if candle.Low, err = strconv.ParseFloat(kline.Low, 64); err != nil {
			return nil, fmt.Errorf("convert candles `low` value: %w", err)
		}

		if kline.Volume == "" {
			return nil, errors.New("convert candles `volume` value is empty")
		}
		if candle.Volume, err = strconv.ParseFloat(kline.Volume, 64); err != nil {
			return nil, fmt.Errorf("convert candles `volume` value: %w", err)
		}

		candles = append(candles, candle)
	}

	return candles, nil
}
