package mappers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/adshao/go-binance/v2"

	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
)

func ConvertPriceEvent(event binance.WsBookTickerEvent) (ask, bid float64, err error) {
	ask, err = strconv.ParseFloat(event.BestAskPrice, 64)
	if err != nil {
		return
	}

	bid, err = strconv.ParseFloat(event.BestBidPrice, 64)
	return
}

func fixTradeEventEndTime(endTime int64) int64 {
	if strings.HasSuffix(strconv.FormatInt(endTime, 10), "999") {
		endTime++
	}
	return endTime
}

func ConvertTradeEvent(
	event binance.WsTradeEvent,
	exchangeTag string,
) (workers.TradeEvent, error) {
	wEvent := workers.TradeEvent{
		ID:            event.TradeID,
		Time:          fixTradeEventEndTime(event.Time),
		Symbol:        event.Symbol,
		ExchangeTag:   exchangeTag,
		BuyerOrderID:  event.BuyerOrderID,
		SellerOrderID: event.SellerOrderID,
	}

	var err error
	wEvent.Price, err = strconv.ParseFloat(event.Price, 64)
	if err != nil {
		return workers.TradeEvent{},
			fmt.Errorf("parse price: %w", err)
	}

	wEvent.Quantity, err = strconv.ParseFloat(event.Quantity, 64)
	if err != nil {
		return workers.TradeEvent{},
			fmt.Errorf("parse qty: %w", err)
	}
	return wEvent, nil
}

func ConvertTradeEventPrivate(event binance.WsUserDataEvent, exchangeTag string) (workers.TradeEventPrivate, error) {
	wEvent := workers.TradeEventPrivate{
		ID:            strconv.FormatInt(event.OrderUpdate.TradeId, 10),
		Time:          fixTradeEventEndTime(event.Time),
		ExchangeTag:   exchangeTag,
		Symbol:        event.OrderUpdate.Symbol,
		OrderID:       strconv.FormatInt(event.OrderUpdate.Id, 10),
		ClientOrderID: event.OrderUpdate.ClientOrderId,
	}

	var err error
	wEvent.Price, err = strconv.ParseFloat(event.OrderUpdate.Price, 64)
	if err != nil {
		return workers.TradeEventPrivate{},
			fmt.Errorf("parse price: %w", err)
	}

	wEvent.Quantity, err = strconv.ParseFloat(event.OrderUpdate.Volume, 64)
	if err != nil {
		return workers.TradeEventPrivate{},
			fmt.Errorf("parse quantity: %w", err)
	}

	return wEvent, nil
}
