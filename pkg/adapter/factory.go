package adapter

import (
	"errors"

	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/wrapper"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bingx"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bybit"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/gate"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
)

func CreateAdapter(exchangeID int) (Adapter, error) {
	switch exchangeID {
	default:
		return nil, errors.New("exchange not found")
	case consts.ExchangeIDbinanceSpot:
		return binance.New(wrapper.NewWrapper()), nil
	case consts.ExchangeIDbybitSpot:
		return bybit.New(), nil
	case consts.ExchangeIDbingx:
		return bingx.New(), nil
	case consts.ExchangeIDgateSpot:
		return gate.New(), nil
	}
}

func CreateAdapters() map[int]Adapter {
	return map[int]Adapter{
		consts.ExchangeIDbinanceSpot: binance.New(wrapper.NewWrapper()),
		consts.ExchangeIDbybitSpot:   bybit.New(),
		consts.ExchangeIDbingx:       bingx.New(),
		consts.ExchangeIDgateSpot:    gate.New(),
	}
}
