package matrixgates

import (
	"context"
	"strconv"

	"github.com/adshao/go-binance/v2"
	sharederrs "github.com/matrixbotio/shared-errors"
)

//BinanceSpotAdapter - bot exchange adapter for BinanceSpot
type BinanceSpotAdapter struct {
	ExchangeAdapter
}

func NewBinanceSpotAdapter() *ExchangeAdapter {
	return newExchangeAdapter("Binance Spot", 1)
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

//GetPairs get all Binance pairs
func (a *BinanceSpotAdapter) GetPairs() ([]*ExchangePairData, *sharederrs.APIError) {
	client := binance.NewClient("", "")
	service := client.NewExchangeInfoService()
	res, err := service.Do(context.Background())
	if err != nil {
		return nil, sharederrs.ServiceDisconnectedErr.
			M("error while connecting to ExchangeInfoService: " + err.Error()).SetTrace()
	}

	pairs := []*ExchangePairData{}
	for _, symbolData := range res.Symbols {
		pairData := ExchangePairData{
			ExchangeID:     a.ExchangeID,
			BaseAsset:      symbolData.BaseAsset,
			BasePrecision:  symbolData.BaseAssetPrecision,
			QuoteAsset:     symbolData.QuoteAsset,
			QuotePrecision: symbolData.QuotePrecision,
			Status:         symbolData.Status,
			Symbol:         symbolData.Symbol,
			MinQty:         0.001,
			MaxQty:         99999.99,
			MinPrice:       0.000001,
			QtyStep:        0.001,
			PriceStep:      0.000001,
			AllowedMargin:  symbolData.IsMarginTradingAllowed,
			AllowedSpot:    symbolData.IsSpotTradingAllowed,
		}

		marketLotSizeFilter := symbolData.MarketLotSizeFilter()
		if marketLotSizeFilter != nil {
			minQtyRaw := marketLotSizeFilter.MinQuantity
			maxQtyRaw := marketLotSizeFilter.MaxQuantity
			minQty, err := strconv.ParseFloat(minQtyRaw, 64)
			if err != nil {
				return nil, sharederrs.DataHandleErr.SetMessage(err.Error())
			}
			if minQty == 0 {
				minQty = 0.001
			}

			pairData.MaxQty, err = strconv.ParseFloat(maxQtyRaw, 64)
			if err != nil {
				return nil, sharederrs.DataHandleErr.SetMessage(err.Error())
			}

			qtyStepRaw := symbolData.MarketLotSizeFilter().StepSize
			pairData.QtyStep, err = strconv.ParseFloat(qtyStepRaw, 64)
			if err != nil {
				return nil, sharederrs.DataHandleErr.SetMessage(err.Error())
			}
			if pairData.QtyStep == 0 {
				pairData.QtyStep = minQty
			}
		}

		priceFilter := symbolData.PriceFilter()
		if priceFilter != nil {
			//add max price?
			minPriceRaw := priceFilter.MinPrice
			pairData.MinPrice, err = strconv.ParseFloat(minPriceRaw, 64)
			if err != nil {
				return nil, sharederrs.DataHandleErr.SetMessage(err.Error())
			}
			if pairData.MinPrice == 0 {
				pairData.MinPrice = 0.000001
			}
			priceStepRaw := symbolData.PriceFilter().TickSize
			pairData.PriceStep, err = strconv.ParseFloat(priceStepRaw, 64)
			if err != nil {
				return nil, sharederrs.DataHandleErr.SetMessage(err.Error())
			}
			if pairData.PriceStep == 0 {
				pairData.PriceStep = pairData.MinPrice
			}
		}
		pairs = append(pairs, &pairData)
	}
	return pairs, nil
}
