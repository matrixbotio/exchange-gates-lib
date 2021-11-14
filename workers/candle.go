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
		w.WsChannels.WsStop <- struct{}{}
	}()
}

// CandleEvent - changes in trading candles for a specific pair
type CandleEvent struct {
	Symbol     string     `json:"symbol"`
	BaseAsset  string     `json:"baseAsset"`
	QuoteAsset string     `json:"quoteAsset"`
	Candle     CandleData `json:"candle"`
	Time       int64      `json:"time"`
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
