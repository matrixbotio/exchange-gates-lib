package bybit

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hirokisan/bybit/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bybit/helpers/accessors"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bybit/helpers/errs"
	order_mappers "github.com/matrixbotio/exchange-gates-lib/internal/adapters/bybit/helpers/mappers/order"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
)

func (a *adapter) GetOrderData(pairSymbol string, orderID int64) (structs.OrderData, error) {
	orderIDFormatted := strconv.FormatInt(orderID, 10)

	return a.getOrderDataByParams(bybit.V5GetHistoryOrdersParam{
		Category: bybit.CategoryV5Spot,
		Symbol:   accessors.GetPairSymbolPointerV5(pairSymbol),
		OrderID:  &orderIDFormatted,
	})
}

func (a *adapter) GetOrderByClientOrderID(pairSymbol, clientOrderID string) (
	structs.OrderData,
	error,
) {
	return a.getOrderDataByParams(bybit.V5GetHistoryOrdersParam{
		Category:    bybit.CategoryV5Spot,
		Symbol:      accessors.GetPairSymbolPointerV5(pairSymbol),
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
		Side:        order_mappers.ConvertOrderSideToBybit(order.Type),
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
	orderID := accessors.GetOrderIDFromHistoryOrdersParam(param)
	pairSymbol := accessors.GetOrderSymbolFromHistoryOrdersParam(param)

	r, err := a.client.V5().Order().GetHistoryOrders(param)
	if err != nil {
		return structs.OrderData{}, fmt.Errorf(
			"get %q order in %q: %w",
			orderID, pairSymbol, err,
		)
	}

	data, err := order_mappers.ParseHistoryOrder(r, orderID, pairSymbol)
	if err != nil {
		return structs.OrderData{}, fmt.Errorf("parse history order: %w", err)
	}
	return data, nil
}

func (a *adapter) CancelPairOrder(pairSymbol string, orderID int64, ctx context.Context) error {
	orderIDFormatted := strconv.FormatInt(orderID, 10)

	_, err := a.client.V5().Order().CancelOrder(bybit.V5CancelOrderParam{
		Category: bybit.CategoryV5Spot,
		Symbol:   bybit.SymbolV5(pairSymbol),
		OrderID:  &orderIDFormatted,
	})
	if err != nil {
		return errs.HandleCancelOrderError(orderIDFormatted, pairSymbol, err)
	}
	return nil
}

func (a *adapter) CancelPairOrderByClientOrderID(
	pairSymbol string,
	clientOrderID string,
	ctx context.Context,
) error {
	_, err := a.client.V5().Order().CancelOrder(bybit.V5CancelOrderParam{
		Category:    bybit.CategoryV5Spot,
		Symbol:      bybit.SymbolV5(pairSymbol),
		OrderLinkID: &clientOrderID,
	})
	if err != nil {
		return errs.HandleCancelOrderError(clientOrderID, pairSymbol, err)
	}
	return nil
}

func (a *adapter) GetOrderExecFee(
	pairSymbol string,
	orderSide string,
	orderID int64,
) (structs.OrderFees, error) {
	if pairSymbol == "" {
		return structs.OrderFees{}, errors.New("pair symbol is not set")
	}
	if orderID == 0 {
		return structs.OrderFees{}, errors.New("order ID is not set")
	}

	orderIDFormatted := strconv.FormatInt(orderID, 10)

	payload := bybit.V5GetExecutionParam{
		Category: bybit.CategoryV5Spot,
		Symbol:   accessors.GetPairSymbolPointerV5(pairSymbol),
		OrderId:  &orderIDFormatted,
	}

	orderExecData, err := a.client.V5().Execution().GetExecutionList(payload)
	if err != nil {
		return structs.OrderFees{}, fmt.Errorf("get order execution history: %w", err)
	}

	fees, err := order_mappers.ParseOrderExecFee(orderExecData.Result, orderSide)
	if err != nil {
		return structs.OrderFees{}, fmt.Errorf("parse order fees: %w", err)
	}
	return fees, nil
}
