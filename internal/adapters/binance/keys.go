package binance

import (
	"context"
	"fmt"

	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/errs"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/pkg/structs"
)

func (a *adapter) VerifyAPIKeys(keyPublic, keySecret string) (
	structs.VerifyKeyStatus,
	error,
) {
	if err := a.Connect(structs.APICredentials{
		Type: structs.APICredentialsTypeKeypair,
		Keypair: structs.APIKeypair{
			Public: keyPublic,
			Secret: keySecret,
		},
	}); err != nil {
		return structs.VerifyKeyStatus{}, fmt.Errorf("binance connect: %w", err)
	}

	accountData, err := a.binanceAPI.GetAccountData(context.Background())
	if err != nil {
		return structs.VerifyKeyStatus{}, fmt.Errorf("invalid api key: %w", err)
	}

	if !accountData.CanTrade {
		return structs.VerifyKeyStatus{}, errs.ErrTradingNotAllowed
	}
	return structs.VerifyKeyStatus{
		Active:      true,
		AccountType: consts.AccountTypeStandart,
	}, nil
}
