package binance

import (
	"testing"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestOrder() *binance.Order {
	return &binance.Order{
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
}

func TestParseOrderOriginalQty(t *testing.T) {
	originalOrder := getTestOrder()

	qty, err := parseOrderOriginalQty(originalOrder)
	require.NoError(t, err)
	assert.Equal(t, float64(0.00081), qty)
}

func TestParseOrderExecutedQty(t *testing.T) {
	originalOrder := getTestOrder()

	qty, err := parseOrderExecutedQty(originalOrder)
	require.NoError(t, err)
	assert.Equal(t, float64(0.00095), qty)
}

func TestParseOrderPrice(t *testing.T) {
	originalOrder := getTestOrder()

	qty, err := parseOrderPrice(originalOrder)
	require.NoError(t, err)
	assert.Equal(t, float64(19236.86), qty)
}
