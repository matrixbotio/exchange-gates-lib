package bingx

import (
	"fmt"
	"strings"
	"time"

	bingxgo "github.com/Sagleft/go-bingx"

	adp "github.com/matrixbotio/exchange-gates-lib/internal/adapters"
	baseadp "github.com/matrixbotio/exchange-gates-lib/internal/adapters/base"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bingx/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
	"github.com/matrixbotio/go-common-lib/pkg/nano"
)

const (
	adapterName           = "BingX Spot"
	limitOrder            = "LIMIT"
	limitOrderTimeInForce = "GTC"
	brokerSourceKey       = "Matrixbot"
	symbolFormat          = "%s-%s"
	clientOrderIDLength   = 32
)

type adapter struct {
	baseadp.AdapterBase

	client bingxgo.SpotClient
}

func New() adp.Adapter {
	return &adapter{
		AdapterBase: baseadp.NewAdapterBase(
			consts.ExchangeIDbingx,
			adapterName,
			consts.BingXAdapterTag,
		),
	}
}

func (a *adapter) GenClientOrderID() string {
	return strings.ToLower(nano.ID(clientOrderIDLength))
}

func (a *adapter) GetPairSymbol(
	baseTicker string,
	quoteTicker string,
) string {
	return fmt.Sprintf(symbolFormat, baseTicker, quoteTicker)
}

func (a *adapter) GetLimits() pkgStructs.ExchangeLimits {
	return pkgStructs.ExchangeLimits{
		MaxConnectionsPerBatch:   10,
		MaxConnectionsInDuration: time.Second,
		MaxTopicsPerWebsocket:    200,
	}
}

func (a *adapter) Connect(credentials pkgStructs.APICredentials) error {
	a.client = bingxgo.NewSpotClient(bingxgo.NewClient(
		credentials.Keypair.Public,
		credentials.Keypair.Secret,
	).SetBrokerSourceKey(brokerSourceKey))
	return nil
}

func (a *adapter) VerifyAPIKeys(keyPublic, keySecret string) error {
	if err := a.Connect(pkgStructs.APICredentials{
		Type: pkgStructs.APICredentialsTypeKeypair,
		Keypair: pkgStructs.APIKeypair{
			Public: keyPublic,
			Secret: keySecret,
		},
	}); err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	_, err := a.client.GetBalance()
	return err
}

func (a *adapter) GetCandles(
	limit int,
	symbol string,
	interval string,
) ([]workers.CandleData, error) {
	klines, err := a.client.GetHistoricalKlines(symbol, interval, int64(limit))
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return mappers.ConvertKlines(klines), nil
}
