package matrixgates

/*
BotOrder - structure containing information about the order placed by the bot.
Used when auto-resuming trades
*/
type BotOrder struct {
	Type    string  `json:"type"`
	Qty     float64 `json:"qty"`
	Price   float64 `json:"price"`
	Deposit float64 `json:"deposit"`
}

//TradeEventData - container for bot trading new event data
type TradeEventData struct {
	OrderID        int64
	OrderAwaitQty  float64 //initial order qty
	OrderFilledQty float64 //event executed qty
	Status         string  //used in bot.getOrderData
}

//CreateOrderResponse ..
type CreateOrderResponse struct {
	OrderID       int64
	ClientOrderID string
	OrigQuantity  float64
	Price         float64
}
