package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
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

	return roundFloatToDecimal(val, precision).InexactFloat64(), nil
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

func PrintObject(o any) {
	data, err := json.MarshalIndent(o, "", "	")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(data))
}

func StringPointer(val string) *string {
	return &val
}
