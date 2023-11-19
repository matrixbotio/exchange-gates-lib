package binance

import (
	"context"
	"fmt"

	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/errs"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
)

func (a *adapter) GetPairLastPrice(pairSymbol string) (float64, error) {
	prices, err := a.binanceAPI.GetPrices(context.Background(), pairSymbol)
	if err != nil {
		return 0, fmt.Errorf("get pair last price: %w", err)
	}

	lastPrice, err := mappers.GetPairPrice(prices, pairSymbol)
	if err != nil {
		return 0, fmt.Errorf("convert pair price: %w", err)
	}
	return lastPrice, nil
}

func (a *adapter) CancelPairOrder(pairSymbol string, orderID int64, ctx context.Context) error {
	err := a.binanceAPI.CancelOrderByID(
		context.Background(),
		pairSymbol,
		orderID,
	)

	return mappers.MapCancelOrderError(err)
}

func (a *adapter) CancelPairOrderByClientOrderID(
	pairSymbol string,
	clientOrderID string,
	ctx context.Context,
) error {
	err := a.binanceAPI.CancelOrderByClientOrderID(
		context.Background(),
		pairSymbol,
		clientOrderID,
	)

	return mappers.MapCancelOrderError(err)
}

func (a *adapter) GetPairData(pairSymbol string) (structs.ExchangePairData, error) {
	exchangeInfo, err := a.binanceAPI.GetExchangeInfo(context.Background(), pairSymbol)
	if err != nil {
		return structs.ExchangePairData{}, fmt.Errorf("get exchange info: %w", err)
	}

	for _, symbolData := range exchangeInfo.Symbols {
		if symbolData.Symbol == pairSymbol {
			return mappers.ConvertExchangePairData(symbolData, a.ExchangeID)
		}
	}

	return structs.ExchangePairData{},
		fmt.Errorf("data for %q pair not found", pairSymbol)
}

func (a *adapter) GetPairOpenOrders(pairSymbol string) ([]structs.OrderData, error) {
	ordersRaw, err := a.binanceAPI.GetOpenOrders(context.Background(), pairSymbol)
	if err != nil {
		return nil, fmt.Errorf("get open orders: %w", err)
	}

	orders, err := mappers.ConvertOrders(ordersRaw)
	if err != nil {
		return nil, fmt.Errorf("convert orders: %w", err)
	}
	return orders, nil
}

func (a *adapter) GetPairs() ([]structs.ExchangePairData, error) {
	pairsResponse, err := a.binanceAPI.GetExchangeInfo(context.Background(), "")
	if err != nil {
		return nil, fmt.Errorf("get pairs: %w", err)
	}

	if pairsResponse == nil {
		return nil, errs.ErrPairResponseEmpty
	}

	return mappers.ConvertExchangePairsData(*pairsResponse, a.ExchangeID)
}
