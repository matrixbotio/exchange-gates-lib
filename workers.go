package matrixgates

import sharederrs "github.com/matrixbotio/shared-errors"

// MarketDataWorker - a worker interface based on data from a specific market, such as quotes
type MarketDataWorker struct {
	WsChannels *WorkerChannels
}

// WorkerChannels - channels to control the worker
type WorkerChannels struct {
	WsBookDone chan struct{}
	WsBookStop chan struct{}
}

// IMarketDataWorker - MarketDataWorker interface
type IMarketDataWorker interface {
	SubscribeToBookEvents(
		eventCallback func(event MDWBookEvent),
		errorHandler func(err *sharederrs.APIError),
	) *sharederrs.APIError
}

// MDWBookEvent - data on changes in trade data in the market
type MDWBookEvent struct {
	UpdateID     int64  `json:"u"`
	Symbol       string `json:"s"`
	BestBidPrice string `json:"b"`
	BestBidQty   string `json:"B"`
	BestAskPrice string `json:"a"`
	BestAskQty   string `json:"A"`
}

// Placeholders
