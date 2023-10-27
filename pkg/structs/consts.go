package structs

import "github.com/matrixbotio/exchange-gates-lib/internal/consts"

const (
	OrderStatusNew                      = consts.OrderStatusNew
	OrderStatusPartiallyFilled          = consts.OrderStatusPartiallyFilled
	OrderStatusPartiallyFilledCancelled = consts.OrderStatusPartiallyFilledCancelled
	OrderStatusFilled                   = consts.OrderStatusFilled
	OrderStatusCancelled                = consts.OrderStatusCancelled
	OrderStatusPendingCancel            = consts.OrderStatusPendingCancel
	OrderStatusRejected                 = consts.OrderStatusRejected
	OrderStatusExpired                  = consts.OrderStatusExpired
	OrderStatusUnknown                  = consts.OrderStatusUnknown
	OrderStatusUntriggered              = consts.OrderStatusUntriggered
	OrderStatusTriggered                = consts.OrderStatusTriggered
	OrderStatusDeactivated              = consts.OrderStatusDeactivated
)

const (
	OrderTypeBuy  = "buy"
	OrderTypeSell = "sell"
)

const (
	BotStrategyLong  BotStrategy = "long"
	BotStrategyShort BotStrategy = "short"
)
