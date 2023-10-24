package bybit

import (
	"fmt"
	"time"

	"github.com/hirokisan/bybit/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bybit/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
)

func (a *adapter) GetCandles(
	limit int,
	symbol string,
	interval string,
) ([]workers.CandleData, error) {
	bybitInterval, isExists := mappers.CandleIntervalsToBybit[interval]
	if !isExists {
		return nil, fmt.Errorf("interval %q not available", interval)
	}

	timeTo := time.Now()
	periodDuration := bybitInterval.Duration * time.Duration(limit)
	timeFrom := timeTo.Add(-periodDuration)

	fromTimestamp := int(timeFrom.UnixMilli())
	toTimestamp := int(timeTo.UnixMilli())

	response, err := a.client.V5().Market().GetKline(bybit.V5GetKlineParam{
		Category: bybit.CategoryV5Spot,
		Symbol:   bybit.SymbolV5(symbol),
		Interval: bybit.Interval(bybitInterval.Code),
		Start:    &fromTimestamp,
		End:      &toTimestamp,
		Limit:    &limit,
	})
	if err != nil {
		return nil, fmt.Errorf("get candles: %w", err)
	}

	var events []workers.CandleData
	for _, data := range response.Result.List {
		event, err := mappers.ConvertHistoricalCandle(
			symbol,
			data,
			bybitInterval.Duration,
			interval,
		)
		if err != nil {
			return nil, fmt.Errorf("convert candle: %w", err)
		}

		events = append(events, event)
	}
	return events, nil
}
