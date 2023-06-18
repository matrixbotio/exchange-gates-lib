package structs

// BotOrder - structure containing information about the order calculated by the bot
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

func (o BotOrder) IsEmpty() bool {
	return o.PairSymbol == "" && o.Qty == 0
}
