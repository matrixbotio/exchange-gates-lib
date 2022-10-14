package test

import (
	"context"
	"time"

	"github.com/matrixbotio/exchange-gates-lib/pkg/consts"
	"github.com/matrixbotio/exchange-gates-lib/pkg/structs"
	"github.com/matrixbotio/exchange-gates-lib/pkg/utils"
	"github.com/matrixbotio/exchange-gates-lib/pkg/workers"
)

type adapter struct {
	ExchangeID int
	Name       string
	Tag        string
}

func New() *adapter {
	return &adapter{
		ExchangeID: consts.TestExchangeID,
		Name:       "Test Exchange",
		Tag:        "test-exchange",
	}
}

// GetName - get exchange adapter name
func (a *adapter) GetName() string {
	return a.Name
}

// GetTag - get exchange adapter tag
func (a *adapter) GetTag() string {
	return a.Tag
}

// GetID - get exchange adapter name
func (a *adapter) GetID() int {
	return a.ExchangeID
}

// Placeholders

// Connect to exchange
func (a *adapter) Connect(credentials structs.APICredentials) error {
	return nil
}

// PlaceOrder - place order on exchange
func (a *adapter) PlaceOrder(
	ctx context.Context, order structs.BotOrderAdjusted,
) (structs.CreateOrderResponse, error) {
	return structs.CreateOrderResponse{
		OrderID:       1,
		ClientOrderID: "test",
		OrigQuantity:  0.1,
	}, nil
}

// GetAccountData ..
func (a *adapter) GetAccountData() (structs.AccountData, error) {
	return structs.AccountData{
		CanTrade: true,
		Balances: []structs.Balance{
			{
				Asset:  consts.PairDefaultAsset,
				Free:   0,
				Locked: 0,
			},
		},
	}, nil
}

// GetPairLastPrice ..
func (a *adapter) GetPairLastPrice(pairSymbol string) (float64, error) {
	return 0, nil
}

// CancelPairOrder ..
func (a *adapter) CancelPairOrder(pairSymbol string, orderID int64, ctx context.Context) error {
	return nil
}

func (a *adapter) CancelPairOrderByClientOrderID(
	pairSymbol string,
	clientOrderID string,
	ctx context.Context,
) error {
	return nil
}

// GetOrderData - get test order data
func (a *adapter) GetOrderData(pairSymbol string, orderID int64) (structs.OrderData, error) {
	return structs.OrderData{
		OrderID:       orderID,
		ClientOrderID: "",
		Status:        consts.OrderStatusNew,
		AwaitQty:      100,
		FilledQty:     10,
		Price:         500,
		Symbol:        pairSymbol,
		Type:          consts.OrderTypeBuy,
		CreatedTime:   time.Now().UnixMilli(),
		UpdatedTime:   time.Now().UnixMilli(),
	}, nil
}

// GetOrderByClientOrderID ..
func (a *adapter) GetOrderByClientOrderID(pairSymbol string, clientOrderID string) (structs.OrderData, error) {
	return a.GetOrderData(pairSymbol, 0)
}

// GetPairOpenOrders ..
func (a *adapter) GetPairOpenOrders(pairSymbol string) ([]structs.OrderData, error) {
	return nil, nil
}

// GetPairs get all Binance pairs
func (a *adapter) GetPairs() ([]structs.ExchangePairData, error) {
	return nil, nil
}

// VerifyAPIKeys ..
func (a *adapter) VerifyAPIKeys(keyPublic, keySecret string) error {
	return nil
}

// GetTradeEventsWorker - create empty trade data worker
func (a *adapter) GetTradeEventsWorker() workers.ITradeEventWorker {
	w := workers.TradeEventWorker{}
	w.ExchangeTag = a.GetTag()
	return &w
}

// GetPriceWorker - create empty market data worker
func (a *adapter) GetPriceWorker(callback workers.PriceEventCallback) workers.IPriceWorker {
	w := workers.PriceWorker{}
	w.ExchangeTag = a.GetTag()
	return &w
}

// GetCandleWorker - create empty market candle worker
func (a *adapter) GetCandleWorker() workers.ICandleWorker {
	w := workers.CandleWorker{}
	w.ExchangeTag = a.GetTag()
	return &w
}

func (a *adapter) GetPairBalance(pair structs.PairSymbolData) (structs.PairBalance, error) {
	return structs.PairBalance{
		BaseAsset: &structs.AssetBalance{
			Ticker: pair.BaseTicker,
			Free:   10000,
		},
		QuoteAsset: &structs.AssetBalance{
			Ticker: pair.QuoteTicker,
			Free:   10000,
		},
	}, nil
}

func (a *adapter) GetPairData(pairSymbol string) (structs.ExchangePairData, error) {
	return utils.GetDefaultPairData(), nil
}

func (a *adapter) GetPairOrdersHistory(task structs.GetOrdersHistoryTask) ([]structs.OrderData, error) {
	return []structs.OrderData{}, nil
}

func (a *adapter) GetPrices() ([]structs.SymbolPrice, error) {
	return []structs.SymbolPrice{}, nil
}
