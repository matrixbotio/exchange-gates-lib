package workers

import sharederrs "github.com/matrixbotio/shared-errors"

// CandleWorker - worker for subscribtion to exchange candle events
type CandleWorker struct {
	WsChannels *WorkerChannels
}

// ICandleWorker - interface for CandleWorker
type ICandleWorker interface {
	SubscribeToCandleEvents(
		pairSymbols []string,
		eventCallback func(event CandleEvent),
		errorHandler func(err *sharederrs.APIError),
	) *sharederrs.APIError
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
