package utils

import (
	"strconv"
	"testing"

	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
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
