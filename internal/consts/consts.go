package consts

import "time"

const (
	TestExchangeID        = 0
	PairDefaultExchangeID = 1
	PairDefaultMinQty     = 0.001
	PairDefaultMaxQty     = 99999.99
	PairDefaultMinPrice   = 0.000001
	PairDefaultQtyStep    = 0.001
	PairDefaultPriceStep  = 0.000001
	PairMinDeposit        = 10
	PairDefaultBaseAsset  = "BTC"
	PairDefaultQuoteAsset = "BUSD"
	PairDefaultStatus     = "TRADING"
	PairDefaultAsset      = PairDefaultBaseAsset + PairDefaultQuoteAsset
)

const (
	OrderStatusNew             = "NEW"
	OrderStatusPartiallyFilled = "PARTIALLY_FILLED"
	OrderStatusFilled          = "FILLED"
	OrderStatusCancelled       = "CANCELED"
	OrderStatusPendingCancel   = "PENDING_CANCEL"
	OrderStatusRejected        = "REJECTED"
	OrderStatusExpired         = "EXPIRED"
	OrderStatusUnknown         = "UNKNOWN"
	OrderStatusUntriggered     = "UNTRIGGERED"
	OrderStatusTriggered       = "TRIGGERED"
	OrderStatusDeactivated     = "DEACTIVATED"
)

const (
	CandlesInterval          = "1m"
	ExchangeSetupConnTimeout = 3500 // ms
	ReadTimeout              = time.Second * 5
)

const (
	PingRetryAttempts = 3
	PingRetryWaitTime = time.Second * 2
)

const (
	ExchangeIDbinanceSpot = 1
	ExchangeIDbybitSpot   = 2
)
