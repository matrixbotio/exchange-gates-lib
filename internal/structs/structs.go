package structs

import (
	"context"

	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/shopspring/decimal"
)

// BotOrderAdjusted - the same as BotOrder, only with the given values for the trading pair
type BotOrderAdjusted struct {
	// required
	PairSymbol string           `json:"pair"`
	Type       consts.OrderSide `json:"type"`
	Qty        string           `json:"qty"`
	Price      string           `json:"price"`
	Deposit    string           `json:"deposit"`

	// optional
	ClientOrderID string `json:"clientOrderID"`
	IsMarketOrder bool   `json:"isMarket"`

	// calculated
	MinQty           float64 `json:"minQty"`
	MinQtyPassed     bool    `json:"minQtyPassed"`
	MinDeposit       float64 `json:"minDeposit"`
	MinDepositPassed bool    `json:"minDepositPassed"`
}

func (o BotOrderAdjusted) IsEmpty() bool {
	return o.Qty == "" && o.Price == ""
}

// CreateOrderResponse - response from the exchange about the placed order
type CreateOrderResponse struct {
	OrderID       int64              `json:"orderID"`
	ClientOrderID string             `json:"clientOrderID"`
	OrigQuantity  float64            `json:"originalQty"`
	Price         float64            `json:"price"`
	Symbol        string             `json:"symbol"`
	Type          consts.OrderSide   `json:"orderRes"`
	CreatedTime   int64              `json:"createdTime"` // unix timestamp ms
	Status        consts.OrderStatus `json:"status"`
}

// Balance - Trading pair balance
type Balance struct {
	Asset  string  `json:"asset"`
	Free   float64 `json:"free"`
	Locked float64 `json:"locked"`
}

// AccountData & balances
type AccountData struct {
	CanTrade bool      `json:"canTrade"`
	Balances []Balance `json:"balances"`
}

// PairBalance - data on the balance of a trading pair for each of the two currencies
type PairBalance struct {
	BaseAsset  *AssetBalance `json:"base"`
	QuoteAsset *AssetBalance `json:"quote"`
}

// AssetBalance - is a wraper for asset balance data
type AssetBalance struct {
	Ticker string  `json:"ticker"`
	Free   float64 `json:"free"`
	Locked float64 `json:"locked"`
}

// PairSymbolData - contains pair symbol data
type PairSymbolData struct {
	BaseTicker  string `json:"base"`   // ETH
	QuoteTicker string `json:"quote"`  // USDT
	Symbol      string `json:"symbol"` // ETHUSDT
}

// APIPassword - password authentication
type APIPassword string

// APIEmail - email authentication
type APIEmail string

// GetOrdersHistoryTask - data for GetPairOrdersHistory request
type GetOrdersHistoryTask struct {
	// required
	PairSymbol string
	StartTime  int64 // unix timestamp ms

	// optional
	EndTime int64 // unix timestamp ms
	Ctx     context.Context
}

// SymbolPrice define symbol and price pair
type SymbolPrice struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}

type OrderFees struct {
	BaseAsset  decimal.Decimal `json:"base"`
	QuoteAsset decimal.Decimal `json:"quote"`
}
