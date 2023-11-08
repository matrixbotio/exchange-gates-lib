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

func TestCalcTPOrderErrorEmptyStrategy(t *testing.T) {
	// given
	proc := NewCalcTPOrderProcessor().CoinsQty(0)

	// when
	_, err := proc.Do()

	// then
	require.Contains(t, err.Error(), "strategy is not set")
}

func TestCalcTPOrderShortNoFees(t *testing.T) {
	// given
	pairData := structs.ExchangePairData{
		Symbol:    "BTCBUSD",
		QtyStep:   0.00001,
		PriceStep: 0.01,
	}
	zeroProfitPrice := testTPDepositSpent / testTPCoinsQty

	proc := NewCalcTPOrderProcessor().CoinsQty(testTPCoinsQty).
		Profit(testTPProfitPercent).
		DepositSpent(testTPDepositSpent).
		Strategy(pkgStructs.BotStrategyShort).
		PairData(pairData)

	// when
	order, err := proc.Do()

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(25833.5), order.Price)
	assert.Equal(t, float64(0.00126), order.Qty)
	assert.Less(t, order.Price, zeroProfitPrice)
	assert.Equal(t, pkgStructs.OrderTypeBuy, order.Type)
	assert.NotEmpty(t, order.ClientOrderID)
}

func TestCalcTPOrderShortNoFeesBigQty(t *testing.T) {
	// given
	pairData := structs.ExchangePairData{
		Symbol:     "TWTBUSD",
		MinQty:     1,
		MinDeposit: 10,
		QtyStep:    1,
		MinPrice:   0.0001,
		PriceStep:  0.0001,
	}

	srv := NewCalcTPOrderProcessor().
		CoinsQty(348).
		Profit(0.1).
		DepositSpent(262.33212).
		Strategy(pkgStructs.BotStrategyShort).
		PairData(pairData)

	// when
	order, err := srv.Do()

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(0.753), order.Price)
	assert.Equal(t, float64(348), order.Qty)
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

	proc := NewCalcTPOrderProcessor().CoinsQty(testTPCoinsQty).
		Profit(testTPProfitPercent).
		DepositSpent(testTPDepositSpent).
		Strategy(pkgStructs.BotStrategyShort).
		PairData(pairData).
		Fees(fees)

	// when
	order, err := proc.Do()

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

	proc := NewCalcTPOrderProcessor().
		CoinsQty(testTPCoinsQty).
		Profit(testTPProfitPercent).
		DepositSpent(testTPDepositSpent).
		Strategy(pkgStructs.BotStrategyLong).
		PairData(pairData)

	// when
	order, err := proc.Do()

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

	proc := NewCalcTPOrderProcessor().
		CoinsQty(testTPCoinsQty).
		Profit(testTPProfitPercent).
		DepositSpent(testTPDepositSpent).
		Strategy(pkgStructs.BotStrategyLong).
		PairData(pairData).
		Fees(fees)

	// when
	order, err := proc.Do()

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(26014.75), order.Price)
	assert.Equal(t, float64(0.00126), order.Qty)
	assert.Greater(t, order.Price, zeroProfitPrice)
	assert.Equal(t, pkgStructs.OrderTypeSell, order.Type)
	assert.NotEmpty(t, order.ClientOrderID)
}
