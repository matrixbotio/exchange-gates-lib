package mappers

import (
	"fmt"
	"time"

	bingxgo "github.com/Sagleft/go-bingx"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
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

func ConvertWsKline(kline *bingxgo.WsKlineEvent) workers.CandleEvent {
	intervalData, err := ConvertBingXInterval(kline.Interval)
	if err != nil {
		intervalData.Code = intervalMinute
	}

	return workers.CandleEvent{
		Symbol: kline.Symbol,
		Candle: workers.CandleData{
			StartTime: int64(kline.StartTime),
			EndTime:   int64(kline.EndTime),
			Interval:  intervalData.Code,
			Open:      kline.Open,
			Close:     kline.Close,
			High:      kline.High,
			Low:       kline.Low,
			Volume:    kline.Volume,
		},
		Time:       int64(kline.EndTime),
		IsFinished: kline.Completed,
	}
}
