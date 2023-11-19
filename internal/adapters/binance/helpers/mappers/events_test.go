package mappers

import (
	"testing"

	"github.com/adshao/go-binance/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testExchangeTag = "binance-spot"

func getTestTradeEvent() binance.WsTradeEvent {
	return binance.WsTradeEvent{
		Time:          1700945872999,
		Symbol:        "LTCUSDT",
		Price:         "65.614",
		Quantity:      "1.1742",
		BuyerOrderID:  100,
		SellerOrderID: 101,
	}
}

func TestConvertPriceEventSuccess(t *testing.T) {
	// given
	event := binance.WsBookTickerEvent{
		Symbol:       "BTCUSDT",
		BestBidPrice: "20000",
		BestAskPrice: "20100",
		BestBidQty:   "0.1",
		BestAskQty:   "0,2",
	}

	// when
	ask, bid, err := ConvertPriceEvent(event)

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(20100), ask)
	assert.Equal(t, float64(20000), bid)
}

func TestConvertPriceEventError(t *testing.T) {
	// given
	event := binance.WsBookTickerEvent{
		Symbol:       "BTCUSDT",
		BestBidPrice: "wtf",
		BestAskPrice: "omg",
	}

	// when
	_, _, err := ConvertPriceEvent(event)

	// then
	require.ErrorContains(t, err, "invalid syntax")
}

func TestConvertTradeEventSuccess(t *testing.T) {
	// given
	event := getTestTradeEvent()

	// when
	result, err := ConvertTradeEvent(event, testExchangeTag)

	// then
	require.NoError(t, err)
	assert.Equal(t, int64(1700945873000), result.Time)
	assert.Equal(t, event.Symbol, result.Symbol)
	assert.Equal(t, float64(65.614), result.Price)
	assert.Equal(t, float64(1.1742), result.Quantity)
	assert.Equal(t, event.BuyerOrderID, result.BuyerOrderID)
	assert.Equal(t, event.SellerOrderID, result.SellerOrderID)
}

func TestConvertTradeEventParsePriceError(t *testing.T) {
	// given
	event := getTestTradeEvent()
	event.Price = "broken data"

	// when
	_, err := ConvertTradeEvent(event, testExchangeTag)

	// then
	require.ErrorContains(t, err, "parse price: strconv.ParseFloat")
}

func TestConvertTradeEventParseQtyError(t *testing.T) {
	// given
	event := getTestTradeEvent()
	event.Quantity = "broken data"

	// when
	_, err := ConvertTradeEvent(event, testExchangeTag)

	// then
	require.ErrorContains(t, err, "parse qty: strconv.ParseFloat")
}
