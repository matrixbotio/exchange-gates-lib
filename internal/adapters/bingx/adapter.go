package bingx

import (
	"fmt"
	"strings"
	"time"

	bingxgo "github.com/matrixbotio/go-bingx"

	adp "github.com/matrixbotio/exchange-gates-lib/internal/adapters"
	baseadp "github.com/matrixbotio/exchange-gates-lib/internal/adapters/base"
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
	idReplaceFrom         = "-"
	idReplaceTo           = "_"
)

type adapter struct {
	baseadp.AdapterBase

	client bingxgo.SpotClient
	creds  pkgStructs.APICredentials
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
	return strings.ReplaceAll(
		strings.ToLower(nano.ID(clientOrderIDLength)),
		idReplaceFrom, idReplaceTo,
	)
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
	a.creds = credentials
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
	interval consts.Interval,
) ([]workers.CandleData, error) {
	bingxInterval, err := ConvertIntervalToBingXRest(interval)
	if err != nil {
		return nil, fmt.Errorf("convert interval: %w", err)
	}

	klines, err := a.client.GetHistoricalKlines(
		symbol,
		string(bingxInterval),
		int64(limit),
	)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return ConvertKlinesRest(klines)
}
