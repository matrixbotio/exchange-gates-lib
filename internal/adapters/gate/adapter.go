package gate

import (
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
}

func New() adp.Adapter {
	return &adapter{
		ExchangeID: consts.ExchangeIDbybitSpot,
		Name:       "Gate.io Spot (Beta)",
		Tag:        "gate-spot",
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
	// TODO
	return nil
}

func (a *adapter) CanTrade() (bool, error) {
	return true, nil // TODO
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
