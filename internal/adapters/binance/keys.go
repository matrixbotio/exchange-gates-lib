package binance

import (
	"context"
	"fmt"

	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/errs"
	"github.com/matrixbotio/exchange-gates-lib/pkg/structs"
)

func (a *adapter) VerifyAPIKeys(keyPublic, keySecret string) error {
	if err := a.Connect(structs.APICredentials{
		Type: structs.APICredentialsTypeKeypair,
		Keypair: structs.APIKeypair{
			Public: keyPublic,
			Secret: keySecret,
		},
	}); err != nil {
		return fmt.Errorf("binance connect: %w", err)
	}

	accountData, err := a.binanceAPI.GetAccountData(context.Background())
	if err != nil {
		return fmt.Errorf("invalid api key: %w", err)
	}

	if !accountData.CanTrade {
		return errs.ErrTradingNotAllowed
	}
	return nil
}
