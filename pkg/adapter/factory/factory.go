package factory

import (
	"errors"

	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	adp "github.com/matrixbotio/exchange-gates-lib/pkg/adapter"
	"github.com/matrixbotio/exchange-gates-lib/pkg/adapter/binance"
	"github.com/matrixbotio/exchange-gates-lib/pkg/adapter/test"
)

func CreateAdapter(exchangeID int) (adp.Adapter, error) {
	switch exchangeID {
	default:
		return nil, errors.New("exchange not found")
	case consts.ExchangeIDbinanceSpot:
		return binance.New(), nil
	case consts.TestExchangeID:
		return test.New(), nil
	}
}

func CreateAdapters() map[int]adp.Adapter {
	return map[int]adp.Adapter{
		consts.ExchangeIDbinanceSpot: binance.New(),
	}
}
