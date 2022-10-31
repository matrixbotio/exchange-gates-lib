package binance

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/go-stack/stack"

	adp "github.com/matrixbotio/exchange-gates-lib/internal/adapters"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
	"github.com/matrixbotio/exchange-gates-lib/pkg/utils"
)

type adapter struct {
	ExchangeID int
	Name       string
	Tag        string

	binanceAPI *binance.Client
}

func New() adp.Adapter {
	stack.Caller(0)
	a := adapter{}
	a.Name = "Binance Spot"
	a.Tag = "binance-spot"
	a.ExchangeID = consts.ExchangeIDbinanceSpot
	return &a
}

func (a *adapter) GetTag() string {
	return a.Tag
}

func (a *adapter) GetID() int {
	return a.ExchangeID
}

func (a *adapter) GetName() string {
	return a.Name
}

// Connect to exchange
func (a *adapter) Connect(credentials pkgStructs.APICredentials) error {
	if credentials.Type != pkgStructs.APICredentialsTypeKeypair {
		return errors.New("invalid credentials to connect to Binance")
	}

	a.binanceAPI = binance.NewClient(credentials.Keypair.Public, credentials.Keypair.Secret)
	err := a.ping()
	if err != nil {
		return fmt.Errorf("ping binance: %w", err)
	}

	// sync time
	a.sync()
	return nil
}

func (a *adapter) sync() {
	a.binanceAPI.NewSetServerTimeService().Do(context.Background())
}

func (a *adapter) GetPrices() ([]structs.SymbolPrice, error) {
	prices, err := a.binanceAPI.NewListPricesService().Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get prices: %w", err)
	}

	r := []structs.SymbolPrice{}
	for _, priceData := range prices {
		pairPrice, err := strconv.ParseFloat(priceData.Price, 64)
		if err != nil {
			return nil, fmt.Errorf("parse price for symbol %q: %w", priceData.Symbol, err)
		}
		r = append(r, structs.SymbolPrice{Price: pairPrice, Symbol: priceData.Symbol})
	}
	return r, nil
}

func (a *adapter) getOrderFromService(
	pairSymbol string, orderID int64, clientOrderID string,
) (*binance.Order, error) {

	tradeData := structs.OrderData{}
	s := a.binanceAPI.NewGetOrderService().Symbol(pairSymbol)
	if orderID > 0 {
		tradeData.OrderID = orderID
		s.OrderID(orderID)
	} else {
		tradeData.ClientOrderID = clientOrderID
		s.OrigClientOrderID(clientOrderID)
	}

	orderResponse, err := s.Do(context.Background())
	if err != nil {
		if strings.Contains(err.Error(), "Order does not exist") {
			tradeData.Status = pkgStructs.OrderStatusUnknown
			return nil, nil
		}
		return nil, err
	}

	return orderResponse, nil
}

func (a *adapter) getOrderData(
	pairSymbol string, orderID int64, clientOrderID string,
) (structs.OrderData, error) {

	if orderID == 0 && clientOrderID == "" {
		return structs.OrderData{}, errors.New("orderID & client order ID is not set")
	}

	// get service & set order ID
	orderResponse, err := a.getOrderFromService(pairSymbol, orderID, clientOrderID)
	if err != nil {
		return structs.OrderData{}, err
	}

	return convertOrder(orderResponse)
}

// GetOrderData - get order data
func (a *adapter) GetOrderData(pairSymbol string, orderID int64) (structs.OrderData, error) {
	return a.getOrderData(pairSymbol, orderID, "")
}

// GetClientOrderData - get order data by client order ID
func (a *adapter) GetOrderByClientOrderID(pairSymbol string, clientOrderID string) (structs.OrderData, error) {
	return a.getOrderData(pairSymbol, 0, clientOrderID)
}

