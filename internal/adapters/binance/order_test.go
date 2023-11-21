package binance

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/bmizerany/assert"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/errs"
	pkgErrs "github.com/matrixbotio/exchange-gates-lib/pkg/errs"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
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

func TestGetOrderDataErrorIDNotSet(t *testing.T) {
	// given
	a := createAdapter(NewMockBinanceAPIWrapper(t))

	// when
	_, err := a.GetOrderData(testPairSymbol, 0)

	// then
	require.ErrorIs(t, err, errs.ErrOrderIDNotSet)
}

func TestGetOrderDataSuccess(t *testing.T) {
	// given
	w := NewMockBinanceAPIWrapper(t)
	a := createAdapter(w)
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
	w := NewMockBinanceAPIWrapper(t)
	a := createAdapter(w)

	w.EXPECT().GetOrderDataByOrderID(mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New(errs.UnknownOrderMsg))

	// when
	_, err := a.GetOrderData(testPairSymbol, testOrderID)

	// then
	require.ErrorIs(t, err, pkgErrs.OrderNotFound)
}

func TestGetOrderDataErrorUnknown(t *testing.T) {
	// given
	var w = NewMockBinanceAPIWrapper(t)
	var a = createAdapter(w)
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
	w := NewMockBinanceAPIWrapper(t)
	a := createAdapter(w)
	order := getTestOrderData()
	order.Price = "strange data"

	w.EXPECT().GetOrderDataByOrderID(mock.Anything, mock.Anything, mock.Anything).
		Return(&order, nil)

	// when
	_, err := a.GetOrderData(testPairSymbol, testOrderID)

	// then
	require.ErrorContains(t, err, "invalid syntax")
}
