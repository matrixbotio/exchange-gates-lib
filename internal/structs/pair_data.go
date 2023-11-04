package structs

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

func (data ExchangePairData) IsEmpty() bool {
	return data.Symbol == "" && data.ExchangeID == 0 && data.Status == ""
}
