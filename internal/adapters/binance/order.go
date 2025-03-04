package binance

import (
	"context"
	"fmt"
	"strings"

	"github.com/adshao/go-binance/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/errs"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	pkgErrs "github.com/matrixbotio/exchange-gates-lib/pkg/errs"
)

func (a *adapter) GetOrderData(
	pairSymbol string,
	orderID int64,
) (structs.OrderData, error) {
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
			return structs.OrderData{}, pkgErrs.ErrOrderNotFound
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
			return structs.OrderData{}, pkgErrs.ErrOrderNotFound
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

	var orderResponse *binance.CreateOrderResponse
	if order.IsMarketOrder {
		orderResponse, err = a.binanceAPI.PlaceMarketOrder(
			ctx,
			order.PairSymbol,
			orderSide,
			order.Qty,
			order.Price,
			order.ClientOrderID,
		)
	} else {
		orderResponse, err = a.binanceAPI.PlaceLimitOrder(
			ctx,
			order.PairSymbol,
			orderSide,
			order.Qty,
			order.Price,
			order.ClientOrderID,
		)
	}

	if err != nil {
		if strings.Contains(err.Error(), errs.ErrMsgOrderDuplicate) {
			return structs.CreateOrderResponse{}, pkgErrs.ErrOrderDuplicate
		}

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

func (a *adapter) GetHistoryOrder(
	pairSymbol string,
	orderID int64,
) (structs.OrderHistory, error) {
	// not emplemented yet
	return structs.OrderHistory{}, nil
}

func (a *adapter) GetOrderExecFee(
	baseAssetTicker string,
	quoteAssetTicker string,
	_ consts.OrderSide,
	orderID int64,
) (structs.OrderFees, error) {
	pairSymbol := baseAssetTicker + quoteAssetTicker

	trades, err := a.binanceAPI.GetOrderTradeHistory(
		context.Background(),
		orderID,
		pairSymbol,
	)
	if err != nil {
		return structs.OrderFees{}, fmt.Errorf("get order trade history: %w", err)
	}

	fees, err := mappers.GetFeesFromTradeList(
		trades,
		baseAssetTicker,
		quoteAssetTicker,
		orderID,
	)
	if err != nil {
		return structs.OrderFees{}, fmt.Errorf("convert fees: %w", err)
	}
	return fees, nil
}
