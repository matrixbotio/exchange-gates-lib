package binance

import (
	"context"
	"errors"
	"fmt"
)

func (a *adapter) VerifyAPIKeys(keyPublic, keySecret string) error {
	accountData, err := a.binanceAPI.GetAccountData(context.Background())
	if err != nil {
		return fmt.Errorf("invalid api key: %w", err)
	}

	if !accountData.CanTrade {
		return errors.New("your API key does not have permission to trade," +
			" change its restrictions")
	}
	return nil
}
