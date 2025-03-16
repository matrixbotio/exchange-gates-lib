package mappers

import (
	"fmt"
	"strconv"

	bingxgo "github.com/matrixbotio/go-bingx"

	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
)

func ConvertOrderEvent(o *bingxgo.WsOrder) (workers.TradeEventPrivate, error) {
	orderPrice, err := strconv.ParseFloat(o.Price, 64)
	if err != nil {
		return workers.TradeEventPrivate{},
			fmt.Errorf("parse price: %w", err)
	}

	orderQty, err := strconv.ParseFloat(o.Quantity, 64)
	if err != nil {
		return workers.TradeEventPrivate{},
			fmt.Errorf("parse qty: %w", err)
	}

	return workers.TradeEventPrivate{
		ID:            o.TransactionID,
		Time:          int64(o.Timestamp),
		ExchangeTag:   consts.BingXAdapterTag,
		Symbol:        o.Symbol,
		OrderID:       strconv.Itoa(o.OrderID),
		ClientOrderID: o.ClientOrderID,
		Price:         orderPrice,
		Quantity:      orderQty,
	}, nil
}
