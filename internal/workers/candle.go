//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package workers

// CandleEvent - changes in trading candles for a specific pair
type CandleEvent struct {
	Symbol     string     `json:"symbol"`
	BaseAsset  string     `json:"baseAsset"`
	QuoteAsset string     `json:"quoteAsset"`
	Candle     CandleData `json:"candle"`
	Time       int64      `json:"time"`

	IsFinished bool
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

// ICandleWorker - interface for CandleWorker
type ICandleWorker interface {
	/*
		SubscribeToCandle - websocket subscription to change trade candles
		on the exchange per one pair
	*/
	SubscribeToCandle(
		pairSymbol string,
		eventCallback func(event CandleEvent),
		errorHandler func(err error),
	) error

	/*
		SubscribeToCandlesList - websocket subscription to change trade candles
		on the exchange per specific pairs
	*/
	SubscribeToCandlesList(
		intervalsPerPair map[string]string,
		eventCallback func(event CandleEvent),
		errorHandler func(err error),
	) error

	// GetExchangeTag - get worker exchange tag from exchange adapter
	GetExchangeTag() string

	// Stop listening ws events
	Stop()
}

// CandleWorker - worker for subscribtion to exchange candle events
type CandleWorker struct {
	workerBase
	ExchangeTag string
}

func (w *CandleWorker) SubscribeToCandle(
	pairSymbol string,
	eventCallback func(event CandleEvent),
	errorHandler func(err error),
) error {
	return nil
}

func (w *CandleWorker) SubscribeToCandlesList(
	intervalsPerPair map[string]string,
	eventCallback func(event CandleEvent),
	errorHandler func(err error),
) error {
	return nil
}

func (w *CandleWorker) GetExchangeTag() string {
	return w.ExchangeTag
}
