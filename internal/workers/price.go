package workers

import "github.com/matrixbotio/exchange-gates-lib/pkg/structs"

// PriceEventCallback - callback to handle price event
type PriceEventCallback func(event PriceEvent)

// PriceWorker - a worker interface based on data from a specific market, such as quotes
type PriceWorker struct {
	ExchangeTag         string
	WsChannels          *structs.WorkerChannels
	HandleEventCallback PriceEventCallback
}

// IPriceWorker - interface for PriceWorker
type IPriceWorker interface {
	SubscribeToPriceEvents(
		pairSymbols []string,
		errorHandler func(err error),
	) (map[string]structs.WorkerChannels, error)

	GetExchangeTag() string

	Stop()
}

// SubscribeToPriceEvents - websocket subscription to change quotes and ask-, bid-qty on the exchange (placeholder)
func (w *PriceWorker) SubscribeToPriceEvents(
	pairSymbols []string,
	errorHandler func(err error),
) (map[string]structs.WorkerChannels, error) {
	// placeholder
	return map[string]structs.WorkerChannels{}, nil
}

// GetExchangeTag - get worker exchange tag from exchange adapter
func (w *PriceWorker) GetExchangeTag() string {
	return w.ExchangeTag
}

// Stop listening ws events
func (w *PriceWorker) Stop() {
	go func() {
		w.WsChannels.WsStop <- struct{}{}
	}()
}

// PriceEvent - data on changes in trade data in the market
type PriceEvent struct {
	ExchangeTag string  `json:"exchangeTag"`
	Symbol      string  `json:"symbol"`
	Ask         float64 `json:"ask"`
	Bid         float64 `json:"bid"`
}
