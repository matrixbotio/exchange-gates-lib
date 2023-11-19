package helpers

import (
	"fmt"

	"github.com/adshao/go-binance/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
)

func GetTradeEventsHandler(
	exchangeTag string,
	eventCallback func(event workers.TradeEvent),
	errorHandler func(err error),
) func(event *binance.WsTradeEvent) {
	return func(event *binance.WsTradeEvent) {
		if event == nil {
			return
		}

		wEvent, err := mappers.ConvertTradeEvent(*event, exchangeTag)
		if err != nil {
			errorHandler(fmt.Errorf("convert trade event: %w", err))
			return
		}

		eventCallback(wEvent)
	}
}
