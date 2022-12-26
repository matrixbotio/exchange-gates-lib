package mappers

import (
	"fmt"
	"strconv"

	"github.com/hirokisan/bybit/v2"

	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bybit/helpers/conditions"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	"github.com/matrixbotio/exchange-gates-lib/pkg/utils"
)

func ConvertPairsData(pairs *bybit.V5GetInstrumentsInfoSpotResult, exchangeID int) (
	[]structs.ExchangePairData,
	error,
) {
	var result []structs.ExchangePairData

	/* TBD: https://github.com/matrixbotio/exchange-gates-lib/issues/152
	we can put the code inside the loop into a separate
	function when the author of the lib fix the field types
	of the structure: https://github.com/hirokisan/bybit/pull/138
	*/
	for _, rawPairData := range pairs.List {
		pairSymbol := string(rawPairData.Symbol)

		pairData := structs.ExchangePairData{
			ExchangeID:    exchangeID,
			BaseAsset:     string(rawPairData.BaseCoin),
			QuoteAsset:    string(rawPairData.QuoteCoin),
			MinPrice:      consts.PairDefaultMinPrice,
			Symbol:        pairSymbol,
			Status:        consts.PairDefaultStatus,
			AllowedMargin: false,
			AllowedSpot:   conditions.IsSpotTradingAvailable(rawPairData.Status),
			InUse:         true,
		}

		minBaseAmount, err := strconv.ParseFloat(rawPairData.LotSizeFilter.BasePrecision, 64)
		if err != nil {
			return nil, fmt.Errorf("parse %q base precision: %w", pairSymbol, err)
		}

		minQuoteAmount, err := strconv.ParseFloat(rawPairData.LotSizeFilter.QuotePrecision, 64)
		if err != nil {
			return nil, fmt.Errorf("parse %q quote precision: %w", pairSymbol, err)
		}

		pairData.BasePrecision = utils.GetFloatPrecision(minBaseAmount)
		pairData.QuotePrecision = utils.GetFloatPrecision(minQuoteAmount)

		pairData.MinQty, err = strconv.ParseFloat(rawPairData.LotSizeFilter.MinOrderQty, 64)
		if err != nil {
			return nil, fmt.Errorf("parse %q min qty: %w", pairSymbol, err)
		}
		pairData.QtyStep = utils.GetValueStep(pairData.MinQty)

		pairData.MaxQty, err = strconv.ParseFloat(rawPairData.LotSizeFilter.MaxOrderQty, 64)
		if err != nil {
			return nil, fmt.Errorf("parse %q max qty: %w", pairSymbol, err)
		}

		pairData.OriginalMinDeposit, err = strconv.ParseFloat(
			rawPairData.LotSizeFilter.MinOrderAmt,
			64,
		)
		if err != nil {
			return nil, fmt.Errorf("parse %q min deposit: %w", pairSymbol, err)
		}
		pairData.MinDeposit = pairData.OriginalMinDeposit

		pairData.PriceStep, err = strconv.ParseFloat(rawPairData.PriceFilter.TickSize, 64)
		if err != nil {
			return nil, fmt.Errorf("parse price precision: %w", err)
		}
		if pairData.PriceStep == 0 {
			return nil, fmt.Errorf("%q price step is empty", pairSymbol)
		}

		result = append(result, pairData)
	}

	return result, nil
}