// PlaceOrder - place order on exchange
func (a *adapter) PlaceOrder(
	ctx context.Context,
	order pkgStructs.BotOrderAdjusted,
) (structs.CreateOrderResponse, error) {
	r := structs.CreateOrderResponse{}
	orderSide := binance.SideType(pkgStructs.OrderTypeBuy)
	switch order.Type {
	default:
		return r, errors.New("data invalid error: unknown strategy given for order")
	case pkgStructs.OrderTypeBuy:
		orderSide = binance.SideTypeBuy
	case pkgStructs.OrderTypeSell:
		orderSide = binance.SideTypeSell
	}

	a.sync() // sync client

	// setup order
	orderService := a.binanceAPI.NewCreateOrderService().Symbol(order.PairSymbol).
		Side(orderSide).Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeGTC).Quantity(order.Qty).
		Price(order.Price)

	// set order ID
	if order.ClientOrderID != "" {
		orderService.NewClientOrderID(order.ClientOrderID)
	}

	// place order
	orderRes, err := orderService.Do(ctx)
	if err != nil {
		return r, fmt.Errorf("create order: %w", err)
	}

	// parse qty & price from order response
	orderResOrigQty, convErr := strconv.ParseFloat(orderRes.OrigQuantity, 64)
	if convErr != nil {
		return r, fmt.Errorf("parse order origQty: %w", err)
	}
	orderResPrice, convErr := strconv.ParseFloat(orderRes.Price, 64)
	if convErr != nil {
		return r, fmt.Errorf("parse order price: %w", err)
	}

	return structs.CreateOrderResponse{
		OrderID:       orderRes.OrderID,
		ClientOrderID: orderRes.ClientOrderID,
		OrigQuantity:  orderResOrigQty,
		Price:         orderResPrice,
		Symbol:        orderRes.Symbol,
		Type:          a.getOrderType(orderRes.Side),
	}, nil
}

// convert order side to bot order type
func (a *adapter) getOrderType(orderSide binance.SideType) string {
	return strings.ToLower(string(orderSide))
}

// GetAccountData - get account data ^ↀᴥↀ^
func (a *adapter) GetAccountData() (structs.AccountData, error) {
	binanceAccountData, clientErr := a.binanceAPI.NewGetAccountService().Do(context.Background())
	if clientErr != nil {
		return structs.AccountData{}, errors.New("send request to trade, " + clientErr.Error())
	}
	accountDataResult := structs.AccountData{
		CanTrade: binanceAccountData.CanTrade,
	}

	balances := []structs.Balance{}
	for _, binanceBalanceData := range binanceAccountData.Balances {
		// convert strings to float64
		balanceFree, convErr := strconv.ParseFloat(binanceBalanceData.Free, 64)
		if convErr != nil {
			balanceFree = 0
		}
		balanceLocked, convErr := strconv.ParseFloat(binanceBalanceData.Locked, 64)
		if convErr != nil {
			balanceLocked = 0
		}
		if balanceFree != 0 || balanceLocked != 0 {
			balances = append(balances, structs.Balance{
				Asset:  binanceBalanceData.Asset,
				Free:   balanceFree,
				Locked: balanceLocked,
			})
		}
	}
	accountDataResult.Balances = balances
	return accountDataResult, nil
}

// GetPairLastPrice - get pair last price ^ↀᴥↀ^
func (a *adapter) GetPairLastPrice(pairSymbol string) (float64, error) {
	tickerService := a.binanceAPI.NewListPricesService()
	prices, err := tickerService.Symbol(pairSymbol).Do(context.Background())
	if err != nil {
		return 0, fmt.Errorf("last price: %w", err)
	}

	// until just brute force. need to be done faster
	var price float64 = 0
	var parseErr error
	for _, p := range prices {
		if p.Symbol == pairSymbol {
			price, parseErr = strconv.ParseFloat(p.Price, 64)
			if parseErr != nil {
				return 0, fmt.Errorf("parse price %q: %w", p.Price, err)
			}
			break
		}
	}
	return price, nil
}

// CancelPairOrder - cancel one exchange pair order by ID
func (a *adapter) CancelPairOrder(pairSymbol string, orderID int64, ctx context.Context) error {
	_, clientErr := a.binanceAPI.NewCancelOrderService().Symbol(pairSymbol).
		OrderID(orderID).Do(ctx)
	if clientErr != nil {
		if !a.isErrorAboutUnknownOrder(clientErr) {
			return clientErr
		}
	}
	return nil
}

