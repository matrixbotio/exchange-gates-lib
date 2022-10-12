package test

import (
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/pkg/adapters"
)

func GetAdapter() adapters.Adapter {
	return &TestAdapter{
		ExchangeID: consts.TestExchangeID,
		Name:       "Test Exchange",
		Tag:        "test-exchange",
	}
}
