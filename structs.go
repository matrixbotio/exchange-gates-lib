package matrixgates

import "context"

/*
BotOrder - structure containing information about the order placed by the bot.
Used when auto-resuming trades
*/
type BotOrder struct {
	// required
	PairSymbol string  `json:"pair"`
	Type       string  `json:"type"`
	Qty        float64 `json:"qty"`
	Price      float64 `json:"price"`
	Deposit    float64 `json:"deposit"`

	// optional
	ClientOrderID string `json:"clientOrderID"`
}

// BotOrderAdjusted - the same as BotOrder, only with the given values for the trading pair
type BotOrderAdjusted struct {
	// required
	PairSymbol string `json:"pair"`
	Type       string `json:"type"`
	Qty        string `json:"qty"`
	Price      string `json:"price"`
	Deposit    string `json:"deposit"`

	// optional
	ClientOrderID string `json:"clientOrderID"`

	// calculated
	MinQty           float64 `json:"minQty"`
	MinQtyPassed     bool    `json:"minQtyPassed"`
	MinDeposit       float64 `json:"minDeposit"`
	MinDepositPassed bool    `json:"minDepositPassed"`
}

// CreateOrderResponse ..
type CreateOrderResponse struct {
	OrderID       int64   `json:"orderID"`
	ClientOrderID string  `json:"clientOrderID"`
	OrigQuantity  float64 `json:"origQty"`
	Price         float64 `json:"price"`
	Symbol        string  `json:"symbol"`
	Type          string  `json:"orderRes"`
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

// OrderData - placed order data
type OrderData struct {
	OrderID       int64   `json:"orderID"`
	ClientOrderID string  `json:"clientOrderID"`
	Status        string  `json:"status"`    // used in bot.getOrderData
	AwaitQty      float64 `json:"awaitQty"`  // initial order qty
	FilledQty     float64 `json:"filledQty"` // event executed qty
	Price         float64 `json:"price"`
	Symbol        string  `json:"symbol"`
	Type          string  `json:"type"`        // "buy" or "sell"
	CreatedTime   int64   `json:"createdTime"` // unix ms
	UpdatedTime   int64   `json:"updatedTime"` // unix ms
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

// ExchangePairData contains information about a trading pair, data about order limits
type ExchangePairData struct {
	ID                 int     `json:"id"`
	ExchangeID         int     `json:"exchangeID"`     // 1
	BaseAsset          string  `json:"baseAsset"`      // ETH
	BasePrecision      int     `json:"basePrecision"`  // 4
	QuoteAsset         string  `json:"quoteAsset"`     // USDT
	QuotePrecision     int     `json:"quotePrecision"` // 2
	Status             string  `json:"status"`         // TRADING
	Symbol             string  `json:"symbol"`         // ETHUSDT
	MinQty             float64 `json:"minQty"`
	MaxQty             float64 `json:"maxQty"`
	OriginalMinDeposit float64 `json:"origMinDeposit"`
	MinDeposit         float64 `json:"minDeposit"`
	MinPrice           float64 `json:"minPrice"`
	QtyStep            float64 `json:"qtyStep"`
	PriceStep          float64 `json:"priceStep"`
	AllowedMargin      bool    `json:"allowedMargin"`
	AllowedSpot        bool    `json:"allowedSpot"`
	InUse              bool    `json:"inUse"`
}

// APICredentialsType - API credentials type ^ↀᴥↀ^
type APICredentialsType string

// APIKeypair - data for authorization via public and private keys
type APIKeypair struct {
	Public string `json:"public"`
	Secret string `json:"secret"`
}

// APIPassword - password authentication
type APIPassword string

// APIEmail - email authentication
type APIEmail string

// APICredentialsTypeKeypair - public and private key pair
var APICredentialsTypeKeypair APICredentialsType = "keypair"

// APICredentials - data for authorization to the exchange API
type APICredentials struct {
	Type APICredentialsType `json:"type"`

	Keypair  APIKeypair  `json:"keypair"`
	Password APIPassword `json:"password"`
	Email    APIEmail    `json:"email"`
}

// CheckOrdersResponse - data on checked and restored orders
type CheckOrdersResponse struct {
	ExecutedOrders  []*OrderData
	CancelledOrders []int64 // order IDs
}

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
