package binance

import (
	"fmt"
	"strconv"

	"github.com/adshao/go-binance/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
)

// convert binance.Symbol to ExchangePairData
func getExchangePairData(symbolData binance.Symbol, exchangeID int) (
	structs.ExchangePairData,
	error,
) {
	pairData := structs.ExchangePairData{
		ExchangeID:     exchangeID,
		BaseAsset:      symbolData.BaseAsset,
		BasePrecision:  symbolData.BaseAssetPrecision,
		QuoteAsset:     symbolData.QuoteAsset,
		QuotePrecision: symbolData.QuotePrecision,
		Status:         symbolData.Status,
		Symbol:         symbolData.Symbol,
		MinQty:         consts.PairDefaultMinQty,
		MaxQty:         consts.PairDefaultMaxQty,
		MinDeposit:     consts.PairMinDeposit,
		MinPrice:       consts.PairDefaultMinPrice,
		QtyStep:        consts.PairDefaultQtyStep,
		PriceStep:      consts.PairDefaultPriceStep,
		AllowedMargin:  symbolData.IsMarginTradingAllowed,
		AllowedSpot:    symbolData.IsSpotTradingAllowed,
	}

	var optionalErr error
	if err := binanceParseLotSizeFilter(symbolData, &pairData); err != nil {
		optionalErr = err
	}

	if err := binanceParsePriceFilter(symbolData, &pairData); err != nil {
		optionalErr = err
	}

	if err := binanceParseMinNotionalFilter(symbolData, &pairData); err != nil {
		optionalErr = err
	}

	return pairData, optionalErr
}

func binanceParseMinNotionalFilter(symbolData binance.Symbol, pairData *structs.ExchangePairData) error {
	var err error
	minNotionalFilter := symbolData.MinNotionalFilter()
	if minNotionalFilter == nil {
		return fmt.Errorf("min notional filter not available for pair %q", symbolData.Symbol)
	}

	pairData.OriginalMinDeposit, err = strconv.ParseFloat(minNotionalFilter.MinNotional, 64)
	if err != nil {
		return fmt.Errorf("parse float: %w", err)
	}
	pairData.MinDeposit = pairData.OriginalMinDeposit
	return nil
}

func binanceParsePriceFilter(symbolData binance.Symbol, pairData *structs.ExchangePairData) error {
	var err error
	priceFilter := symbolData.PriceFilter()
	if priceFilter == nil {
		return fmt.Errorf("get price filter for %q", symbolData.Symbol)
	}

	minPriceRaw := priceFilter.MinPrice
	pairData.MinPrice, err = strconv.ParseFloat(minPriceRaw, 64)
	if err != nil {
		return fmt.Errorf("data handle error: %w", err)
	}
	if pairData.MinPrice == 0 {
		pairData.MinPrice = consts.PairDefaultMinPrice
	}

	priceStepRaw := priceFilter.TickSize
	pairData.PriceStep, err = strconv.ParseFloat(priceStepRaw, 64)
	if err != nil {
		return fmt.Errorf("data handle error: %w", err)
	}
	if pairData.PriceStep == 0 {
		pairData.PriceStep = pairData.MinPrice
	}
	return nil
}

func binanceParseLotSizeFilter(symbolData binance.Symbol, pairData *structs.ExchangePairData) error {
	lotSizeFilter := symbolData.LotSizeFilter()
	if lotSizeFilter == nil {
		return fmt.Errorf("lot size filter for symbol %q not found", symbolData.Symbol)
	}
	minQtyRaw := lotSizeFilter.MinQuantity
	maxQtyRaw := lotSizeFilter.MaxQuantity

	var err error
	pairData.MinQty, err = strconv.ParseFloat(minQtyRaw, 64)
	if err != nil {
		return fmt.Errorf("parse pair min qty: %w", err)
	}
	if pairData.MinQty == 0 {
		pairData.MinQty = consts.PairDefaultMinQty
	}

	pairData.MaxQty, err = strconv.ParseFloat(maxQtyRaw, 64)
	if err != nil {
		return fmt.Errorf("parse pair max qty: %w", err)
	}

	qtyStepRaw := lotSizeFilter.StepSize
	pairData.QtyStep, err = strconv.ParseFloat(qtyStepRaw, 64)
	if err != nil {
		return fmt.Errorf("parse pair qty step: %w", err)
	}
	if pairData.QtyStep == 0 {
		pairData.QtyStep = pairData.MinQty
	}
	return nil
}
