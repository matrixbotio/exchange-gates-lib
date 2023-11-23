package helpers

import (
	"strconv"
	"strings"

	"github.com/adshao/go-binance/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
	"github.com/matrixbotio/exchange-gates-lib/pkg/utils"
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

		// fix event.Time
		if strings.HasSuffix(strconv.FormatInt(event.Time, 10), "999") {
			event.Time++
		}
		wEvent := workers.TradeEvent{
			ID:            event.TradeID,
			Time:          event.Time,
			Symbol:        event.Symbol,
			ExchangeTag:   exchangeTag,
			BuyerOrderID:  event.BuyerOrderID,
			SellerOrderID: event.SellerOrderID,
		}
		errs := make([]error, 2)
		wEvent.Price, errs[0] = strconv.ParseFloat(event.Price, 64)
		wEvent.Quantity, errs[0] = strconv.ParseFloat(event.Quantity, 64)
		if utils.LogNotNilError(errs) {
			return
		}
		eventCallback(wEvent)
	}
}
