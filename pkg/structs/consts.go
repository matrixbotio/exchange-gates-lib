package structs

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
	OrderTypeBuy  = "buy"
	OrderTypeSell = "sell"
)

const (
	BotStrategyLong  BotStrategy = "long"
	BotStrategyShort BotStrategy = "short"
)
