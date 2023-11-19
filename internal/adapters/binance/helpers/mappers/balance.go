package mappers

import (
	"fmt"
	"strconv"

	"github.com/adshao/go-binance/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
)

func ConvertAssetBalance(data binance.Balance) (structs.Balance, error) {
	balanceFree, err := strconv.ParseFloat(data.Free, 64)
	if err != nil {
		return structs.Balance{},
			fmt.Errorf("parse %q free balance: %w", data.Asset, err)
	}

	balanceLocked, err := strconv.ParseFloat(data.Locked, 64)
	if err != nil {
		return structs.Balance{},
			fmt.Errorf("parse %q locked balance: %w", data.Asset, err)
	}

	return structs.Balance{
		Asset:  data.Asset,
		Free:   balanceFree,
		Locked: balanceLocked,
	}, nil
}

func FindAssetBalances(
	accountData structs.AccountData,
	pair structs.PairSymbolData,
) structs.PairBalance {
	pairBalanceData := structs.PairBalance{}
	for _, balanceData := range accountData.Balances {
		if balanceData.Asset == pair.BaseTicker {
			// base asset found
			pairBalanceData.BaseAsset = &structs.AssetBalance{
				Ticker: balanceData.Asset,
				Free:   balanceData.Free,
				Locked: balanceData.Locked,
			}
		}
		if balanceData.Asset == pair.QuoteTicker {
			// quote asset found
			pairBalanceData.QuoteAsset = &structs.AssetBalance{
				Ticker: balanceData.Asset,
				Free:   balanceData.Free,
				Locked: balanceData.Locked,
			}
		}
		if pairBalanceData.BaseAsset != nil && pairBalanceData.QuoteAsset != nil {
			// found
			break
		}
	}
	if pairBalanceData.BaseAsset == nil {
		pairBalanceData.BaseAsset = &structs.AssetBalance{
			Ticker: pair.BaseTicker,
			Free:   0,
			Locked: 0,
		}
	}
	if pairBalanceData.QuoteAsset == nil {
		pairBalanceData.QuoteAsset = &structs.AssetBalance{
			Ticker: pair.QuoteTicker,
			Free:   0,
			Locked: 0,
		}
	}
	return pairBalanceData
}

func ConvertAccountBalances(binanceAccountData binance.Account) (structs.AccountData, error) {
	accountDataResult := structs.AccountData{
		CanTrade: binanceAccountData.CanTrade,
	}

	balances := []structs.Balance{}
	for _, data := range binanceAccountData.Balances {
		assetBalance, err := ConvertAssetBalance(data)
		if err != nil {
			return structs.AccountData{}, fmt.Errorf("parse asset balance: %w", err)
		}

		if assetBalance.Free == 0 && assetBalance.Locked == 0 {
			continue
		}

		balances = append(balances, assetBalance)
	}

	accountDataResult.Balances = balances
	return accountDataResult, nil
}
