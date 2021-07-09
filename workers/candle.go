package workers

// CandleWorker - worker for subscribtion to exchange candle events
type CandleWorker struct {
	ExchangeTag string
	WsChannels  *WorkerChannels
}

// ICandleWorker - interface for CandleWorker
type ICandleWorker interface {
	SubscribeToCandleEvents(
		pairSymbols []string,
		eventCallback func(event CandleEvent),
		errorHandler func(err error),
	) error
	GetExchangeTag() string
	Stop()
}

// SubscribeToCandleEvents - websocket subscription to change trade candles on the exchange (placeholder)
func (w *CandleWorker) SubscribeToCandleEvents(
	pairSymbols []string,
	eventCallback func(event CandleEvent),
	errorHandler func(err error),
) error {
	// placeholder
	return nil
}

// GetExchangeTag - get worker exchange tag from exchange adapter
func (w *CandleWorker) GetExchangeTag() string {
	return w.ExchangeTag
}

// Stop listening ws events
func (w *CandleWorker) Stop() {
	go func() {
		<-w.WsChannels.WsStop
		close(w.WsChannels.WsDone)
	}()
}

// CandleEvent - changes in trading candles for a specific pair
type CandleEvent struct {
	Symbol string     `json:"symbol"`
	Candle CandleData `json:"candle"`
}

// CandleData - trading candle
type CandleData struct {
	StartTime int64   `json:"startTime"`
	EndTime   int64   `json:"endTime"`
	Interval  string  `json:"interval"`
	Open      float64 `json:"open"`
	Close     float64 `json:"close"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Volume    float64 `json:"volume"`
}
