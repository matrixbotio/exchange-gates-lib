package mappers

import (
	"testing"

	"github.com/adshao/go-binance/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPairPriceSuccess(t *testing.T) {
	// given
	pairSymbol := "LTCUSDT"
	prices := []*binance.SymbolPrice{
		{
			Symbol: "USDCUSDT",
			Price:  "1.001",
		},
		{
			Symbol: pairSymbol,
			Price:  "65",
		},
	}

	// when
	lastPrice, err := GetPairPrice(prices, pairSymbol)

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(65), lastPrice)
}

func TestGetPairPriceParseError(t *testing.T) {
	// given
	pairSymbol := "USDCUSDT"
	prices := []*binance.SymbolPrice{
		{
			Symbol: pairSymbol,
			Price:  "1-001",
		},
	}

	// when
	_, err := GetPairPrice(prices, pairSymbol)

	// then
	require.Error(t, err)
}

func TestGetPairPriceNotFound(t *testing.T) {
	// given
	pairSymbol := "MTXBUSDC"
	prices := []*binance.SymbolPrice{
		{
			Symbol: "USDCUSDT",
			Price:  "1.001",
		},
	}

	// when
	_, err := GetPairPrice(prices, pairSymbol)

	// then
	require.Error(t, err)
}
