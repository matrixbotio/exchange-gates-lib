package matrixgates

/*
BotOrder - structure containing information about the order placed by the bot.
Used when auto-resuming trades
*/
type BotOrder struct {
	PairSymbol string  `json:"pair"`
	Type       string  `json:"type"`
	Qty        float64 `json:"qty"`
	Price      float64 `json:"price"`
	Deposit    float64 `json:"deposit"`
}

// TradeEventData - container for bot trading new event data
type TradeEventData struct {
	OrderID        int64
	OrderAwaitQty  float64 //initial order qty
	OrderFilledQty float64 //event executed qty
	Status         string  //used in bot.getOrderData
}

// CreateOrderResponse ..
type CreateOrderResponse struct {
	OrderID       int64
	ClientOrderID string
	OrigQuantity  float64
	Price         float64
}

// Balance - Trading pair balance
type Balance struct {
	Asset  string
	Free   float64
	Locked float64
}

// AccountData & balances
type AccountData struct {
	CanTrade bool
	Balances []Balance
}

// Order data
type Order struct {
	OrderID       int64
	ClientOrderID string
	Status        string
}

// PairBalance - data on the balance of a trading pair for each of the two currencies
type PairBalance struct {
	BaseAsset  *AssetBalance
	QuoteAsset *AssetBalance
}

// AssetBalance - is a wraper for asset balance data
type AssetBalance struct {
	Ticker string  `json:"ticker"`
	Free   float64 `json:"free"`
	Locked float64 `json:"locked"`
}

// PairSymbolData - contains pair symbol data
type PairSymbolData struct {
	BaseTicker  string // ETH
	QuoteTicker string // USDT
	Symbol      string // ETHUSDT
}

// ExchangePairData contains information about a trading pair, data about order limits
type ExchangePairData struct {
	ID             int     `json:"id"`
	ExchangeID     int     `json:"exchangeID"`     // 1
	BaseAsset      string  `json:"baseAsset"`      // ETH
	BasePrecision  int     `json:"basePrecision"`  // 4
	QuoteAsset     string  `json:"quoteAsset"`     // USDT
	QuotePrecision int     `json:"quotePrecision"` // 2
	Status         string  `json:"status"`         // TRADING
	Symbol         string  `json:"symbol"`         // ETHUSDT
	MinQty         float64 `json:"minQty"`
	MaxQty         float64 `json:"maxQty"`
	MinPrice       float64 `json:"minPrice"`
	QtyStep        float64 `json:"qtyStep"`
	PriceStep      float64 `json:"priceStep"`
	AllowedMargin  bool    `json:"allowedMargin"`
	AllowedSpot    bool    `json:"allowedSpot"`
	InUse          bool    `json:"inUse"`
}

// APICredentialsType - API credentials type ^ↀᴥↀ^
type APICredentialsType string

// APIKeypair - data for authorization via public and private keys
type APIKeypair struct {
	Public string
	Secret string
}

// APIPassword - password authentication
type APIPassword string

// APIEmail - email authentication
type APIEmail string

// APICredentialsTypeKeypair - public and private key pair
var APICredentialsTypeKeypair APICredentialsType = "keypair"

// APICredentials - data for authorization to the exchange API
type APICredentials struct {
	Type APICredentialsType

	Keypair  APIKeypair
	Password APIPassword
	Email    APIEmail
}
