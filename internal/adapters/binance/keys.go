package binance

import (
	"context"
	"fmt"

	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/errs"
)

func (a *adapter) VerifyAPIKeys(keyPublic, keySecret string) error {
	accountData, err := a.binanceAPI.GetAccountData(context.Background())
	if err != nil {
		return fmt.Errorf("invalid api key: %w", err)
	}

	if !accountData.CanTrade {
		return errs.ErrTradingNotAllowed
	}
	return nil
}
