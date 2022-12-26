package mappers

import (
	"testing"

	"github.com/bmizerany/assert"
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
