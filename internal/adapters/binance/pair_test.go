package binance

import (
	"context"
	"testing"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/errs"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/wrapper"
	"github.com/matrixbotio/exchange-gates-lib/pkg/structs"
)

func TestGetPairLastPriceSuccess(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)

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
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)

	w.EXPECT().GetPrices(mock.Anything, testPairSymbol).
		Return(nil, errTestException)

	// when
	_, err := a.GetPairLastPrice(testPairSymbol)

	// then
	require.ErrorIs(t, err, errTestException)
}

func TestGetPairLastPriceConvertError(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)

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

func TestCancelPairOrderSuccess(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)

	w.EXPECT().CancelOrderByID(mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	// when
	err := a.CancelPairOrder(testPairSymbol, testOrderID, context.Background())

	// then
	require.NoError(t, err)
}

func TestCancelPairOrderByClientOrderIDSuccess(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)

	w.EXPECT().CancelOrderByClientOrderID(
		mock.Anything, mock.Anything, mock.Anything,
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
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)

	baseAsset := "MTXB"
	quoteAsset := "USDC"

	w.EXPECT().GetExchangeInfo(mock.Anything, mock.Anything).
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
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)

	w.EXPECT().GetExchangeInfo(mock.Anything, mock.Anything).
		Return(nil, errTestException)

	// when
	_, err := a.GetPairData(testPairSymbol)

	// then
	require.ErrorIs(t, err, errTestException)
}

func TestGetPairDataNotFound(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)

	w.EXPECT().GetExchangeInfo(mock.Anything, mock.Anything).
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

func TestGetPairOpenOrdersSucess(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)

	ordersResponse := []*binance.Order{
		{
			Symbol:           "LTCUSDT",
			OrderID:          100,
			Price:            "65.418",
			OrigQuantity:     "0.1205",
			ExecutedQuantity: "0.0000",
			Status:           binance.OrderStatusTypeNew,
			Side:             binance.SideTypeSell,
			Type:             binance.OrderTypeLimit,
			Time:             time.Now().UnixMilli(),
		},
		{
			Symbol:           "LTCUSDT",
			OrderID:          101,
			Price:            "66.918",
			OrigQuantity:     "0.1425",
			ExecutedQuantity: "0.0000",
			Status:           binance.OrderStatusTypeNew,
			Side:             binance.SideTypeSell,
			Type:             binance.OrderTypeLimit,
			Time:             time.Now().UnixMilli(),
		},
	}

	w.EXPECT().GetOpenOrders(mock.Anything, mock.Anything).
		Return(ordersResponse, nil)

	// when
	orders, err := a.GetPairOpenOrders(testPairSymbol)

	// then
	require.NoError(t, err)
	require.Len(t, orders, 2)
	assert.Equal(t, float64(65.418), orders[0].Price)
	assert.Equal(t, float64(0.1205), orders[0].AwaitQty)
	assert.Equal(t, structs.OrderTypeSell, orders[0].Type)
	assert.Equal(t, float64(66.918), orders[1].Price)
	assert.Equal(t, float64(0.1425), orders[1].AwaitQty)
	assert.Equal(t, structs.OrderTypeSell, orders[1].Type)
}

func TestGetPairOpenOrdersError(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)

	w.EXPECT().GetOpenOrders(mock.Anything, mock.Anything).
		Return(nil, errTestException)

	// when
	_, err := a.GetPairOpenOrders(testPairSymbol)

	// then
	require.ErrorIs(t, err, errTestException)
}

func TestGetPairOpenOrdersConvertError(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)

	ordersResponse := []*binance.Order{
		{
			Symbol:           "LTCUSDT",
			OrderID:          100,
			Price:            "65.418",
			OrigQuantity:     "broken data",
			ExecutedQuantity: "0.0000",
			Status:           binance.OrderStatusTypeNew,
			Side:             binance.SideTypeSell,
			Type:             binance.OrderTypeLimit,
			Time:             time.Now().UnixMilli(),
		},
	}

	w.EXPECT().GetOpenOrders(mock.Anything, mock.Anything).
		Return(ordersResponse, nil)

	// when
	_, err := a.GetPairOpenOrders(testPairSymbol)

	// then
	require.ErrorContains(t, err, "invalid syntax")
}

func TestGetPairs(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)

	w.EXPECT().GetExchangeInfo(mock.Anything, mock.Anything).
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
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)

	w.EXPECT().GetExchangeInfo(mock.Anything, mock.Anything).
		Return(nil, errTestException)

	// when
	_, err := a.GetPairs()

	// then
	require.ErrorIs(t, err, errTestException)
}

func TestGetPairsResponseEmpty(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)

	w.EXPECT().GetExchangeInfo(mock.Anything, mock.Anything).
		Return(nil, nil)

	// when
	_, err := a.GetPairs()

	// then
	require.ErrorIs(t, err, errs.ErrPairResponseEmpty)
}
