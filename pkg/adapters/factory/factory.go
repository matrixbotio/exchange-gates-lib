package factory

import (
	"errors"

	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/pkg/adapters"
	"github.com/matrixbotio/exchange-gates-lib/pkg/adapters/binance"
	"github.com/matrixbotio/exchange-gates-lib/pkg/adapters/test"
)

// GetAdapter - get supported exchange adapter with interface
func GetAdapter(exchangeID int) (adapters.Adapter, error) {
	switch exchangeID {
	default:
		return nil, errors.New("exchange not found")
	case consts.ExchangeIDbinanceSpot:
		return binance.NewBinanceSpotAdapter(), nil
	case consts.TestExchangeID:
		return test.GetAdapter(), nil
	}
}

// GetAdapters - get all supported exchange adapters
func GetAdapters() map[int]adapters.Adapter {
	return map[int]adapters.Adapter{
		consts.ExchangeIDbinanceSpot: binance.NewBinanceSpotAdapter(),
	}
}
