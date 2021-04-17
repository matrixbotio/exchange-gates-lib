package matrixgates

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/adshao/go-binance/v2"
	sharederrs "github.com/matrixbotio/shared-errors"
)

//BinanceSpotAdapter - bot exchange adapter for BinanceSpot
type BinanceSpotAdapter struct {
	ExchangeAdapter
	binanceAPI *binance.Client
}

//NewBinanceSpotAdapter - create binance exchange adapter
func NewBinanceSpotAdapter() *ExchangeAdapter {
	adapter := newExchangeAdapter("Binance Spot", 1)
	return adapter
}

//Connect to exchange
func (a *BinanceSpotAdapter) Connect(credentials APICredentials) *sharederrs.APIError {
	switch credentials.Type {
	default:
		return sharederrs.DataInvalidErr
	case APICredentialsTypeKeypair:
		a.binanceAPI = binance.NewClient(credentials.Keypair.Public, credentials.Keypair.Secret)
		break
	}
	return a.ping()
}

//GetOrderData ..
func (a *BinanceSpotAdapter) GetOrderData(pairSymbol string, orderID int64) (*TradeEventData, *sharederrs.APIError) {
	tradeData := TradeEventData{
		OrderID: orderID,
	}
	//order status: NEW, PARTIALLY_FILLED, FILLED, CANCELED, PENDING_CANCEL, REJECTED, EXPIRED
	orderResponse, err := a.binanceAPI.NewGetOrderService().Symbol(pairSymbol).
		OrderID(orderID).Do(context.Background())

	if err != nil {
		if strings.Contains(err.Error(), "Order does not exist") {
			fmt.Println("[CHECK] order " + strconv.FormatInt(orderID, 10) + " doesn't exists")
			tradeData.Status = "UNKNOWN"
			return &tradeData, nil
		}
		return nil, sharederrs.ServiceReqFailedErr.SetMessage(err.Error()).SetTrace()
	}
	orderFilledQty, convErr := strconv.ParseFloat(orderResponse.ExecutedQuantity, 64)
	if convErr != nil {
		return nil, sharederrs.DataHandleErr.M("failed to parse order filled qty: " + convErr.Error()).SetTrace()
	}
	tradeData.OrderAwaitQty = orderFilledQty
	tradeData.Status = string(orderResponse.Status)
	return &tradeData, nil
}

//PlaceOrder ..
func (a *BinanceSpotAdapter) PlaceOrder(order BotOrder) (*CreateOrderResponse, *sharederrs.APIError) {
	//TODO
	return nil, nil
}

//GetAccountData - get account data ^ↀᴥↀ^
func (a *BinanceSpotAdapter) GetAccountData() (*AccountData, *sharederrs.APIError) {
	binanceAccountData, clientErr := a.binanceAPI.NewGetAccountService().Do(context.Background())
	if clientErr != nil {
		return nil, sharederrs.DataInvalidErr.
			M("failed to send request to trade, " + clientErr.Error()).SetTrace()
	}
	accountDataResult := AccountData{
		CanTrade: binanceAccountData.CanTrade,
	}
	balances := []Balance{}
	for _, binanceBalanceData := range binanceAccountData.Balances {
		//convert strings to float64
		balanceFree, convErr := strconv.ParseFloat(binanceBalanceData.Free, 64)
		if convErr != nil {
			balanceFree = 0
			//log.Println("failed to parse free balance: " + convErr.Error())
		}
		balanceLocked, convErr := strconv.ParseFloat(binanceBalanceData.Locked, 64)
		if convErr != nil {
			balanceLocked = 0
			//log.Println("failed to parse locked balance: " + convErr.Error())
		}
		balances = append(balances, Balance{
			Asset:  binanceBalanceData.Asset,
			Free:   balanceFree,
			Locked: balanceLocked,
		})
	}
	accountDataResult.Balances = balances
	return &accountDataResult, nil
}

//GetPairLastPrice - get pair last price ^ↀᴥↀ^
func (a *BinanceSpotAdapter) GetPairLastPrice(pairSymbol string) (float64, *sharederrs.APIError) {
	tickerService := a.binanceAPI.NewListPricesService()
	prices, srvErr := tickerService.Symbol(pairSymbol).Do(context.Background())
	if srvErr != nil {
		return 0, sharederrs.ServiceReqFailedErr.
			SetMessage("failed to request last price, " + srvErr.Error()).SetTrace()
	}
	//until just brute force. need to be done faster
	var price float64 = 0
	var parseErr error
	for _, p := range prices {
		if p.Symbol == pairSymbol {
			price, parseErr = strconv.ParseFloat(p.Price, 64)
			if parseErr != nil {
				return 0, sharederrs.DataHandleErr.
					M("failed to parse " + p.Price + " as float").SetTrace()
			}
			break
		}
	}
	return price, nil
}

//CancelPairOrder ..
func (a *BinanceSpotAdapter) CancelPairOrder(pairSymbol string, orderID int64) *sharederrs.APIError {
	//TODO
	return nil
}

//CancelPairOrders - cancel pair all orders
func (a *BinanceSpotAdapter) CancelPairOrders(pairSymbol string) *sharederrs.APIError {
	_, clientErr := a.binanceAPI.NewCancelOpenOrdersService().
		Symbol(pairSymbol).Do(context.Background())
	if clientErr != nil {
		//log.Println("failed to cancel all orders, " + clientErr.Error())
		//log.Println("let's try cancel orders manualy..")
		//handle error
		if strings.Contains(clientErr.Error(), "Unknown order sent") {
			/*canceling all orders failed,
			let's try to request a list of them and cancel them individually*/
			orders, err := a.GetPairOpenOrders(pairSymbol)
			if err != nil {
				// =(
				fmt.Println("[DEBUG] error while b.getOpenOrders(): " + err.Message)
				return err
			}
			if len(orders) == 0 {
				//orders already cancelled
				return nil
			}
			for _, order := range orders {
				err := a.CancelPairOrder(pairSymbol, order.OrderID)
				if err != nil {
					return err
				}
			}
			return nil
		}
		//log.Println("[DEBUG] service error: " + clientErr.Error())
		return sharederrs.ServiceReqFailedErr.SetMessage(clientErr.Error()).SetTrace()
	}
	return nil
}

//GetPairOpenOrders ..
func (a *BinanceSpotAdapter) GetPairOpenOrders(pairSymbol string) ([]*Order, *sharederrs.APIError) {
	//TODO
	return nil, nil
}

func (a *BinanceSpotAdapter) ping() *sharederrs.APIError {
	err := a.binanceAPI.NewPingService().Do(context.Background())
	if err != nil {
		return sharederrs.ServiceDisconnectedErr.
			SetMessage("failed to connect to the trade, please try again later").SetTrace()
	}
	return nil
}

//VerifyAPIKeys ..
func (a *BinanceSpotAdapter) VerifyAPIKeys(keyPublic, keySecret string) *sharederrs.APIError {
	accountService, err := a.binanceAPI.NewGetAccountService().Do(context.Background())
	if err != nil {
		return sharederrs.ServiceReqFailedErr.
			SetMessage(err.Error()).SetTrace()
	}
	if !accountService.CanTrade {
		return sharederrs.ServiceNoAccess.
			SetMessage("Your API key does not have permission to trade, change its restrictions").SetTrace()
	}
	return a.ping()
}

//GetPairs get all Binance pairs
func (a *BinanceSpotAdapter) GetPairs() ([]*ExchangePairData, *sharederrs.APIError) {
	service := a.binanceAPI.NewExchangeInfoService()
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
