package pkg

import (
	"testing"

	"github.com/adshao/go-binance/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBinanceOrderConvert(t *testing.T) {
	a := NewBinanceSpotAdapter()

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
	orders, err := a.convertOrders(ordersRaw)
	if err != nil {
		t.Fatal(err)
	}

	if len(orders) == 0 {
		t.Fatal("0 orders converted")
	}
}

func TestBinanceAdapter(t *testing.T) {
	a := NewBinanceSpotAdapter()
	exchangeID := a.GetID()
	require.Equal(t, exchangeID, exchangeIDbinanceSpot)
}

func TestBinanceConvertOrderSide(t *testing.T) {
	a := NewBinanceSpotAdapter()
	orderSide, err := a.convertOrderSide(binance.SideTypeBuy)
	assert.Nil(t, err)
	assert.Equal(t, orderSide, OrderTypeBuy)

	orderSide, err = a.convertOrderSide(binance.SideTypeSell)
	assert.Nil(t, err)
	assert.Equal(t, orderSide, OrderTypeSell)

	_, err = a.convertOrderSide("wtf")
	assert.NotNil(t, err)
}
