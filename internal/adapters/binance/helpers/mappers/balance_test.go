package mappers

import (
	"testing"

	"github.com/adshao/go-binance/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseAssetBalance(t *testing.T) {
	// given
	rawData := binance.Balance{
		Asset:  "BTC",
		Free:   "0.01",
		Locked: "0",
	}

	// when
	assetBalance, err := ConvertAssetBalance(rawData)

	// then
	require.NoError(t, err)
	assert.Equal(t, rawData.Asset, assetBalance.Asset)
	assert.Equal(t, float64(0.01), assetBalance.Free)
	assert.Equal(t, float64(0), assetBalance.Locked)
}

func TestParseAssetBalanceFreeEmpty(t *testing.T) {
	// given
	rawData := binance.Balance{
		Asset:  "BTC",
		Free:   "",
		Locked: "0",
	}

	// when
	_, err := ConvertAssetBalance(rawData)

	// then
	require.ErrorContains(t, err, "invalid syntax")
}

func TestParseAssetBalanceLockedEmpty(t *testing.T) {
	// given
	rawData := binance.Balance{
		Asset:  "BTC",
		Free:   "0.01",
		Locked: "",
	}

	// when
	_, err := ConvertAssetBalance(rawData)

	// then
	require.ErrorContains(t, err, "invalid syntax")
}

func TestFindAssetBalancesSuccess(t *testing.T) {
	// given
	baseAsset := "BTC"
	quoteAsset := "BUSD"
	baseAssetFree := float64(0.1)
	quoteAssetFree := float64(95.12)

	accountData := structs.AccountData{
		Balances: []structs.Balance{
			{
				Asset: "LTC",
				Free:  10,
			},
			{
				Asset: baseAsset,
				Free:  baseAssetFree,
			},
			{
				Asset: quoteAsset,
				Free:  quoteAssetFree,
			},
		},
	}
	pairSymbolData := structs.PairSymbolData{
		BaseTicker:  baseAsset,
		QuoteTicker: quoteAsset,
		Symbol:      baseAsset + quoteAsset,
	}

	// when
	pairBalance := FindAssetBalances(accountData, pairSymbolData)

	// then
	assert.Equal(t, baseAsset, pairBalance.BaseAsset.Ticker)
	assert.Equal(t, quoteAsset, pairBalance.QuoteAsset.Ticker)
	assert.Equal(t, baseAssetFree, pairBalance.BaseAsset.Free)
	assert.Equal(t, float64(0), pairBalance.BaseAsset.Locked)
	assert.Equal(t, quoteAssetFree, pairBalance.QuoteAsset.Free)
	assert.Equal(t, float64(0), pairBalance.QuoteAsset.Locked)
}

func TestFindAssetBalancesNotFound(t *testing.T) {
	// given
	baseAsset := "MTXB"
	quoteAsset := "USDC"

	accountData := structs.AccountData{
		Balances: []structs.Balance{
			{
				Asset: "LTC",
				Free:  10,
			},
			{
				Asset: "USDT",
				Free:  0.5,
			},
		},
	}
	pairSymbolData := structs.PairSymbolData{
		BaseTicker:  baseAsset,
		QuoteTicker: quoteAsset,
		Symbol:      baseAsset + quoteAsset,
	}

	// when
	pairBalance := FindAssetBalances(accountData, pairSymbolData)

	// then
	assert.Equal(t, baseAsset, pairBalance.BaseAsset.Ticker)
	assert.Equal(t, quoteAsset, pairBalance.QuoteAsset.Ticker)
	assert.Equal(t, float64(0), pairBalance.BaseAsset.Free)
	assert.Equal(t, float64(0), pairBalance.BaseAsset.Locked)
	assert.Equal(t, float64(0), pairBalance.QuoteAsset.Free)
	assert.Equal(t, float64(0), pairBalance.QuoteAsset.Locked)
}

func TestConvertAccountBalances(t *testing.T) {
	// given
	binanceAccountData := binance.Account{
		CanTrade: true,
		Balances: []binance.Balance{
			{
				Asset:  "MTXB",
				Free:   "100500.000",
				Locked: "0.000000000",
			},
			{
				Asset:  "USDT",
				Free:   "10.000000001",
				Locked: "5.0200000000",
			},
		},
	}

	// when
	accountData, err := ConvertAccountBalances(binanceAccountData)

	// then
	require.NoError(t, err)
	assert.Equal(t, true, accountData.CanTrade)
	require.Len(t, accountData.Balances, 2)
	assert.Equal(t, "MTXB", accountData.Balances[0].Asset)
	assert.Equal(t, float64(100500), accountData.Balances[0].Free)
	assert.Equal(t, float64(0), accountData.Balances[0].Locked)
	assert.Equal(t, "USDT", accountData.Balances[1].Asset)
	assert.Equal(t, float64(10.000000001), accountData.Balances[1].Free)
	assert.Equal(t, float64(5.02), accountData.Balances[1].Locked)
}
