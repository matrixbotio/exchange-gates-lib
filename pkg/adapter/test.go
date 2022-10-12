package adapter

import (
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/test"
)

func NewTestAdapter() Adapter {
	return test.New()
}
