package binance

import (
	"testing"

	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/wrapper"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const testInterval = consts.Interval15min

func TestGetCandlesSuccess(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	w := wrapper.NewMockBinanceAPIWrapper(ctrl)
	a := New(w)
	limit := 5
	pairSymbol := "LTCUSDT"
	interval := testInterval

	w.EXPECT().GetKlines(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mappers.GetTestKlines(), nil)

	// when
	candles, err := a.GetCandles(limit, pairSymbol, interval)

	// then
	require.NoError(t, err)
	require.Len(t, candles, 2)
	assert.Equal(t, interval, candles[0].Interval)
	assert.Equal(t, 1000.0, candles[0].Open)
	assert.Equal(t, 2000.0, candles[0].Close)
	assert.Equal(t, 3000.0, candles[0].High)
	assert.Equal(t, 500.0, candles[0].Low)
	assert.Equal(t, 10000.0, candles[0].Volume)
}

func TestGetCandlesGetKlinesError(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	w := wrapper.NewMockBinanceAPIWrapper(ctrl)
	a := New(w)
	limit := 5
	pairSymbol := "LTCUSDT"
	interval := testInterval

	w.EXPECT().GetKlines(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, errTestException)

	// when
	_, err := a.GetCandles(limit, pairSymbol, interval)

	// then
	require.ErrorIs(t, err, errTestException)
}

func TestGetCandlesConvertError(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	w := wrapper.NewMockBinanceAPIWrapper(ctrl)
	a := New(w)
	limit := 5
	pairSymbol := "LTCUSDT"
	interval := testInterval

	klines := mappers.GetTestKlines()
	klines[0].Low = "strange data"

	w.EXPECT().GetKlines(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(klines, nil)

	// when
	_, err := a.GetCandles(limit, pairSymbol, interval)

	// then
	require.ErrorContains(t, err, "convert candles")
}
