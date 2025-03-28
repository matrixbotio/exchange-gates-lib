//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package workers

// TradeEventWorker - a worker interface based on pair trade events
type TradeEventWorker struct {
	workerBase
	ExchangeTag string
}

type TradeEventPrivateCallback func(event TradeEventPrivate)

// GetExchangeTag - get worker exchange tag from exchange adapter
func (w *TradeEventWorker) GetExchangeTag() string {
	return w.ExchangeTag
}

type TradeEventPrivate struct {
	ID            string  `json:"id,omitempty"`
	Time          int64   `json:"time,omitempty"`
	ExchangeTag   string  `json:"exchangeTag,omitempty"`
	Symbol        string  `json:"symbol,omitempty"`
	OrderID       string  `json:"orderID,omitempty"`
	ClientOrderID string  `json:"clientOrderID,omitempty"`
	Price         float64 `json:"price,omitempty"`
	Quantity      float64 `json:"quantity,omitempty"`
}

type OrderEvent struct {
	// required
	APIKeyID string            `json:"apiKeyID"`
	Data     TradeEventPrivate `json:"data"`

	// optional
	BotID string `json:"botID"`
}
