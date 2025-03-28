package utils

import "github.com/matrixbotio/exchange-gates-lib/internal/structs"

func FindPairBalance(
	balances []structs.Balance,
	pair structs.PairSymbolData,
) structs.PairBalance {
	var result structs.PairBalance

	for _, balance := range balances {
		if balance.Asset == pair.BaseTicker {
			result.BaseAsset = &structs.AssetBalance{
				Ticker: balance.Asset,
				Free:   balance.Free,
				Locked: balance.Locked,
			}
		}

		if balance.Asset == pair.QuoteTicker {
			result.QuoteAsset = &structs.AssetBalance{
				Ticker: balance.Asset,
				Free:   balance.Free,
				Locked: balance.Locked,
			}
		}

		if result.BaseAsset != nil && result.QuoteAsset != nil {
			return result
		}
	}

	if result.BaseAsset == nil {
		result.BaseAsset = &structs.AssetBalance{
			Ticker: pair.BaseTicker,
		}
	}
	if result.QuoteAsset == nil {
		result.QuoteAsset = &structs.AssetBalance{
			Ticker: pair.BaseTicker,
		}
	}

	return result
}
