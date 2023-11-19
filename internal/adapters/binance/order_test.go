package binance

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/errs"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/wrapper"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	pkgErrs "github.com/matrixbotio/exchange-gates-lib/pkg/errs"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	testPairSymbol          = "MTXBUSDC"
	testOrderID       int64 = 100
	testClientOrderID       = "test"
)

func getTestOrderData() binance.Order {
	return binance.Order{
		Symbol:           testPairSymbol,
		OrderID:          testOrderID,
		ClientOrderID:    testClientOrderID,
		Price:            "1.001",
		OrigQuantity:     "1230.213",
		ExecutedQuantity: "102.1203220001",
		Status:           binance.OrderStatusTypeFilled,
		Type:             binance.OrderTypeLimit,
		Side:             binance.SideTypeBuy,
		Time:             time.Now().UnixMilli(),
	}
}

func getTestBotOrder() structs.BotOrderAdjusted {
	return structs.BotOrderAdjusted{
		PairSymbol:    testPairSymbol,
		Type:          pkgStructs.OrderTypeBuy,
		Qty:           "0.1005",
		Price:         "102.1924",
		Deposit:       "10.2703",
		ClientOrderID: testClientOrderID,
	}
}

func TestGetOrderDataErrorIDNotSet(t *testing.T) {
	// given
	a := New(wrapper.NewMockBinanceAPIWrapper(t))

	// when
	_, err := a.GetOrderData(testPairSymbol, 0)

	// then
	require.ErrorIs(t, err, errs.ErrOrderIDNotSet)
}

func TestGetOrderDataSuccess(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)
	order := getTestOrderData()

	w.EXPECT().GetOrderDataByOrderID(
		context.Background(),
		testPairSymbol,
		testOrderID,
	).Return(&order, nil)

	// when
	orderData, err := a.GetOrderData(testPairSymbol, testOrderID)

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(1.001), orderData.Price)
	assert.Equal(t, float64(1230.213), orderData.AwaitQty)
	assert.Equal(t, float64(102.1203220001), orderData.FilledQty)
	assert.Equal(t, testOrderID, orderData.OrderID)
	assert.Equal(t, testClientOrderID, orderData.ClientOrderID)
	assert.Equal(t, pkgStructs.OrderStatusFilled, orderData.Status)
	assert.Equal(t, pkgStructs.OrderTypeBuy, orderData.Type)
	assert.Equal(t, order.Time, orderData.CreatedTime)
}

func TestGetOrderDataErrorOrderUnknown(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)

	w.EXPECT().GetOrderDataByOrderID(mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New(errs.UnknownOrderMsg))

	// when
	_, err := a.GetOrderData(testPairSymbol, testOrderID)

	// then
	require.ErrorIs(t, err, pkgErrs.ErrOrderNotFound)
}

func TestGetOrderDataErrorUnknown(t *testing.T) {
	// given
	var w = wrapper.NewMockBinanceAPIWrapper(t)
	var a = New(w)
	var testErr = errors.New("some exception")

	w.EXPECT().GetOrderDataByOrderID(mock.Anything, mock.Anything, mock.Anything).
		Return(nil, testErr)

	// when
	_, err := a.GetOrderData(testPairSymbol, testOrderID)

	// then
	require.ErrorIs(t, err, testErr)
}

func TestGetOrderDataCovertError(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)
	order := getTestOrderData()
	order.Price = "strange data"

	w.EXPECT().GetOrderDataByOrderID(mock.Anything, mock.Anything, mock.Anything).
		Return(&order, nil)

	// when
	_, err := a.GetOrderData(testPairSymbol, testOrderID)

	// then
	require.ErrorContains(t, err, "invalid syntax")
}

func TestGetOrderByClientOrderIDSuccess(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)

	testOrderData := &binance.Order{
		Symbol:           testPairSymbol,
		ClientOrderID:    testClientOrderID,
		Price:            "1.012",
		OrigQuantity:     "125.1564",
		ExecutedQuantity: "0.12",
		Side:             binance.SideTypeBuy,
		Type:             binance.OrderTypeLimit,
		Status:           binance.OrderStatusTypePartiallyFilled,
	}

	w.EXPECT().GetOrderDataByClientOrderID(
		mock.Anything,
		testOrderData.Symbol,
		testClientOrderID,
	).Return(testOrderData, nil)

	// when
	orderData, err := a.GetOrderByClientOrderID(
		testOrderData.Symbol,
		testOrderData.ClientOrderID,
	)

	// then
	require.NoError(t, err)
	assert.Equal(t, testPairSymbol, orderData.Symbol)
	assert.Equal(t, float64(1.012), orderData.Price)
	assert.Equal(t, float64(125.1564), orderData.AwaitQty)
	assert.Equal(t, float64(0.12), orderData.FilledQty)
}

