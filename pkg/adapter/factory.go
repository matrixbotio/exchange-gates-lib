package adapter

import (
	"errors"

	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/test"
	"github.com/matrixbotio/exchange-gates-lib/pkg/consts"
)

func CreateAdapter(exchangeID int) (Adapter, error) {
	switch exchangeID {
	default:
		return nil, errors.New("exchange not found")
	case consts.ExchangeIDbinanceSpot:
		return binance.New(), nil
	case consts.TestExchangeID:
		return test.New(), nil
	}
}

func CreateAdapters() map[int]Adapter {
	return map[int]Adapter{
		consts.ExchangeIDbinanceSpot: binance.New(),
	}
}
