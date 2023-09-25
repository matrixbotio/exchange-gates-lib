package binance

func (a *adapter) GetOrderExecFee(pairSymbol string, orderID int64) (float64, error) {
	// TBD: https://github.com/matrixbotio/exchange-gates-lib/issues/167
	return 0, nil
}
