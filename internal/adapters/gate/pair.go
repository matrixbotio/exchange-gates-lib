package gate

import (
	"context"
	"fmt"

	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/gate/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
)

func (a *adapter) GetPairData(pairSymbol string) (structs.ExchangePairData, error) {
	// TODO
	return structs.ExchangePairData{}, nil
}

func (a *adapter) GetPairLastPrice(pairSymbol string) (float64, error) {
	// TODO
	return 0, nil
}

func (a *adapter) CancelPairOrder(
	pairSymbol string,
	orderID int64,
	ctx context.Context,
) error {
	// TODO
	return nil
}

func (a *adapter) CancelPairOrderByClientOrderID(
	pairSymbol string,
	clientOrderID string,
	ctx context.Context,
) error {
	// TODO
	return nil
}

func (a *adapter) GetPairs() ([]structs.ExchangePairData, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), requestTimeout)
	defer ctxCancel()

	pairs, _, err := a.client.SpotApi.ListCurrencyPairs(ctx)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	result, err := mappers.ConvertPairData(pairs)
	if err != nil {
		return nil, fmt.Errorf("convert: %w", err)
	}
	return result, nil
}

func (a *adapter) GetPairBalance(pair structs.PairSymbolData) (structs.PairBalance, error) {
	// TODO
	return structs.PairBalance{}, nil
}
