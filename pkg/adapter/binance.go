package adapter

import (
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance"
)

func NewBinanceAdapter() Adapter {
	return binance.New()
}
