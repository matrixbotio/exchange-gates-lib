package gate

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gateio/gateapi-go/v6"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
)

func (a *adapter) GetOrderData(pairSymbol string, orderID int64) (structs.OrderData, error) {
	data, _, err := a.client.SpotApi.GetOrder(
		a.auth,
		strconv.FormatInt(orderID, 10),
		pairSymbol,
		&gateapi.GetOrderOpts{},
	)
	if err != nil {
		return structs.OrderData{}, fmt.Errorf("get order: %w", err)
	}

	// TODO
	return structs.OrderData{
		OrderID:       orderID,
		ClientOrderID: data.Text,
		// TODO
	}, nil
}

func (a *adapter) GetOrderByClientOrderID(
	pairSymbol string,
	clientOrderID string,
) (structs.OrderData, error) {
	// TODO
	return structs.OrderData{}, nil
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
