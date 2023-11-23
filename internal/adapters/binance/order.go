package binance

import (
	"context"
	"fmt"

	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/errs"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	pkgErrs "github.com/matrixbotio/exchange-gates-lib/pkg/errs"
	"github.com/shopspring/decimal"
)

func (a *adapter) GetOrderData(pairSymbol string, orderID int64) (structs.OrderData, error) {
	if orderID == 0 {
		return structs.OrderData{}, errs.ErrOrderIDNotSet
	}

	order, err := a.binanceAPI.GetOrderDataByOrderID(
		context.Background(),
		pairSymbol,
		orderID,
	)
	if err != nil {
		if errs.IsErrorAboutUnknownOrder(err) {
			return structs.OrderData{}, pkgErrs.OrderNotFound
		}

		return structs.OrderData{}, err
	}

	result, err := mappers.ConvertOrderData(order)
	if err != nil {
		return structs.OrderData{}, fmt.Errorf("convert order: %w", err)
	}
	return result, nil
}

func (a *adapter) GetOrderByClientOrderID(pairSymbol string, clientOrderID string) (
	structs.OrderData,
	error,
) {
	if clientOrderID == "" {
		return structs.OrderData{}, errs.ErrClientOrderIDNotSet
	}

	order, err := a.binanceAPI.GetOrderDataByClientOrderID(
		context.Background(),
		pairSymbol,
		clientOrderID,
	)
	if err != nil {
		if errs.IsErrorAboutUnknownOrder(err) {
			return structs.OrderData{}, pkgErrs.OrderNotFound
		}

		return structs.OrderData{}, err
	}

	result, err := mappers.ConvertOrderData(order)
	if err != nil {
		return structs.OrderData{}, fmt.Errorf("convert order: %w", err)
	}
	return result, nil
}

func (a *adapter) PlaceOrder(ctx context.Context, order structs.BotOrderAdjusted) (
	structs.CreateOrderResponse,
	error,
) {
	orderSide, err := mappers.GetBinanceOrderSide(order.Type)
	if err != nil {
		return structs.CreateOrderResponse{}, fmt.Errorf("get order side: %w", err)
	}

	orderResponse, err := a.binanceAPI.PlaceLimitOrder(
		ctx,
		order.PairSymbol,
		orderSide,
		order.Qty,
		order.Price,
		order.ClientOrderID,
	)
	if err != nil {
		return structs.CreateOrderResponse{}, fmt.Errorf("create order: %w", err)
	}

	if orderResponse == nil {
		return structs.CreateOrderResponse{}, errs.ErrOrderResponseEmpty
	}

	orderConverted, err := mappers.ConvertPlacedOrder(*orderResponse)
	if err != nil {
		return structs.CreateOrderResponse{}, fmt.Errorf("convert order: %w", err)
	}
	return orderConverted, nil
}

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
