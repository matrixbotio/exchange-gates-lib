//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package workers

// TradeEventWorker - a worker interface based on pair trade events
type TradeEventWorker struct {
	workerBase
	ExchangeTag string
}

type TradeEventCallback func(event TradeEvent)

type TradeEventPrivateCallback func(event TradeEventPrivate)

// ITradeEventWorker - interface for PriceWorker
type ITradeEventWorker interface {
	SubscribeToTradeEvents(
		symbol string,
		eventCallback TradeEventCallback,
		errorHandler func(err error),
	) error
	SubscribeToTradeEventsPrivate(
		eventCallback TradeEventPrivateCallback,
		errorHandler func(err error),
	) error
	GetExchangeTag() string
	Stop()
}

// SubscribeToTradeEvents - websocket subscription to pair trade events
func (w *TradeEventWorker) SubscribeToTradeEvents(
	_ string,
	_ TradeEventCallback,
	_ func(err error),
) error {
	// placeholder
	return nil
}

// GetExchangeTag - get worker exchange tag from exchange adapter
func (w *TradeEventWorker) GetExchangeTag() string {
	return w.ExchangeTag
}

// TradeEvent - data on a executed order in a trading pair
type TradeEvent struct {
	ID            int64   `json:"id"`
	Time          int64   `json:"time"`
	Symbol        string  `json:"symbol"`
	Price         float64 `json:"price"`
	Quantity      float64 `json:"quantity"`
	ExchangeTag   string  `json:"exchangeTag"`
	BuyerOrderID  int64   `json:"buyerID"`
	SellerOrderID int64   `json:"sellerID"`
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
