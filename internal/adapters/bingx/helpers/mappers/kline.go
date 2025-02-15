package mappers

import (
	"fmt"
	"time"

	bingxgo "github.com/Sagleft/go-bingx"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
	"github.com/shopspring/decimal"
)

const intervalMinute = "1m"

type IntervalData struct {
	Code     string
	Duration time.Duration
}

var intervalBingxToOur = map[bingxgo.Interval]IntervalData{
	bingxgo.Interval1:  {intervalMinute, time.Minute},
	bingxgo.Interval3:  {"3m", time.Minute * 3},
	bingxgo.Interval5:  {"5m", time.Minute * 5},
	bingxgo.Interval15: {"15m", time.Minute * 15},
	bingxgo.Interval30: {"30m", time.Minute * 30},
	bingxgo.Interval60: {"1h", time.Hour},
	bingxgo.Interval2h: {"2h", time.Hour * 2},
	bingxgo.Interval4h: {"4h", time.Hour * 4},
	bingxgo.Interval6h: {"6h", time.Hour * 6},
	bingxgo.Interval1d: {"D", time.Hour * 24},
	bingxgo.Interval1w: {"W", time.Hour * 24 * 7},
	bingxgo.Interval1M: {"M", time.Hour * 24 * 30},
}

var ourIntervalToBingX = func() map[string]bingxgo.Interval {
	r := map[string]bingxgo.Interval{}
	for interval, data := range intervalBingxToOur {
		r[data.Code] = interval
	}
	return r
}()

func ConvertBingXInterval(interval bingxgo.Interval) (IntervalData, error) {
	result, isExists := intervalBingxToOur[interval]
	if !isExists {
		return IntervalData{}, fmt.Errorf("unknown interval: %q", interval)
	}

	return result, nil
}

func ConvertKlines(klines []bingxgo.KlineData) []workers.CandleData {
	var result []workers.CandleData

	for _, kline := range klines {
		result = append(result, workers.CandleData{
			StartTime: kline.StartTime,
			EndTime:   kline.EndTime,
			Interval:  kline.Interval,
			Open:      kline.Open,
			Close:     kline.Close,
			High:      kline.High,
			Low:       kline.Low,
			Volume:    kline.Volume,
		})
	}

	return result
}

func ConvertWsKline(kline bingxgo.KlineEvent) (workers.CandleEvent, error) {
	intervalData, err := ConvertBingXInterval(kline.Interval)
	if err != nil {
		intervalData.Code = intervalMinute
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
			Interval:  intervalData.Code,
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
