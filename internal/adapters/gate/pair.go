package gate

import (
	"context"

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

func (a *adapter) GetPairOpenOrders(pairSymbol string) ([]structs.OrderData, error) {
	// TODO
	return nil, nil
}

func (a *adapter) GetPairs() ([]structs.ExchangePairData, error) {
	// TODO
	return nil, nil
}

func (a *adapter) GetPairBalance(pair structs.PairSymbolData) (structs.PairBalance, error) {
	// TODO
	return structs.PairBalance{}, nil
}
