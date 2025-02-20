package bingx

import (
	"context"
	"fmt"

	bingxgo "github.com/Sagleft/go-bingx"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bingx/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	"github.com/shopspring/decimal"
)

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
		return structs.OrderFees{},
			fmt.Errorf("get order data: %w", err)
	}

	fees := structs.OrderFees{
		BaseAsset:  decimal.Zero,
		QuoteAsset: decimal.Zero,
	}

	feeVal, err := decimal.NewFromString(data.Fee)
	if err != nil {
		return fees, fmt.Errorf("parse: %w", err)
	}

	if data.FeeAsset == baseAssetTicker {
		fees.BaseAsset = feeVal
	}
	if data.FeeAsset == quoteAssetTicker {
		fees.QuoteAsset = feeVal
	}

	return fees, nil
}

func (a *adapter) GetOrderData(
	pairSymbol string,
	orderID int64,
) (structs.OrderData, error) {
	data, err := a.client.GetOrder(pairSymbol, orderID)
	if err != nil {
		return structs.OrderData{}, fmt.Errorf("get: %w", err)
	}

	return mappers.ConvertBingXOrderData(data)
}

func (a *adapter) GetOrderByClientOrderID(
	pairSymbol string,
	clientOrderID string,
) (structs.OrderData, error) {
	data, err := a.client.GetOrderByClientOrderID(
		pairSymbol, clientOrderID,
	)
	if err != nil {
		return structs.OrderData{}, fmt.Errorf("get: %w", err)
	}

	return mappers.ConvertBingXOrderData(data)
}
