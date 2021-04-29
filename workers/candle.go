package workers

import sharederrs "github.com/matrixbotio/shared-errors"

// CandleWorker - worker for subscribtion to exchange candle events
type CandleWorker struct {
	WsChannels *WorkerChannels
}

// ICandleWorker - interface for CandleWorker
type ICandleWorker interface {
	SubscribeToCandleEvents(
		pairs []string,
		eventCallback func(event CandleEvent),
		errorHandler func(err *sharederrs.APIError),
	) *sharederrs.APIError
}

// CandleEvent - changes in trading candles for a specific pair
type CandleEvent struct {
	Event  string `json:"e"`
	Time   int64  `json:"E"`
	Symbol string `json:"s"`
	Kline  Candle `json:"k"`
}

// Candle - trading candle
type Candle struct {
	StartTime            int64  `json:"t"`
	EndTime              int64  `json:"T"`
	Symbol               string `json:"s"`
	Interval             string `json:"i"`
	FirstTradeID         int64  `json:"f"`
	LastTradeID          int64  `json:"L"`
	Open                 string `json:"o"`
	Close                string `json:"c"`
	High                 string `json:"h"`
	Low                  string `json:"l"`
	Volume               string `json:"v"`
	TradeNum             int64  `json:"n"`
	IsFinal              bool   `json:"x"`
	QuoteVolume          string `json:"q"`
	ActiveBuyVolume      string `json:"V"`
	ActiveBuyQuoteVolume string `json:"Q"`
}
