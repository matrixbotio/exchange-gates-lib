package matrixgates

import sharederrs "github.com/matrixbotio/shared-errors"

// PriceWorker - a worker interface based on data from a specific market, such as quotes
type PriceWorker struct {
	WsChannels *WorkerChannels
}

// WorkerChannels - channels container to control the worker
type WorkerChannels struct {
	WsDone chan struct{}
	WsStop chan struct{}
}

// IPriceWorker - MarketDataWorker interface
type IPriceWorker interface {
	SubscribeToPriceEvents(
		eventCallback func(event PriceEvent),
		errorHandler func(err *sharederrs.APIError),
	) *sharederrs.APIError
}

// PriceEvent - data on changes in trade data in the market
type PriceEvent struct {
	UpdateID     int64  `json:"u"`
	Symbol       string `json:"s"`
	BestBidPrice string `json:"b"`
	BestBidQty   string `json:"B"`
	BestAskPrice string `json:"a"`
	BestAskQty   string `json:"A"`
}
