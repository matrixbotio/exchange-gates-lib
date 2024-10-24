package gate

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gateio/gateapi-go/v6"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/gate/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
)

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
	data, _, err := a.client.SpotApi.GetOrder(
		a.auth,
		clientOrderID,
		pairSymbol,
		&gateapi.GetOrderOpts{},
	)
	if err != nil {
		return structs.OrderData{}, fmt.Errorf("get order: %w", err)
	}

	orderID, err := strconv.ParseInt(data.Id, 10, 64)
	if err != nil {
		return structs.OrderData{}, fmt.Errorf("parse orderID: %w", err)
	}

	orderData := structs.OrderData{
		OrderID:       orderID,
		ClientOrderID: data.Text,
		Symbol:        pairSymbol,
		Type:          data.Side,
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
	ctx context.Context,
	order structs.BotOrderAdjusted,
) (structs.CreateOrderResponse, error) {
	// TODO
	return structs.CreateOrderResponse{}, nil
}

func (a *adapter) GetOrderExecFee(
	baseAssetTicker string,
	quoteAssetTicker string,
	orderSide string,
	orderID int64,
) (structs.OrderFees, error) {
	// TODO
	return structs.OrderFees{}, nil
}
