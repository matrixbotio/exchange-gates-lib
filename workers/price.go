package workers

// PriceWorker - a worker interface based on data from a specific market, such as quotes
type PriceWorker struct {
	ExchangeTag string
	WsChannels  *WorkerChannels
}

// IPriceWorker - interface for PriceWorker
type IPriceWorker interface {
	SubscribeToPriceEvents(
		eventCallback func(event PriceEvent),
		errorHandler func(err error),
	) error
	GetExchangeTag() string
}

// SubscribeToPriceEvents - websocket subscription to change quotes and ask-, bid-qty on the exchange (placeholder)
func (w *PriceWorker) SubscribeToPriceEvents(
	eventCallback func(event PriceEvent),
	errorHandler func(err error),
) error {
	// placeholder
	return nil
}

// GetExchangeTag - get worker exchange tag from exchange adapter
func (w *PriceWorker) GetExchangeTag() string {
	return w.ExchangeTag
}

// PriceEvent - data on changes in trade data in the market
type PriceEvent struct {
	Symbol string  `json:"symbol"`
	Ask    float64 `json:"ask"`
	Bid    float64 `json:"bid"`
}
