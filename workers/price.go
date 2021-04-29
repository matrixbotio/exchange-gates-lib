package workers

import sharederrs "github.com/matrixbotio/shared-errors"

// PriceWorker - a worker interface based on data from a specific market, such as quotes
type PriceWorker struct {
	WsChannels *WorkerChannels
}

// IPriceWorker - interface for PriceWorker
type IPriceWorker interface {
	SubscribeToPriceEvents(
		eventCallback func(event PriceEvent),
		errorHandler func(err *sharederrs.APIError),
	) *sharederrs.APIError
}

// PriceEvent - data on changes in trade data in the market
type PriceEvent struct {
	Symbol string `json:"s"`
	Ask    string `json:"a"`
	Bid    string `json:"b"`

	//BestBidQty   string `json:"B"`
	//BestAskQty   string `json:"A"`
}
