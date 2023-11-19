package binance

import (
	"context"
	"fmt"
)

func (a *adapter) CanTrade() (bool, error) {
	data, err := a.binanceAPI.GetAccountData(context.Background())
	if err != nil {
		return false, fmt.Errorf("get account data: %w", err)
	}

	return data.CanTrade, nil
}
