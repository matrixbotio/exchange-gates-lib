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

	proc := NewCalcTPOrderProcessor().CoinsQty(testTPCoinsQty).
		Profit(testTPProfitPercent).
		DepositSpent(testTPDepositSpent).
		Strategy(pkgStructs.BotStrategyShort).
		PairData(pairData)

	// when
	result, err := proc.Do()

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(25833.5), result.TPOrder.Price)
	assert.Equal(t, float64(0.00126), result.TPOrder.Qty)
	assert.LessOrEqual(t, result.TPOrder.Price, zeroProfitPrice.InexactFloat64())
	assert.Equal(t, pkgStructs.OrderTypeBuy, result.TPOrder.Type)
	assert.NotEmpty(t, result.TPOrder.ClientOrderID)
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
	result, err := srv.Do()

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(0.753), result.TPOrder.Price)
	assert.Equal(t, float64(348), result.TPOrder.Qty)
	assert.Equal(t, pkgStructs.OrderTypeBuy, result.TPOrder.Type)
	assert.NotEmpty(t, result.TPOrder.ClientOrderID)
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
	result, err := proc.Do()

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(25807.68), result.TPOrder.Price)
	assert.Equal(t, float64(0.00126), result.TPOrder.Qty)
	assert.LessOrEqual(t, result.TPOrder.Price, zeroProfitPriceFloat)
	assert.Equal(t, pkgStructs.OrderTypeBuy, result.TPOrder.Type)
	assert.NotEmpty(t, result.TPOrder.ClientOrderID)
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
	result, err := proc.Do()

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(25988.74), result.TPOrder.Price)
	assert.Equal(t, float64(0.00126), result.TPOrder.Qty)
	assert.GreaterOrEqual(t, result.TPOrder.Price, zeroProfitPriceFloat)
	assert.Equal(t, pkgStructs.OrderTypeSell, result.TPOrder.Type)
	assert.NotEmpty(t, result.TPOrder.ClientOrderID)
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
	result, err := proc.Do()

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(26014.75), result.TPOrder.Price)
	assert.Equal(t, float64(0.00125), result.TPOrder.Qty)
	assert.GreaterOrEqual(t, result.TPOrder.Price, zeroProfitPriceFloat)
	assert.Equal(t, pkgStructs.OrderTypeSell, result.TPOrder.Type)
	assert.NotEmpty(t, result.TPOrder.ClientOrderID)
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
	result, err := proc.Do()

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(69.9), result.TPOrder.Price)
	assert.Equal(t, float64(0.088), result.TPOrder.Qty)
	assert.Equal(t, pkgStructs.OrderTypeSell, result.TPOrder.Type)
	assert.NotEmpty(t, result.TPOrder.ClientOrderID)
}

func TestCalcTPOrderLongRemainsBigQtyStep(t *testing.T) {
	// given
	pairData := structs.ExchangePairData{
		Symbol:    "FIREUSDT",
		QtyStep:   1,
		PriceStep: 0.0001,
	}
	coinsQty := float64(2)
	depoSpent := decimal.NewFromFloat(2.6694)

	fees := structs.OrderFees{
		BaseAsset:  decimal.NewFromFloat(0.000088),
		QuoteAsset: decimal.Zero,
	}

	zeroProfitPrice := depoSpent.InexactFloat64() / coinsQty

	accBase := decimal.NewFromFloat(0.9966683)
	accQuote := decimal.Zero

	proc := NewCalcTPOrderProcessor().
		CoinsQty(coinsQty).
		Profit(0.7).
		DepositSpent(depoSpent).
		Strategy(pkgStructs.BotStrategyLong).
		PairData(pairData).
		Remains(accBase, accQuote).
		Fees(fees)

	// when
	result, err := proc.Do()

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(1.3441), result.TPOrder.Price)
	assert.Equal(t, float64(2), result.TPOrder.Qty)
	assert.Equal(t, pkgStructs.OrderTypeSell, result.TPOrder.Type)
	assert.NotEmpty(t, result.TPOrder.ClientOrderID)
	assert.Greater(t, result.TPOrder.Price, zeroProfitPrice)
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
	result, err := NewCalcTPOrderProcessor().CoinsQty(coinsQty).
		Profit(profit).
		DepositSpent(depoSpent).
		Strategy(pkgStructs.BotStrategyShort).
		PairData(pairData).
		Remains(accBase, accQuote).
		Fees(fees).
		Do()

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(58752.97), result.TPOrder.Price)
	assert.Equal(t, float64(0.000489), result.TPOrder.Qty)
	assert.Equal(t, pkgStructs.OrderTypeBuy, result.TPOrder.Type)
	assert.NotEmpty(t, result.TPOrder.ClientOrderID)
}

