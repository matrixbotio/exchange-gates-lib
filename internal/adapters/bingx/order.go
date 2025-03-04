package bingx

import (
	"context"
	"errors"
	"fmt"
	"strings"

	bingxgo "github.com/Sagleft/go-bingx"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bingx/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	"github.com/matrixbotio/exchange-gates-lib/pkg/errs"
	"github.com/shopspring/decimal"
)

const errOrderNotActualMessage = "the order is FILLED or CANCELLED already before"

func (a *adapter) PlaceOrder(
	ctx context.Context,
	order structs.BotOrderAdjusted,
) (structs.CreateOrderResponse, error) {
	orderSide, err := mappers.GetBingXOrderSide(order.Type)
	if err != nil {
		return structs.CreateOrderResponse{},
			fmt.Errorf("get order side: %w", err)
	}

	orderQty, err := decimal.NewFromString(order.Qty)
	if err != nil {
		return structs.CreateOrderResponse{},
			fmt.Errorf("parse qty: %w", err)
	}

	orderPrice, err := decimal.NewFromString(order.Price)
	if err != nil {
		return structs.CreateOrderResponse{},
			fmt.Errorf("parse price: %w", err)
	}

	response, err := a.client.CreateOrder(bingxgo.SpotOrderRequest{
		Symbol:        order.PairSymbol,
		Side:          orderSide,
		Type:          limitOrder,
		Quantity:      orderQty.InexactFloat64(),
		Price:         orderPrice.InexactFloat64(),
		TimeInForce:   limitOrderTimeInForce,
		ClientOrderID: order.ClientOrderID,
	})
	if err != nil {
		return structs.CreateOrderResponse{},
			fmt.Errorf("create: %w", err)
	}

	result, err := mappers.ConvertOrderResponse(response)
	if err != nil {
		return structs.CreateOrderResponse{},
			fmt.Errorf("convert: %w", err)
	}
	return result, nil
}

func (a *adapter) GetOrderExecFee(
	baseAssetTicker string,
	quoteAssetTicker string,
	orderSide consts.OrderSide,
	orderID int64,
) (structs.OrderFees, error) {
	pairSymbol := a.GetPairSymbol(baseAssetTicker, quoteAssetTicker)

	data, err := a.client.GetOrder(pairSymbol, orderID)
	if err != nil {
		if strings.Contains(err.Error(), errOrderNotActualMessage) {
			return structs.OrderFees{}, errs.ErrOrderDataNotActual
		}

		return structs.OrderFees{},
			fmt.Errorf("get order data: %w", err)
	}

	if data == nil {
		return structs.OrderFees{}, errors.New("order data not set")
	}

	return getFeesFromOrderData(orderSide, data.Fee), nil
}

func getFeesFromOrderData(
	orderSide consts.OrderSide,
	fee float64,
) structs.OrderFees {
	fees := structs.OrderFees{
		BaseAsset:  decimal.Zero,
		QuoteAsset: decimal.Zero,
	}

	feeVal := decimal.NewFromFloat(fee).Abs()

	if orderSide == consts.OrderSideBuy {
		fees.BaseAsset = feeVal
	}
	if orderSide == consts.OrderSideSell {
		fees.QuoteAsset = feeVal
	}
	return fees
}

func (a *adapter) GetOrderData(
	pairSymbol string,
	orderID int64,
) (structs.OrderData, error) {
	data, err := a.client.GetOrder(pairSymbol, orderID)
	if err != nil {
		if strings.Contains(err.Error(), errOrderNotActualMessage) {
			return structs.OrderData{}, errs.ErrOrderDataNotActual
		}

		return structs.OrderData{}, fmt.Errorf("get: %w", err)
	}

	return mappers.ConvertBingXOrderData(data)
}

func (a *adapter) GetHistoryOrder(
	baseAssetTicker string,
	quoteAssetTicker string,
	orderID int64,
) (structs.OrderHistory, error) {
	pairSymbol := a.GetPairSymbol(baseAssetTicker, quoteAssetTicker)

	order, err := a.client.GetHistoryOrder(pairSymbol, orderID)
	if err != nil {
		return structs.OrderHistory{}, fmt.Errorf("get: %w", err)
	}

	orderData, err := mappers.ConvertBingXOrderData(&order)
	if err != nil {
		return structs.OrderHistory{}, fmt.Errorf("convert: %w", err)
	}

	return structs.OrderHistory{
		OrderData: orderData,
		Fees:      getFeesFromOrderData(orderData.Side, order.Fee),
	}, nil
}

func (a *adapter) GetOrderByClientOrderID(
	pairSymbol string,
	clientOrderID string,
) (structs.OrderData, error) {
	data, err := a.client.GetOrderByClientOrderID(
		pairSymbol, clientOrderID,
	)
	if err != nil {
		if strings.Contains(err.Error(), errOrderNotActualMessage) {
			return structs.OrderData{}, errs.ErrOrderDataNotActual
		}

		return structs.OrderData{}, fmt.Errorf("get: %w", err)
	}

	return mappers.ConvertBingXOrderData(data)
}
