package bingx

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bingx/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	"github.com/matrixbotio/exchange-gates-lib/pkg/utils"
)

func (a *adapter) GetPairData(pairSymbol string) (structs.ExchangePairData, error) {
	symbols, err := a.client.GetSymbols(pairSymbol)
	if err != nil {
		return structs.ExchangePairData{},
			fmt.Errorf("get symbols: %w", err)
	}

	tickers, err := a.client.GetTickers()
	if err != nil {
		return structs.ExchangePairData{}, fmt.Errorf("get tickers: %w", err)
	}

	pairs, err := mappers.ConvertPairs(symbols, tickers)
	if err != nil {
		return structs.ExchangePairData{}, fmt.Errorf("convert: %w", err)
	}
	if len(pairs) == 0 {
		return structs.ExchangePairData{}, errors.New("data not found")
	}

	return pairs[0], nil
}

func (a *adapter) GetPairLastPrice(pairSymbol string) (float64, error) {
	tickers, err := a.client.GetTickers(pairSymbol)
	if err != nil {
		return 0, fmt.Errorf("get tickers: %w", err)
	}

	lastPrice, isExists := tickers[pairSymbol]
	if !isExists {
		return 0, fmt.Errorf("%q last price not found", pairSymbol)
	}
	return lastPrice, nil
}

func (a *adapter) CancelPairOrder(
	pairSymbol string,
	orderID int64,
	ctx context.Context,
) error {
	return a.client.CancelOrder(
		pairSymbol,
		strconv.FormatInt(orderID, 10),
	)
}

func (a *adapter) CancelPairOrderByClientOrderID(
	pairSymbol string,
	clientOrderID string,
	ctx context.Context,
) error {
	return a.client.CancelOrderByClientOrderID(pairSymbol, clientOrderID)
}

func (a *adapter) GetPairs() ([]structs.ExchangePairData, error) {
	symbols, err := a.client.GetSymbols()
	if err != nil {
		return nil, fmt.Errorf("get symbols: %w", err)
	}

	tickers, err := a.client.GetTickers()
	if err != nil {
		return nil, fmt.Errorf("get tickers: %w", err)
	}

	pairs, err := mappers.ConvertPairs(symbols, tickers)
	if err != nil {
		return nil, fmt.Errorf("convert: %w", err)
	}
	return pairs, nil
}

func (a *adapter) GetPairBalance(pair structs.PairSymbolData) (
	structs.PairBalance,
	error,
) {
	balances, err := a.GetAccountBalance()
	if err != nil {
		return structs.PairBalance{}, fmt.Errorf("get: %w", err)
	}

	return utils.FindPairBalance(balances, pair), nil
}
