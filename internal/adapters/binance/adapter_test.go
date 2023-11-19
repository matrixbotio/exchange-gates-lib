package binance

import (
	"testing"

	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/stretchr/testify/require"
)

func TestBinanceAdapter(t *testing.T) {
	a := New()
	exchangeID := a.GetID()
	require.Equal(t, exchangeID, consts.ExchangeIDbinanceSpot)
}
