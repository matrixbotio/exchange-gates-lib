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
	assert.Nil(t, err)
	assert.Equal(t, orderSide, structs.OrderTypeBuy)

}

func TestConvertOrderSideSell(t *testing.T) {
	// given
	var exchangeOrderSide = binance.SideTypeSell

	// when
	orderSide, err := ConvertOrderSide(exchangeOrderSide)

	// then
	assert.Nil(t, err)
	assert.Equal(t, orderSide, structs.OrderTypeSell)
}

func TestConvertOrderSideUnknown(t *testing.T) {
	// given
	var exchangeOrderSide = binance.SideType("wtf")

	// when
	_, err := ConvertOrderSide(exchangeOrderSide)

	// then
	assert.NotNil(t, err)
}