// CancelPairOrder - cancel one exchange pair order by client order ID
func (a *adapter) CancelPairOrderByClientOrderID(
	pairSymbol string,
	clientOrderID string,
	ctx context.Context,
) error {
	_, clientErr := a.binanceAPI.NewCancelOrderService().Symbol(pairSymbol).
		OrigClientOrderID(clientOrderID).Do(ctx)
	if clientErr != nil {
		if !a.isErrorAboutUnknownOrder(clientErr) {
			return clientErr
		}
	}
	return nil
}

func (a *adapter) isErrorAboutUnknownOrder(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "Unknown order sent")
}

// GetPairData - get pair data & limits
func (a *adapter) GetPairData(pairSymbol string) (structs.ExchangePairData, error) {
	exchangeInfo, err := a.binanceAPI.NewExchangeInfoService().Symbol(pairSymbol).Do(context.Background())
	if err != nil {
		return structs.ExchangePairData{}, err
	}

	// find pairSymbol
	for _, symbolData := range exchangeInfo.Symbols {
		return a.getExchangePairData(symbolData)
	}

	return structs.ExchangePairData{}, errors.New("data for " + pairSymbol + " pair not found")
}

//GetPairOpenOrders - get open orders array
func (a *adapter) GetPairOpenOrders(pairSymbol string) ([]structs.OrderData, error) {
	ordersRaw, err := a.binanceAPI.NewListOpenOrdersService().Symbol(pairSymbol).Do(context.Background())
	if err != nil {
		return nil, err
	}

	return convertOrders(ordersRaw)
}

func (a *adapter) ping() error {
	var err error
	for attemptNumber := 1; attemptNumber <= consts.PingRetryAttempts; attemptNumber++ {
		err := a.binanceAPI.NewPingService().Do(context.Background())
		if err == nil {
			return nil
		}

		time.Sleep(consts.PingRetryWaitTime)
	}

	return fmt.Errorf("ping exchange: %w", err)
}

// VerifyAPIKeys - create new exchange client & attempt to get account data
func (a *adapter) VerifyAPIKeys(keyPublic, keySecret string) error {
	newClient := binance.NewClient(keyPublic, keySecret)
	accountService, err := newClient.NewGetAccountService().Do(context.Background())
	if err != nil {
		return fmt.Errorf("invalid api key: %w", err)
	}
	if !accountService.CanTrade {
		return errors.New("your API key does not have permission to trade, change its restrictions")
	}
	return nil
}

// convert binance.Symbol to ExchangePairData
func (a *adapter) getExchangePairData(symbolData binance.Symbol) (structs.ExchangePairData, error) {
	pairData := structs.ExchangePairData{
		ExchangeID:     a.ExchangeID,
		BaseAsset:      symbolData.BaseAsset,
		BasePrecision:  symbolData.BaseAssetPrecision,
		QuoteAsset:     symbolData.QuoteAsset,
		QuotePrecision: symbolData.QuotePrecision,
		Status:         symbolData.Status,
		Symbol:         symbolData.Symbol,
		MinQty:         consts.PairDefaultMinQty,
		MaxQty:         consts.PairDefaultMaxQty,
		MinDeposit:     consts.PairMinDeposit,
		MinPrice:       consts.PairDefaultMinPrice,
		QtyStep:        consts.PairDefaultQtyStep,
		PriceStep:      consts.PairDefaultPriceStep,
		AllowedMargin:  symbolData.IsMarginTradingAllowed,
		AllowedSpot:    symbolData.IsSpotTradingAllowed,
	}

	var optionalErr error
	err := binanceParseLotSizeFilter(&symbolData, &pairData)
	if err != nil {
		optionalErr = err
	}

	err = binanceParsePriceFilter(&symbolData, &pairData)
	if err != nil {
		optionalErr = err
	}

	err = binanceParseMinNotionalFilter(&symbolData, &pairData)
	if err != nil {
		optionalErr = err
	}

	return pairData, optionalErr
}

