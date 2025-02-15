package mappers

import (
	"fmt"
	"strings"

	bingxgo "github.com/Sagleft/go-bingx"
	"github.com/shopspring/decimal"

	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	"github.com/matrixbotio/exchange-gates-lib/pkg/utils"
)

const pairSymbolDelimiter = "-"

func ConvertPairData(
	data bingxgo.SymbolInfo,
	lastPrice float64,
) (structs.ExchangePairData, error) {
	if lastPrice == 0 {
		return structs.ExchangePairData{},
			fmt.Errorf("%q price not set", data.Symbol)
	}

	symbolParts := strings.Split(data.Symbol, pairSymbolDelimiter)

	lastPriceDec := decimal.NewFromFloat(lastPrice)
	minQty := decimal.NewFromFloat(data.MinNotional).Div(lastPriceDec)
	maxQty := decimal.NewFromFloat(data.MaxNotional).Div(lastPriceDec)

	return structs.ExchangePairData{
		ExchangeID:         consts.ExchangeIDbingx,
		BaseAsset:          symbolParts[0],
		QuoteAsset:         symbolParts[1],
		BasePrecision:      utils.GetFloatPrecision(data.StepSize),
		QuotePrecision:     utils.GetFloatPrecision(data.MinNotional),
		Status:             ConvertPairStatus(data.Status),
		Symbol:             data.Symbol,
		MinQty:             minQty.InexactFloat64(),
		MaxQty:             maxQty.InexactFloat64(),
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
		lastPrice, isExists := tickers[pair.Symbol]
		if !isExists {
			return nil, fmt.Errorf("%q last price not found", pair.Symbol)
		}

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
