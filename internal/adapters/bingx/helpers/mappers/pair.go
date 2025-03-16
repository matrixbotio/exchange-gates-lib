package mappers

import (
	"strings"

	bingxgo "github.com/matrixbotio/go-bingx"
	"github.com/shopspring/decimal"

	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	"github.com/matrixbotio/exchange-gates-lib/pkg/utils"
)

const (
	pairSymbolDelimiter   = "-"
	defaultQuotePrecision = 2
)

func ConvertPairData(
	data bingxgo.SymbolInfo,
	lastPrice float64,
) (structs.ExchangePairData, error) {
	symbolParts := strings.Split(data.Symbol, pairSymbolDelimiter)
	lastPriceDec := decimal.NewFromFloat(lastPrice)

	var minQty float64
	var maxQty float64
	if data.MinQty > 0 {
		// use exchange deprecated field when available
		minQty = data.MinQty
	} else if lastPrice > 0 {
		// or recalc min qty based on min order
		minQty = decimal.NewFromFloat(data.MinNotional).
			Div(lastPriceDec).InexactFloat64()
	}

	if data.MaxQty > 0 {
		maxQty = data.MaxQty
	} else if lastPrice > 0 {
		maxQty = decimal.NewFromFloat(data.MaxNotional).
			Div(lastPriceDec).InexactFloat64()
	}

	quotePricision := utils.GetFloatPrecision(data.MinNotional)
	if quotePricision == 0 {
		quotePricision = defaultQuotePrecision
	}

	return structs.ExchangePairData{
		ExchangeID:         consts.ExchangeIDbingx,
		BaseAsset:          symbolParts[0],
		QuoteAsset:         symbolParts[1],
		BasePrecision:      utils.GetFloatPrecision(data.StepSize),
		QuotePrecision:     quotePricision,
		Status:             ConvertPairStatus(data.Status),
		Symbol:             data.Symbol,
		MinQty:             minQty,
		MaxQty:             maxQty,
		OriginalMinDeposit: data.MinNotional,
		MinDeposit:         data.MinNotional,
		MinPrice:           0, // TBD
		QtyStep:            data.StepSize,
		PriceStep:          data.TickSize,
		AllowedMargin:      false,
		AllowedSpot:        true,
		InUse:              true,
	}, nil
}

func ConvertPairs(
	pairs []bingxgo.SymbolInfo,
	tickers bingxgo.Tickers,
) ([]structs.ExchangePairData, error) {
	var result []structs.ExchangePairData

	for _, pair := range pairs {
		lastPrice := tickers[pair.Symbol]

		data, err := ConvertPairData(pair, lastPrice)
		if err != nil {
			return nil, err
		}

		result = append(result, data)
	}

	return result, nil
}

// exchange status -> our const
var pairStatusConverter = map[int]string{
	0:  consts.PairStatusOffline,
	1:  consts.PairStatusTrading,
	5:  consts.PairStatusPreOpen,
	25: consts.PairStatusSuspended,
}

func ConvertPairStatus(status int) string {
	result, isExists := pairStatusConverter[status]
	if !isExists {
		return consts.PairStatusTrading
	}
	return result
}
