package mappers

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/hirokisan/bybit/v2"
	"github.com/stretchr/testify/require"
)

func TestConvertCoinData(t *testing.T) {
	// given
	tickerTag := "LTC"
	coinData := bybit.V5WalletBalanceCoin{
		Free:   "10",
		Locked: "0",
		Coin:   bybit.Coin(tickerTag),
	}

	// when
	balanceData, err := convertCoinData(coinData, tickerTag)

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(10), balanceData.Free)
	assert.Equal(t, float64(0), balanceData.Locked)
}

func TestFindAndConvertAssetBalance(t *testing.T) {
	// given
	tickerTag := "LTC"
	coins := []bybit.V5WalletBalanceCoin{
		{
			Free:   "10",
			Locked: "0",
			Coin:   bybit.Coin(tickerTag),
		},
	}

	// when
	balanceData, err := ConvertBalances(coins, tickerTag)

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(10), balanceData.Free)
	assert.Equal(t, float64(0), balanceData.Locked)
}

func TestFindAndConvertAssetBalanceNotFound(t *testing.T) {
	// given
	tickerTag := "wtf"
	coins := []bybit.V5WalletBalanceCoin{
		{
			Free:   "10",
			Locked: "0",
			Coin:   bybit.Coin("BTCUSDT"),
		},
	}

	// when
	_, err := ConvertBalances(coins, tickerTag)

	// then
	require.Error(t, err)
}

func TestConvertCoinDataParseError(t *testing.T) {
	// given
	tickerTag := "LTC"
	coinData := bybit.V5WalletBalanceCoin{
		Free:   "",
		Locked: "",
	}

	// when
	_, err := convertCoinData(coinData, tickerTag)

	// then
	require.Error(t, err)
}

func TestConvertCoinDataParseError2(t *testing.T) {
	// given
	tickerTag := "LTC"
	coinData := bybit.V5WalletBalanceCoin{
		Free:   "10",
		Locked: "",
	}

	// when
	_, err := convertCoinData(coinData, tickerTag)

	// then
	require.Error(t, err)
}
