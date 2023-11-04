package utils

import (
	"errors"

	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
	"github.com/shopspring/decimal"
)

type CalcTPService struct {
	strategy     pkgStructs.BotStrategy
	coinsQty     float64
	profit       float64
	depositSpent float64
	fees         structs.OrderFees
	pairData     structs.ExchangePairData
}

func NewCalcTPOrderService() *CalcTPService {
	return &CalcTPService{}
}

func (s *CalcTPService) Strategy(strategy pkgStructs.BotStrategy) *CalcTPService {
	s.strategy = strategy
	return s
}

func (s *CalcTPService) CoinsQty(coinsQty float64) *CalcTPService {
	s.coinsQty = coinsQty
	return s
}

func (s *CalcTPService) Profit(profit float64) *CalcTPService {
	s.profit = profit
	return s
}

func (s *CalcTPService) DepositSpent(depositSpent float64) *CalcTPService {
	s.depositSpent = depositSpent
	return s
}

func (s *CalcTPService) PairData(pairData structs.ExchangePairData) *CalcTPService {
	s.pairData = pairData
	return s
}

func (s *CalcTPService) Fees(fees structs.OrderFees) *CalcTPService {
	s.fees = fees
	return s
}

func (s *CalcTPService) checkParams() error {
	if s.strategy == "" {
		return errors.New("strategy is not set")
	}
	if s.coinsQty == 0 {
		return errors.New("invalid coins qty (0)")
	}
	if s.profit == 0 {
		return errors.New("invalid profit value (0)")
	}
	if s.depositSpent == 0 {
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

func (s *CalcTPService) checkFees() {
	if !s.fees.BaseAsset.IsZero() && !s.fees.QuoteAsset.IsZero() {
		return
	}

	s.fees = structs.OrderFees{
		BaseAsset:  decimal.NewFromInt(0),
		QuoteAsset: decimal.NewFromInt(0),
	}
}

func (s *CalcTPService) Do() (pkgStructs.BotOrder, error) {
	if err := s.checkParams(); err != nil {
		return pkgStructs.BotOrder{}, nil
	}

	s.checkFees()

	if s.strategy == pkgStructs.BotStrategyShort {
		return s.calcShortTPOrder(), nil
	}
	return s.calcLongOrder(), nil
}

func (s *CalcTPService) calcShortTPOrder() pkgStructs.BotOrder {
	// subtract fees from depo spent in quote asset (from default SELL orders)
	// example: when pair is LTCUSDT, fees summed up for SELL orders in USDT
	depositSpentDec := decimal.NewFromFloat(s.depositSpent).Sub(s.fees.QuoteAsset)

	coinsQtyDec := decimal.NewFromFloat(s.coinsQty)
	profitDec := decimal.NewFromFloat(s.profit)
	profitDelta := decimal.NewFromFloat(1).Add(profitDec.Div(decimal.NewFromInt(100)))
	// qty = (1 + profit/100) * coinsQty
	tpQty := coinsQtyDec.Mul(profitDelta)

	// price = depositSpent / tpQty
	tpPrice := depositSpentDec.Div(tpQty)

	qtyPrecision := GetFloatPrecision(s.pairData.QtyStep)
	tpQtyFloat, _ := tpQty.RoundFloor(int32(qtyPrecision)).Float64()

	pricePrecision := GetFloatPrecision(s.pairData.PriceStep)
	tpPriceFloat, _ := tpPrice.RoundFloor(int32(pricePrecision)).Float64()

	return pkgStructs.BotOrder{
		PairSymbol:    s.pairData.Symbol,
		Type:          GetTPOrderType(pkgStructs.BotStrategyShort),
		Qty:           tpQtyFloat,
		Price:         tpPriceFloat,
		Deposit:       s.depositSpent,
		ClientOrderID: GenerateUUID(),
	}
}

func (s *CalcTPService) calcLongOrder() pkgStructs.BotOrder {
	// subtract fees from coins qty in base asset (from default BUY orders)
	// example: when pair is LTCUSDT, fees summed up for BUY orders in LTC
	coinsQtyDec := decimal.NewFromFloat(s.coinsQty).Sub(s.fees.BaseAsset)

	depositSpentDec := decimal.NewFromFloat(s.depositSpent)
	profitDec := decimal.NewFromFloat(s.profit)
	profitDelta := decimal.NewFromFloat(1).Add(profitDec.Div(decimal.NewFromInt(100)))

	// deposit = (1 + profit/100) * depositSpent
	tpDeposit := profitDelta.Mul(depositSpentDec)

	tpPrice := tpDeposit.Div(coinsQtyDec)

	tpDepositFloat, _ := tpDeposit.Float64()
	pricePrecision := GetFloatPrecision(s.pairData.PriceStep)
	tpPriceFloat, _ := tpPrice.RoundFloor(int32(pricePrecision)).Float64()

	return pkgStructs.BotOrder{
		PairSymbol:    s.pairData.Symbol,
		Type:          GetTPOrderType(pkgStructs.BotStrategyLong),
		Qty:           s.coinsQty,
		Price:         tpPriceFloat,
		Deposit:       tpDepositFloat,
		ClientOrderID: GenerateUUID(),
	}
}
