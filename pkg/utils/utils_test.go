package utils

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetFloatPrecision(t *testing.T) {
	floatVal := 56.13954
	precisionExpected := 5
	precision := GetFloatPrecision(floatVal)
	if precision != precisionExpected {
		t.Fatalf("count float value precision. Received " +
			strconv.Itoa(precision) + ", expected " + strconv.Itoa(precisionExpected))
	}
}

func TestGetFloatPrecision2(t *testing.T) {
	// given
	var val float64 = 30
	var precisionExpected int = 0

	// when
	var precision = GetFloatPrecision(val)

	// then
	assert.Equal(t, precisionExpected, precision)
}

func TestGetFloatPrecision3(t *testing.T) {
	// given
	var val float64 = 0.000048
	var precisionExpected int = 6

	// when
	var precision = GetFloatPrecision(val)

	// then
	assert.Equal(t, precisionExpected, precision)
}

func TestOrderResponseToBotOrder(t *testing.T) {
	fromOrder := structs.CreateOrderResponse{}

	toOrder := OrderResponseToBotOrder(fromOrder)

	if toOrder.ClientOrderID != fromOrder.ClientOrderID {
		t.Fatal("ClientOrderID is not equal in orders")
	}
	if toOrder.PairSymbol != fromOrder.Symbol {
		t.Fatal("PairSymbol is not equal in orders")
	}
	if toOrder.Type != fromOrder.Type {
		t.Fatal("Type is not equal in orders")
	}
	if toOrder.Qty != fromOrder.OrigQuantity {
		t.Fatal("Qty is not equal in orders")
	}
	if toOrder.Price != fromOrder.Price {
		t.Fatal("Price is not equal in orders")
	}
}

func TestRoundPairOrderValues(t *testing.T) {
	originalOrder := pkgStructs.BotOrder{
		Qty:           0.666666666666,
		Price:         100.66666666666666,
		Deposit:       67.111111111044,
		ClientOrderID: "12345",
	}

	pairData := structs.ExchangePairData{
		BaseAsset:          "ETH",
		QuoteAsset:         "USDT",
		BasePrecision:      8,
		QuotePrecision:     8,
		Symbol:             "ETHUSDT",
		MinQty:             0.0001,
		MaxQty:             9000,
		OriginalMinDeposit: 10,
		MinDeposit:         11,
		MinPrice:           0.01,
		QtyStep:            0.0001,
		PriceStep:          0.01,
	}

	roundedOrder, err := RoundPairOrderValues(originalOrder, pairData)
	require.Nil(t, err)

	parsedOrder, err := ParseAdjustedOrder(roundedOrder)
	require.Nil(t, err)

	assert.Equal(t, 0.6666, parsedOrder.Qty)
	assert.Equal(t, 100.66, parsedOrder.Price)
	assert.LessOrEqual(t, parsedOrder.Deposit, originalOrder.Deposit)
	assert.NotEmpty(t, parsedOrder.ClientOrderID)
}

func TestRoundPairOrderValues2(t *testing.T) {
	originalOrder := pkgStructs.BotOrder{
		Qty:   0.191,
		Price: 70,
	}

	pairData := structs.ExchangePairData{
		MinPrice:  0.01,
		PriceStep: 0.01,
		MinQty:    0.001,
		QtyStep:   0.001,
		MaxQty:    99999,
	}

	// when
	roundedOrder, err := RoundPairOrderValues(originalOrder, pairData)

	// then
	require.Nil(t, err)
	assert.Equal(t, "70", roundedOrder.Price)
	assert.Equal(t, "0.191", roundedOrder.Qty)
}

func TestFormatFloatFloor(t *testing.T) {
	// given
	qty := float64(0.00056)
	precision := int(5)
	qtyFormatedExpected := "0.00056"

	// when
	qtyFormated, err := formatFloatFloor(qty, precision)

	// then
	require.NoError(t, err)
	assert.Equal(t, qtyFormatedExpected, qtyFormated)
}

func TestFormatFloatFloor2(t *testing.T) {
	// given
	qty := float64(30)
	precision := 0
	qtyFormatedExpected := "30"

	// when
	qtyFormated, err := formatFloatFloor(qty, precision)

	// then
	require.NoError(t, err)
	assert.Equal(t, qtyFormatedExpected, qtyFormated)
}

func TestFormatFloatFloor3(t *testing.T) {
	// given
	qty := float64(0.0)
	precision := 0
	qtyFormatedExpected := "0"

	// when
	qtyFormated, err := formatFloatFloor(qty, precision)

	// then
	require.NoError(t, err)
	assert.Equal(t, qtyFormatedExpected, qtyFormated)
}

func TestRoundFloatToDecimal(t *testing.T) {
	// given
	val := float64(70)
	precision := int(2)

	// when
	result := roundFloatToDecimal(val, precision)
	f, _ := result.Float64()

	// then
	assert.Equal(t, val, f)
}

