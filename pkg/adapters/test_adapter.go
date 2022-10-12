package adapters

import (
	"context"
	"time"

	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/utils"
	"github.com/matrixbotio/exchange-gates-lib/pkg/structs"
	workers2 "github.com/matrixbotio/exchange-gates-lib/pkg/workers"
)

// TestAdapter - abstract universal exchange adapter
type TestAdapter struct {
	ExchangeID int
	Name       string
	Tag        string
}

// GetName - get exchange adapter name
func (a *TestAdapter) GetName() string {
	return a.Name
}

// GetTag - get exchange adapter tag
func (a *TestAdapter) GetTag() string {
	return a.Tag
}

// GetID - get exchange adapter name
func (a *TestAdapter) GetID() int {
	return a.ExchangeID
}

// Placeholders

// Connect to exchange
func (a *TestAdapter) Connect(credentials structs.APICredentials) error {
	return nil
}

// PlaceOrder - place order on exchange
func (a *TestAdapter) PlaceOrder(
	ctx context.Context, order structs.BotOrderAdjusted,
) (structs.CreateOrderResponse, error) {
	return structs.CreateOrderResponse{
		OrderID:       1,
		ClientOrderID: "test",
		OrigQuantity:  0.1,
	}, nil
}

// GetAccountData ..
func (a *TestAdapter) GetAccountData() (structs.AccountData, error) {
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
func (a *TestAdapter) GetPairLastPrice(pairSymbol string) (float64, error) {
	return 0, nil
}

// CancelPairOrder ..
func (a *TestAdapter) CancelPairOrder(pairSymbol string, orderID int64, ctx context.Context) error {
	return nil
}

func (a *TestAdapter) CancelPairOrderByClientOrderID(
	pairSymbol string,
	clientOrderID string,
	ctx context.Context,
) error {
	return nil
}

// GetOrderData - get test order data
func (a *TestAdapter) GetOrderData(pairSymbol string, orderID int64) (structs.OrderData, error) {
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
func (a *TestAdapter) GetOrderByClientOrderID(pairSymbol string, clientOrderID string) (structs.OrderData, error) {
	return a.GetOrderData(pairSymbol, 0)
}

// GetPairOpenOrders ..
func (a *TestAdapter) GetPairOpenOrders(pairSymbol string) ([]structs.OrderData, error) {
	return nil, nil
}

// GetPairs get all Binance pairs
func (a *TestAdapter) GetPairs() ([]structs.ExchangePairData, error) {
	return nil, nil
}

// VerifyAPIKeys ..
func (a *TestAdapter) VerifyAPIKeys(keyPublic, keySecret string) error {
	return nil
}

// GetTradeEventsWorker - create empty trade data worker
func (a *TestAdapter) GetTradeEventsWorker() workers2.ITradeEventWorker {
	w := workers2.TradeEventWorker{}
	w.ExchangeTag = a.GetTag()
	return &w
}

// GetPriceWorker - create empty market data worker
func (a *TestAdapter) GetPriceWorker(callback workers2.PriceEventCallback) workers2.IPriceWorker {
	w := workers2.PriceWorker{}
	w.ExchangeTag = a.GetTag()
	return &w
}

// GetCandleWorker - create empty market candle worker
func (a *TestAdapter) GetCandleWorker() workers2.ICandleWorker {
	w := workers2.CandleWorker{}
	w.ExchangeTag = a.GetTag()
	return &w
}

func (a *TestAdapter) GetPairBalance(pair structs.PairSymbolData) (structs.PairBalance, error) {
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

func (a *TestAdapter) GetPairData(pairSymbol string) (structs.ExchangePairData, error) {
	return utils.GetDefaultPairData(), nil
}

func (a *TestAdapter) GetPairOrdersHistory(task structs.GetOrdersHistoryTask) ([]structs.OrderData, error) {
	return []structs.OrderData{}, nil
}

func (a *TestAdapter) GetPrices() ([]structs.SymbolPrice, error) {
	return []structs.SymbolPrice{}, nil
}
