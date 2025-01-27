package worker

import "github.com/matrixbotio/exchange-gates-lib/internal/workers"

type (
	CandleEvent = workers.CandleEvent

	MockCandleWorker = workers.MockICandleWorker
	MockTradeWorker  = workers.MockITradeEventWorker
	MockPriceWorker  = workers.MockIPriceWorker
)
