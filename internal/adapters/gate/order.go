package gate

import (
	"context"

	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
)

func (a *adapter) GetOrderData(pairSymbol string, orderID int64) (structs.OrderData, error) {
	// TODO
	return structs.OrderData{}, nil
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
