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
	testTPDepositSpent          = decimal.NewFromFloat(32.64787)
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
	zeroProfitPrice := testTPDepositSpent.Div(decimal.NewFromFloat(testTPCoinsQty))
	zeroProfitPriceFloat, _ := zeroProfitPrice.Float64()

	proc := NewCalcTPOrderProcessor().CoinsQty(testTPCoinsQty).
		Profit(testTPProfitPercent).
		DepositSpent(testTPDepositSpent).
		Strategy(pkgStructs.BotStrategyShort).
		PairData(pairData)

	// when
	order, err := proc.Do()

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(25833.507414265143), order.Price)
	assert.Equal(t, float64(0.00126378), order.Qty)
	assert.LessOrEqual(t, order.Price, zeroProfitPriceFloat)
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
		DepositSpent(decimal.NewFromFloat(262.33212)).
		Strategy(pkgStructs.BotStrategyShort).
		PairData(pairData)

	// when
	order, err := srv.Do()

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(0.7530748561783045), order.Price)
	assert.Equal(t, float64(348.348), order.Qty)
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
	zeroProfitPrice := testTPDepositSpent.Div(decimal.NewFromFloat(testTPCoinsQty))
	zeroProfitPriceFloat, _ := zeroProfitPrice.Float64()

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
	assert.Equal(t, float64(25807.680134200575), order.Price)
	assert.Equal(t, float64(0.00126378), order.Qty)
	assert.LessOrEqual(t, order.Price, zeroProfitPriceFloat)
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
	zeroProfitPrice := testTPDepositSpent.Div(decimal.NewFromFloat(testTPCoinsQty))
	zeroProfitPriceFloat, _ := zeroProfitPrice.Float64()

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
	assert.Equal(t, float64(25988.74096031746), order.Price)
	assert.Equal(t, float64(0.00126), order.Qty)
	assert.GreaterOrEqual(t, order.Price, zeroProfitPriceFloat)
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
	zeroProfitPrice := testTPDepositSpent.Div(decimal.NewFromFloat(testTPCoinsQty))
	zeroProfitPriceFloat, _ := zeroProfitPrice.Float64()

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
	assert.Equal(t, float64(26014.755716033495), order.Price)
	assert.Equal(t, float64(0.00125874), order.Qty)
	assert.GreaterOrEqual(t, order.Price, zeroProfitPriceFloat)
	assert.Equal(t, pkgStructs.OrderTypeSell, order.Type)
	assert.NotEmpty(t, order.ClientOrderID)
}

func TestCalcTPOrderLongRemains(t *testing.T) {
	// given
	pairData := structs.ExchangePairData{
		Symbol:    "QNTUSDT",
		QtyStep:   0.001,
		PriceStep: 0.1,
	}
	coinsQty := 0.088
	depoSpent := decimal.NewFromFloat(6.1205)

	fees := structs.OrderFees{
		BaseAsset:  decimal.NewFromFloat(0.000088),
		QuoteAsset: decimal.Zero,
	}
	accBase := decimal.NewFromFloat(0.001)
	accQuote := decimal.Zero

	proc := NewCalcTPOrderProcessor().
		CoinsQty(coinsQty).
		Profit(0.53).
		DepositSpent(depoSpent).
		Strategy(pkgStructs.BotStrategyLong).
		PairData(pairData).
		Remains(accBase, accQuote).
		Fees(fees)

	// when
	order, err := proc.Do()

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(69.2025671450423), order.Price)
	assert.Equal(t, float64(0.088912), order.Qty)
	assert.Equal(t, pkgStructs.OrderTypeSell, order.Type)
	assert.NotEmpty(t, order.ClientOrderID)
}

func TestCalcTPOrderShortRemains(t *testing.T) {
	// given
	pairData := structs.ExchangePairData{
		Symbol:    "BTCUSDT",
		QtyStep:   0.000001,
		PriceStep: 0.01,
	}
	coinsQty := 0.000484
	depoSpent := decimal.NewFromFloat(28.80802933)
	profit := float64(1)

	fees := structs.OrderFees{
		BaseAsset:  decimal.NewFromInt(0),
		QuoteAsset: decimal.NewFromFloat(0.02880802933),
	}

	accBase := decimal.NewFromFloat(0)
	accQuote := decimal.NewFromFloat(0.05853762067)

	// when
	order, err := NewCalcTPOrderProcessor().CoinsQty(coinsQty).
		Profit(profit).
		DepositSpent(depoSpent).
		Strategy(pkgStructs.BotStrategyShort).
		PairData(pairData).
		Remains(accBase, accQuote).
		Fees(fees).
		Do()

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(58992.22428880615), order.Price)
	assert.Equal(t, float64(0.00048884), order.Qty)
	assert.Equal(t, pkgStructs.OrderTypeBuy, order.Type)
	assert.NotEmpty(t, order.ClientOrderID)
}

func TestCalcTPOrderLongFeesBigQty(t *testing.T) {
	// given
	pairData := structs.ExchangePairData{
		Symbol:    "XRPUSDT",
		QtyStep:   0.01,
		PriceStep: 0.001,
	}

	averagePrice := float64(2.5)
	coinsQty := float64(125.16)
	depoSpent := coinsQty * averagePrice
	zeroProfitPrice := depoSpent / coinsQty

	fees := structs.OrderFees{
		BaseAsset:  decimal.NewFromFloat(coinsQty * 0.001),
		QuoteAsset: decimal.NewFromFloat(0),
	}

	proc := NewCalcTPOrderProcessor().
		CoinsQty(coinsQty).
		Profit(0.5).
		DepositSpent(decimal.NewFromFloat(depoSpent)).
		Strategy(pkgStructs.BotStrategyLong).
		PairData(pairData).
		Fees(fees)

	// when
	order, err := proc.Do()

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(2.515015015015015), order.Price)
	assert.Equal(t, float64(125.03484), order.Qty)
	assert.GreaterOrEqual(t, order.Price, zeroProfitPrice)
	assert.Equal(t, pkgStructs.OrderTypeSell, order.Type)
	assert.NotEmpty(t, order.ClientOrderID)
}
