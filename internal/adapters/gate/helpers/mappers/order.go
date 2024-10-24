package mappers

import "github.com/matrixbotio/exchange-gates-lib/internal/consts"

func ConvertOrderStatus(gateOrderStatus string) string {
	switch gateOrderStatus {
	default:
		return consts.OrderStatusUnknown
	case "open":
		return consts.OrderStatusNew
	case "closed":
		return consts.OrderStatusFilled
	case "cancelled":
		return consts.OrderStatusCancelled
	}
}
