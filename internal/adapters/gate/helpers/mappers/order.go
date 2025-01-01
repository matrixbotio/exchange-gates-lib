package mappers

import (
	"fmt"
	"strings"

	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
)

func ConvertOrderStatus(gateOrderStatus string) consts.OrderStatus {
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

func ConvertOrderSide(gateOrderSide string) (consts.OrderSide, error) {
	switch strings.ToLower(gateOrderSide) {
	default:
		return "", fmt.Errorf("unknown side: %q", gateOrderSide)
	case "buy":
		return consts.OrderSideBuy, nil
	case "sell":
		return consts.OrderSideSell, nil
	}
}