func TestCalcTPOrderShortRemainsBigQtyStep(t *testing.T) {
	// given
	pairData := structs.ExchangePairData{
		Symbol:    "FIREUSDT",
		QtyStep:   1,
		PriceStep: 0.0001,
	}
	coinsQty := float64(2)
	depoSpent := decimal.NewFromFloat(2.6694)
	profit := float64(1)

	zeroProfitPrice := depoSpent.InexactFloat64() / coinsQty

	fees := structs.OrderFees{
		BaseAsset:  decimal.NewFromInt(0),
		QuoteAsset: decimal.NewFromFloat(0.02880802933),
	}

	accBase := decimal.NewFromFloat(0)
	accQuote := decimal.NewFromFloat(0.00585)

	// when
	result, err := NewCalcTPOrderProcessor().CoinsQty(coinsQty).
		Profit(profit).
		DepositSpent(depoSpent).
		Strategy(pkgStructs.BotStrategyShort).
		PairData(pairData).
		Remains(accBase, accQuote).
		Fees(fees).
		Do()

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(1.3043), result.TPOrder.Price)
	assert.Equal(t, float64(2), result.TPOrder.Qty)
	assert.Equal(t, pkgStructs.OrderTypeBuy, result.TPOrder.Type)
	assert.NotEmpty(t, result.TPOrder.ClientOrderID)
	assert.Less(t, result.TPOrder.Price, zeroProfitPrice)
	assert.LessOrEqual(
		t, result.AccQuoteUsed.InexactFloat64(),
		fees.QuoteAsset.InexactFloat64(),
	)
}

func TestCalcTPOrderShortRemainsBigQtyStep2(t *testing.T) {
	// given
	pairData := structs.ExchangePairData{
		ExchangeID:     2,
		BaseAsset:      "SUI",
		QuoteAsset:     "USDT",
		Symbol:         "SUIUSDT",
		BasePrecision:  4,
		QuotePrecision: 4,
		MinQty:         1,
		MaxQty:         9999999,
		QtyStep:        1,
		MinPrice:       0.000001,
		PriceStep:      0.0001,
		MinDeposit:     1,
	}
	coinsQty := float64(1)
	depoSpent := decimal.NewFromFloat(1.9068)
	profit := float64(0.15)

	zeroProfitPrice := depoSpent.InexactFloat64() / coinsQty

	fees := structs.OrderFees{
		BaseAsset:  decimal.NewFromInt(0),
		QuoteAsset: decimal.NewFromFloat(0.0019068),
	}

	accBase := decimal.NewFromFloat(0)
	accQuote := decimal.NewFromFloat(0.0057874)

	// when
	result, err := NewCalcTPOrderProcessor().
		CoinsQty(coinsQty).
		Profit(profit).
		DepositSpent(depoSpent).
		Strategy(pkgStructs.BotStrategyShort).
		PairData(pairData).
		Remains(accBase, accQuote).
		Fees(fees).
		Do()

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(1.8962), result.TPOrder.Price)
	assert.Equal(t, float64(1), result.TPOrder.Qty)
	assert.Equal(t, pkgStructs.OrderTypeBuy, result.TPOrder.Type)
	assert.NotEmpty(t, result.TPOrder.ClientOrderID)
	assert.Less(t, result.TPOrder.Price, zeroProfitPrice)
	assert.LessOrEqual(
		t, result.AccQuoteUsed.InexactFloat64(),
		accQuote.InexactFloat64(),
	)
}

func TestCalcTPOrderLongRemainsBig(t *testing.T) {
	// given
	pairData := structs.ExchangePairData{
		Symbol:    "FIREUSDT",
		QtyStep:   1,
		PriceStep: 0.0001,
	}
	coinsQty := float64(2.9)
	depoSpent := decimal.NewFromFloat(4.986434)
	profit := float64(1)

	zeroProfitPrice := depoSpent.InexactFloat64() / coinsQty

	fees := structs.OrderFees{
		BaseAsset:  decimal.NewFromFloat(0.00522),
		QuoteAsset: decimal.NewFromFloat(0),
	}

	accBase := decimal.NewFromFloat(0.558992)
	accQuote := decimal.NewFromFloat(0)

	// when
	result, err := NewCalcTPOrderProcessor().CoinsQty(coinsQty).
		Profit(profit).
		DepositSpent(depoSpent).
		Strategy(pkgStructs.BotStrategyLong).
		PairData(pairData).
		Remains(accBase, accQuote).
		Fees(fees).
		Do()

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(1.7397), result.TPOrder.Price)
	assert.Equal(t, float64(3), result.TPOrder.Qty)
	assert.Equal(t, pkgStructs.OrderTypeSell, result.TPOrder.Type)
	assert.NotEmpty(t, result.TPOrder.ClientOrderID)
	assert.Greater(t, result.TPOrder.Price, zeroProfitPrice)
	assert.LessOrEqual(
		t, result.AccBaseUsed.InexactFloat64(),
		accBase.InexactFloat64(),
	)
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
	result, err := proc.Do()

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(2.515), result.TPOrder.Price)
	assert.Equal(t, float64(125.03), result.TPOrder.Qty)
	assert.GreaterOrEqual(t, result.TPOrder.Price, zeroProfitPrice)
	assert.Equal(t, pkgStructs.OrderTypeSell, result.TPOrder.Type)
	assert.NotEmpty(t, result.TPOrder.ClientOrderID)
	assert.LessOrEqual(
		t, result.AccQuoteUsed.InexactFloat64(),
		fees.QuoteAsset.InexactFloat64(),
	)
}

