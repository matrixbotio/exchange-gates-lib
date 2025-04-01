package gate

import (
	"context"
	"fmt"
	"time"

	"github.com/antihax/optional"
	"github.com/gateio/gateapi-go/v6"

	adp "github.com/matrixbotio/exchange-gates-lib/internal/adapters"
	baseadp "github.com/matrixbotio/exchange-gates-lib/internal/adapters/base"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/gate/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
	"github.com/matrixbotio/exchange-gates-lib/pkg/errs"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
	"github.com/matrixbotio/exchange-gates-lib/pkg/utils"
)

const (
	adapterName         = "Gate.io Spot (Beta)"
	clientOrderIDFormat = "t-%s"
	spotAccountType     = "spot"
	channelID           = "matrixbot"
	requestTimeout      = time.Second * 15
)

type adapter struct {
	baseadp.AdapterBase

	creds  pkgStructs.APICredentials
	client *gateapi.APIClient
	auth   context.Context

	candleWorker GateCandleWorker
	tradeWorker  GateTradeWorker
}

func New() adp.Adapter {
	cfg := gateapi.NewConfiguration()
	cfg.AddDefaultHeader("X-Gate-Channel-Id", channelID)

	return &adapter{
		AdapterBase: baseadp.NewAdapterBase(
			consts.ExchangeIDgateSpot,
			adapterName,
			consts.GateAdapterTag,
		),
		client: gateapi.NewAPIClient(cfg),
	}
}

func (a *adapter) GetPairSymbol(baseTicker, quoteTicker string) string {
	return mappers.GetPairSymbol(baseTicker, quoteTicker)
}

func (a *adapter) GenClientOrderID() string {
	return fmt.Sprintf(clientOrderIDFormat, utils.GenClientOrderID())
}

func (a *adapter) Connect(credentials pkgStructs.APICredentials) error {
	a.creds = credentials
	a.tradeWorker.creds = credentials

	a.auth = context.WithValue(
		context.Background(),
		gateapi.ContextGateAPIV4,
		gateapi.GateAPIV4{
			Key:    credentials.Keypair.Public,
			Secret: credentials.Keypair.Secret,
		},
	)
	return nil
}

func (a *adapter) getUID() (int64, error) {
	ctx, ctxCancel := context.WithTimeout(a.auth, requestTimeout)
	defer ctxCancel()

	data, _, err := a.client.AccountApi.GetAccountDetail(ctx)
	if err != nil {
		return 0, fmt.Errorf("get account data: %w", err)
	}
	return data.UserId, nil
}

func (a *adapter) CanTrade() (bool, error) {
	if _, err := a.getUID(); err != nil {
		return false, fmt.Errorf("uid: %w", err)
	}
	return true, nil
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

	_, err := a.CanTrade()
	return err
}

func (a *adapter) GetAccountBalance() ([]structs.Balance, error) {
	if !a.creds.Keypair.IsSet() {
		return nil, errs.ErrAPIKeyNotSet
	}

	ctx, ctxCancel := context.WithTimeout(a.auth, requestTimeout)
	defer ctxCancel()

	data, _, err := a.client.SpotApi.ListSpotAccounts(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("list accounts: %w", err)
	}

	result, err := mappers.ConvertBalances(data)
	if err != nil {
		return nil, fmt.Errorf("convert: %w", err)
	}
	return result, nil
}

func (a *adapter) GetCandles(
	limit int,
	symbol string,
	interval consts.Interval,
) ([]workers.CandleData, error) {
	intervalGate, err := mappers.ConvertIntervalToGate(interval)
	if err != nil {
		return nil, fmt.Errorf("convert interval: %w", err)
	}

	data, _, err := a.client.SpotApi.ListCandlesticks(
		a.auth, symbol, &gateapi.ListCandlesticksOpts{
			Limit:    optional.NewInt32(int32(limit)),
			Interval: optional.NewString(intervalGate),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	result, err := mappers.ConvertCandles(data, interval)
	if err != nil {
		return nil, fmt.Errorf("convert: %w", err)
	}
	return result, nil
}

func (a *adapter) GetLimits() pkgStructs.ExchangeLimits {
	return pkgStructs.ExchangeLimits{
		MaxConnectionsPerBatch:   50,
		MaxConnectionsInDuration: time.Second,
		MaxTopicsPerWebsocket:    30,
	}
}
