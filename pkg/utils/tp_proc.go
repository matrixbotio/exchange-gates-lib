package utils

import (
	"errors"
	"fmt"

	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
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

func (s *CalcTPProcessor) Do() (pkgStructs.BotOrder, error) {
	if err := s.checkParams(); err != nil {
		return pkgStructs.BotOrder{}, fmt.Errorf("check params: %w", err)
	}

	if s.strategy == pkgStructs.BotStrategyShort {
		return s.calcShortTPOrder(), nil
	}
	return s.calcLongOrder(), nil
}

func (s *CalcTPProcessor) calcShortTPOrder() pkgStructs.BotOrder {
	// subtract fees from depo spent in quote asset (from default SELL orders)
	// example: when pair is LTCUSDT, fees summed up for SELL orders in USDT
	depositSpentWithFee := s.depositSpent.
		Sub(s.fees.QuoteAsset).
		Add(s.accQuote)

	coinsQtyDec := decimal.NewFromFloat(s.coinsQty)
	profitDec := decimal.NewFromFloat(s.profit)
	profitDelta := decimal.NewFromFloat(1).Add(profitDec.Div(decimal.NewFromInt(100)))
	// qty = (1 + profit/100) * coinsQty
	tpQty := coinsQtyDec.Mul(profitDelta)

	// price = depositSpent / tpQty
	tpPriceWithoutFee := depositSpentWithFee.Div(tpQty)

	qtyPrecision := GetFloatPrecision(s.pairData.QtyStep)
	tpQtyFloat, _ := tpQty.RoundFloor(int32(qtyPrecision)).Float64()

	pricePrecision := GetFloatPrecision(s.pairData.PriceStep)
	tpPriceFloat, _ := tpPriceWithoutFee.RoundFloor(int32(pricePrecision)).Float64()

	depoSpentFloat, _ := depositSpentWithFee.Float64()

	return pkgStructs.BotOrder{
		PairSymbol:    s.pairData.Symbol,
		Type:          GetTPOrderType(pkgStructs.BotStrategyShort),
		Qty:           tpQtyFloat,
		Price:         tpPriceFloat,
		Deposit:       depoSpentFloat,
		ClientOrderID: GenerateUUID(),
	}
}

func (s *CalcTPProcessor) calcLongOrder() pkgStructs.BotOrder {
	// subtract fees from coins qty in base asset (from default BUY orders)
	// example: when pair is LTCUSDT, fees summed up for BUY orders in LTC
	coinsQtyDec := decimal.NewFromFloat(s.coinsQty).
		Sub(s.fees.BaseAsset).
		Add(s.accBase)

	profitDec := decimal.NewFromFloat(s.profit)
	profitDelta := decimal.NewFromFloat(1).Add(profitDec.Div(decimal.NewFromInt(100)))

	// deposit = (1 + profit/100) * depositSpent
	tpDeposit := profitDelta.Mul(s.depositSpent)

	tpPrice := tpDeposit.Div(coinsQtyDec)

	tpDepositFloat, _ := tpDeposit.Float64()
	pricePrecision := GetFloatPrecision(s.pairData.PriceStep)
	tpPriceFloat, _ := tpPrice.RoundFloor(int32(pricePrecision)).Float64()

	qtyPrecision := GetFloatPrecision(s.pairData.QtyStep)
	tpCoinsQty, _ := coinsQtyDec.RoundFloor(int32(qtyPrecision)).Float64()

	return pkgStructs.BotOrder{
		PairSymbol:    s.pairData.Symbol,
		Type:          GetTPOrderType(pkgStructs.BotStrategyLong),
		Qty:           tpCoinsQty,
		Price:         tpPriceFloat,
		Deposit:       tpDepositFloat,
		ClientOrderID: GenerateUUID(),
	}
}
