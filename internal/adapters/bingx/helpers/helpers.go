package helpers

import (
	bingxgo "github.com/Sagleft/go-bingx"

	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bingx/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
)

func GetBingXCandleEventsHandler(
	eventCallback func(event workers.CandleEvent),
) func(event *bingxgo.WsKlineEvent) {
	return func(event *bingxgo.WsKlineEvent) {
		if event == nil {
			return
		}

		if !event.Completed {
			return
		}

		eventCallback(mappers.ConvertWsKline(event))
	}
}
