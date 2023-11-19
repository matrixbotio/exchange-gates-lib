package binance

import (
	"context"
	"fmt"
)

func (a *adapter) CanTrade() (bool, error) {
	binanceAccountData, err := a.binanceAPI.NewGetAccountService().
		Do(context.Background())
	if err != nil {
		return false, fmt.Errorf("get account data: %w", err)
	}
	return binanceAccountData.CanTrade, nil
}