func binanceParseMinNotionalFilter(symbolData *binance.Symbol, pairData *structs.ExchangePairData) error {
	var err error
	minNotionalFilter := symbolData.MinNotionalFilter()
	pairData.OriginalMinDeposit, err = strconv.ParseFloat(minNotionalFilter.MinNotional, 64)
	if err != nil {
		return fmt.Errorf("parse float: %w", err)
	}
	pairData.MinDeposit = utils.RoundMinDeposit(pairData.OriginalMinDeposit)
	return nil
}

func binanceParsePriceFilter(symbolData *binance.Symbol, pairData *structs.ExchangePairData) error {
	var err error
	priceFilter := symbolData.PriceFilter()
	if priceFilter == nil {
		return errors.New("get price filter for symbol data")
	}

	minPriceRaw := priceFilter.MinPrice
	pairData.MinPrice, err = strconv.ParseFloat(minPriceRaw, 64)
	if err != nil {
		return fmt.Errorf("data handle error: %w", err)
	}
	if pairData.MinPrice == 0 {
		pairData.MinPrice = consts.PairDefaultMinPrice
	}

	priceStepRaw := priceFilter.TickSize
	pairData.PriceStep, err = strconv.ParseFloat(priceStepRaw, 64)
	if err != nil {
		return fmt.Errorf("data handle error: %w", err)
	}
	if pairData.PriceStep == 0 {
		pairData.PriceStep = pairData.MinPrice
	}
	return nil
}

func binanceParseLotSizeFilter(symbolData *binance.Symbol, pairData *structs.ExchangePairData) error {
	lotSizeFilter := symbolData.LotSizeFilter()
	if lotSizeFilter == nil {
		return fmt.Errorf("lot size filter for symbol %q not found", symbolData.Symbol)
	}
	minQtyRaw := lotSizeFilter.MinQuantity
	maxQtyRaw := lotSizeFilter.MaxQuantity

	var err error
	pairData.MinQty, err = strconv.ParseFloat(minQtyRaw, 64)
	if err != nil {
		return fmt.Errorf("parse pair min qty: %w", err)
	}
	if pairData.MinQty == 0 {
		pairData.MinQty = consts.PairDefaultMinQty
	}

	pairData.MaxQty, err = strconv.ParseFloat(maxQtyRaw, 64)
	if err != nil {
		return fmt.Errorf("parse pair max qty: %w", err)
	}

	qtyStepRaw := lotSizeFilter.StepSize
	pairData.QtyStep, err = strconv.ParseFloat(qtyStepRaw, 64)
	if err != nil {
		return fmt.Errorf("parse pair qty step: %w", err)
	}
	if pairData.QtyStep == 0 {
		pairData.QtyStep = pairData.MinQty
	}
	return nil
}

// GetPairs get all Binance pairs
func (a *adapter) GetPairs() ([]structs.ExchangePairData, error) {
	service := a.binanceAPI.NewExchangeInfoService()
	res, err := service.Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	var lastError error
	pairs := []structs.ExchangePairData{}
	for _, symbolData := range res.Symbols {
		pairData, err := a.getExchangePairData(symbolData)
		if err != nil {
			lastError = err
		} else {
			pairs = append(pairs, pairData)
		}
	}
	return pairs, lastError
}

func (a *adapter) GetPairOrdersHistory(task structs.GetOrdersHistoryTask) ([]structs.OrderData, error) {
	// check data
	if task.PairSymbol == "" {
		return nil, errors.New("pair symbol is not set")
	}
	if task.StartTime == 0 {
		return nil, errors.New("start time is not set")
	}

	// create request service
	service := a.binanceAPI.NewListOrdersService().StartTime(task.StartTime).
		Symbol(task.PairSymbol)
	if task.EndTime > 0 {
		service.EndTime(task.EndTime)
	}

	// set context
	if task.Ctx == nil {
		task.Ctx = context.Background()
	}

	// send request
	ordersRaw, err := service.Do(task.Ctx)
	if err != nil {
		return nil, fmt.Errorf("get orders history: %w", err)
	}
	if ordersRaw == nil {
		return nil, errors.New("request orders history: orders not set")
	}

	// convert orders
	return convertOrders(ordersRaw)
}

