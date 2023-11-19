package mappers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/adshao/go-binance/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
)

func fixCandleEndTime(endTime int64) int64 {
	if strings.HasSuffix(strconv.FormatInt(endTime, 10), "999") {
		return endTime - 59999
	}
	return endTime
}

func ConvertBinanceCandleEvent(event *binance.WsKlineEvent) (workers.CandleEvent, error) {
	e := workers.CandleEvent{
		Symbol: event.Symbol,
		Candle: workers.CandleData{
			StartTime: event.Kline.StartTime,
			EndTime:   fixCandleEndTime(event.Kline.EndTime),
			Interval:  event.Kline.Interval,
		},
		Time: event.Time,
	}

	var err error
	if event.Kline.Open == "" {
		return workers.CandleEvent{}, errors.New("candle `open` value is empty")
	}
	if e.Candle.Open, err = strconv.ParseFloat(event.Kline.Open, 64); err != nil {
		return workers.CandleEvent{}, fmt.Errorf("parse candle `open` value: %w", err)
	}

	if event.Kline.Close == "" {
		return workers.CandleEvent{}, errors.New("candle `close` value is empty")
	}
	if e.Candle.Close, err = strconv.ParseFloat(event.Kline.Close, 64); err != nil {
		return workers.CandleEvent{}, fmt.Errorf("parse candle `close` value: %w", err)
	}

	if event.Kline.High == "" {
		return workers.CandleEvent{}, errors.New("candle `high` value is empty")
	}
	if e.Candle.High, err = strconv.ParseFloat(event.Kline.High, 64); err != nil {
		return workers.CandleEvent{}, fmt.Errorf("parse candle `high` value: %w", err)
	}

	if event.Kline.Low == "" {
		return workers.CandleEvent{}, errors.New("candle `low` value is empty")
	}
	if e.Candle.Low, err = strconv.ParseFloat(event.Kline.Low, 64); err != nil {
		return workers.CandleEvent{}, fmt.Errorf("parse candle `low` value: %w", err)
	}

	if event.Kline.Volume == "" {
		return workers.CandleEvent{}, errors.New("candle `volume` value is empty")
	}
	if e.Candle.Volume, err = strconv.ParseFloat(event.Kline.Volume, 64); err != nil {
		return workers.CandleEvent{}, fmt.Errorf("parse candle `volume` value: %w", err)
	}

	return e, nil
}

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
