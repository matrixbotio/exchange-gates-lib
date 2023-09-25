package order_mappers

import (
	"strconv"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/hirokisan/bybit/v2"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
	"github.com/stretchr/testify/require"
)

func TestParseHistoryOrder(t *testing.T) {
	// given
	orderID := int64(12345)
	orderIDStr := strconv.FormatInt(orderID, 10)
	pairSymbol := "LTCUSDT"
	ordersResponse := &bybit.V5GetOrdersResponse{
		Result: bybit.V5GetOrdersResult{
			List: []bybit.V5GetOrder{
				{
					Symbol:      bybit.SymbolV5(pairSymbol),
					OrderID:     orderIDStr,
					Qty:         "0.5",
					CumExecQty:  "0.1",
					Price:       "82",
					UpdatedTime: "1692119310600",
					Side:        bybit.SideSell,
					OrderStatus: bybit.OrderStatusActive,
				},
			},
		},
	}

	// when
	orderData, err := ParseHistoryOrder(ordersResponse, orderIDStr, pairSymbol)

	// then
	require.NoError(t, err)
	assert.Equal(t, orderID, orderData.OrderID)
	assert.Equal(t, pairSymbol, orderData.Symbol)
	assert.Equal(t, pkgStructs.OrderTypeSell, orderData.Type)
	assert.Equal(t, pkgStructs.OrderStatusNew, orderData.Status)
	assert.Equal(t, float64(0.5), orderData.AwaitQty)
	assert.Equal(t, float64(0.1), orderData.FilledQty)
	assert.Equal(t, float64(82), orderData.Price)
	assert.Equal(t, int64(1692119310600), orderData.UpdatedTime)
}

func TestParseHistoryOrderNotFound(t *testing.T) {
	// given
	ordersResponse := &bybit.V5GetOrdersResponse{
		Result: bybit.V5GetOrdersResult{
			List: []bybit.V5GetOrder{},
		},
	}

	// when
	_, err := ParseHistoryOrder(ordersResponse, "", "")

	// then
	require.Error(t, err)
}
