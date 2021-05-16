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
}

// PriceEvent - data on changes in trade data in the market
type PriceEvent struct {
	Symbol string  `json:"symbol"`
	Ask    float64 `json:"ask"`
	Bid    float64 `json:"bid"`
}
