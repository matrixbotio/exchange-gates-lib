package consts

import "time"

const (
	TestExchangeID                = 0
	PairDefaultExchangeID         = 1
	PairDefaultMinQty             = 0.001
	PairDefaultMaxQty             = 99999.99
	PairDefaultMinPrice           = 0.000001
	PairDefaultQtyStep            = 0.001
	PairDefaultPriceStep          = 0.000001
	PairMinDeposit                = 10
	PairDefaultBaseAsset          = "BTC"
	PairDefaultQuoteAsset         = "BUSD"
	PairDefaultAsset              = PairDefaultBaseAsset + PairDefaultQuoteAsset
	MinDepositFix         float64 = 10 // percent

	CandlesInterval          = "1m"
	ExchangeSetupConnTimeout = 3500 // ms

	ExchangeIDbinanceSpot = 1
	PingRetryAttempts     = 3
	PingRetryWaitTime     = time.Second * 2
)
