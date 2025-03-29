package mappers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gateio/gateapi-go/v6"
	gate "github.com/gateio/gatews/go"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
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

func ParseOrderEvent(event gate.SpotUserTradesMsg) (
	workers.TradeEventPrivate,
	error,
) {
	timestamp, err := strconv.ParseInt(event.CreateTimeMs, 10, 64)
	if err != nil {
		fmt.Printf("timestamp: %s\n", err.Error())
		timestamp = time.Now().UnixMilli()
	}

	price, err := decimal.NewFromString(event.Price)
	if err != nil {
		return workers.TradeEventPrivate{}, fmt.Errorf("price: %w", err)
	}

	qty, err := decimal.NewFromString(event.Amount)
	if err != nil {
		return workers.TradeEventPrivate{}, fmt.Errorf("qty: %w", err)
	}

	return workers.TradeEventPrivate{
		Time:          timestamp,
		ExchangeTag:   consts.GateAdapterTag,
		Symbol:        event.CurrencyPair,
		OrderID:       event.OrderId,
		ClientOrderID: event.Text,
		Price:         price.InexactFloat64(),
		Quantity:      qty.InexactFloat64(),
	}, nil
}
