package bingx

import (
	"fmt"
	"strconv"

	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
)

func (a *adapter) CanTrade() (bool, error) {
	_, err := a.client.GetBalance()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a *adapter) GetAccountBalance() ([]structs.Balance, error) {
	balances, err := a.client.GetBalance()
	if err != nil {
		return nil, fmt.Errorf("get balance: %w", err)
	}

	var result []structs.Balance
	for _, balanceData := range balances {
		assetFree, err := strconv.ParseFloat(balanceData.Free, 64)
		if err != nil {
			return nil, fmt.Errorf("parse %q free: %w", err)
		}

		assetLocked, err := strconv.ParseFloat(balanceData.Locked, 64)
		if err != nil {
			return nil, fmt.Errorf("parse %q locked: %w", err)
		}

		result = append(result, structs.Balance{
			Asset:  balanceData.Asset,
			Free:   assetFree,
			Locked: assetLocked,
		})
	}
	return result, nil
}
