package factory

import (
	"errors"

	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/pkg/adapters"
	"github.com/matrixbotio/exchange-gates-lib/pkg/adapters/binance"
	"github.com/matrixbotio/exchange-gates-lib/pkg/adapters/test"
)

func CreateAdapter(exchangeID int) (adapters.Adapter, error) {
	switch exchangeID {
	default:
		return nil, errors.New("exchange not found")
	case consts.ExchangeIDbinanceSpot:
		return binance.New(), nil
	case consts.TestExchangeID:
		return test.New(), nil
	}
}

func CreateAdapters() map[int]adapters.Adapter {
	return map[int]adapters.Adapter{
		consts.ExchangeIDbinanceSpot: binance.New(),
	}
}
