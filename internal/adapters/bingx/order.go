package bingx

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bingx/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	"github.com/matrixbotio/exchange-gates-lib/pkg/errs"
	bingxgo "github.com/matrixbotio/go-bingx"
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

	feeValue, err := decimal.NewFromString(data.Fee)
	if err != nil {
		return structs.OrderFees{}, fmt.Errorf("parse fee: %w", err)
	}

	return getFeesFromOrderData(
		orderSide,
		feeValue.Abs(),
	), nil
}

func getFeesFromOrderData(
	orderSide consts.OrderSide,
	fee decimal.Decimal,
) structs.OrderFees {
	fees := structs.OrderFees{
		BaseAsset:  decimal.Zero,
		QuoteAsset: decimal.Zero,
	}

	if orderSide == consts.OrderSideBuy {
		fees.BaseAsset = fee
	}
	if orderSide == consts.OrderSideSell {
		fees.QuoteAsset = fee
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
	pairSymbol string,
	orderID int64,
) (structs.OrderHistory, error) {
	order, err := a.client.GetHistoryOrder(pairSymbol, orderID)
	if err != nil {
		return structs.OrderHistory{}, fmt.Errorf("get: %w", err)
	}

	orderData, err := mappers.ConvertBingXHistoryOrder(order)
	if err != nil {
		return structs.OrderHistory{}, fmt.Errorf("convert: %w", err)
	}

	return structs.OrderHistory{
		OrderData: orderData,
		Fees: getFeesFromOrderData(
			orderData.Side,
			decimal.NewFromFloat(order.Fee),
		),
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
