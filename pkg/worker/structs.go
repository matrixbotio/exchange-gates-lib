package worker

import "github.com/matrixbotio/exchange-gates-lib/internal/workers"

type (
	CandleEvent = workers.CandleEvent

	MockCandleWorker = workers.MockICandleWorker
	MockTradeWorker  = workers.MockITradeEventWorker
	MockPriceWorker  = workers.MockIPriceWorker
)

var (
	NewMockCandleWorker = workers.NewMockICandleWorker
	NewMockTradeWorker  = workers.NewMockITradeEventWorker
	NewMockPriceWorker  = workers.NewMockIPriceWorker
)
