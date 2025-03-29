package binance

import (
	"context"
	"fmt"
	"time"

	adp "github.com/matrixbotio/exchange-gates-lib/internal/adapters"
	baseadp "github.com/matrixbotio/exchange-gates-lib/internal/adapters/base"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/errs"
	binanceworkers "github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/workers"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/wrapper"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
	"github.com/matrixbotio/exchange-gates-lib/pkg/utils"
)

const (
	adapterName = "Binance Spot"
)

type adapter struct {
	baseadp.AdapterBase
	binanceAPI wrapper.BinanceAPIWrapper

	tradeWorker  *binanceworkers.TradeEventWorkerBinance
	candleWorker *CandleWorkerBinance
}

func New(wrapper wrapper.BinanceAPIWrapper) adp.Adapter {
	return &adapter{
		AdapterBase: baseadp.NewAdapterBase(
			consts.ExchangeIDbinanceSpot,
			adapterName,
			consts.BinanceAdapterTag,
		),
		binanceAPI:   wrapper,
		candleWorker: NewCandleWorker(wrapper),
		tradeWorker:  binanceworkers.NewTradeEventsWorker(wrapper),
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
