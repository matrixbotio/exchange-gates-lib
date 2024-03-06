package utils

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
	"github.com/shopspring/decimal"
)

const (
	defaultCheckOrdersTimeout = time.Second * 30
)

func GetCheckOrdersTimeout(exchangeID int) time.Duration {
	switch exchangeID {
	default:
		return defaultCheckOrdersTimeout
	case consts.ExchangeIDbinanceSpot:
		return consts.CheckOrdersTimeoutBinance
	case consts.ExchangeIDbybitSpot:
		return consts.CheckOrdersTimeoutBybit
	}
}

func GenerateUUID() string {
	return uuid.New().String()
}

// GetFloatPrecision returns the number of decimal places in a float
func GetFloatPrecision(value float64) int {
	precision := 0
	for {
		rounded := math.Round(value*math.Pow(10, float64(precision))) / math.Pow(10, float64(precision))
		if value == rounded {
			break
		}
		precision++
	}
	return precision
}

func roundFloatToDecimal(val float64, precision int) decimal.Decimal {
	return decimal.NewFromFloat(val).RoundFloor(int32(precision))
}

func RoundFloatFloor(val float64, precision int) (float64, error) {
	if math.IsNaN(val) {
		return 0, errors.New("value is NaN")
	}
	if math.IsInf(val, 0) {
		return 0, errors.New("value is Inf")
	}

	f, _ := roundFloatToDecimal(val, precision).Float64()
	return f, nil
}

func formatFloatFloor(val float64, precision int) (string, error) {
	if precision == 0 {
		return strconv.FormatFloat(val, 'f', 0, 64), nil
	}

	valRounded, err := RoundFloatFloor(val, precision)
	if err != nil {
		return "", fmt.Errorf("round value: %w", err)
	}

	return formatAndTrimFloat(valRounded, precision), nil
}

func formatAndTrimFloat(val float64, precision int) string {
	f := strconv.FormatFloat(val, 'f', precision, 64)
	v := strings.TrimRight(strings.TrimRight(f, "0"), ".")
	if v == "" {
		return "0"
	}
	return v
}

/*
RoundPairOrderValues - adjusts the order values in accordance
with the trading pair parameters
*/
func RoundPairOrderValues(
	order pkgStructs.BotOrder,
	pairLimits structs.ExchangePairData,
) (structs.BotOrderAdjusted, error) {
	result := structs.BotOrderAdjusted{
		PairSymbol:       order.PairSymbol,
		Type:             order.Type,
		ClientOrderID:    order.ClientOrderID,
		MinQty:           pairLimits.MinQty,
		MinQtyPassed:     true, // by default
		MinDeposit:       pairLimits.OriginalMinDeposit,
		MinDepositPassed: true, // by default
	}

	if order.Qty == 0 {
		return structs.BotOrderAdjusted{}, errors.New("order qty is not set")
	}
	if order.Qty < pairLimits.MinQty {
		result.MinQtyPassed = false
	}
	if order.Qty > pairLimits.MaxQty {
		return structs.BotOrderAdjusted{}, errors.New("too much amount to open an order in this pair. " +
			"order qty: " + strconv.FormatFloat(order.Qty, 'f', 8, 32) +
			" max: " + strconv.FormatFloat(pairLimits.MaxQty, 'f', 8, 32))
	}

	if order.Price == 0 {
		return structs.BotOrderAdjusted{}, errors.New("order price is not set")
	}
	if order.Price < pairLimits.MinPrice {
		return structs.BotOrderAdjusted{}, errors.New("insufficient price to open an order in this pair. " +
			"order price: " + strconv.FormatFloat(order.Price, 'f', 8, 32) +
			" min: " + strconv.FormatFloat(pairLimits.MinPrice, 'f', 8, 32))
	}

	// check min deposit
	orderDeposit := order.Qty * order.Price
	if orderDeposit < pairLimits.OriginalMinDeposit {
		result.MinDepositPassed = false
	}

	// round order values
	var err error
	result.Qty, err = formatFloatFloor(order.Qty, GetFloatPrecision(pairLimits.QtyStep))
	if err != nil {
		return structs.BotOrderAdjusted{}, fmt.Errorf("format qty: %w", err)
	}
	result.Price, err = formatFloatFloor(order.Price, GetFloatPrecision(pairLimits.PriceStep))
	if err != nil {
		return structs.BotOrderAdjusted{}, fmt.Errorf("format price: %w", err)
	}

	qtyRounded, err := strconv.ParseFloat(result.Qty, 64)
	if err != nil {
		return result, fmt.Errorf("parse order qty: %w", err)
	}
	priceRounded, err := strconv.ParseFloat(result.Price, 64)
	if err != nil {
		return structs.BotOrderAdjusted{}, fmt.Errorf("parse order price: %w", err)
	}

	depositRounded := qtyRounded * priceRounded
	result.Deposit = strconv.FormatFloat(depositRounded, 'f', GetFloatPrecision(orderDeposit), 64)
	return result, nil
}

// RoundDeposit - round deposit. tickerValueStep - minimum value step
func RoundDeposit(deposit, tickerValueStep float64) (float64, error) {
	depositRoundedStr := strconv.FormatFloat(
		deposit,
		'f',
		GetFloatPrecision(tickerValueStep),
		64,
	)
	depositRounded, err := strconv.ParseFloat(depositRoundedStr, 64)
	if err != nil {
		return 0, fmt.Errorf("round deposit: %w", err)
	}
	return depositRounded, nil
}

// ParseAdjustedOrder - parse rounded order to bot order
func ParseAdjustedOrder(order structs.BotOrderAdjusted) (pkgStructs.BotOrder, error) {
	resultOrder := pkgStructs.BotOrder{
		PairSymbol:    order.PairSymbol,
		Type:          order.Type,
		ClientOrderID: order.ClientOrderID,
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
		ExchangeID:    consts.PairDefaultExchangeID,
		BaseAsset:     consts.PairDefaultBaseAsset,
		QuoteAsset:    consts.PairDefaultQuoteAsset,
		MinQty:        consts.PairDefaultMinQty,
		MaxQty:        consts.PairDefaultMaxQty,
		MinDeposit:    consts.PairMinDeposit,
		MinPrice:      consts.PairDefaultMinPrice,
		QtyStep:       consts.PairDefaultQtyStep,
		PriceStep:     consts.PairDefaultPriceStep,
		AllowedMargin: true,
		AllowedSpot:   true,
	}
}

func GetValueStep(minValue float64) float64 {
	precision := GetFloatPrecision(minValue)
	divisor := decimal.NewFromFloat(math.Pow(10, float64(precision)))
	valueStep, _ := decimal.NewFromInt(1).Div(divisor).Float64()
	return valueStep
}
