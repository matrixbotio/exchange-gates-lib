package utils

import (
	"testing"

	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalcTPOrderShort(t *testing.T) {
	// given
	coinsQty := float64(0.00126)
	profitPercent := float64(0.3)
	depositSpent := float64(32.64787)
	pairData := structs.ExchangePairData{
		Symbol:    "BTCBUSD",
		QtyStep:   0.00001,
		PriceStep: 0.01,
	}
	zeroProfitPrice := depositSpent / coinsQty

	srv := NewCalcTPOrderService().CoinsQty(coinsQty).Profit(profitPercent).
		DepositSpent(depositSpent).Strategy(pkgStructs.BotStrategyShort).
		PairData(pairData)

	// when
	order, err := srv.Do()

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(25833.5), order.Price)
	assert.Equal(t, float64(0.00126), order.Qty)
	assert.Less(t, order.Price, zeroProfitPrice)
	assert.Equal(t, pkgStructs.OrderTypeBuy, order.Type)
	assert.NotEmpty(t, order.ClientOrderID)
}

func TestCalcTPOrderLong(t *testing.T) {
	// given
	coinsQty := float64(0.00126)
	profitPercent := float64(0.3)
	depositSpent := float64(32.64787)
	pairData := structs.ExchangePairData{
		Symbol:    "BTCBUSD",
		QtyStep:   0.00001,
		PriceStep: 0.01,
	}
	zeroProfitPrice := depositSpent / coinsQty

	srv := NewCalcTPOrderService().CoinsQty(coinsQty).Profit(profitPercent).
		DepositSpent(depositSpent).Strategy(pkgStructs.BotStrategyLong).
		PairData(pairData)

	// when
	order, err := srv.Do()

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(25988.74), order.Price)
	assert.Equal(t, float64(0.00126), order.Qty)
	assert.Greater(t, order.Price, zeroProfitPrice)
	assert.Equal(t, pkgStructs.OrderTypeSell, order.Type)
	assert.NotEmpty(t, order.ClientOrderID)
}