// GetPairBalance - get pair balance: ticker, quote asset balance for pair symbol
func (a *adapter) GetPairBalance(pair structs.PairSymbolData) (structs.PairBalance, error) {
	accountData, err := a.GetAccountData()
	if err != nil {
		return structs.PairBalance{}, err
	}

	pairBalanceData := structs.PairBalance{}
	for _, balanceData := range accountData.Balances {
		if balanceData.Asset == pair.BaseTicker {
			// base asset found
			pairBalanceData.BaseAsset = &structs.AssetBalance{
				Ticker: balanceData.Asset,
				Free:   balanceData.Free,
				Locked: balanceData.Locked,
			}
		}
		if balanceData.Asset == pair.QuoteTicker {
			// quote asset found
			pairBalanceData.QuoteAsset = &structs.AssetBalance{
				Ticker: balanceData.Asset,
				Free:   balanceData.Free,
				Locked: balanceData.Locked,
			}
		}
		if pairBalanceData.BaseAsset != nil && pairBalanceData.QuoteAsset != nil {
			// found
			break
		}
	}
	if pairBalanceData.BaseAsset == nil {
		pairBalanceData.BaseAsset = &structs.AssetBalance{
			Ticker: pair.BaseTicker,
			Free:   0,
			Locked: 0,
		}
	}
	if pairBalanceData.QuoteAsset == nil {
		pairBalanceData.QuoteAsset = &structs.AssetBalance{
			Ticker: pair.QuoteTicker,
			Free:   0,
			Locked: 0,
		}
	}
	return pairBalanceData, nil
}

/*
                    _
                   | |
__      _____  _ __| | _____ _ __ ___
\ \ /\ / / _ \| '__| |/ / _ \ '__/ __|
 \ V  V / (_) | |  |   <  __/ |  \__ \
  \_/\_/ \___/|_|  |_|\_\___|_|  |___/

*/

// PriceWorkerBinance - MarketDataWorker for binance
type PriceWorkerBinance struct {
	workers.PriceWorker
}

// GetPriceWorker - create new market data worker
func (a *adapter) GetPriceWorker(callback workers.PriceEventCallback) workers.IPriceWorker {
	w := PriceWorkerBinance{}
	w.PriceWorker.ExchangeTag = a.Tag
	w.PriceWorker.HandleEventCallback = callback
	return &w
}

func (w *PriceWorkerBinance) handlePriceEvent(event *binance.WsBookTickerEvent) {
	if event == nil {
		return
	}

	eventAsk, convErr := strconv.ParseFloat(event.BestAskPrice, 64)
	if convErr != nil {
		return // ignore event
	}

	eventBid, convErr := strconv.ParseFloat(event.BestBidPrice, 64)
	if convErr != nil {
		return // ignore event
	}

	w.HandleEventCallback(workers.PriceEvent{
		ExchangeTag: w.ExchangeTag,
		Symbol:      event.Symbol,
		Ask:         eventAsk,
		Bid:         eventBid,
	})
}

// SubscribeToPriceEvents - websocket subscription to change quotes and ask-, bid-qty on the exchange
// returns map[pair symbol] -> worker channels
func (w *PriceWorkerBinance) SubscribeToPriceEvents(
	pairSymbols []string,
	errorHandler func(err error),
) (map[string]pkgStructs.WorkerChannels, error) {
	result := map[string]pkgStructs.WorkerChannels{}

	// event handler func
	w.WsChannels = new(pkgStructs.WorkerChannels)

	var openWsErr error
	for _, pairSymbol := range pairSymbols {
		newChannels := pkgStructs.WorkerChannels{}
		newChannels.WsDone, newChannels.WsStop, openWsErr = binance.WsBookTickerServe(pairSymbol, w.handlePriceEvent, errorHandler)
		if openWsErr != nil {
			return result, fmt.Errorf("subscribe to %q price: %w", pairSymbol, openWsErr)
		}

		result[pairSymbol] = newChannels
	}

	return result, nil
}

// CandleWorkerBinance - MarketDataWorker for binance
type CandleWorkerBinance struct {
	workers.CandleWorker
}