func TestFormatAndTrimFloat(t *testing.T) {
	// given
	values := []struct {
		valRaw      float64
		valExpected string
	}{
		{valRaw: 70, valExpected: "70"},
		{valRaw: 60.1, valExpected: "60.1"},
		{valRaw: 15.32, valExpected: "15.32"},
		{valRaw: 1250, valExpected: "1250"},
		{valRaw: 0.918, valExpected: "0.92"},
		{valRaw: 0.01, valExpected: "0.01"},
	}
	precision := int(2)

	for _, v := range values {
		// when
		valFormatted := formatAndTrimFloat(v.valRaw, precision)

		// then
		assert.Equal(t, v.valExpected, valFormatted)
	}
}

func TestFormatFloatFloor6(t *testing.T) {
	// given
	price := float64(70)
	priceStep := float64(0.01)
	pricePrecision := GetFloatPrecision(priceStep)
	valFormatedExpected := "70"

	// when
	valFormated, err := formatFloatFloor(price, pricePrecision)

	// then
	require.NoError(t, err)
	assert.Equal(t, valFormatedExpected, valFormated)
}

func TestRoundFloatFloor(t *testing.T) {
	// given
	val := float64(0.00053)
	precision := int(5)

	// when
	valRounded, err := RoundFloatFloor(val, precision)

	// then
	require.NoError(t, err)
	assert.Equal(t, val, valRounded)
}

func TestRoundFloatFloor2(t *testing.T) {
	// given
	val := float64(0.00056)
	precision := int(5)

	// when
	valRounded, err := RoundFloatFloor(val, precision)

	// then
	require.NoError(t, err)
	assert.Equal(t, val, valRounded)
}

func TestRoundFloatFloor3(t *testing.T) {
	// given
	val := float64(0.666666666666)
	valRoundedExpected := float64(0.6666)
	precision := int(4)

	// when
	valRounded, err := RoundFloatFloor(val, precision)

	// then
	require.NoError(t, err)
	assert.Equal(t, valRoundedExpected, valRounded)
}

func TestFormatFloatFloor4(t *testing.T) {
	// given
	qty := float64(4.0000000000000295e-07)
	precision := 3
	qtyFormatedExpected := "0"

	// when
	qtyFormated, err := formatFloatFloor(qty, precision)

	// then
	require.NoError(t, err)
	assert.Equal(t, qtyFormatedExpected, qtyFormated)
}

func TestFormatFloatFloor5(t *testing.T) {
	// given
	qty := float64(0.00555)
	qtyStep := 0.00001
	precision := GetFloatPrecision(qtyStep)
	qtyFormatedExpected := "0.00555"

	// when
	qtyFormated, err := formatFloatFloor(qty, precision)

	// then
	require.NoError(t, err)
	assert.Equal(t, qtyFormatedExpected, qtyFormated)
}

func TestRoundOrderQty(t *testing.T) {
	var orderQty float64 = 0.00056
	var qtyStep float64 = 0.00001
	var qtyPrecision int = GetFloatPrecision(qtyStep)
	assert.Equal(t, 5, qtyPrecision)

	roundedQty, err := formatFloatFloor(orderQty, qtyPrecision)
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("%v", orderQty), roundedQty)
}

func TestGetFloatPrecisionPriceStep(t *testing.T) {
	assert.Equal(t, 5, GetFloatPrecision(0.00001))
}

func TestFormatPrice(t *testing.T) {
	// given
	var precision int = 5
	var price float64 = 0.067302

	// when
	valFormatted, err := formatFloatFloor(price, precision)

	// then
	require.NoError(t, err)
	assert.Equal(t, "0.0673", valFormatted)
}

func TestFormatFractionalPart(t *testing.T) {
	// given
	var precision int = 4
	var price float64 = 1.00002

	// when
	valFormatted, err := formatFloatFloor(price, precision)

	// then
	require.NoError(t, err)
	assert.Equal(t, "1", valFormatted)
}

func TestRoundDeposit(t *testing.T) {
	// given
	deposit := float64(100.1)
	depositStep := float64(0.00001)

	// when
	roundedDeposit, err := RoundDeposit(deposit, depositStep)

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(100.1), roundedDeposit)
}

func TestRoundDeposit2(t *testing.T) {
	// given
	deposit := float64(100.1)
	depositStep := float64(0.1)

	// when
	roundedDeposit, err := RoundDeposit(deposit, depositStep)

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(100.1), roundedDeposit)
}

func TestRoundDeposit3(t *testing.T) {
	// given
	deposit := float64(100.1)
	depositStep := float64(0)

	// when
	roundedDeposit, err := RoundDeposit(deposit, depositStep)

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(100), roundedDeposit)
}

func TestGetQtyStep(t *testing.T) {
	// given
	var minQty = 0.000048
	var qtyStepExpected = 0.000001

	// when
	qtyStep := GetValueStep(minQty)

	// then
	assert.Equal(t, qtyStepExpected, qtyStep)
}

func TestGetQtyStep2(t *testing.T) {
	// given
	var minQty = 0.00001
	var qtyStepExpected = 0.00001

	// when
	qtyStep := GetValueStep(minQty)

	// then
	assert.Equal(t, qtyStepExpected, qtyStep)
}

func TestGetValueStep_Successful(t *testing.T) {
	// given
	minOrderQTY := 821.02

	// when
	result := GetValueStep(minOrderQTY)

	// then
	assert.Equal(t, 0.01, result)
}
