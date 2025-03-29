package gate

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gateio/gateapi-go/v6"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/gate/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	"github.com/matrixbotio/exchange-gates-lib/pkg/errs"
	"github.com/shopspring/decimal"
)

const orderTypeLimit = "limit"

func (a *adapter) GetOrderData(pairSymbol string, orderID int64) (structs.OrderData, error) {
	if !a.creds.Keypair.IsSet() {
		return structs.OrderData{}, errs.ErrAPIKeyNotSet
	}

	return a.GetOrderByClientOrderID(
		pairSymbol,
		strconv.FormatInt(orderID, 10),
	)
}

func (a *adapter) GetOrderByClientOrderID(
	pairSymbol string,
	clientOrderID string,
) (structs.OrderData, error) {
	if !a.creds.Keypair.IsSet() {
		return structs.OrderData{}, errs.ErrAPIKeyNotSet
	}

	ctx, ctxCancel := context.WithTimeout(a.auth, requestTimeout)
	defer ctxCancel()

	data, _, err := a.client.SpotApi.GetOrder(
		ctx,
		clientOrderID,
		pairSymbol,
		nil,
	)
	if err != nil {
		//if strings.Contains(err.Error(), mappers.ErrOrderNotActualMessage) {
		//	return structs.OrderData{}, errs.ErrOrderDataNotActual
		//}

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

	if orderData.FilledQty > 0 && orderData.FilledQty < orderData.AwaitQty {
		orderData.Status = consts.OrderStatusPartiallyFilled
	}

	return orderData, nil
}

func (a *adapter) PlaceOrder(
	_ context.Context,
	order structs.BotOrderAdjusted,
) (structs.CreateOrderResponse, error) {
	if !a.creds.Keypair.IsSet() {
		return structs.CreateOrderResponse{}, errs.ErrAPIKeyNotSet
	}

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
	if !a.creds.Keypair.IsSet() {
		return structs.OrderFees{}, errs.ErrAPIKeyNotSet
	}

	ctx, ctxCancel := context.WithTimeout(a.auth, requestTimeout)
	defer ctxCancel()

	pairSymbol := a.GetPairSymbol(baseAssetTicker, quoteAssetTicker)

	data, _, err := a.client.SpotApi.GetOrder(
		ctx,
		strconv.FormatInt(orderID, 10),
		pairSymbol,
		nil,
	)
	if err != nil {
		if strings.Contains(err.Error(), mappers.ErrOrderNotActualMessage) {
			return structs.OrderFees{}, errs.ErrOrderDataNotActual
		}

		return structs.OrderFees{}, fmt.Errorf("get: %w", err)
	}

	fees, err := mappers.GetOrderFees(data, baseAssetTicker, quoteAssetTicker)
	if err != nil {
		return structs.OrderFees{}, fmt.Errorf("convert: %w", err)
	}
	return fees, nil
}

func (a *adapter) GetHistoryOrder(
	pairSymbol string,
	orderID int64,
) (structs.OrderHistory, error) {
	/*if !a.creds.Keypair.IsSet() {
		return structs.OrderHistory{}, errs.ErrAPIKeyNotSet
	}

	ctx, ctxCancel := context.WithTimeout(a.auth, requestTimeout)
	defer ctxCancel()

	events, _, err := a.client.SpotApi.ListMyTrades(ctx, &gateapi.ListMyTradesOpts{
		CurrencyPair: optional.NewString(pairSymbol),
		OrderId:      optional.NewString(strconv.FormatInt(orderID, 10)),
	})
	if err != nil {
		return structs.OrderHistory{}, fmt.Errorf("get: %w", err)
	}

	result, err := mappers.ConvertTradesToOrderHistory(events)
	if err != nil {
		return structs.OrderHistory{}, fmt.Errorf("convert: %w", err)
	}
	return result, nil*/

	// not ready yet
	return structs.OrderHistory{}, errors.New("not implemented")
}
