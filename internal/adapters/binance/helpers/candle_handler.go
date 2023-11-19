package helpers

import (
	"github.com/adshao/go-binance/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
)

func GetCandleEventsHandler(
	eventCallback func(event workers.CandleEvent),
	errorHandler func(err error),
) func(event *binance.WsKlineEvent) {
	return func(event *binance.WsKlineEvent) {
		if event == nil {
			return
		}

		wEvent, err := mappers.ConvertBinanceCandleEvent(event)
		if err != nil {
			errorHandler(err)
			return
		}

		eventCallback(wEvent)
	}
}