func TestCalcTPOrderShortFeesInBaseAsset(t *testing.T) {
	// given
	pairData := structs.ExchangePairData{
		Symbol:    "BNBUSDT",
		QtyStep:   0.001,
		PriceStep: 0.1,
	}

	strategy := pkgStructs.BotStrategyShort
	profit := float64(1)
	gridOrderPrice := decimal.NewFromFloat(596)
	gridOrderCoinsQty := decimal.NewFromFloat(0.009)
	depoSpent := gridOrderCoinsQty.Mul(gridOrderPrice)

	fees := structs.OrderFees{
		BaseAsset:  decimal.NewFromFloat(0.00000675),
		QuoteAsset: decimal.NewFromFloat(0),
	}

	accBase := decimal.Zero
	accQuote := decimal.NewFromFloat(0.0013)

	proc := NewCalcTPOrderProcessor().
		CoinsQty(gridOrderCoinsQty.InexactFloat64()).
		Profit(profit).
		DepositSpent(depoSpent).
		Strategy(strategy).
		PairData(pairData).
		Fees(fees).
		Remains(accBase, accQuote)

	// when
	result, err := proc.Do()

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(590.3), result.TPOrder.Price)
	assert.Equal(t, float64(0.009), result.TPOrder.Qty)
	assert.Less(t, result.TPOrder.Price, gridOrderPrice.InexactFloat64())
	assert.Equal(t, pkgStructs.OrderTypeBuy, result.TPOrder.Type)
	assert.NotEmpty(t, result.TPOrder.ClientOrderID)
	assert.LessOrEqual(
		t, result.AccQuoteUsed.InexactFloat64(),
		accQuote.InexactFloat64(),
	)
}

func TestCalcTPShortUsedRemains(t *testing.T) {
	// given
	pairData := structs.ExchangePairData{
		Symbol:    "ARBUSDT",
		QtyStep:   1,
		PriceStep: 0.001,
	}

	strategy := pkgStructs.BotStrategyShort
	profit := float64(0.59)
	lapCoinsQty := decimal.NewFromFloat(75)
	depoSpent := decimal.NewFromFloat(83.725)

	fees := structs.OrderFees{
		BaseAsset:  decimal.NewFromFloat(0),
		QuoteAsset: decimal.NewFromFloat(0.083725),
	}

	accBase := decimal.Zero
	accQuote := decimal.NewFromFloat(1.533)

	proc := NewCalcTPOrderProcessor().
		CoinsQty(lapCoinsQty.InexactFloat64()).
		Profit(profit).
		DepositSpent(depoSpent).
		Strategy(strategy).
		PairData(pairData).
		Fees(fees).
		Remains(accBase, accQuote)

	// when
	result, err := proc.Do()

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(1.088), result.TPOrder.Price)
	assert.Equal(t, float64(76), result.TPOrder.Qty)
	assert.Equal(t, pkgStructs.OrderTypeBuy, result.TPOrder.Type)
	assert.NotEmpty(t, result.TPOrder.ClientOrderID)
	assert.Equal(t, 0.0, result.AccBaseOriginal.InexactFloat64())
	assert.Equal(t, accQuote.InexactFloat64(), result.AccQuoteOriginal.InexactFloat64())
	assert.Equal(t, 0.0, result.AccBaseUsed.InexactFloat64())
	assert.Equal(t, 0.635149, result.AccQuoteUsed.RoundFloor(6).InexactFloat64())
	assert.LessOrEqual(
		t, result.AccQuoteUsed.InexactFloat64(),
		accQuote.InexactFloat64(),
	)
}

func TestRoundQtyDown(t *testing.T) {
	// given
	pairData := structs.ExchangePairData{
		Symbol:    "BNBUSDT",
		QtyStep:   0.001,
		PriceStep: 0.1,
	}

	val := decimal.NewFromFloat(0.00899325)

	proc := NewCalcTPOrderProcessor().
		PairData(pairData)

		// when
	result := proc.roundQtyDown(val)

	// then
	assert.Equal(t, 0.008, result.InexactFloat64())
}
