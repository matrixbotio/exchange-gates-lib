package mappers

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/hirokisan/bybit/v2"
	"github.com/stretchr/testify/require"
)

func TestParsePriceEvent(t *testing.T) {
	// given
	pairSymbol := "LTCUSDT"
	rawEvent := bybit.V5WebsocketPublicTickerResponse{
		Data: bybit.V5WebsocketPublicTickerData{
			Spot: &bybit.V5WebsocketPublicTickerSpotResult{
				Symbol:    bybit.SymbolV5(pairSymbol),
				LastPrice: "84.5",
			},
		},
	}
	exchangeTag := "binance-spot"

	// when
	event, err := ParsePriceEvent(rawEvent, exchangeTag)

	// then
	require.NoError(t, err)
	assert.Equal(t, exchangeTag, event.ExchangeTag)
	assert.Equal(t, float64(84.5), event.Ask)
}

func TestParseTradeEventSuccessful(t *testing.T) {
	// given
	mockData := bybit.V5WebsocketPrivateOrderData{
		BlockTradeID: "12345",
		Symbol:       "BTCUSD",
		OrderID:      "67890",
		OrderLinkID:  "abcdef",
		Price:        "45000.50",
		Qty:          "1.5",
		CumExecQty:   "1.0",
	}

	eventTime := int64(1643616000)
	exchangeTag := "bybit-spot"

	// Call the function to be tested
	result, err := ParseTradeEvent(mockData, eventTime, exchangeTag)

	// Assertions
	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, "12345", result.ID)
	assert.Equal(t, eventTime, result.Time)
	assert.Equal(t, exchangeTag, result.ExchangeTag)
	assert.Equal(t, "BTCUSD", result.Symbol)
	assert.Equal(t, "67890", result.OrderID)
	assert.Equal(t, "abcdef", result.ClientOrderID)
	assert.Equal(t, 45000.50, result.Price)
	assert.Equal(t, 1.5, result.Quantity)
	assert.Equal(t, 1.0, result.FilledQuantity)
}
