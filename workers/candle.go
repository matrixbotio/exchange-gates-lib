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

// CandleEvent - changes in trading candles for a specific pair
type CandleEvent struct {
	//Event  string     `json:"e"`
	Symbol string     `json:"s"`
	Candle CandleData `json:"k"`
}

// CandleData - trading candle
type CandleData struct {
	StartTime int64   `json:"t"`
	EndTime   int64   `json:"T"`
	Interval  string  `json:"i"`
	Open      float64 `json:"o"`
	Close     float64 `json:"c"`
	High      float64 `json:"h"`
	Low       float64 `json:"l"`
	Volume    float64 `json:"v"`
	//IsFinal              bool   `json:"x"`
}
