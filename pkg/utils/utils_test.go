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
		Qty:     0.666666666666,
		Price:   100.66666666666666,
		Deposit: 67.111111111044,
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

func TestRoundFloatFloor1(t *testing.T) {
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

func TestRoundOrderQty(t *testing.T) {
	var orderQty float64 = 0.00056
	var qtyStep float64 = 0.00001
	var qtyPrecision int = GetFloatPrecision(qtyStep)
	assert.Equal(t, 5, qtyPrecision)

	roundedQty, err := formatFloatFloor(orderQty, qtyPrecision)
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("%v", orderQty), roundedQty)
}
