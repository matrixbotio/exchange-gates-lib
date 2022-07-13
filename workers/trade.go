package workers

// TradeEventWorker - a worker interface based on pair trade events
type TradeEventWorker struct {
	ExchangeTag string
	WsChannels  *WorkerChannels
}

// ITradeEventWorker - interface for PriceWorker
type ITradeEventWorker interface {
	SubscribeToTradeEvents(
		symbol string,
		eventCallback func(event TradeEvent),
		errorHandler func(err error),
	) error
	GetExchangeTag() string
	Stop()
}

// SubscribeToTradeEvents - websocket subscription to pair trade events
func (w *TradeEventWorker) SubscribeToTradeEvents(
	eventCallback func(event TradeEvent),
	errorHandler func(err error),
) error {
	// placeholder
	return nil
}

// GetExchangeTag - get worker exchange tag from exchange adapter
func (w *TradeEventWorker) GetExchangeTag() string {
	return w.ExchangeTag
}

// Stop listening ws events
func (w *TradeEventWorker) Stop() {
	go func() {
		w.WsChannels.WsStop <- struct{}{}
	}()
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
