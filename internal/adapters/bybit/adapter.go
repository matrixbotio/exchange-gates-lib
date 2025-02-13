package bybit

import (
	"fmt"
	"time"

	"github.com/hirokisan/bybit/v2"

	adp "github.com/matrixbotio/exchange-gates-lib/internal/adapters"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
	"github.com/matrixbotio/exchange-gates-lib/pkg/utils"
)

type adapter struct {
	ExchangeID int
	Name       string
	Tag        string

	client   *bybit.Client
	wsClient *bybit.WebSocketClient
}

func New() adp.Adapter {
	return &adapter{
		ExchangeID: consts.ExchangeIDbybitSpot,
		Name:       "ByBit Spot",
		Tag:        "bybit-spot",
		client:     bybit.NewClient(),
		wsClient:   bybit.NewWebsocketClient(),
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

func (a *adapter) GetLimits() pkgStructs.ExchangeLimits {
	return pkgStructs.ExchangeLimits{
		MaxConnectionsPerBatch:   499,
		MaxConnectionsInDuration: 5 * time.Minute,
		MaxTopicsPerWebsocket:    10,
	}
}

func (a *adapter) GetPairSymbol(baseTicker string, quoteTicker string) string {
	return fmt.Sprintf("%s%s", baseTicker, quoteTicker)
}

func (a *adapter) GenClientOrderID() string {
	return utils.GenClientOrderID()
}

func (a *adapter) Connect(credentials pkgStructs.APICredentials) error {
	a.client.WithAuth(credentials.Keypair.Public, credentials.Keypair.Secret)
	a.wsClient.WithAuth(credentials.Keypair.Public, credentials.Keypair.Secret)

	if err := a.client.SyncServerTime(); err != nil {
		return fmt.Errorf("sync time: %w", err)
	}
	return nil
}

func (a *adapter) getAccountType() bybit.AccountTypeV5 {
	return bybit.AccountTypeV5UNIFIED
}

func (a *adapter) CanTrade() (bool, error) {
	response, err := a.client.V5().User().GetAPIKey()
	if err != nil {
		return false, fmt.Errorf("get API key info: %w", err)
	}

	for _, permission := range response.Result.Permissions.Spot {
		if permission == "SpotTrade" {
			return true, nil
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
