package pkg

import (
	"context"
	"time"

	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	workers2 "github.com/matrixbotio/exchange-gates-lib/pkg/workers"
)

// ExchangeAdapter - abstract universal exchange adapter
type ExchangeAdapter struct {
	ExchangeID int
	Name       string
	Tag        string
}

// GetName - get exchange adapter name
func (a *ExchangeAdapter) GetName() string {
	return a.Name
}

// GetTag - get exchange adapter tag
func (a *ExchangeAdapter) GetTag() string {
	return a.Tag
}

// GetID - get exchange adapter name
func (a *ExchangeAdapter) GetID() int {
	return a.ExchangeID
}

// Placeholders

// Connect to exchange
func (a *ExchangeAdapter) Connect(credentials APICredentials) error {
	return nil
}

// PlaceOrder - place order on exchange
func (a *ExchangeAdapter) PlaceOrder(
	ctx context.Context, order BotOrderAdjusted,
) (CreateOrderResponse, error) {
	return CreateOrderResponse{
		OrderID:       1,
		ClientOrderID: "test",
		OrigQuantity:  0.1,
	}, nil
}

// GetAccountData ..
func (a *ExchangeAdapter) GetAccountData() (AccountData, error) {
	return AccountData{
		CanTrade: true,
		Balances: []Balance{
			{
				Asset:  consts.PairDefaultAsset,
				Free:   0,
				Locked: 0,
			},
		},
	}, nil
}

// GetPairLastPrice ..
func (a *ExchangeAdapter) GetPairLastPrice(pairSymbol string) (float64, error) {
	return 0, nil
}

// CancelPairOrder ..
func (a *ExchangeAdapter) CancelPairOrder(pairSymbol string, orderID int64, ctx context.Context) error {
	return nil
}

func (a *ExchangeAdapter) CancelPairOrderByClientOrderID(
	pairSymbol string,
	clientOrderID string,
	ctx context.Context,
) error {
	return nil
}

// GetOrderData - get test order data
func (a *ExchangeAdapter) GetOrderData(pairSymbol string, orderID int64) (OrderData, error) {
	return OrderData{
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
func (a *ExchangeAdapter) GetOrderByClientOrderID(pairSymbol string, clientOrderID string) (OrderData, error) {
	return a.GetOrderData(pairSymbol, 0)
}

// GetPairOpenOrders ..
func (a *ExchangeAdapter) GetPairOpenOrders(pairSymbol string) ([]OrderData, error) {
	return nil, nil
}

// GetPairs get all Binance pairs
func (a *ExchangeAdapter) GetPairs() ([]ExchangePairData, error) {
	return nil, nil
}

// VerifyAPIKeys ..
func (a *ExchangeAdapter) VerifyAPIKeys(keyPublic, keySecret string) error {
	return nil
}

// GetTradeEventsWorker - create empty trade data worker
func (a *ExchangeAdapter) GetTradeEventsWorker() workers2.ITradeEventWorker {
	w := workers2.TradeEventWorker{}
	w.ExchangeTag = a.GetTag()
	return &w
}

// GetPriceWorker - create empty market data worker
func (a *ExchangeAdapter) GetPriceWorker(callback workers2.PriceEventCallback) workers2.IPriceWorker {
	w := workers2.PriceWorker{}
	w.ExchangeTag = a.GetTag()
	return &w
}

// GetCandleWorker - create empty market candle worker
func (a *ExchangeAdapter) GetCandleWorker() workers2.ICandleWorker {
	w := workers2.CandleWorker{}
	w.ExchangeTag = a.GetTag()
	return &w
}

func (a *ExchangeAdapter) GetPairBalance(pair PairSymbolData) (PairBalance, error) {
	return PairBalance{
		BaseAsset: &AssetBalance{
			Ticker: pair.BaseTicker,
			Free:   10000,
		},
		QuoteAsset: &AssetBalance{
			Ticker: pair.QuoteTicker,
			Free:   10000,
		},
	}, nil
}

func (a *ExchangeAdapter) GetPairData(pairSymbol string) (ExchangePairData, error) {
	return GetDefaultPairData(), nil
}

func (a *ExchangeAdapter) GetPairOrdersHistory(task GetOrdersHistoryTask) ([]OrderData, error) {
	return []OrderData{}, nil
}

func (a *ExchangeAdapter) GetPrices() ([]SymbolPrice, error) {
	return []SymbolPrice{}, nil
}
