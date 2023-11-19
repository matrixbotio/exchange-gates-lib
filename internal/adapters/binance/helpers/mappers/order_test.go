package mappers

import (
	"testing"

	"github.com/adshao/go-binance/v2"
	"github.com/matrixbotio/exchange-gates-lib/pkg/structs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBinanceOrderConvert(t *testing.T) {
	// given
	pairSymbol := "LTCBUSD"
	ordersRaw := []*binance.Order{
		{
			Symbol:           pairSymbol,
			OrderID:          1498236,
			Price:            "125.1",
			OrigQuantity:     "0.1",
			ExecutedQuantity: "0",
			Status:           binance.OrderStatusTypeNew,
			Type:             binance.OrderTypeLimit,
			Side:             binance.SideTypeBuy,
		},
	}

	// when
	orders, err := ConvertOrders(ordersRaw)

	// then
	require.NoError(t, err)
	assert.NotEqual(t, 0, len(orders))
}

func TestBinanceConvertOrderSide(t *testing.T) {
	// when
	var exchangeOrderSide = binance.SideTypeBuy

	// when
	orderSide, err := ConvertOrderSide(exchangeOrderSide)

	// then
	require.NoError(t, err)
	assert.Equal(t, orderSide, structs.OrderTypeBuy)
}

func TestConvertOrderSideSell(t *testing.T) {
	// given
	var exchangeOrderSide = binance.SideTypeSell

	// when
	orderSide, err := ConvertOrderSide(exchangeOrderSide)

	// then
	require.NoError(t, err)
	assert.Equal(t, orderSide, structs.OrderTypeSell)
}

func TestConvertOrderSideUnknown(t *testing.T) {
	// given
	var exchangeOrderSide = binance.SideType("wtf")

	// when
	_, err := ConvertOrderSide(exchangeOrderSide)

	// then
	assert.ErrorContains(t, err, "unknown order side")
}

func TestGetBinanceOrderSideBuy(t *testing.T) {
	// given
	var botOrderSide = structs.OrderTypeBuy

	// when
	orderType, err := GetBinanceOrderSide(botOrderSide)

	// then
	require.NoError(t, err)
	assert.Equal(t, binance.SideTypeBuy, orderType)
}

func TestGetBinanceOrderSideSell(t *testing.T) {
	// given
	var botOrderSide = structs.OrderTypeSell

	// when
	orderType, err := GetBinanceOrderSide(botOrderSide)

	// then
	require.NoError(t, err)
	assert.Equal(t, binance.SideTypeSell, orderType)
}

func TestGetBinanceOrderSideUnknown(t *testing.T) {
	// given
	var botOrderSide = "wtf"

	// when
	_, err := GetBinanceOrderSide(botOrderSide)

	// then
	require.ErrorContains(t, err, "unknown order side")
}

func TestConvertBinanceToBotOrderSuccess(t *testing.T) {
	// given
	orderResponse := binance.CreateOrderResponse{
		Symbol:           "LTCUSDC",
		OrderID:          100,
		ClientOrderID:    "test",
		Price:            "65.108",
		OrigQuantity:     "1.219",
		ExecutedQuantity: "0.025",
		Status:           binance.OrderStatusTypeFilled,
		Type:             binance.OrderTypeLimit,
		Side:             binance.SideTypeBuy,
	}

	// when
	order, err := ConvertPlacedOrder(orderResponse)

	// then
	require.NoError(t, err)
	assert.Equal(t, orderResponse.Symbol, order.Symbol)
	assert.Equal(t, orderResponse.OrderID, order.OrderID)
	assert.Equal(t, orderResponse.ClientOrderID, order.ClientOrderID)
	assert.Equal(t, float64(65.108), order.Price)
	assert.Equal(t, float64(1.219), order.OrigQuantity)
	assert.Equal(t, structs.OrderTypeBuy, order.Type)
}

func TestConvertBinanceToBotOrderInvalidQty(t *testing.T) {
	// given
	orderResponse := binance.CreateOrderResponse{
		Price:        "65.108",
		OrigQuantity: "wtf",
	}

	// when
	_, err := ConvertPlacedOrder(orderResponse)

	// then
	require.ErrorContains(t, err, "parse order origQty")
}

func TestConvertBinanceToBotOrderInvalidPrice(t *testing.T) {
	// given
	orderResponse := binance.CreateOrderResponse{
		Price:        "wtf",
		OrigQuantity: "1.108",
	}

	// when
	_, err := ConvertPlacedOrder(orderResponse)

	// then
	require.ErrorContains(t, err, "parse order price")
}

func TestConvertBinanceToBotOrderInvalidSide(t *testing.T) {
	// given
	orderResponse := binance.CreateOrderResponse{
		Price:        "65.118",
		OrigQuantity: "1.108",
		Side:         binance.SideType("wtf"),
	}

	// when
	_, err := ConvertPlacedOrder(orderResponse)

	// then
	require.ErrorContains(t, err, "unknown order side")
}
