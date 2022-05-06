package matrixgates

import (
	"testing"

	"github.com/adshao/go-binance/v2"
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
