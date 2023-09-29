package binance

import "github.com/shopspring/decimal"

func (a *adapter) GetOrderExecFee(pairSymbol string, orderID int64) (decimal.Decimal, error) {
	// TBD: https://github.com/matrixbotio/exchange-gates-lib/issues/167
	return decimal.NewFromInt(0), nil
}
