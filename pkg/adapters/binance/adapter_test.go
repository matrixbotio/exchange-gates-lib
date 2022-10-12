package binance

import (
	"testing"

	"github.com/adshao/go-binance/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
)

func TestBinanceOrderConvert(t *testing.T) {
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
	orders, err := convertOrders(ordersRaw)
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
	require.Equal(t, exchangeID, consts.ExchangeIDbinanceSpot)
}

func TestBinanceConvertOrderSide(t *testing.T) {
	orderSide, err := convertOrderSide(binance.SideTypeBuy)
	assert.Nil(t, err)
	assert.Equal(t, orderSide, consts.OrderTypeBuy)

	orderSide, err = convertOrderSide(binance.SideTypeSell)
	assert.Nil(t, err)
	assert.Equal(t, orderSide, consts.OrderTypeSell)

	_, err = convertOrderSide("wtf")
	assert.NotNil(t, err)
}
