package gate

import (
	"context"
	"fmt"
	"time"

	"github.com/gateio/gateapi-go/v6"

	adp "github.com/matrixbotio/exchange-gates-lib/internal/adapters"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
	"github.com/matrixbotio/exchange-gates-lib/pkg/utils"
)

const clientOrderIDFormat = "t-%s"
const requestTimeout = time.Second * 10
const spotAccountType = "spot"

type adapter struct {
	ExchangeID int
	Name       string
	Tag        string

	accountType consts.AccountType
	keyPublic   string
	client      *gateapi.APIClient
	auth        context.Context
}

func New() adp.Adapter {
	return &adapter{
		ExchangeID: consts.ExchangeIDbybitSpot,
		Name:       "Gate.io Spot (Beta)",
		Tag:        "gate-spot",
		client:     gateapi.NewAPIClient(gateapi.NewConfiguration()),
	}
}

func (a *adapter) GetTag() string {
	return a.Tag
}

func (a *adapter) GetID() int {
	return a.ExchangeID
}

func (a *adapter) GetName() string {
	return a.Name
}

func (a *adapter) GetPairSymbol(baseTicker string, quoteTicker string) string {
	return fmt.Sprintf("%s_%s", baseTicker, quoteTicker)
}

func (a *adapter) GenClientOrderID() string {
	return fmt.Sprintf(clientOrderIDFormat, utils.GenClientOrderID())
}

func (a *adapter) Connect(credentials pkgStructs.APICredentials) error {
	a.keyPublic = credentials.Keypair.Public

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
	data, _, err := a.client.AccountApi.GetAccountDetail(a.auth)
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

	keyData, _, err := a.client.SubAccountApi.GetSubAccountKey(ctx, int32(uid), a.keyPublic)
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

func (a *adapter) VerifyAPIKeys(keyPublic, keySecret string) (
	pkgStructs.VerifyKeyStatus,
	error,
) {
	if err := a.Connect(pkgStructs.APICredentials{
		Type: pkgStructs.APICredentialsTypeKeypair,
		Keypair: pkgStructs.APIKeypair{
			Public: keyPublic,
			Secret: keySecret,
		},
	}); err != nil {
		return pkgStructs.VerifyKeyStatus{}, fmt.Errorf("connect: %w", err)
	}

	_, err := a.CanTrade()
	if err != nil {
		return pkgStructs.VerifyKeyStatus{}, err
	}

	return pkgStructs.VerifyKeyStatus{
		Active:      true,
		AccountType: consts.AccountTypeStandart, // TODO
	}, nil
}

func (a *adapter) GetAccountBalance() ([]structs.Balance, error) {
	//a.client.SpotApi.ListSpotAccounts(a.auth, &gateapi.ListSpotAccountsOpts{})

	// TODO
	return nil, nil
}

func (a *adapter) GetCandles(
	limit int,
	symbol string,
	interval string,
) ([]workers.CandleData, error) {
	// TODO
	return nil, nil
}

func (a *adapter) GetLimits() pkgStructs.ExchangeLimits {
	return pkgStructs.ExchangeLimits{
		MaxConnectionsPerBatch:   50,
		MaxConnectionsInDuration: time.Second,
	}
}

func (a *adapter) SetAccountType(accountType consts.AccountType) {
	a.accountType = accountType
}
