package utils

import (
	"strconv"
	"testing"

	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
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

func TestGetFloatPrecision4(t *testing.T) {
	// given
	var val float64 = 1
	var precisionExpected int = 0

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

func TestGetFloatPrecisionPriceStep(t *testing.T) {
	assert.Equal(t, 5, GetFloatPrecision(0.00001))
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
