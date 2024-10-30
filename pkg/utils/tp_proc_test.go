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
	order, err := proc.Do()

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(25833.5), order.Price)
	assert.Equal(t, float64(0.00126), order.Qty)
	assert.LessOrEqual(t, order.Price, zeroProfitPrice.InexactFloat64())
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
	assert.Equal(t, float64(25807.68), order.Price)
	assert.Equal(t, float64(0.00126), order.Qty)
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
	assert.Equal(t, float64(25988.74), order.Price)
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
	assert.Equal(t, float64(26014.75), order.Price)
	assert.Equal(t, float64(0.00125), order.Qty)
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
	assert.Equal(t, float64(69.9), order.Price)
	assert.Equal(t, float64(0.088), order.Qty)
	assert.Equal(t, pkgStructs.OrderTypeSell, order.Type)
	assert.NotEmpty(t, order.ClientOrderID)
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
	order, err := proc.Do()

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(1.3441), order.Price)
	assert.Equal(t, float64(2), order.Qty)
	assert.Equal(t, pkgStructs.OrderTypeSell, order.Type)
	assert.NotEmpty(t, order.ClientOrderID)
	assert.Greater(t, order.Price, zeroProfitPrice)
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
	assert.Equal(t, float64(58752.97), order.Price)
	assert.Equal(t, float64(0.000489), order.Qty)
	assert.Equal(t, pkgStructs.OrderTypeBuy, order.Type)
	assert.NotEmpty(t, order.ClientOrderID)
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
	assert.Equal(t, float64(1.3043), order.Price)
	assert.Equal(t, float64(2), order.Qty)
	assert.Equal(t, pkgStructs.OrderTypeBuy, order.Type)
	assert.NotEmpty(t, order.ClientOrderID)
	assert.Less(t, order.Price, zeroProfitPrice)
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
	order, err := NewCalcTPOrderProcessor().
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
	assert.Equal(t, float64(1.8962), order.Price)
	assert.Equal(t, float64(1), order.Qty)
	assert.Equal(t, pkgStructs.OrderTypeBuy, order.Type)
	assert.NotEmpty(t, order.ClientOrderID)
	assert.Less(t, order.Price, zeroProfitPrice)
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
	order, err := NewCalcTPOrderProcessor().CoinsQty(coinsQty).
		Profit(profit).
		DepositSpent(depoSpent).
		Strategy(pkgStructs.BotStrategyLong).
		PairData(pairData).
		Remains(accBase, accQuote).
		Fees(fees).
		Do()

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(1.7397), order.Price)
	assert.Equal(t, float64(3), order.Qty)
	assert.Equal(t, pkgStructs.OrderTypeSell, order.Type)
	assert.NotEmpty(t, order.ClientOrderID)
	assert.Greater(t, order.Price, zeroProfitPrice)
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
	assert.Equal(t, float64(2.515), order.Price)
	assert.Equal(t, float64(125.03), order.Qty)
	assert.GreaterOrEqual(t, order.Price, zeroProfitPrice)
	assert.Equal(t, pkgStructs.OrderTypeSell, order.Type)
	assert.NotEmpty(t, order.ClientOrderID)
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
	order, err := proc.Do()

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(590.3), order.Price)
	assert.Equal(t, float64(0.009), order.Qty)
	assert.Less(t, order.Price, gridOrderPrice.InexactFloat64())
	assert.Equal(t, pkgStructs.OrderTypeBuy, order.Type)
	assert.NotEmpty(t, order.ClientOrderID)
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

func TestRoundQtyUp(t *testing.T) {
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
	result := proc.roundQtyUp(val)

	// then
	assert.Equal(t, 0.009, result.InexactFloat64())
}
