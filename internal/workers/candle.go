//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package workers

import (
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
)

// CandleEvent - changes in trading candles for a specific pair
type CandleEvent struct {
	// required
	Symbol     string     `json:"symbol"`
	Candle     CandleData `json:"candle"`
	Time       int64      `json:"time"`
	IsFinished bool

	// optional
	BaseAsset  string `json:"baseAsset"`
	QuoteAsset string `json:"quoteAsset"`
}

// CandleData - trading candle
type CandleData struct {
	StartTime int64           `json:"startTime"`
	EndTime   int64           `json:"endTime"`
	Interval  consts.Interval `json:"interval"`
	Open      float64         `json:"open"`
	Close     float64         `json:"close"`
	High      float64         `json:"high"`
	Low       float64         `json:"low"`
	Volume    float64         `json:"volume"`
}

// CandleWorker - worker for subscribtion to exchange candle events
type CandleWorker struct {
	workerBase
	ExchangeTag string
}

func (w *CandleWorker) SubscribeToCandle(
	pairSymbol string,
	interval consts.Interval,
	eventCallback func(event CandleEvent),
	errorHandler func(err error),
) error {
	return nil
}

func (w *CandleWorker) GetExchangeTag() string {
	return w.ExchangeTag
}
