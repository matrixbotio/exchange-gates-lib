package utils

import (
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
	"github.com/shopspring/decimal"
)

func CalcTPOrder(
	strategy pkgStructs.BotStrategy,
	coinsQty float64,
	profit float64,
	depositSpent float64,
	pairLimits structs.ExchangePairData,
) pkgStructs.BotOrder {
	if strategy == pkgStructs.BotStrategyShort {
		return calcShortTPOrder(coinsQty, profit, depositSpent, pairLimits)
	}
	return calcLongOrder(coinsQty, profit, depositSpent, pairLimits)
}

func GetTPOrderType(strategy pkgStructs.BotStrategy) string {
	if strategy == pkgStructs.BotStrategyLong {
		return pkgStructs.OrderTypeSell
	}
	return pkgStructs.OrderTypeBuy
}

func calcShortTPOrder(
	coinsQty float64,
	profit float64,
	depositSpent float64,
	pairLimits structs.ExchangePairData,
) pkgStructs.BotOrder {
	coinsQtyDec := decimal.NewFromFloat(coinsQty)
	profitDec := decimal.NewFromFloat(profit)
	profitDelta := decimal.NewFromFloat(1).Add(profitDec.Div(decimal.NewFromInt(100)))
	// qty = (1 + profit/100) * coinsQty
	tpQty := coinsQtyDec.Mul(profitDelta)

	// price = depositSpent / tpQty
	tpPrice := decimal.NewFromFloat(depositSpent).Div(tpQty)

	qtyPrecision := GetFloatPrecision(pairLimits.QtyStep)
	tpQtyFloat, _ := tpQty.RoundFloor(int32(qtyPrecision)).Float64()

	pricePrecision := GetFloatPrecision(pairLimits.PriceStep)
	tpPriceFloat, _ := tpPrice.RoundFloor(int32(pricePrecision)).Float64()

	return pkgStructs.BotOrder{
		PairSymbol:    pairLimits.Symbol,
		Type:          GetTPOrderType(pkgStructs.BotStrategyShort),
		Qty:           tpQtyFloat,
		Price:         tpPriceFloat,
		Deposit:       depositSpent,
		ClientOrderID: GenerateUUID(),
	}
}

func calcLongOrder(
	coinsQty float64,
	profit float64,
	depositSpent float64,
	pairLimits structs.ExchangePairData,
) pkgStructs.BotOrder {
	coinsQtyDec := decimal.NewFromFloat(coinsQty)
	depositSpentDec := decimal.NewFromFloat(depositSpent)
	profitDec := decimal.NewFromFloat(profit)
	profitDelta := decimal.NewFromFloat(1).Add(profitDec.Div(decimal.NewFromInt(100)))

	// deposit = (1 + profit/100) * depositSpent
	tpDeposit := profitDelta.Mul(depositSpentDec)

	tpPrice := tpDeposit.Div(coinsQtyDec)

	tpDepositFloat, _ := tpDeposit.Float64()
	pricePrecision := GetFloatPrecision(pairLimits.PriceStep)
	tpPriceFloat, _ := tpPrice.RoundFloor(int32(pricePrecision)).Float64()

	return pkgStructs.BotOrder{
		PairSymbol:    pairLimits.Symbol,
		Type:          GetTPOrderType(pkgStructs.BotStrategyLong),
		Qty:           coinsQty,
		Price:         tpPriceFloat,
		Deposit:       tpDepositFloat,
		ClientOrderID: GenerateUUID(),
	}
}
