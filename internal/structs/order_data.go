package structs

import "github.com/matrixbotio/exchange-gates-lib/internal/consts"

// OrderData - the result of checking the data of the placed order
type OrderData struct {
	OrderID       int64   `json:"orderID"`
	ClientOrderID string  `json:"clientOrderID"`
	Status        string  `json:"status"`      // used in bot.getOrderData
	AwaitQty      float64 `json:"originalQty"` // initial order qty
	FilledQty     float64 `json:"filledQty"`   // event executed qty
	Price         float64 `json:"price"`
	Symbol        string  `json:"symbol"`
	Type          string  `json:"type"`        // "buy" or "sell"
	CreatedTime   int64   `json:"createdTime"` // unix ms
	UpdatedTime   int64   `json:"updatedTime"` // unix ms
}

func (data OrderData) IsPendingCancel() bool {
	return data.Status == consts.OrderStatusPendingCancel
}
