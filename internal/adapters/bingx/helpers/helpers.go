package helpers

import (
	bingxgo "github.com/Sagleft/go-bingx"

	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bingx/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
)

func GetBingXCandleEventsHandler(
	eventCallback func(event workers.CandleEvent),
	errorCallback func(error),
) func(event bingxgo.KlineEvent) {
	return func(event bingxgo.KlineEvent) {
		candle, err := mappers.ConvertWsKline(event)
		if err != nil {
			errorCallback(err)
			return
		}

		eventCallback(candle)
	}
}
