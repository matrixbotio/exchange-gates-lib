package matrixgates

const (
	PairDefaultExchangeID = 1
	PairDefaultMinQty     = 0.001
	PairDefaultMaxQty     = 99999.99
	PairDefaultMinPrice   = 0.000001
	PairDefaultQtyStep    = 0.001
	PairDefaultPriceStep  = 0.000001
	PairMinDeposit        = 10
	minDepositFix         = 10 // percent

	candlesInterval          = "1m"
	exchangeSetupConnTimeout = 3500 // ms

	OrderStatusNew             = "NEW"
	OrderStatusPartiallyFilled = "PARTIALLY_FILLED"
	OrderStatusFilled          = "FILLED"
	OrderStatusCancelled       = "CANCELED"
	OrderStatusPendingCancel   = "PENDING_CANCEL"
	OrderStatusRejected        = "REJECTED"
	OrderStatusExpired         = "EXPIRED"
	OrderStatusUnknown         = "UNKNOWN"
)
