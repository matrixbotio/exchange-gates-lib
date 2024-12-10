package utils

import (
	"errors"
	"fmt"

	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	pkgErrors "github.com/matrixbotio/exchange-gates-lib/pkg/errs"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
	"github.com/shopspring/decimal"
)

type CalcTPProcessor struct {
	strategy     pkgStructs.BotStrategy
	coinsQty     float64
	profit       float64
	depositSpent decimal.Decimal
	fees         structs.OrderFees
	pairData     structs.ExchangePairData

	accBase  decimal.Decimal
	accQuote decimal.Decimal
}

func NewCalcTPOrderProcessor() *CalcTPProcessor {
	return &CalcTPProcessor{}
}

func (s *CalcTPProcessor) Remains(
	accBase decimal.Decimal,
	accQuote decimal.Decimal,
) *CalcTPProcessor {
	s.accBase = accBase
	s.accQuote = accQuote

	if s.accBase.IsNegative() {
		s.accBase = decimal.Zero
	}
	if s.accQuote.IsNegative() {
		s.accQuote = decimal.Zero
	}
	return s
}

func (s *CalcTPProcessor) Strategy(strategy pkgStructs.BotStrategy) *CalcTPProcessor {
	s.strategy = strategy
	return s
}

func (s *CalcTPProcessor) CoinsQty(coinsQty float64) *CalcTPProcessor {
	s.coinsQty = coinsQty
	return s
}

func (s *CalcTPProcessor) Profit(profit float64) *CalcTPProcessor {
	s.profit = profit
	return s
}

func (s *CalcTPProcessor) DepositSpent(depositSpent decimal.Decimal) *CalcTPProcessor {
	s.depositSpent = depositSpent
	return s
}

func (s *CalcTPProcessor) PairData(pairData structs.ExchangePairData) *CalcTPProcessor {
	s.pairData = pairData
	return s
}

func (s *CalcTPProcessor) Fees(fees structs.OrderFees) *CalcTPProcessor {
	s.fees = fees
	return s
}

func (s *CalcTPProcessor) checkParams() error {
	if s.strategy == "" {
		return errors.New("strategy is not set")
	}
	if s.coinsQty == 0 {
		return errors.New("invalid coins qty (0)")
	}
	if s.profit == 0 {
		return errors.New("invalid profit value (0)")
	}
	if s.depositSpent.IsZero() {
		return errors.New("invalid depositSpent value (0)")
	}
	if s.pairData.IsEmpty() {
		return errors.New("pair data is not set")
	}
	if s.pairData.Symbol == "" {
		return errors.New("pair symbol is not set in pair data")
	}
	if s.pairData.QtyStep == 0 {
		return errors.New("invalid qty step value (0)")
	}
	if s.pairData.PriceStep == 0 {
		return errors.New("invalid price step value (0)")
	}
	return nil
}

type CalcTPResult struct {
	TPOrder pkgStructs.BotOrder

	// remains
	AccBaseOriginal  decimal.Decimal
	AccQuoteOriginal decimal.Decimal

	// remains used in TP
	AccBaseUsed  decimal.Decimal
	AccQuoteUsed decimal.Decimal
}

func (s *CalcTPProcessor) Do() (CalcTPResult, error) {
	if err := s.checkParams(); err != nil {
		return CalcTPResult{}, fmt.Errorf("check params: %w", err)
	}

	if s.strategy == pkgStructs.BotStrategyShort {
		return s.calcShortTPOrder()
	}
	return s.calcLongOrder()
}

func (s *CalcTPProcessor) getMinQtyError(qty decimal.Decimal) error {
	return fmt.Errorf(
		"%w: not enough coins (%s %s with a minimum of %v %s)",
		pkgErrors.ErrMinimumTP,
		qty.String(),
		s.pairData.BaseAsset,
		s.pairData.MinQty,
		s.pairData.BaseAsset,
	)
}

func (s *CalcTPProcessor) getMaxQtyError(qty decimal.Decimal) error {
	return fmt.Errorf(
		"too many coins (%s %s with a max of %v %s)",
		qty.String(),
		s.pairData.BaseAsset,
		s.pairData.MaxQty,
		s.pairData.BaseAsset,
	)
}

