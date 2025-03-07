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
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
	"github.com/matrixbotio/exchange-gates-lib/pkg/utils"
)

const (
	adapterName = "Gate.io Spot (Beta)"
	adapterTag  = "gate-spot"

	clientOrderIDFormat = "t-%s"
	spotAccountType     = "spot"
	requestTimeout      = time.Second * 15
)

type adapter struct {
	baseadp.AdapterBase

	creds  pkgStructs.APICredentials
	client *gateapi.APIClient
	auth   context.Context
}

func New() adp.Adapter {
	return &adapter{
		AdapterBase: baseadp.NewAdapterBase(
			consts.ExchangeIDgateSpot,
			adapterName,
			adapterTag,
		),
		client: gateapi.NewAPIClient(gateapi.NewConfiguration()),
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
	uid, err := a.getUID()
	if err != nil {
		return false, fmt.Errorf("uid: %w", err)
	}

	ctx, ctxCancel := context.WithTimeout(a.auth, requestTimeout)
	defer ctxCancel()

	keyData, _, err := a.client.SubAccountApi.GetSubAccountKey(
		ctx,
		int32(uid),
		a.creds.Keypair.Public,
	)
	if err != nil {
		return false, fmt.Errorf("get key data: %w", err)
	}

	for _, permissions := range keyData.Perms {
		if permissions.Name == spotAccountType {
			return !permissions.ReadOnly, nil
		}
	}
	return false, nil
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
