package binance

import (
	"context"
	"errors"
	"fmt"

	adp "github.com/matrixbotio/exchange-gates-lib/internal/adapters"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
)

const (
	adapterName = "Binance Spot"
	adapterTag  = "binance-spot"
)

type adapter struct {
	ExchangeID int
	Name       string
	Tag        string

	binanceAPI BinanceAPIWrapper
}

func New() adp.Adapter {
	return createAdapter(NewWrapper())
}

func createAdapter(wrapper BinanceAPIWrapper) *adapter {
	return &adapter{
		ExchangeID: consts.ExchangeIDbinanceSpot,
		Name:       adapterName,
		Tag:        adapterTag,
		binanceAPI: wrapper,
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
	if credentials.Type != pkgStructs.APICredentialsTypeKeypair {
		return errs.ErrInvalidCredentials
	}

	if err := a.binanceAPI.Connect(
		context.Background(),
		credentials.Keypair.Public,
		credentials.Keypair.Secret,
	); err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	a.binanceAPI.Sync(context.Background())
	return nil
}
