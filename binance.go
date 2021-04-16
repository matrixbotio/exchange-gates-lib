package matrixgates

import sharederrs "github.com/matrixbotio/shared-errors"

//BinanceSpotAdapter - bot exchange adapter for BinanceSpot
type BinanceSpotAdapter struct {
	ExchangeAdapter
}

//GetOrderData ..
func (a *BinanceSpotAdapter) GetOrderData() (*TradeEventData, *sharederrs.APIError) {
	//TODO
	return nil, nil
}
