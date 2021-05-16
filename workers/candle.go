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
}

// CandleEvent - changes in trading candles for a specific pair
type CandleEvent struct {
	//Event  string     `json:"e"`
	Symbol string     `json:"symbol"`
	Candle CandleData `json:"candle"`
}

// CandleData - trading candle
type CandleData struct {
	StartTime int64   `json:"start"`
	EndTime   int64   `json:"end"`
	Interval  string  `json:"interval"`
	Open      float64 `json:"open"`
	Close     float64 `json:"close"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Volume    float64 `json:"volume"`
	//IsFinal              bool   `json:"x"`
}
