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

//TradeEventData - container for bot trading new event data
type TradeEventData struct {
	OrderID        int64   `json:"orderID"`
	OrderAwaitQty  float64 `json:"await_qty"`  //initial order qty
	OrderFilledQty float64 `json:"filled_qty"` //event executed qty
	Status         string  `json:"status"`     //used in bot.getOrderData
}

//CreateOrderResponse ..
type CreateOrderResponse struct {
	OrderID       int64   `json:"orderID"`
	ClientOrderID string  `json:"clientOrderID"`
	OrigQuantity  float64 `json:"orig_qty"`
	Price         float64 `json:"price"`
}

//Balance - Trading pair balance
type Balance struct {
	Asset  string  `json:"asset"`
	Free   float64 `json:"free"`
	Locked float64 `json:"locked"`
}

//AccountData & balances
type AccountData struct {
	CanTrade bool      `json:"canTrade"`
	Balances []Balance `json:"balances"`
}

//Order data
type Order struct {
	OrderID       int64  `json:"orderID"`
	ClientOrderID string `json:"clientOrderID"`
	Status        string `json:"status"`
}

//ExchangePairData contains information about a trading pair, data about order limits
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

//APICredentialsType - API credentials type ^ↀᴥↀ^
type APICredentialsType string

//APIKeypair - data for authorization via public and private keys
type APIKeypair struct {
	Public string `json:"public"`
	Secret string `json:"secret"`
}

//APIPassword - password authentication
type APIPassword string

//APIEmail - email authentication
type APIEmail string

//APICredentialsTypeKeypair - public and private key pair
var APICredentialsTypeKeypair APICredentialsType = "keypair"

//APICredentials - data for authorization to the exchange API
type APICredentials struct {
	Type APICredentialsType `json:"type"`

	Keypair  APIKeypair  `json:"keypair"`
	Password APIPassword `json:"password"`
	Email    APIEmail    `json:"email"`
}
