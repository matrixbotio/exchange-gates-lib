package mappers

import (
	"fmt"
	"strconv"

	"github.com/hirokisan/bybit/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
)

func ParseTradeEventPrivate(
	event bybit.V5WebsocketPrivateExecutionData,
	eventTime int64,
	exchangeTag string,
) (
	workers.TradeEventPrivate,
	error,
) {
	wEvent := workers.TradeEventPrivate{
		ID:            event.BlockTradeID,
		Time:          eventTime,
		ExchangeTag:   exchangeTag,
		Symbol:        string(event.Symbol),
		OrderID:       event.OrderID,
		ClientOrderID: event.OrderLinkID,
	}

	var err error
	wEvent.Price, err = strconv.ParseFloat(event.ExecPrice, 64)
	if err != nil {
		return workers.TradeEventPrivate{},
			fmt.Errorf("parse price: %w", err)
	}

	wEvent.Quantity, err = strconv.ParseFloat(event.ExecQty, 64)
	if err != nil {
		return workers.TradeEventPrivate{},
			fmt.Errorf("parse quantity: %w", err)
	}

	return wEvent, nil
}