func (s *CalcTPProcessor) getMinAmountError(amount decimal.Decimal) error {
	return fmt.Errorf(
		"%w: not enough amount (%s %s with a minimum of %v %s)",
		pkgErrors.ErrMinimumTP,
		amount.String(),
		s.pairData.QuoteAsset,
		s.pairData.MinDeposit,
		s.pairData.QuoteAsset,
	)
}

func (s *CalcTPProcessor) getMinPriceError(price decimal.Decimal) error {
	return fmt.Errorf(
		"too low price (%s with a minimum of %v)",
		price.String(),
		s.pairData.MinPrice,
	)
}

func (s *CalcTPProcessor) calcShortTPQty(coinsQtyDec, amountAvailable decimal.Decimal) (
	decimal.Decimal,
	error,
) {
	// increase qty by profit %
	profitDec := decimal.NewFromFloat(s.profit)
	profitDelta := decimal.NewFromFloat(1).Add(profitDec.Div(decimal.NewFromInt(100)))

	// qty = coinsQty * (1 + profit/100)
	tpQty := coinsQtyDec.Mul(profitDelta)

	// check max qty
	if s.pairData.MaxQty > 0 &&
		tpQty.GreaterThan(decimal.NewFromFloat(s.pairData.MaxQty)) {
		return decimal.Zero, s.getMaxQtyError(tpQty)
	}

	if !s.accQuote.IsPositive() {
		// remains not set
		return tpQty, nil
	}

	// Let's try to calculate how much remains amount we can
	// convert to qty to add to the order
	zeroProfitPrice := amountAvailable.Div(tpQty)
	remainsQty := s.accQuote.Div(zeroProfitPrice)
	qtyWithRemains := tpQty.Add(remainsQty)

	// round qty with remains
	return qtyWithRemains, nil
}

func (s *CalcTPProcessor) roundQtyDown(qty decimal.Decimal) decimal.Decimal {
	qtyPrecision := GetFloatPrecision(s.pairData.QtyStep)
	return qty.RoundFloor(int32(qtyPrecision))
}

func (s *CalcTPProcessor) roundPrice(price decimal.Decimal) decimal.Decimal {
	pricePrecision := GetFloatPrecision(s.pairData.PriceStep)
	return price.RoundFloor(int32(pricePrecision))
}

func (s *CalcTPProcessor) roundAmount(amount decimal.Decimal) decimal.Decimal {
	return RoundAmount(
		amount,
		string(s.strategy),
		s.pairData.BasePrecision,
		s.pairData.QuotePrecision,
	)
}

func (s *CalcTPProcessor) calcShortTPOrder() (CalcTPResult, error) {
	// coins qty - fees
	coinsQtyDec := decimal.NewFromFloat(s.coinsQty).
		Sub(s.fees.BaseAsset)

	amountAvailable := s.depositSpent.Sub(s.fees.QuoteAsset)
	tpQtyWithRemains, err := s.calcShortTPQty(coinsQtyDec, amountAvailable)
	if err != nil {
		return CalcTPResult{}, err
	}

	tpQtyRounded := s.roundQtyDown(tpQtyWithRemains)

	// check min qty
	if tpQtyRounded.LessThan(decimal.NewFromFloat(s.pairData.MinQty)) {
		return CalcTPResult{}, s.getMinQtyError(tpQtyRounded)
	}

	// price = depositSpent / tpQty
	tpPrice := amountAvailable.Div(tpQtyWithRemains)

	// check price
	if tpPrice.LessThan(decimal.NewFromFloat(s.pairData.MinPrice)) {
		return CalcTPResult{}, s.getMinPriceError(tpPrice)
	}

	tpPrice = s.roundPrice(tpPrice)

	// calc remains used
	accQuoteUsed := tpQtyRounded.Sub(tpQtyWithRemains).Mul(tpPrice).Add(s.accQuote)

	// recalc amount
	tpAmount := s.roundAmount(tpQtyRounded.Mul(tpPrice))

	// check min amount
	if tpAmount.LessThan(decimal.NewFromFloat(s.pairData.MinDeposit)) {
		return CalcTPResult{}, s.getMinAmountError(tpAmount)
	}

	order := pkgStructs.BotOrder{
		PairSymbol:    s.pairData.Symbol,
		Type:          GetTPOrderType(pkgStructs.BotStrategyShort),
		Qty:           tpQtyRounded.InexactFloat64(),
		Price:         tpPrice.InexactFloat64(),
		Deposit:       tpAmount.InexactFloat64(),
		ClientOrderID: GenerateUUID(),
	}

	// let's check that the TP order will not close in the minus
	zeroProfitPrice := s.depositSpent.Div(coinsQtyDec)
	if tpPrice.GreaterThan(zeroProfitPrice) {
		return CalcTPResult{},
			fmt.Errorf(
				"invalid TP calc: order: %s",
				order.String(),
			)
	}

	result := CalcTPResult{
		TPOrder:         order,
		AccBaseOriginal: decimal.Zero,
		AccBaseUsed:     decimal.Zero,
	}
	if s.accQuote.IsPositive() {
		result.AccQuoteOriginal = s.accQuote
		result.AccQuoteUsed = accQuoteUsed
	}

	return result, nil
}

