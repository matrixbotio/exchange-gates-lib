package bybit

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hirokisan/bybit/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bybit/helpers/accessors"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bybit/helpers/errs"
	order_mappers "github.com/matrixbotio/exchange-gates-lib/internal/adapters/bybit/helpers/mappers/order"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	pkgErrs "github.com/matrixbotio/exchange-gates-lib/pkg/errs"
	"github.com/matrixbotio/exchange-gates-lib/pkg/utils"
)

func (a *adapter) GetOrderData(pairSymbol string, orderID int64) (structs.OrderData, error) {
	orderIDFormatted := strconv.FormatInt(orderID, 10)

	data, err := a.getOrderDataByParams(bybit.V5GetHistoryOrdersParam{
		Category: bybit.CategoryV5Spot,
		Symbol:   accessors.GetPairSymbolPointerV5(pairSymbol),
		OrderID:  &orderIDFormatted,
	})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return structs.OrderData{}, pkgErrs.ErrOrderNotFound
		}
		return structs.OrderData{}, err
	}

	return data, nil
}

func (a *adapter) GetOrderByClientOrderID(pairSymbol, clientOrderID string) (
	structs.OrderData,
	error,
) {
	data, err := a.getOrderDataByParams(bybit.V5GetHistoryOrdersParam{
		Category:    bybit.CategoryV5Spot,
		Symbol:      accessors.GetPairSymbolPointerV5(pairSymbol),
		OrderLinkID: &clientOrderID,
	})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return structs.OrderData{}, pkgErrs.ErrOrderNotFound
		}
		return structs.OrderData{}, err
	}

	return data, nil
}

func (a *adapter) PlaceOrder(
	ctx context.Context,
	order structs.BotOrderAdjusted,
) (structs.CreateOrderResponse, error) {
	data := bybit.V5CreateOrderParam{
		Category:    bybit.CategoryV5Spot,
		Symbol:      bybit.SymbolV5(order.PairSymbol),
		Side:        order_mappers.ConvertOrderSideToBybit(order.Type),
		OrderType:   bybit.OrderTypeLimit,
		Qty:         order.Qty,
		Price:       &order.Price,
		OrderLinkID: &order.ClientOrderID,
	}

	if order.IsMarketOrder {
		data.OrderType = bybit.OrderTypeMarket
		data.Price = nil
	}

	response, err := a.client.V5().Order().CreateOrder(data)
	if err != nil {
		// if the order has already been placed, we will receive and return its data
		if strings.Contains(err.Error(), errs.ErrMsgOrderDuplicate) {
			return structs.CreateOrderResponse{}, pkgErrs.ErrOrderDuplicate
		}

		return structs.CreateOrderResponse{}, fmt.Errorf("create: %w", err)
	}

	orderID, err := strconv.ParseInt(response.Result.OrderID, 10, 64)
	if err != nil {
		return structs.CreateOrderResponse{}, fmt.Errorf("parse order ID: %w", err)
	}

	orderData, err := a.getOrderDataByParams(bybit.V5GetHistoryOrdersParam{
		Category: bybit.CategoryV5Spot,
		Symbol:   accessors.GetPairSymbolPointerV5(order.PairSymbol),
		OrderID:  utils.StringPointer(strconv.FormatInt(orderID, 10)),
	})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			// order not found, return original order data
			return utils.OrderToOrderResponse(order, orderID)
		}

		return structs.CreateOrderResponse{},
			fmt.Errorf("get order data after place order: %w", err)
	}
	return utils.OrderDataToCreateOrderResponse(orderData, orderID), nil
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
		if strings.Contains(err.Error(), "not found") {
			// history order not available, let's find in opened orders
			return a.getOpenedOrder(param)
		}

		return structs.OrderData{}, fmt.Errorf("parse history order: %w", err)
	}
	return data, nil
}

func (a *adapter) getOpenedOrder(param bybit.V5GetHistoryOrdersParam) (
	structs.OrderData,
	error,
) {
	orderID := accessors.GetOrderIDFromHistoryOrdersParam(param)
	pairSymbol := accessors.GetOrderSymbolFromHistoryOrdersParam(param)

	r, err := a.client.V5().Order().GetOpenOrders(bybit.V5GetOpenOrdersParam{
		Category:    param.Category,
		Symbol:      param.Symbol,
		OrderID:     param.OrderID,
		OrderLinkID: param.OrderLinkID,
	})
	if err != nil {
		return structs.OrderData{}, fmt.Errorf("get open orders: %w", err)
	}

	orderData, err := order_mappers.ParseHistoryOrder(r, orderID, pairSymbol)
	if err != nil {
		return structs.OrderData{}, fmt.Errorf("parse: %w", err)
	}
	return orderData, nil
}

func (a *adapter) CancelPairOrder(pairSymbol string, orderID int64, ctx context.Context) error {
	orderIDFormatted := strconv.FormatInt(orderID, 10)

	_, err := a.client.V5().Order().CancelOrder(bybit.V5CancelOrderParam{
		Category: bybit.CategoryV5Spot,
		Symbol:   bybit.SymbolV5(pairSymbol),
		OrderID:  &orderIDFormatted,
	})
	if err != nil {
		return errs.MapCancelOrderError(orderIDFormatted, pairSymbol, err)
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
		return errs.MapCancelOrderError(clientOrderID, pairSymbol, err)
	}
	return nil
}

func (a *adapter) GetOrderExecFee(
	baseAssetTicker string,
	quoteAssetTicker string,
	orderSide consts.OrderSide,
	orderID int64,
) (structs.OrderFees, error) {
	if baseAssetTicker == "" {
		return structs.OrderFees{}, errors.New("base asset ticker is not set")
	}
	if quoteAssetTicker == "" {
		return structs.OrderFees{}, errors.New("quote asset ticker is not set")
	}

	pairSymbol := baseAssetTicker + quoteAssetTicker
	if orderID == 0 {
		return structs.OrderFees{}, errors.New("order ID is not set")
	}

	orderIDFormatted := strconv.FormatInt(orderID, 10)

	payload := bybit.V5GetExecutionParam{
		Category: bybit.CategoryV5Spot,
		Symbol:   accessors.GetPairSymbolPointerV5(pairSymbol),
		OrderID:  &orderIDFormatted,
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

func (a *adapter) GetHistoryOrder(
	pairSymbol string,
	orderID int64,
) (structs.OrderHistory, error) {
	// not emplemented yet
	return structs.OrderHistory{}, nil
}
