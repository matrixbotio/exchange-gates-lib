package mappers

import (
	"testing"

	"github.com/adshao/go-binance/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
