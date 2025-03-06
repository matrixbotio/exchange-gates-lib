package bingx

import (
	"fmt"
	"time"

	bingxgo "github.com/Sagleft/go-bingx"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
	"github.com/shopspring/decimal"
)

type IntervalData struct {
	Interval consts.Interval
	Duration time.Duration
}

var intervalBingxWsToOur = map[bingxgo.Interval]IntervalData{
	bingxgo.Interval1:   {consts.Interval1min, time.Minute},
	bingxgo.Interval5:   {consts.Interval5min, time.Minute * 5},
	bingxgo.Interval15:  {consts.Interval15min, time.Minute * 15},
	bingxgo.Interval30:  {consts.Interval30min, time.Minute * 30},
	bingxgo.Interval60:  {consts.Interval1hour, time.Hour},
	bingxgo.Interval4h:  {consts.Interval4hour, time.Hour * 4},
	bingxgo.Interval6h:  {consts.Interval6hour, time.Hour * 6},
	bingxgo.Interval12h: {consts.Interval12hour, time.Hour * 12},
	bingxgo.Interval1d:  {consts.Interval1day, time.Hour * 24},
}

var intervalBingxRestToOur = map[bingxgo.Interval]IntervalData{
	"1m":  {consts.Interval1min, time.Minute},
	"5m":  {consts.Interval5min, time.Minute * 5},
	"15m": {consts.Interval15min, time.Minute * 15},
	"30m": {consts.Interval30min, time.Minute * 30},
	"1h":  {consts.Interval1hour, time.Hour},
	"4h":  {consts.Interval4hour, time.Hour * 4},
	"6h":  {consts.Interval6hour, time.Hour * 6},
	"12h": {consts.Interval12hour, time.Hour * 12},
	"1d":  {consts.Interval1day, time.Hour * 24},
}

var ourIntervalToBingXWs = func() map[consts.Interval]bingxgo.Interval {
	r := map[consts.Interval]bingxgo.Interval{}
	for interval, data := range intervalBingxWsToOur {
		r[data.Interval] = interval
	}
	return r
}()

var ourIntervalToBingXRest = func() map[consts.Interval]bingxgo.Interval {
	r := map[consts.Interval]bingxgo.Interval{}
	for interval, data := range intervalBingxRestToOur {
		r[data.Interval] = interval
	}
	return r
}()

// Convert inteval in websocket format
func ConvertIntervalToBingXWs(interval consts.Interval) (bingxgo.Interval, error) {
	result, isExists := ourIntervalToBingXWs[interval]
	if !isExists {
		return "", fmt.Errorf("unknown interval: %q", interval)
	}

	return result, nil
}

// Convert inteval in REST format
func ConvertIntervalToBingXRest(interval consts.Interval) (bingxgo.Interval, error) {
	result, isExists := ourIntervalToBingXRest[interval]
	if !isExists {
		return "", fmt.Errorf("unknown interval: %q", interval)
	}

	return result, nil
}

func ConvertBingXWsInterval(interval bingxgo.Interval) (IntervalData, error) {
	result, isExists := intervalBingxWsToOur[interval]
	if !isExists {
		return IntervalData{}, fmt.Errorf("unknown interval: %q", interval)
	}

	return result, nil
}

func ConvertBingXRestInterval(interval bingxgo.Interval) (IntervalData, error) {
	result, isExists := intervalBingxRestToOur[interval]
	if !isExists {
		return IntervalData{}, fmt.Errorf("unknown interval: %q", interval)
	}

	return result, nil
}

func ConvertKlinesRest(klines []bingxgo.KlineData) ([]workers.CandleData, error) {
	var result []workers.CandleData

	for _, kline := range klines {
		interval, err := ConvertBingXRestInterval(bingxgo.Interval(kline.Interval))
		if err != nil {
			return nil, fmt.Errorf("convert interval: %w", err)
		}

		result = append(result, workers.CandleData{
			StartTime: kline.StartTime,
			EndTime:   kline.EndTime,
			Interval:  interval.Interval,
			Open:      kline.Open,
			Close:     kline.Close,
			High:      kline.High,
			Low:       kline.Low,
			Volume:    kline.Volume,
		})
	}

	return result, nil
}

func ConvertWsKline(kline bingxgo.KlineEvent) (workers.CandleEvent, error) {
	intervalData, err := ConvertBingXWsInterval(kline.Interval)
	if err != nil {
		return workers.CandleEvent{}, fmt.Errorf("convert interval: %w", err)
	}

	klineOpen, err := decimal.NewFromString(kline.Open)
	if err != nil {
		return workers.CandleEvent{},
			fmt.Errorf("parse open price: %w", err)
	}

	klineClose, err := decimal.NewFromString(kline.Close)
	if err != nil {
		return workers.CandleEvent{},
			fmt.Errorf("parse close price: %w", err)
	}

	klineHigh, err := decimal.NewFromString(kline.High)
	if err != nil {
		return workers.CandleEvent{},
			fmt.Errorf("parse high price: %w", err)
	}

	klineLow, err := decimal.NewFromString(kline.Low)
	if err != nil {
		return workers.CandleEvent{},
			fmt.Errorf("parse low price: %w", err)
	}

	klineVolume, err := decimal.NewFromString(kline.Volume)
	if err != nil {
		return workers.CandleEvent{},
			fmt.Errorf("parse volume: %w", err)
	}

	return workers.CandleEvent{
		Symbol: kline.Symbol,
		Candle: workers.CandleData{
			StartTime: kline.StartTime,
			EndTime:   kline.EndTime,
			Interval:  intervalData.Interval,
			Open:      klineOpen.InexactFloat64(),
			Close:     klineClose.InexactFloat64(),
			High:      klineHigh.InexactFloat64(),
			Low:       klineLow.InexactFloat64(),
			Volume:    klineVolume.InexactFloat64(),
		},
		Time:       kline.EventTime,
		IsFinished: kline.Completed,
	}, nil
}

func GetBingXCandleEventsHandler(
	eventCallback func(event workers.CandleEvent),
	errorCallback func(error),
) func(event bingxgo.KlineEvent) {
	return func(event bingxgo.KlineEvent) {
		candle, err := ConvertWsKline(event)
		if err != nil {
			errorCallback(err)
			return
		}

		eventCallback(candle)
	}
}
