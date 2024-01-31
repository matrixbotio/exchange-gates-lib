package mappers

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hirokisan/bybit/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
)

func ParsePriceEvent(
	r bybit.V5WebsocketPublicTickerResponse,
	exchangeTag string,
) (workers.PriceEvent, error) {
	if r.Data.Spot == nil {
		return workers.PriceEvent{}, errors.New("price event data is empty")
	}

	if r.Data.Spot.LastPrice == "" {
		return workers.PriceEvent{}, errors.New("last price is empty")
	}
	lastPrice, err := strconv.ParseFloat(r.Data.Spot.LastPrice, 64)
	if err != nil {
		return workers.PriceEvent{}, fmt.Errorf("parse last price: %w", err)
	}

	return workers.PriceEvent{
		ExchangeTag: exchangeTag,
		Symbol:      string(r.Data.Spot.Symbol),
		Ask:         lastPrice,
		Bid:         lastPrice,
	}, nil
}

func ParseTradeEvent(event bybit.V5WebsocketPrivateOrderData, eventTime int64, exchangeTag string) (
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
	wEvent.Price, err = strconv.ParseFloat(event.Price, 64)
	if err != nil {
		return workers.TradeEventPrivate{},
			fmt.Errorf("parse price: %w", err)
	}

	wEvent.Quantity, err = strconv.ParseFloat(event.Qty, 64)
	if err != nil {
		return workers.TradeEventPrivate{},
			fmt.Errorf("parse quantity: %w", err)
	}

	wEvent.FilledQuantity, err = strconv.ParseFloat(event.CumExecQty, 64)
	if err != nil {
		return workers.TradeEventPrivate{},
			fmt.Errorf("parse filledQuantity: %w", err)
	}

	return wEvent, nil
}
