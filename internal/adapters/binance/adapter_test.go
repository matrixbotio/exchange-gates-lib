package binance

import (
	"testing"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/matrixbotio/exchange-gates-lib/pkg/consts"
	"github.com/matrixbotio/exchange-gates-lib/pkg/utils"
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
	a := New()
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

func TestParseOrderOriginalQty(t *testing.T) {
	order := binance.Order{
		Symbol:           "LTCBUSD",
		OrderID:          102140140,
		ClientOrderID:    "123-456-789-ABC",
		Price:            "19236.86",
		OrigQuantity:     "0.00081",
		ExecutedQuantity: "0.00095",
		Status:           binance.OrderStatusTypeNew,
		Type:             binance.OrderTypeLimit,
		Side:             binance.SideTypeSell,
		Time:             time.Now().UnixMilli(),
		UpdateTime:       time.Now().UnixMilli(),
	}

	qty, err := utils.ParseStringToFloat64(order.OrigQuantity, "await qty")
	require.NoError(t, err)
	assert.Equal(t, float64(0.00081), qty)

	qty, err = utils.ParseStringToFloat64(order.ExecutedQuantity, "executed qty")
	require.NoError(t, err)
	assert.Equal(t, float64(0.00095), qty)

	price, err := utils.ParseStringToFloat64(order.Price, "price")
	require.NoError(t, err)
	assert.Equal(t, float64(19236.86), price)
}
