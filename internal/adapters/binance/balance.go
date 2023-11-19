package binance

import (
	"context"
	"fmt"

	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/errs"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
)

func (a *adapter) GetAccountBalance() ([]structs.Balance, error) {
	accountBalances, err := a.getAccountBalances()
	if err != nil {
		return nil, fmt.Errorf("get account balances: %w", err)
	}

	return accountBalances.Balances, nil
}

func (a *adapter) getAccountBalances() (structs.AccountData, error) {
	data, err := a.binanceAPI.GetAccountData(context.Background())
	if err != nil {
		return structs.AccountData{}, fmt.Errorf("get account data: %w", err)
	}

	if data == nil {
		return structs.AccountData{}, errs.ErrAccountDataEmpty
	}

	accountData, err := mappers.ConvertAccountBalances(*data)
	if err != nil {
		return structs.AccountData{},
			fmt.Errorf("convert account balance: %w", err)
	}

	return accountData, nil
}

func (a *adapter) GetPairBalance(pair structs.PairSymbolData) (structs.PairBalance, error) {
	accountData, err := a.getAccountBalances()
	if err != nil {
		return structs.PairBalance{}, fmt.Errorf("get pair balance: %w", err)
	}

	return mappers.FindAssetBalances(accountData, pair), nil
}
