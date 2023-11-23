package binance

import (
	"testing"

	"github.com/adshao/go-binance/v2"
	"github.com/bmizerany/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetPairLastPriceSuccess(t *testing.T) {
	// given
	w := NewMockBinanceAPIWrapper(t)
	a := createAdapter(w)

	w.EXPECT().GetPrices(mock.Anything, testPairSymbol).
		Return([]*binance.SymbolPrice{
			{
				Symbol: testPairSymbol,
				Price:  "65.01294",
			},
			{
				Symbol: "BTCUSDT",
				Price:  "35000",
			},
		}, nil)

	// when
	lastPrice, err := a.GetPairLastPrice(testPairSymbol)

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(65.01294), lastPrice)
}

func TestGetPairLastPriceError(t *testing.T) {
	// given
	w := NewMockBinanceAPIWrapper(t)
	a := createAdapter(w)

	w.EXPECT().GetPrices(mock.Anything, testPairSymbol).
		Return(nil, errTestException)

	// when
	_, err := a.GetPairLastPrice(testPairSymbol)

	// then
	require.ErrorIs(t, err, errTestException)
}

func TestGetPairLastPriceConvertError(t *testing.T) {
	// given
	w := NewMockBinanceAPIWrapper(t)
	a := createAdapter(w)

	w.EXPECT().GetPrices(mock.Anything, testPairSymbol).
		Return([]*binance.SymbolPrice{
			{
				Symbol: testPairSymbol,
				Price:  "broken data",
			},
		}, nil)

	// when
	_, err := a.GetPairLastPrice(testPairSymbol)

	// then
	require.ErrorContains(t, err, "invalid syntax")
}
