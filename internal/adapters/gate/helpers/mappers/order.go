package mappers

import (
	"fmt"
	"strings"

	"github.com/gateio/gateapi-go/v6"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	"github.com/shopspring/decimal"
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

func GetOrderFees(
	order gateapi.Order,
	baseTicker string,
	quoteTicker string,
) (structs.OrderFees, error) {
	fees := structs.OrderFees{
		BaseAsset:  decimal.Zero,
		QuoteAsset: decimal.Zero,
	}

	feeValue, err := decimal.NewFromString(order.Fee)
	if err != nil {
		return structs.OrderFees{}, fmt.Errorf("parse: %w", err)
	}

	if order.FeeCurrency == baseTicker {
		fees.BaseAsset = feeValue
	}
	if order.FeeCurrency == quoteTicker {
		fees.QuoteAsset = feeValue
	}
	return fees, nil
}
