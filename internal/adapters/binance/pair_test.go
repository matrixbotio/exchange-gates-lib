package binance

import (
	"context"
	"testing"

	"github.com/adshao/go-binance/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/errs"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/wrapper"
)

func TestGetPairLastPriceSuccess(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	w := wrapper.NewMockBinanceAPIWrapper(ctrl)
	a := New(w)

	w.EXPECT().GetPrices(gomock.Any(), testPairSymbol).
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
	ctrl := gomock.NewController(t)
	w := wrapper.NewMockBinanceAPIWrapper(ctrl)
	a := New(w)

	w.EXPECT().GetPrices(gomock.Any(), testPairSymbol).
		Return(nil, errTestException)

	// when
	_, err := a.GetPairLastPrice(testPairSymbol)

	// then
	require.ErrorIs(t, err, errTestException)
}

func TestGetPairLastPriceConvertError(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	w := wrapper.NewMockBinanceAPIWrapper(ctrl)
	a := New(w)

	w.EXPECT().GetPrices(gomock.Any(), testPairSymbol).
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

func TestCancelPairOrderSuccess(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	w := wrapper.NewMockBinanceAPIWrapper(ctrl)
	a := New(w)

	w.EXPECT().CancelOrderByID(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	// when
	err := a.CancelPairOrder(testPairSymbol, testOrderID, context.Background())

	// then
	require.NoError(t, err)
}

func TestCancelPairOrderByClientOrderIDSuccess(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	w := wrapper.NewMockBinanceAPIWrapper(ctrl)
	a := New(w)

	w.EXPECT().CancelOrderByClientOrderID(
		gomock.Any(), gomock.Any(), gomock.Any(),
	).Return(nil)

	// when
	err := a.CancelPairOrderByClientOrderID(
		testPairSymbol,
		testClientOrderID,
		context.Background(),
	)

	// then
	require.NoError(t, err)
}

func TestGetPairDataSuccess(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	w := wrapper.NewMockBinanceAPIWrapper(ctrl)
	a := New(w)

	baseAsset := "MTXB"
	quoteAsset := "USDC"

	w.EXPECT().GetExchangeInfo(gomock.Any(), gomock.Any()).
		Return(&binance.ExchangeInfo{
			Symbols: []binance.Symbol{
				{
					Symbol:             testPairSymbol,
					Status:             "TRADING",
					BaseAsset:          baseAsset,
					BaseAssetPrecision: 4,
					QuoteAsset:         quoteAsset,
					QuotePrecision:     4,
					Filters:            mappers.GetTestPairDataFilters(),
				},
			},
		}, nil)

	// when
	pairData, err := a.GetPairData(testPairSymbol)

	// then
	require.NoError(t, err)
	assert.Equal(t, testPairSymbol, pairData.Symbol)
	assert.Equal(t, baseAsset, pairData.BaseAsset)
	assert.Equal(t, quoteAsset, pairData.QuoteAsset)
}

func TestGetPairDataError(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	w := wrapper.NewMockBinanceAPIWrapper(ctrl)
	a := New(w)

	w.EXPECT().GetExchangeInfo(gomock.Any(), gomock.Any()).
		Return(nil, errTestException)

	// when
	_, err := a.GetPairData(testPairSymbol)

	// then
	require.ErrorIs(t, err, errTestException)
}

func TestGetPairDataNotFound(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	w := wrapper.NewMockBinanceAPIWrapper(ctrl)
	a := New(w)

	w.EXPECT().GetExchangeInfo(gomock.Any(), gomock.Any()).
		Return(&binance.ExchangeInfo{
			Symbols: []binance.Symbol{
				{Symbol: "LTCBUSD", Filters: mappers.GetTestPairDataFilters()},
				{Symbol: "BTCUSDC", Filters: mappers.GetTestPairDataFilters()},
			},
		}, nil)

	// when
	_, err := a.GetPairData(testPairSymbol)

	// then
	require.ErrorContains(t, err, "pair not found")
}

func TestGetPairs(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	w := wrapper.NewMockBinanceAPIWrapper(ctrl)
	a := New(w)

	w.EXPECT().GetExchangeInfo(gomock.Any(), gomock.Any()).
		Return(&binance.ExchangeInfo{
			Symbols: []binance.Symbol{
				{
					Symbol:             testPairSymbol,
					Status:             "TRADING",
					BaseAsset:          "MTXB",
					BaseAssetPrecision: 4,
					QuoteAsset:         "USDC",
					QuotePrecision:     4,
					Filters:            mappers.GetTestPairDataFilters(),
				},
				{
					Symbol:             "LTCUSDT",
					Status:             "TRADING",
					BaseAsset:          "LTC",
					BaseAssetPrecision: 8,
					QuoteAsset:         "USDT",
					QuotePrecision:     4,
					Filters:            mappers.GetTestPairDataFilters(),
				},
			},
		}, nil)

	// when
	pairs, err := a.GetPairs()

	// then
	require.NoError(t, err)
	assert.Len(t, pairs, 2)
}

func TestGetPairsError(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	w := wrapper.NewMockBinanceAPIWrapper(ctrl)
	a := New(w)

	w.EXPECT().GetExchangeInfo(gomock.Any(), gomock.Any()).
		Return(nil, errTestException)

	// when
	_, err := a.GetPairs()

	// then
	require.ErrorIs(t, err, errTestException)
}

func TestGetPairsResponseEmpty(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	w := wrapper.NewMockBinanceAPIWrapper(ctrl)
	a := New(w)

	w.EXPECT().GetExchangeInfo(gomock.Any(), gomock.Any()).
		Return(nil, nil)

	// when
	_, err := a.GetPairs()

	// then
	require.ErrorIs(t, err, errs.ErrPairResponseEmpty)
}
