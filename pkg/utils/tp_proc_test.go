package utils

import (
	"testing"

	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testTPCoinsQty      float64 = 0.00126
	testTPProfitPercent float64 = 0.3
	testTPDepositSpent  float64 = 32.64787
)

func TestCalcTPOrderShortNoFees(t *testing.T) {
	// given
	pairData := structs.ExchangePairData{
		Symbol:    "BTCBUSD",
		QtyStep:   0.00001,
		PriceStep: 0.01,
	}
	zeroProfitPrice := testTPDepositSpent / testTPCoinsQty

	srv := NewCalcTPOrderService().CoinsQty(testTPCoinsQty).
		Profit(testTPProfitPercent).
		DepositSpent(testTPDepositSpent).
		Strategy(pkgStructs.BotStrategyShort).
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

func TestCalcTPOrderShortWithFees(t *testing.T) {
	// given
	pairData := structs.ExchangePairData{
		Symbol:    "BTCBUSD",
		QtyStep:   0.00001,
		PriceStep: 0.01,
	}
	zeroProfitPrice := testTPDepositSpent / testTPCoinsQty

	fees := structs.OrderFees{
		BaseAsset:  decimal.NewFromInt(0),
		QuoteAsset: decimal.NewFromFloat(0.03264),
	}

	srv := NewCalcTPOrderService().CoinsQty(testTPCoinsQty).
		Profit(testTPProfitPercent).
		DepositSpent(testTPDepositSpent).
		Strategy(pkgStructs.BotStrategyShort).
		PairData(pairData).
		Fees(fees)

	// when
	order, err := srv.Do()

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(25807.68), order.Price)
	assert.Equal(t, float64(0.00126), order.Qty)
	assert.Less(t, order.Price, zeroProfitPrice)
	assert.Equal(t, pkgStructs.OrderTypeBuy, order.Type)
	assert.NotEmpty(t, order.ClientOrderID)
}

func TestCalcTPOrderLongNoFees(t *testing.T) {
	// given
	pairData := structs.ExchangePairData{
		Symbol:    "BTCBUSD",
		QtyStep:   0.00001,
		PriceStep: 0.01,
	}
	zeroProfitPrice := testTPDepositSpent / testTPCoinsQty

	srv := NewCalcTPOrderService().
		CoinsQty(testTPCoinsQty).
		Profit(testTPProfitPercent).
		DepositSpent(testTPDepositSpent).
		Strategy(pkgStructs.BotStrategyLong).
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

func TestCalcTPOrderLongFees(t *testing.T) {
	// given
	pairData := structs.ExchangePairData{
		Symbol:    "BTCBUSD",
		QtyStep:   0.00001,
		PriceStep: 0.01,
	}
	zeroProfitPrice := testTPDepositSpent / testTPCoinsQty

	fees := structs.OrderFees{
		BaseAsset:  decimal.NewFromFloat(testTPCoinsQty * 0.001),
		QuoteAsset: decimal.NewFromFloat(0),
	}

	srv := NewCalcTPOrderService().
		CoinsQty(testTPCoinsQty).
		Profit(testTPProfitPercent).
		DepositSpent(testTPDepositSpent).
		Strategy(pkgStructs.BotStrategyLong).
		PairData(pairData).
		Fees(fees)

	// when
	order, err := srv.Do()

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(26014.75), order.Price)
	assert.Equal(t, float64(0.00126), order.Qty)
	assert.Greater(t, order.Price, zeroProfitPrice)
	assert.Equal(t, pkgStructs.OrderTypeSell, order.Type)
	assert.NotEmpty(t, order.ClientOrderID)
}
