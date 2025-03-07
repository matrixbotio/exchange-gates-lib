package mappers

import (
	"fmt"

	"github.com/gateio/gateapi-go/v6"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	"github.com/matrixbotio/exchange-gates-lib/pkg/utils"
	"github.com/shopspring/decimal"
)

const pairSymbolFormat = "%s_%s"

func GetPairSymbol(baseTicker, quoteTicker string) string {
	return fmt.Sprintf(pairSymbolFormat, baseTicker, quoteTicker)
}

func ConvertPairs(
	data []gateapi.CurrencyPair,
) ([]structs.ExchangePairData, error) {
	r := []structs.ExchangePairData{}

	for _, pairData := range data {
		pairParsed, err := ConvertPair(pairData)
		if err != nil {
			return nil, fmt.Errorf("parse: %w", err)
		}

		r = append(r, pairParsed)
	}

	return r, nil
}

func ConvertPair(
	data gateapi.CurrencyPair,
) (structs.ExchangePairData, error) {
	r := structs.ExchangePairData{
		ExchangeID:     consts.ExchangeIDgateSpot,
		BaseAsset:      data.Base,
		QuoteAsset:     data.Quote,
		BasePrecision:  0, // TODO
		QuotePrecision: int(data.AmountPrecision),
		Status:         parsePairStatus(data.TradeStatus),
		Symbol:         GetPairSymbol(data.Base, data.Quote),
		AllowedSpot:    true,
		AllowedMargin:  false,
	}

	minQty, err := decimal.NewFromString(data.MinBaseAmount)
	if err != nil {
		return structs.ExchangePairData{},
			fmt.Errorf("parse min qty: %w", err)
	}
	r.MinQty = minQty.InexactFloat64()

	maxQty, err := decimal.NewFromString(data.MaxBaseAmount)
	if err != nil {
		return structs.ExchangePairData{},
			fmt.Errorf("parse max qty: %w", err)
	}
	r.MaxQty = maxQty.InexactFloat64()

	minAmount, err := decimal.NewFromString(data.MinQuoteAmount)
	if err != nil {
		return structs.ExchangePairData{},
			fmt.Errorf("parse min amount: %w", err)
	}
	r.MinDeposit = minAmount.InexactFloat64()
	r.OriginalMinDeposit = r.MinDeposit

	minQtyPrecision := utils.GetFloatPrecision(minQty.InexactFloat64())
	r.QtyStep = getValueStep(int32(minQtyPrecision))

	r.PriceStep = getValueStep(data.Precision)
	r.MinPrice = r.PriceStep

	return r, nil
}

func getValueStep(precision int32) float64 {
	// 1 / 10^precision
	return decimal.NewFromInt(1).
		Sub(decimal.NewFromInt(10).
			Pow(decimal.NewFromInt32(precision)),
		).InexactFloat64()
}

func parsePairStatus(gateStatus string) string {
	switch gateStatus {
	default:
		return consts.PairStatusUnknown
	case "untradable":
		return consts.PairStatusOffline
	case "buyable", "sellable":
		return consts.PairStatusSuspended
	case "tradable":
		return consts.PairStatusTrading
	}
}
