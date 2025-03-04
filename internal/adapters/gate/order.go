package gate

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gateio/gateapi-go/v6"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/gate/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	"github.com/shopspring/decimal"
)

const orderTypeLimit = "limit"

func (a *adapter) GetOrderData(pairSymbol string, orderID int64) (structs.OrderData, error) {
	return a.GetOrderByClientOrderID(
		pairSymbol,
		strconv.FormatInt(orderID, 10),
	)
}

func (a *adapter) GetOrderByClientOrderID(
	pairSymbol string,
	clientOrderID string,
) (structs.OrderData, error) {
	ctx, ctxCancel := context.WithTimeout(a.auth, requestTimeout)
	defer ctxCancel()

	data, _, err := a.client.SpotApi.GetOrder(
		ctx,
		clientOrderID,
		pairSymbol,
		&gateapi.GetOrderOpts{},
	)
	if err != nil {
		return structs.OrderData{}, fmt.Errorf("get order data: %w", err)
	}

	orderID, err := strconv.ParseInt(data.Id, 10, 64)
	if err != nil {
		return structs.OrderData{}, fmt.Errorf("parse orderID: %w", err)
	}

	orderSide, err := mappers.ConvertOrderSide(data.Side)
	if err != nil {
		return structs.OrderData{}, fmt.Errorf("convert: %w", err)
	}

	orderData := structs.OrderData{
		OrderID:       orderID,
		ClientOrderID: data.Text,
		Symbol:        pairSymbol,
		Side:          orderSide,
		CreatedTime:   data.CreateTimeMs,
		UpdatedTime:   data.UpdateTimeMs,
	}

	orderData.AwaitQty, err = strconv.ParseFloat(data.Amount, 64)
	if err != nil {
		return structs.OrderData{}, fmt.Errorf("parse qty: %w", err)
	}

	orderData.FilledQty, err = strconv.ParseFloat(data.FilledAmount, 64)
	if err != nil {
		return structs.OrderData{}, fmt.Errorf("parse filled qty: %w", err)
	}

	orderData.Price, err = strconv.ParseFloat(data.Price, 64)
	if err != nil {
		return structs.OrderData{}, fmt.Errorf("parse price: %w", err)
	}

	orderData.Status = mappers.ConvertOrderStatus(data.Status)

	if orderData.FilledQty > 0 {
		orderData.Status = consts.OrderStatusPartiallyFilled
	}

	return orderData, nil
}

func (a *adapter) PlaceOrder(
	_ context.Context,
	order structs.BotOrderAdjusted,
) (structs.CreateOrderResponse, error) {
	ctx, ctxCancel := context.WithTimeout(a.auth, requestTimeout)
	defer ctxCancel()

	response, _, err := a.client.SpotApi.CreateOrder(ctx, gateapi.Order{
		Text:         order.ClientOrderID,
		CurrencyPair: order.PairSymbol,
		Type:         orderTypeLimit,
		Account:      spotAccountType,
		Side:         string(order.Type),
		Amount:       order.Qty,
		Price:        order.Price,
	}, &gateapi.CreateOrderOpts{})
	if err != nil {
		return structs.CreateOrderResponse{}, fmt.Errorf("create order: %w", err)
	}

	orderID, err := strconv.ParseInt(response.Id, 10, 64)
	if err != nil {
		return structs.CreateOrderResponse{}, fmt.Errorf("parse order ID: %w", err)
	}

	qty, err := decimal.NewFromString(response.Amount)
	if err != nil {
		return structs.CreateOrderResponse{}, fmt.Errorf("parse order qty: %w", err)
	}

	price, err := decimal.NewFromString(response.Price)
	if err != nil {
		return structs.CreateOrderResponse{}, fmt.Errorf("parse order price: %w", err)
	}

	return structs.CreateOrderResponse{
		OrderID:       orderID,
		ClientOrderID: response.Text,
		OrigQuantity:  qty.InexactFloat64(),
		Price:         price.InexactFloat64(),
		Symbol:        order.PairSymbol,
		Type:          order.Type,
		CreatedTime:   response.CreateTimeMs,
		Status:        mappers.ConvertOrderStatus(response.Status),
	}, nil
}

func (a *adapter) GetOrderExecFee(
	baseAssetTicker string,
	quoteAssetTicker string,
	orderSide consts.OrderSide,
	orderID int64,
) (structs.OrderFees, error) {
	/*ctx, ctxCancel := context.WithTimeout(a.auth, requestTimeout)
	defer ctxCancel()

	data, _, err := a.client.SpotApi.GetOrder(
		ctx,
		strconv.FormatInt(orderID, 10),
		a.GetPairSymbol(baseAssetTicker, quoteAssetTicker),
		&gateapi.GetOrderOpts{},
	)
	if err != nil {
		return structs.OrderFees{}, fmt.Errorf("get order data: %w", err)
	}

	// TODO*/
	return structs.OrderFees{}, nil
}

func (a *adapter) GetHistoryOrder(
	baseAssetTicker string,
	quoteAssetTicker string,
	orderID int64,
) (structs.OrderHistory, error) {
	// not emplemented yet
	return structs.OrderHistory{}, nil
}