func (s *CalcTPProcessor) calcLongOrder() (CalcTPResult, error) {
	// subtract fees from coins qty in base asset (from default BUY orders)
	// example: when pair is LTCUSDT, fees summed up for BUY orders in LTC
	coinsQtyDec := decimal.NewFromFloat(s.coinsQty).
		Sub(s.fees.BaseAsset)

	profitDec := decimal.NewFromFloat(s.profit)
	profitDelta := decimal.NewFromFloat(1).Add(profitDec.Div(decimal.NewFromInt(100)))

	// deposit = (1 + profit/100) * depositSpent
	tpAmount := profitDelta.Mul(s.depositSpent)
	tpPrice := tpAmount.Div(coinsQtyDec)

	tpQtyWithRemains := coinsQtyDec.Copy()
	if s.accBase.IsPositive() {
		tpQtyWithRemains = tpQtyWithRemains.Add(s.accBase)
	}

	tpQtyRounded := s.roundQtyDown(tpQtyWithRemains)

	// check min qty
	if tpQtyRounded.LessThan(decimal.NewFromFloat(s.pairData.MinQty)) {
		return CalcTPResult{}, s.getMinQtyError(tpQtyRounded)
	}

	tpPrice = s.roundPrice(tpPrice)

	// recalc amount
	tpAmount = s.roundAmount(tpQtyRounded.Mul(tpPrice))

	// check amount
	if tpAmount.LessThan(decimal.NewFromFloat(s.pairData.MinDeposit)) {
		return CalcTPResult{}, s.getMinAmountError(tpAmount)
	}

	// check price
	if tpPrice.LessThan(decimal.NewFromFloat(s.pairData.MinPrice)) {
		return CalcTPResult{}, s.getMinPriceError(tpPrice)
	}

	// calc remains used
	accBaseUsed := tpQtyRounded.Sub(tpQtyWithRemains).Add(s.accBase)

	order := pkgStructs.BotOrder{
		PairSymbol:    s.pairData.Symbol,
		Type:          GetTPOrderType(pkgStructs.BotStrategyLong),
		Qty:           tpQtyRounded.InexactFloat64(),
		Price:         tpPrice.InexactFloat64(),
		Deposit:       tpAmount.InexactFloat64(),
		ClientOrderID: GenerateUUID(),
	}

	// let's check that the TP order will not close in the minus
	zeroProfitPrice := s.depositSpent.Div(coinsQtyDec)
	if tpPrice.LessThan(zeroProfitPrice) {
		return CalcTPResult{},
			fmt.Errorf(
				"invalid TP calc: order: %s",
				order.String(),
			)
	}

	result := CalcTPResult{
		TPOrder:          order,
		AccQuoteOriginal: decimal.Zero,
		AccQuoteUsed:     decimal.Zero,
	}
	if s.accBase.IsPositive() {
		result.AccBaseOriginal = s.accBase
		result.AccBaseUsed = accBaseUsed
	}

	return result, nil
}
