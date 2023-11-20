package binance

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-stack/stack"

	adp "github.com/matrixbotio/exchange-gates-lib/internal/adapters"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
)

type adapter struct {
	ExchangeID int
	Name       string
	Tag        string

	binanceAPI BinanceAPIWrapper
}

func New() adp.Adapter {
	stack.Caller(0)
	a := adapter{}
	a.Name = "Binance Spot"
	a.Tag = "binance-spot"
	a.ExchangeID = consts.ExchangeIDbinanceSpot
	return &a
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
		return errors.New("invalid credentials to connect to Binance")
	}

	if err := a.binanceAPI.Connect(
		credentials.Keypair.Public,
		credentials.Keypair.Secret,
		context.Background(),
	); err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	a.binanceAPI.Sync(context.Background())
	return nil
}
