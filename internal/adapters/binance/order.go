package binance

import (
	"context"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	"github.com/shopspring/decimal"
)

func (a *adapter) GetOrderExecFee(
	pairSymbol string,
	orderSide string,
	orderID int64,
) (structs.OrderFees, error) {
	// TBD: https://github.com/matrixbotio/exchange-gates-lib/issues/167

	return structs.OrderFees{
		BaseAsset:  decimal.NewFromInt(0),
		QuoteAsset: decimal.NewFromInt(0),
	}, nil
}

// CancelPairOrder - cancel one exchange pair order by ID
func (a *adapter) CancelPairOrder(pairSymbol string, orderID int64, ctx context.Context) error {
	_, err := a.binanceAPI.NewCancelOrderService().Symbol(pairSymbol).
		OrderID(orderID).Do(ctx)
	if err != nil {
		return mappers.MapCancelOrderError(err)
	}
	return nil
}

// CancelPairOrder - cancel one exchange pair order by client order ID
func (a *adapter) CancelPairOrderByClientOrderID(
	pairSymbol string,
	clientOrderID string,
	ctx context.Context,
) error {
	_, err := a.binanceAPI.NewCancelOrderService().Symbol(pairSymbol).
		OrigClientOrderID(clientOrderID).Do(ctx)
	if err != nil {
		return mappers.MapCancelOrderError(err)
	}
	return nil
}
