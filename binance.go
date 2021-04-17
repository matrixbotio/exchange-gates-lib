package matrixgates

import (
	"context"
	"log"

	"github.com/adshao/go-binance/v2"
	sharederrs "github.com/matrixbotio/shared-errors"
)

//BinanceSpotAdapter - bot exchange adapter for BinanceSpot
type BinanceSpotAdapter struct {
	ExchangeAdapter
}

//GetOrderData ..
func (a *BinanceSpotAdapter) GetOrderData() (*TradeEventData, *sharederrs.APIError) {
	//TODO
	return nil, nil
}

//PlaceOrder ..
func (a *BinanceSpotAdapter) PlaceOrder(order BotOrder) (*CreateOrderResponse, *sharederrs.APIError) {
	//TODO
	return nil, nil
}

//GetAccountData ..
func (a *BinanceSpotAdapter) GetAccountData(order BotOrder) (*struct{}, *sharederrs.APIError) {
	//TODO
	return nil, nil
}

//GetPairLastPrice ..
func (a *BinanceSpotAdapter) GetPairLastPrice() (float64, *sharederrs.APIError) {
	//TODO
	return 0, nil
}

//CancelPairOrder ..
func (a *BinanceSpotAdapter) CancelPairOrder() *sharederrs.APIError {
	//TODO
	return nil
}

//CancelPairOrders ..
func (a *BinanceSpotAdapter) CancelPairOrders() *sharederrs.APIError {
	//TODO
	return nil
}

//GetPairOpenOrders ..
func (a *BinanceSpotAdapter) GetPairOpenOrders() ([]*struct{}, *sharederrs.APIError) {
	//TODO
	return nil, nil
}

//VerifyAPIKeys ..
func (a *BinanceSpotAdapter) VerifyAPIKeys() *sharederrs.APIError {
	//TODO
	return nil
}

//GetPairs ..
func (a *BinanceSpotAdapter) GetPairs() *sharederrs.APIError {
	client := binance.NewClient("", "")
	service := client.NewExchangeInfoService()
	res, err := service.Do(context.Background())
	if err != nil {
		return sharederrs.ServiceDisconnectedErr.
			M("error while connecting to ExchangeInfoService: " + err.Error()).SetTrace()
	}
	log.Println(res) //temp
	return nil
}
