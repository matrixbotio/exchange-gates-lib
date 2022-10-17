package utils

import (
	"errors"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/go-stack/stack"

	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
	structs2 "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
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
	if len(stackTrace) == 0 {
		return ""
	}
	return stack.Trace().TrimRuntime().String()
}

// OrderResponseToBotOrder - convert raw order response to bot order
func OrderResponseToBotOrder(response structs.CreateOrderResponse) structs2.BotOrder {
	return structs2.BotOrder{
		PairSymbol:    response.Symbol,
		Type:          response.Type,
		Qty:           response.OrigQuantity,
		Price:         response.Price,
		Deposit:       response.OrigQuantity * response.Price,
		ClientOrderID: response.ClientOrderID,
	}
}

// OrderDataToBotOrder - convert order data to bot order
func OrderDataToBotOrder(order structs.OrderData) structs2.BotOrder {
	return structs2.BotOrder{
		PairSymbol:    order.Symbol,
		Type:          order.Type,
		Qty:           order.AwaitQty,
		Price:         order.Price,
		Deposit:       order.AwaitQty * order.Price,
		ClientOrderID: order.ClientOrderID,
	}
}

// RoundPairOrderValues - adjusts the order values in accordance with the trading pair parameters
func RoundPairOrderValues(order structs2.BotOrder, pairLimits structs.ExchangePairData) (structs.BotOrderAdjusted, error) {
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
	var quantityPrecision int = GetFloatPrecision(pairLimits.QtyStep)
	result.Qty = strconv.FormatFloat(order.Qty, 'f', quantityPrecision, 64)
	var ratePrecision int = GetFloatPrecision(pairLimits.PriceStep)
	result.Price = strconv.FormatFloat(order.Price, 'f', ratePrecision, 64)

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
		return 0, errors.New("failed to round deposit: " + err.Error())
	}
	return depositRounded, nil
}

// ParseAdjustedOrder - parse rounded order to bot order
func ParseAdjustedOrder(order structs.BotOrderAdjusted) (structs2.BotOrder, error) {
	resultOrder := structs2.BotOrder{
		PairSymbol: order.PairSymbol,
		Type:       order.Type,
	}
	// parse qty
	var err error
	resultOrder.Qty, err = strconv.ParseFloat(order.Qty, 64)
	if err != nil {
		return resultOrder, errors.New("failed to parse order qty: " + err.Error())
	}
	// parse price
	resultOrder.Price, err = strconv.ParseFloat(order.Price, 64)
	if err != nil {
		return resultOrder, errors.New("failed to parse order price: " + err.Error())
	}
	// parse deposit
	resultOrder.Deposit, err = strconv.ParseFloat(order.Deposit, 64)
	if err != nil {
		return resultOrder, errors.New("failed to parse order deposit: " + err.Error())
	}
	return resultOrder, nil
}

// RunTimeLimitHandler - func runtime limit handler
type RunTimeLimitHandler struct {
	timeout time.Duration
	runFunc func()
}

// NewRuntimeLimitHandler - create new func runtime limit handler
func NewRuntimeLimitHandler(timeout time.Duration, runFunc func()) *RunTimeLimitHandler {
	return &RunTimeLimitHandler{
		timeout: timeout,
		runFunc: runFunc,
	}
}

// Run - run func & limit runtime.
// returns: bool: true if time is up
func (r *RunTimeLimitHandler) Run() bool {
	timeTo := time.After(r.timeout)
	done := make(chan bool, 1)

	go func() {
		for {
			select {
			case <-timeTo:
				done <- true
				return
			//lint:ignore SA5004 it's meant to be
			default:
				// wait
			}
		}
	}()

	go func() {
		r.runFunc()
		done <- false
	}()

	return <-done
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

	if task.Order.Type == structs2.OrderTypeBuy {
		e.BuyerOrderID = task.Order.OrderID
	} else {
		e.SellerOrderID = task.Order.OrderID
	}

	return e
}
