package utils

import (
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
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

// OrderResponseToBotOrder - convert raw order response to bot order
func OrderResponseToBotOrder(response structs.CreateOrderResponse) pkgStructs.BotOrder {
	return pkgStructs.BotOrder{
		PairSymbol:    response.Symbol,
		Type:          response.Type,
		Qty:           response.OrigQuantity,
		Price:         response.Price,
		Deposit:       response.OrigQuantity * response.Price,
		ClientOrderID: response.ClientOrderID,
	}
}

// OrderDataToBotOrder - convert order data to bot order
func OrderDataToBotOrder(order structs.OrderData) pkgStructs.BotOrder {
	return pkgStructs.BotOrder{
		PairSymbol:    order.Symbol,
		Type:          order.Type,
		Qty:           order.AwaitQty,
		Price:         order.Price,
		Deposit:       order.AwaitQty * order.Price,
		ClientOrderID: order.ClientOrderID,
	}
}

func RoundFloatFloor(val float64, precision int) float64 {
	powLevel := math.Pow10(precision)
	return math.Floor(val*powLevel) / powLevel
}

func formatFloatFloor(val float64, precision int) string {
	return strconv.FormatFloat(RoundFloatFloor(val, precision), 'f', precision, 64)
}

// RoundPairOrderValues - adjusts the order values in accordance with the trading pair parameters
func RoundPairOrderValues(order pkgStructs.BotOrder, pairLimits structs.ExchangePairData) (structs.BotOrderAdjusted, error) {
	result := structs.BotOrderAdjusted{
		PairSymbol:       order.PairSymbol,
		Type:             order.Type,
		ClientOrderID:    order.ClientOrderID,
		MinQty:           pairLimits.MinQty,
		MinQtyPassed:     true, // by default
		MinDeposit:       pairLimits.OriginalMinDeposit,
		MinDepositPassed: true, // by default
	}

	// check lot size
	if order.Qty < pairLimits.MinQty {
		result.MinQtyPassed = false
	}
	if order.Qty > pairLimits.MaxQty {
		return result, errors.New("too much amount to open an order in this pair. " +
			"order qty: " + strconv.FormatFloat(order.Qty, 'f', 8, 32) +
			" max: " + strconv.FormatFloat(pairLimits.MaxQty, 'f', 8, 32))
	}
	if order.Price < pairLimits.MinPrice {
		return result, errors.New("insufficient price to open an order in this pair. " +
			"order price: " + strconv.FormatFloat(order.Price, 'f', 8, 32) +
			" min: " + strconv.FormatFloat(pairLimits.MinPrice, 'f', 8, 32))
	}

	// check min deposit
	orderDeposit := order.Qty * order.Price
	if orderDeposit < pairLimits.OriginalMinDeposit {
		result.MinDepositPassed = false
	}

	// round order values
	result.Qty = formatFloatFloor(order.Qty, GetFloatPrecision(pairLimits.QtyStep))
	result.Price = formatFloatFloor(order.Price, GetFloatPrecision(pairLimits.PriceStep))

	qtyRounded, err := strconv.ParseFloat(result.Qty, 64)
	if err != nil {
		return result, err
	}
	priceRounded, err := strconv.ParseFloat(result.Price, 64)
	if err != nil {
		return result, err
	}
	depositRounded := qtyRounded * priceRounded

	result.Deposit = strconv.FormatFloat(depositRounded, 'f', GetFloatPrecision(orderDeposit), 64)
	return result, nil
}

// RoundDeposit - round deposit for grid by pair limits
func RoundDeposit(deposit float64, pairLimits structs.ExchangePairData) (float64, error) {
	depositStep := pairLimits.PriceStep * pairLimits.QtyStep
	depositRoundedStr := strconv.FormatFloat(deposit, 'f', GetFloatPrecision(depositStep), 64)
	depositRounded, err := strconv.ParseFloat(depositRoundedStr, 64)
	if err != nil {
		return 0, fmt.Errorf("round deposit: %w", err)
	}
	return depositRounded, nil
}

// ParseAdjustedOrder - parse rounded order to bot order
func ParseAdjustedOrder(order structs.BotOrderAdjusted) (pkgStructs.BotOrder, error) {
	resultOrder := pkgStructs.BotOrder{
		PairSymbol: order.PairSymbol,
		Type:       order.Type,
	}
	// parse qty
	var err error
	resultOrder.Qty, err = strconv.ParseFloat(order.Qty, 64)
	if err != nil {
		return resultOrder, fmt.Errorf("parse order qty: %w", err)
	}
	// parse price
	resultOrder.Price, err = strconv.ParseFloat(order.Price, 64)
	if err != nil {
		return resultOrder, fmt.Errorf("parse order price: %w", err)
	}
	// parse deposit
	resultOrder.Deposit, err = strconv.ParseFloat(order.Deposit, 64)
	if err != nil {
		return resultOrder, fmt.Errorf("parse order deposit: %w", err)
	}
	return resultOrder, nil
}

// GetDefaultPairData !
func GetDefaultPairData() structs.ExchangePairData {
	return structs.ExchangePairData{
		ExchangeID: consts.PairDefaultExchangeID,
		BaseAsset:  consts.PairDefaultBaseAsset,
		QuoteAsset: consts.PairDefaultQuoteAsset,
		MinQty:     consts.PairDefaultMinQty,
		MaxQty:     consts.PairDefaultMaxQty,
		MinDeposit: consts.PairMinDeposit,
		MinPrice:   consts.PairDefaultMinPrice,
		QtyStep:    consts.PairDefaultQtyStep,
	}
}

// RoundMinDeposit - update the value of the minimum deposit in accordance with the minimum threshold
func RoundMinDeposit(pairMinDeposit float64) float64 {
	return pairMinDeposit * (1 + consts.MinDepositFix/100)
}

// OrderDataToTradeEvent data
type TradeOrderConvertTask struct {
	Order       structs.OrderData
	ExchangeTag string
}

// OrderDataToTradeEvent - convert order data into a trade event.
func OrderDataToTradeEvent(task TradeOrderConvertTask) workers.TradeEvent {
	e := workers.TradeEvent{
		ID:          0,
		Time:        task.Order.UpdatedTime,
		Symbol:      task.Order.Symbol,
		Price:       task.Order.Price,
		Quantity:    task.Order.FilledQty,
		ExchangeTag: task.ExchangeTag,
	}

	if task.Order.Type == pkgStructs.OrderTypeBuy {
		e.BuyerOrderID = task.Order.OrderID
	} else {
		e.SellerOrderID = task.Order.OrderID
	}

	return e
}