func TestGetOrderByClientOrderIDNotSet(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)

	// when
	_, err := a.GetOrderByClientOrderID("", "")

	// then
	require.ErrorIs(t, err, errs.ErrClientOrderIDNotSet)
}

func TestGetOrderByClientOrderUnknown(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)

	testOrderData := getTestOrderData()

	w.EXPECT().GetOrderDataByClientOrderID(
		mock.Anything,
		testOrderData.Symbol,
		testClientOrderID,
	).Return(nil, errors.New(errs.UnknownOrderMsg))

	// when
	_, err := a.GetOrderByClientOrderID(
		testOrderData.Symbol,
		testOrderData.ClientOrderID,
	)

	// then
	require.ErrorIs(t, err, pkgErrs.ErrOrderNotFound)
}

func TestGetOrderByClientOrderError(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)

	testOrderData := getTestOrderData()

	w.EXPECT().GetOrderDataByClientOrderID(
		mock.Anything,
		testOrderData.Symbol,
		testClientOrderID,
	).Return(nil, errTestException)

	// when
	_, err := a.GetOrderByClientOrderID(
		testOrderData.Symbol,
		testOrderData.ClientOrderID,
	)

	// then
	require.ErrorIs(t, err, errTestException)
}

func TestGetOrderByClientOrderConvertError(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)

	testOrderData := getTestOrderData()
	testOrderData.Price = "broken data"

	w.EXPECT().GetOrderDataByClientOrderID(
		mock.Anything,
		testOrderData.Symbol,
		testClientOrderID,
	).Return(&testOrderData, nil)

	// when
	_, err := a.GetOrderByClientOrderID(
		testOrderData.Symbol,
		testOrderData.ClientOrderID,
	)

	// then
	require.ErrorContains(t, err, "convert order: parse price")
}

func TestPlaceOrderSucess(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)
	order := getTestBotOrder()

	w.EXPECT().PlaceLimitOrder(
		mock.Anything, order.PairSymbol, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything,
	).Return(&binance.CreateOrderResponse{
		Symbol:        order.PairSymbol,
		OrderID:       testOrderID,
		ClientOrderID: testClientOrderID,
		Price:         order.Price,
		OrigQuantity:  order.Qty,
		Status:        binance.OrderStatusTypePartiallyFilled,
		Type:          binance.OrderTypeLimit,
		Side:          binance.SideTypeBuy,
	}, nil)

	// when
	response, err := a.PlaceOrder(context.Background(), order)

	// then
	require.NoError(t, err)
	assert.Equal(t, order.PairSymbol, response.Symbol)
	assert.Equal(t, order.Type, response.Type)
	assert.Equal(t, float64(0.1005), response.OrigQuantity)
	assert.Equal(t, float64(102.1924), response.Price)
}

func TestPlaceOrderInvalidOrderSide(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)
	order := getTestBotOrder()
	order.Type = "strange data"

	// when
	_, err := a.PlaceOrder(context.Background(), order)

	// then
	require.ErrorContains(t, err, "unknown order side")
}

func TestPlaceCreateOrderError(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)
	order := getTestBotOrder()

	w.EXPECT().PlaceLimitOrder(
		mock.Anything, order.PairSymbol, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything,
	).Return(nil, errTestException)

	// when
	_, err := a.PlaceOrder(context.Background(), order)

	// then
	require.ErrorIs(t, err, errTestException)
}

func TestPlaceOrderResponseEmpty(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)
	order := getTestBotOrder()

	w.EXPECT().PlaceLimitOrder(
		mock.Anything, order.PairSymbol, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything,
	).Return(nil, nil)

	// when
	_, err := a.PlaceOrder(context.Background(), order)

	// then
	require.ErrorIs(t, err, errs.ErrOrderResponseEmpty)
}

func TestPlaceOrderConvertError(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)
	order := getTestBotOrder()

	w.EXPECT().PlaceLimitOrder(
		mock.Anything, order.PairSymbol, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything,
	).Return(&binance.CreateOrderResponse{
		Price: "broken data",
	}, nil)

	// when
	_, err := a.PlaceOrder(context.Background(), order)

	// then
	require.ErrorContains(t, err, "invalid syntax")
}

func TestGetOrderExecFee(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)

	// when
	fees, err := a.GetOrderExecFee(
		testPairSymbol,
		pkgStructs.OrderTypeBuy,
		testOrderID,
	)

	// then
	require.NoError(t, err)
	assert.True(t, fees.BaseAsset.Equal(decimal.Zero))
	assert.True(t, fees.QuoteAsset.Equal(decimal.Zero))
}
