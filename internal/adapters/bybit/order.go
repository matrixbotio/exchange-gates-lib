package bybit

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hirokisan/bybit/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bybit/helpers"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bybit/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
)

func (a *adapter) GetOrderData(pairSymbol string, orderID int64) (structs.OrderData, error) {
	orderIDFormatted := strconv.FormatInt(orderID, 10)

	return a.getOrderDataByParams(bybit.V5GetHistoryOrdersParam{
		Category: bybit.CategoryV5Spot,
		Symbol:   helpers.GetPairSymbolPointerV5(pairSymbol),
		OrderID:  &orderIDFormatted,
	})
}

func (a *adapter) GetOrderByClientOrderID(pairSymbol, clientOrderID string) (
	structs.OrderData,
	error,
) {
	return a.getOrderDataByParams(bybit.V5GetHistoryOrdersParam{
		Category:    bybit.CategoryV5Spot,
		Symbol:      helpers.GetPairSymbolPointerV5(pairSymbol),
		OrderLinkID: &clientOrderID,
	})
}

func (a *adapter) PlaceOrder(
	ctx context.Context,
	order structs.BotOrderAdjusted,
) (structs.CreateOrderResponse, error) {
	response, err := a.client.V5().Order().CreateOrder(bybit.V5CreateOrderParam{
		Category:    bybit.CategoryV5Spot,
		Symbol:      bybit.SymbolV5(order.PairSymbol),
		Side:        mappers.ConvertOrderSideToBybit(order.Type),
		OrderType:   bybit.OrderTypeLimit,
		Qty:         order.Qty,
		Price:       &order.Price,
		OrderLinkID: &order.ClientOrderID,
	})
	if err != nil {
		return structs.CreateOrderResponse{}, fmt.Errorf("create: %w", err)
	}

	orderID, err := strconv.ParseInt(response.Result.OrderID, 10, 64)
	if err != nil {
		return structs.CreateOrderResponse{}, fmt.Errorf("parse order ID: %w", err)
	}

	orderData, err := a.GetOrderData(order.PairSymbol, orderID)
	if err != nil {
		return structs.CreateOrderResponse{}, fmt.Errorf("get order data: %w", err)
	}

	return structs.CreateOrderResponse{
		OrderID:       orderID,
		ClientOrderID: orderData.ClientOrderID,
		OrigQuantity:  orderData.AwaitQty,
		Price:         orderData.Price,
		Symbol:        orderData.Symbol,
		Type:          orderData.Type,
	}, nil
}

func (a *adapter) getOrderDataByParams(param bybit.V5GetHistoryOrdersParam) (
	structs.OrderData,
	error,
) {
	orderID := helpers.GetOrderIDFromHistoryOrdersParam(param)
	pairSymbol := helpers.GetOrderSymbolFromHistoryOrdersParam(param)

	r, err := a.client.V5().Order().GetHistoryOrders(param)
	if err != nil {
		return structs.OrderData{}, fmt.Errorf(
			"get %q order in %q: %w",
			orderID, pairSymbol, err,
		)
	}

	data, err := mappers.ParseHistoryOrder(r, orderID, pairSymbol)
	if err != nil {
		return structs.OrderData{}, fmt.Errorf("parse history order: %w", err)
	}
	return data, nil
}
