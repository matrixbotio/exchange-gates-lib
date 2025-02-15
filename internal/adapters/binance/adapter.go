package binance

import (
	"context"
	"fmt"
	"time"

	adp "github.com/matrixbotio/exchange-gates-lib/internal/adapters"
	baseadp "github.com/matrixbotio/exchange-gates-lib/internal/adapters/base"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/errs"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/workers"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/wrapper"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	iWorkers "github.com/matrixbotio/exchange-gates-lib/internal/workers"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
	"github.com/matrixbotio/exchange-gates-lib/pkg/utils"
)

const (
	adapterName = "Binance Spot"
	adapterTag  = "binance-spot"
)

type adapter struct {
	baseadp.AdapterBase

	binanceAPI wrapper.BinanceAPIWrapper
}

func New(wrapper wrapper.BinanceAPIWrapper) adp.Adapter {
	return &adapter{
		AdapterBase: baseadp.NewAdapterBase(
			consts.ExchangeIDbinanceSpot,
			adapterName,
			adapterTag,
		),
		binanceAPI: wrapper,
	}
}

func (a *adapter) GetLimits() pkgStructs.ExchangeLimits {
	return pkgStructs.ExchangeLimits{
		MaxConnectionsPerBatch:   299,
		MaxConnectionsInDuration: 5 * time.Minute,
		MaxTopicsPerWebsocket:    450,
	}
}

func (a *adapter) GetPairSymbol(baseTicker string, quoteTicker string) string {
	return fmt.Sprintf("%s%s", baseTicker, quoteTicker)
}

func (a *adapter) GenClientOrderID() string {
	return utils.GenClientOrderID()
}

func (a *adapter) Connect(credentials pkgStructs.APICredentials) error {
	if credentials.Type != pkgStructs.APICredentialsTypeKeypair {
		return errs.ErrInvalidCredentials
	}

	if err := a.binanceAPI.Connect(
		context.Background(),
		credentials.Keypair.Public,
		credentials.Keypair.Secret,
	); err != nil {
		return fmt.Errorf("binance adapter: connect: %w", err)
	}

	a.binanceAPI.Sync(context.Background())
	return nil
}

func (a *adapter) GetPriceWorker(callback iWorkers.PriceEventCallback) iWorkers.IPriceWorker {
	return workers.NewPriceWorker(
		a.GetTag(),
		a.binanceAPI,
		callback,
	)
}

func (a *adapter) GetTradeEventsWorker() iWorkers.ITradeEventWorker {
	return workers.NewTradeEventsWorker(
		a.GetTag(),
		a.binanceAPI,
	)
}
