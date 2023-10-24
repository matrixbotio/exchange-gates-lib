package binance

import (
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	"github.com/shopspring/decimal"
)

func (a *adapter) GetOrderExecFee(
	pairSymbol string,
	orderSide string,
	orderID int64,
) (structs.OrderFees, error) {
	// TBD: https://github.com/matrixbotio/exchange-gates-lib/issues/167

	return structs.OrderFees{
		BaseAsset:  decimal.NewFromInt(0),
		QuoteAsset: decimal.NewFromInt(0),
	}, nil
}
