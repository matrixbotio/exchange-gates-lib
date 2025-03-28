package utils

import (
	"fmt"
	"strconv"
	"time"

	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
)

// OrderDataToTradeEvent data
type TradeOrderConvertTask struct {
	Order       structs.OrderData
	ExchangeTag string
}

// OrderDataToBotOrder - convert order data to bot order
func OrderDataToBotOrder(order structs.OrderData) pkgStructs.BotOrder {
	return pkgStructs.BotOrder{
		PairSymbol:    order.Symbol,
		Type:          order.Side,
		Qty:           order.AwaitQty,
		Price:         order.Price,
		Deposit:       order.AwaitQty * order.Price,
		ClientOrderID: order.ClientOrderID,
	}
}

// OrderResponseToBotOrder - convert raw order response to bot order
func OrderResponseToBotOrder(response structs.CreateOrderResponse) pkgStructs.BotOrder {
	return pkgStructs.BotOrder{
		PairSymbol:    response.Symbol,
		Type:          response.Type,
		Qty:           response.OrigQuantity,
		Price:         response.Price,
		Deposit:       response.OrigQuantity * response.Price,
		ClientOrderID: response.ClientOrderID,
	}
}

func OrderDataToCreateOrderResponse(
	data structs.OrderData,
	orderID int64,
) structs.CreateOrderResponse {
	return structs.CreateOrderResponse{
		OrderID:       orderID,
		ClientOrderID: data.ClientOrderID,
		OrigQuantity:  data.AwaitQty,
		Price:         data.Price,
		Symbol:        data.Symbol,
		Type:          data.Side,
		CreatedTime:   data.CreatedTime,
		Status:        data.Status,
	}
}

func OrderToOrderResponse(order structs.BotOrderAdjusted, orderID int64) (
	structs.CreateOrderResponse,
	error,
) {
	data := structs.CreateOrderResponse{
		OrderID:       orderID,
		ClientOrderID: order.ClientOrderID,
		Symbol:        order.PairSymbol,
		Type:          order.Type,
		CreatedTime:   time.Now().UnixMilli(),
		Status:        pkgStructs.OrderStatusNew,
	}

	var err error
	data.OrigQuantity, err = strconv.ParseFloat(order.Qty, 64)
	if err != nil {
		return structs.CreateOrderResponse{}, fmt.Errorf("parse qty: %w", err)
	}

	data.Price, err = strconv.ParseFloat(order.Price, 64)
	if err != nil {
		return structs.CreateOrderResponse{}, fmt.Errorf("parse price: %w", err)
	}
	return data, nil
}
