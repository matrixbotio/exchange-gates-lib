package matrixgates

import (
	"errors"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/go-stack/stack"
)

// GetFloatPrecision returns the number of decimal places in a float
func GetFloatPrecision(value float64) int {
	// if you put 15, then the test will fall,
	// because Float is rounded incorrectly
	maxPrecision := 14
	valueFormated := strconv.FormatFloat(math.Abs(value), 'f', maxPrecision, 64)
	valueParts := strings.Split(valueFormated, ".")
	if len(valueParts) <= 1 {
		return 0
	}
	valueLastPartTrimmed := strings.TrimRight(valueParts[1], "0")
	return len(valueLastPartTrimmed)
}

// LogNotNilError logs an array of errors and returns true if an error is found
func LogNotNilError(errs []error) bool {
	for _, err := range errs {
		if err != nil {
			log.Println(err)
			return true
		}
	}
	return false
}

// GetTrace - get stack string
func GetTrace() string {
	stackTrace := stack.Trace()
	if stackTrace == nil || len(stackTrace) == 0 {
		return ""
	}
	return stack.Trace().TrimRuntime().String()
}

// roundPairOrderValues - adjusts the order values in accordance with the trading pair parameters
func roundPairOrderValues(order BotOrder, pairLimits ExchangePairData) (BotOrderAdjusted, error) {
	result := BotOrderAdjusted{
		PairSymbol: order.PairSymbol,
		Type:       order.Type,
	}

	// check lot size
	if order.Qty < pairLimits.MinQty {
		return result, errors.New("bot order invalid error: insufficient amount to open an order in this pair, stack: " + GetTrace())
	}
	if order.Qty > pairLimits.MaxQty {
		return result, errors.New("bot order invalid error: too much amount to open an order in this pair, stack: " + GetTrace())
	}
	if order.Price < pairLimits.MinPrice {
		return result, errors.New("bot order invalid error: insufficient price to open an order in this pair, stack: " + GetTrace())
	}

	// check min deposit
	orderDeposit := order.Qty * order.Price
	if orderDeposit < pairLimits.MinDeposit {
		return result, errors.New("the order deposit is less than the minimum")
	}

	// round order values
	var quantityPrecision int = GetFloatPrecision(pairLimits.QtyStep)
	result.Qty = strconv.FormatFloat(order.Qty, 'f', quantityPrecision, 64)
	var ratePrecision int = GetFloatPrecision(pairLimits.PriceStep)
	result.Price = strconv.FormatFloat(order.Price, 'f', ratePrecision, 32)
	result.Deposit = strconv.FormatFloat(orderDeposit, 'f', GetFloatPrecision(orderDeposit), 32)
	return result, nil
}
