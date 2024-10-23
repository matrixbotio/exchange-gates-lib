package gate

import (
	"context"
	"fmt"

	"github.com/gateio/gateapi-go/v6"

	adp "github.com/matrixbotio/exchange-gates-lib/internal/adapters"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
)

type adapter struct {
	ExchangeID int
	Name       string
	Tag        string

	client *gateapi.APIClient
	auth   context.Context
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

func (a *adapter) Connect(credentials pkgStructs.APICredentials) error {
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

func (a *adapter) CanTrade() (bool, error) {
	_, _, err := a.client.AccountApi.GetAccountDetail(a.auth)
	if err != nil {
		return false, fmt.Errorf("get account data: %w", err)
	}
	return true, nil
}

func (a *adapter) VerifyAPIKeys(keyPublic, keySecret string) error {
	// TODO
	return nil
}

func (a *adapter) GetAccountBalance() ([]structs.Balance, error) {
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
