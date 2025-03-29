package binance

import (
	"testing"

	"github.com/adshao/go-binance/v2"
	"github.com/bmizerany/assert"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/errs"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/wrapper"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func getTestBalances() []binance.Balance {
	return []binance.Balance{
		{
			Asset:  "LTC",
			Free:   "10.1114",
			Locked: "0.0000",
		},
		{
			Asset:  "MTXB",
			Free:   "100500",
			Locked: "24.0201",
		},
	}
}

func TestGetAccountBalanceSuccess(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	w := wrapper.NewMockBinanceAPIWrapper(ctrl)
	a := New(w)

	testBalances := &binance.Account{
		Balances: getTestBalances(),
	}

	w.EXPECT().GetAccountData(gomock.Any()).Return(testBalances, nil)

	// when
	balances, err := a.GetAccountBalance()

	// then
	require.NoError(t, err)
	require.Len(t, balances, 2)
	assert.Equal(t, "LTC", balances[0].Asset)
	assert.Equal(t, "MTXB", balances[1].Asset)
	assert.Equal(t, float64(10.1114), balances[0].Free)
	assert.Equal(t, float64(0), balances[0].Locked)
	assert.Equal(t, float64(100500), balances[1].Free)
	assert.Equal(t, float64(24.0201), balances[1].Locked)
}

func TestGetAccountBalanceError(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	w := wrapper.NewMockBinanceAPIWrapper(ctrl)
	a := New(w)

	w.EXPECT().GetAccountData(gomock.Any()).Return(nil, errTestException)

	// when
	_, err := a.GetAccountBalance()

	// then
	require.ErrorIs(t, err, errTestException)
}

func TestGetAccountBalanceConvertError(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	w := wrapper.NewMockBinanceAPIWrapper(ctrl)
	a := New(w)

	testBalances := &binance.Account{
		Balances: getTestBalances(),
	}
	testBalances.Balances[0].Free = "broken data"

	w.EXPECT().GetAccountData(gomock.Any()).Return(testBalances, nil)

	// when
	_, err := a.GetAccountBalance()

	// then
	require.ErrorContains(t, err, "invalid syntax")
}

func TestGetAccountBalanceErrorEmptyResponse(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	w := wrapper.NewMockBinanceAPIWrapper(ctrl)
	a := New(w)

	w.EXPECT().GetAccountData(gomock.Any()).Return(nil, nil)

	// when
	_, err := a.GetAccountBalance()

	// then
	require.ErrorIs(t, err, errs.ErrAccountDataEmpty)
}

func TestGetPairBalanceSuccess(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	w := wrapper.NewMockBinanceAPIWrapper(ctrl)
	a := New(w)
	pairSymbolData := structs.PairSymbolData{
		BaseTicker:  "LTC",
		QuoteTicker: "MTXB",
		Symbol:      "LTCMTXB",
	}

	w.EXPECT().GetAccountData(gomock.Any()).Return(&binance.Account{
		Balances: getTestBalances(),
	}, nil)

	// when
	pairBalance, err := a.GetPairBalance(pairSymbolData)

	// then
	require.NoError(t, err)
	assert.Equal(t, pairSymbolData.BaseTicker, pairBalance.BaseAsset.Ticker)
	assert.Equal(t, pairSymbolData.QuoteTicker, pairBalance.QuoteAsset.Ticker)
	assert.Equal(t, float64(10.1114), pairBalance.BaseAsset.Free)
	assert.Equal(t, float64(0), pairBalance.BaseAsset.Locked)
	assert.Equal(t, float64(100500), pairBalance.QuoteAsset.Free)
	assert.Equal(t, float64(24.0201), pairBalance.QuoteAsset.Locked)
}

func TestGetPairBalanceError(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	w := wrapper.NewMockBinanceAPIWrapper(ctrl)
	a := New(w)
	pairSymbolData := structs.PairSymbolData{
		BaseTicker:  "LTC",
		QuoteTicker: "MTXB",
		Symbol:      "LTCMTXB",
	}

	w.EXPECT().GetAccountData(gomock.Any()).Return(nil, errTestException)

	// when
	_, err := a.GetPairBalance(pairSymbolData)

	// then
	require.ErrorIs(t, err, errTestException)
}
