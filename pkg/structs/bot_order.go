package structs

import "fmt"

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
	return o.Qty == 0
}

func (o BotOrder) String() string {
	return fmt.Sprintf(
		"order: %q price %v, qty %v, pair %q",
		o.Type, o.Price, o.Qty, o.PairSymbol,
	)
}
