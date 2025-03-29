package bybit

import (
	"fmt"
	"time"

	"github.com/hirokisan/bybit/v2"

	adp "github.com/matrixbotio/exchange-gates-lib/internal/adapters"
	baseadp "github.com/matrixbotio/exchange-gates-lib/internal/adapters/base"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bybit/helpers"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
	"github.com/matrixbotio/exchange-gates-lib/pkg/utils"
)

const (
	adapterName = "ByBit Spot"
	adapterTag  = "bybit-spot"
)

type adapter struct {
	baseadp.AdapterBase

	client   *bybit.Client
	wsClient *bybit.WebSocketClient

	candleWorker *helpers.CandleEventWorkerBybit
	tradeWorker  *TradeEventWorkerBybit
}

func New() adp.Adapter {
	return &adapter{
		AdapterBase: baseadp.NewAdapterBase(
			consts.ExchangeIDbybitSpot,
			adapterName,
			adapterTag,
		),
		client:   bybit.NewClient(),
		wsClient: bybit.NewWebsocketClient(),
	}
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

	a.candleWorker = a.CreateCandleWorker()
	a.tradeWorker = a.CreateTradeEventsWorker()
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