// GetCandleWorker - create new market candle worker
func (a *adapter) GetCandleWorker() workers.ICandleWorker {
	w := CandleWorkerBinance{}
	w.ExchangeTag = a.GetTag()
	return &w
}

// SubscribeToCandleEvents - websocket subscription to candles on the exchange
func (w *CandleWorkerBinance) SubscribeToCandleEvents(
	pairSymbol string,
	eventCallback func(event workers.CandleEvent),
	errorHandler func(err error),
) error {
	wsCandleHandler := func(event *binance.WsKlineEvent) {
		if event != nil {
			// fix endTime
			if strings.HasSuffix(strconv.FormatInt(event.Kline.EndTime, 10), "999") {
				event.Kline.EndTime -= 59999
			}

			wEvent := workers.CandleEvent{
				Symbol: event.Symbol,
				Candle: workers.CandleData{
					StartTime: event.Kline.StartTime,
					EndTime:   event.Kline.EndTime,
					Interval:  event.Kline.Interval,
				},
				Time: event.Time,
			}

			errs := make([]error, 5)
			wEvent.Candle.Open, errs[0] = strconv.ParseFloat(event.Kline.Open, 64)
			wEvent.Candle.Close, errs[1] = strconv.ParseFloat(event.Kline.Close, 64)
			wEvent.Candle.High, errs[2] = strconv.ParseFloat(event.Kline.High, 64)
			wEvent.Candle.Low, errs[3] = strconv.ParseFloat(event.Kline.Low, 64)
			wEvent.Candle.Volume, errs[4] = strconv.ParseFloat(event.Kline.Volume, 64)
			if utils.LogNotNilError(errs) {
				return
			}

			eventCallback(wEvent)
		}
	}
	wsErrHandler := func(err error) {
		errorHandler(err)
	}
	var openWsErr error
	w.WsChannels = new(pkgStructs.WorkerChannels)
	w.WsChannels.WsDone, w.WsChannels.WsStop, openWsErr = binance.WsKlineServe(
		pairSymbol,             // symbol
		consts.CandlesInterval, // interval
		wsCandleHandler,        // event handler
		wsErrHandler,           // error handler
	)
	return openWsErr
}

// TradeEventWorkerBinance - TradeEventWorker for binance
type TradeEventWorkerBinance struct {
	workers.TradeEventWorker
}

// GetTradeEventsWorker - create new market candle worker
func (a *adapter) GetTradeEventsWorker() workers.ITradeEventWorker {
	w := TradeEventWorkerBinance{}
	w.ExchangeTag = a.GetTag()
	return &w
}

// SubscribeToTradeEvents - websocket subscription to change trade candles on the exchange
func (w *TradeEventWorkerBinance) SubscribeToTradeEvents(
	symbol string,
	eventCallback func(event workers.TradeEvent),
	errorHandler func(err error),
) error {

	wsErrHandler := func(err error) {
		errorHandler(err)
	}

	wsTradeHandler := func(event *binance.WsTradeEvent) {
		if event != nil {
			// fix event.Time
			if strings.HasSuffix(strconv.FormatInt(event.Time, 10), "999") {
				event.Time++
			}
			wEvent := workers.TradeEvent{
				ID:            event.TradeID,
				Time:          event.Time,
				Symbol:        event.Symbol,
				ExchangeTag:   w.ExchangeTag,
				BuyerOrderID:  event.BuyerOrderID,
				SellerOrderID: event.SellerOrderID,
			}
			errs := make([]error, 2)
			wEvent.Price, errs[0] = strconv.ParseFloat(event.Price, 64)
			wEvent.Quantity, errs[0] = strconv.ParseFloat(event.Quantity, 64)
			if utils.LogNotNilError(errs) {
				return
			}
			eventCallback(wEvent)
		}
	}

	var err error
	w.WsChannels = new(pkgStructs.WorkerChannels)
	w.WsChannels.WsDone, w.WsChannels.WsStop, err = binance.WsTradeServe(symbol, wsTradeHandler, wsErrHandler)
	if err != nil {
		return fmt.Errorf("subscribe to trade events: %w", err)
	}
	return nil
}
