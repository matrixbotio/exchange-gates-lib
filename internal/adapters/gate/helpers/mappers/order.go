package mappers

import (
	"fmt"
	"math"
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

func ParseTimestamp(rawTs string) int64 {
	timestampRaw, err := strconv.ParseFloat(rawTs, 64)
	if err != nil {
		fmt.Printf("timestamp: %s\n", err.Error())
		return time.Now().UnixMilli()
	} else {
		return int64(math.Floor(timestampRaw))
	}
}

func ParseOrderEvent(event gate.SpotUserTradesMsg) (
	workers.TradeEventPrivate,
	error,
) {
	price, err := decimal.NewFromString(event.Price)
	if err != nil {
		return workers.TradeEventPrivate{}, fmt.Errorf("price: %w", err)
	}

	qty, err := decimal.NewFromString(event.Amount)
	if err != nil {
		return workers.TradeEventPrivate{}, fmt.Errorf("qty: %w", err)
	}

	return workers.TradeEventPrivate{
		Time:          ParseTimestamp(event.CreateTimeMs),
		ExchangeTag:   consts.GateAdapterTag,
		Symbol:        event.CurrencyPair,
		OrderID:       event.OrderId,
		ClientOrderID: event.Text,
		Price:         price.InexactFloat64(),
		Quantity:      qty.InexactFloat64(),
	}, nil
}

/*
NOT IMPLEMENTED YET:
Status        consts.OrderStatus `json:"status"`      // used in bot.getOrderData
AwaitQty      float64            `json:"originalQty"` // initial order qty
FilledQty     float64            `json:"filledQty"`   // event executed qty
Price         float64            `json:"price"`
Side          consts.OrderSide   `json:"type"`        // "buy" or "sell"
*/

func ConvertTradesToOrderHistory(
	baseTicker string,
	quoteTicker string,
	events []gateapi.Trade,
	pairSymbol string,
	orderID int64,
) (structs.OrderHistory, error) {
	r := structs.OrderHistory{
		OrderData: structs.OrderData{
			OrderID: orderID,
			Symbol:  pairSymbol,
		},
		Fees: structs.OrderFees{
			BaseAsset:  decimal.Zero,
			QuoteAsset: decimal.Zero,
		},
	}

	for i, event := range events {

		// sum order fees
		feeValue, err := decimal.NewFromString(event.Fee)
		if err != nil {
			return structs.OrderHistory{}, fmt.Errorf("parse fee: %w", err)
		}
		if event.FeeCurrency == baseTicker {
			r.Fees.BaseAsset = r.Fees.BaseAsset.Add(feeValue)
		}
		if event.FeeCurrency == quoteTicker {
			r.Fees.QuoteAsset = r.Fees.QuoteAsset.Add(feeValue)
		}

		// use some info from last event
		if i == len(events)-1 {
			r.CreatedTime = ParseTimestamp(event.CreateTimeMs)
			r.UpdatedTime = r.CreatedTime // temporary solution
			r.ClientOrderID = event.Text
		}
	}

	return r, nil
}

/*
// Order side
Side string `json:"side,omitempty"`
// Trade role. No value in public endpoints
Role string `json:"role,omitempty"`
// Trade amount
Amount string `json:"amount,omitempty"`
// Order price
Price string `json:"price,omitempty"`
*/
