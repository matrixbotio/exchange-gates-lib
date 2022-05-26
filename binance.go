package matrixgates

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/adshao/go-binance/v2"
	"github.com/go-stack/stack"
	"github.com/matrixbotio/exchange-gates-lib/workers"
)

// BinanceSpotAdapter - bot exchange adapter for BinanceSpot
type BinanceSpotAdapter struct {
	ExchangeAdapter
	binanceAPI *binance.Client
}

// NewBinanceSpotAdapter - create binance exchange adapter
func NewBinanceSpotAdapter() *BinanceSpotAdapter {
	stack.Caller(0)
	a := BinanceSpotAdapter{}
	a.Name = "Binance Spot"
	a.Tag = "binance-spot"
	a.ExchangeID = exchangeIDbinanceSpot
	return &a
}

// Connect to exchange
func (a *BinanceSpotAdapter) Connect(credentials APICredentials) error {
	if credentials.Type != APICredentialsTypeKeypair {
		return errors.New("invalid credentials to connect to Binance")
	}

	a.binanceAPI = binance.NewClient(credentials.Keypair.Public, credentials.Keypair.Secret)
	err := a.ping()
	if err != nil {
		return errors.New("failed to ping binance: " + err.Error())
	}

	// sync time
	a.sync()
	return nil
}

func (a *BinanceSpotAdapter) sync() {
	a.binanceAPI.NewSetServerTimeService().Do(context.Background())
}

// GetOrderData - get order data
func (a *BinanceSpotAdapter) GetOrderData(pairSymbol string, orderID int64) (*OrderData, error) {
	orderResponse, err := a.binanceAPI.NewGetOrderService().Symbol(pairSymbol).
		OrderID(orderID).Do(context.Background())

	tradeData := OrderData{
		OrderID: orderID,
	}

	if err != nil {
		if strings.Contains(err.Error(), "Order does not exist") {
			tradeData.Status = OrderStatusUnknown
			return &tradeData, nil
		}
		return nil, errors.New("service request failed: " + err.Error() + GetTrace())
	}

	return a.convertOrder(orderResponse)
}

// from binance format to our bot order type format
func (a *BinanceSpotAdapter) convertOrderSide(orderSide binance.SideType) (string, error) {
	switch orderSide {
	default:
		return "", errors.New("unknown order type: " + string(orderSide))
	case binance.SideTypeBuy:
		return "buy", nil
	case binance.SideTypeSell:
		return "sell", nil
	}
}

// PlaceOrder - place order on exchange
func (a *BinanceSpotAdapter) PlaceOrder(ctx context.Context, order BotOrderAdjusted) (*CreateOrderResponse, error) {
	var orderSide binance.SideType
	switch order.Type {
	default:
		return nil, errors.New("data invalid error: unknown strategy given for order, stack: " + GetTrace())
	case "buy":
		orderSide = binance.SideTypeBuy
	case "sell":
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
		return nil, errors.New("service request failed: failed to create order, " + err.Error() + ", stack: " + GetTrace())
	}

	// parse qty & price from order response
	orderResOrigQty, convErr := strconv.ParseFloat(orderRes.OrigQuantity, 64)
	if convErr != nil {
		return nil, errors.New("data handle error: failed to parse order origQty, " + convErr.Error() + ", stack: " + GetTrace())
	}
	orderResPrice, convErr := strconv.ParseFloat(orderRes.Price, 64)
	if convErr != nil {
		return nil, errors.New("data handle error: failed to parse order price, " + convErr.Error() + ", stack: " + GetTrace())
	}

	return &CreateOrderResponse{
		OrderID:       orderRes.OrderID,
		ClientOrderID: orderRes.ClientOrderID,
		OrigQuantity:  orderResOrigQty,
		Price:         orderResPrice,
		Symbol:        orderRes.Symbol,
		Type:          a.getOrderType(orderRes.Side),
	}, nil
}

// convert order side to bot order type
func (a *BinanceSpotAdapter) getOrderType(orderSide binance.SideType) string {
	return strings.ToLower(string(orderSide))
}

// GetAccountData - get account data ^ↀᴥↀ^
func (a *BinanceSpotAdapter) GetAccountData() (*AccountData, error) {
	binanceAccountData, clientErr := a.binanceAPI.NewGetAccountService().Do(context.Background())
	if clientErr != nil {
		return nil, errors.New("data invalid error: failed to send request to trade, " + clientErr.Error() + ", stack: " + GetTrace())
	}
	accountDataResult := AccountData{
		CanTrade: binanceAccountData.CanTrade,
	}
	balances := []Balance{}
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
			balances = append(balances, Balance{
				Asset:  binanceBalanceData.Asset,
				Free:   balanceFree,
				Locked: balanceLocked,
			})
		}
	}
	accountDataResult.Balances = balances
	return &accountDataResult, nil
}

// GetPairLastPrice - get pair last price ^ↀᴥↀ^
func (a *BinanceSpotAdapter) GetPairLastPrice(pairSymbol string) (float64, error) {
	tickerService := a.binanceAPI.NewListPricesService()
	prices, srvErr := tickerService.Symbol(pairSymbol).Do(context.Background())
	if srvErr != nil {
		return 0, errors.New("service request failed: failed to request last price, " +
			srvErr.Error() + ", stack: " + GetTrace())
	}
	// until just brute force. need to be done faster
	var price float64 = 0
	var parseErr error
	for _, p := range prices {
		if p.Symbol == pairSymbol {
			price, parseErr = strconv.ParseFloat(p.Price, 64)
			if parseErr != nil {
				return 0, errors.New("data handle error: failed to parse " + p.Price +
					" as float, stack: " + GetTrace())
			}
			break
		}
	}
	return price, nil
}

// CancelPairOrder - cancel one exchange pair order by ID
func (a *BinanceSpotAdapter) CancelPairOrder(pairSymbol string, orderID int64) error {
	_, clientErr := a.binanceAPI.NewCancelOrderService().Symbol(pairSymbol).
		OrderID(orderID).Do(context.Background())
	if clientErr != nil {
		if !a.isErrorAboutUnknownOrder(clientErr) {
			return clientErr
		}
	}
	return nil
}

// CancelPairOrders - cancel pair all orders
func (a *BinanceSpotAdapter) CancelPairOrders(pairSymbol string) error {
	_, clientErr := a.binanceAPI.NewCancelOpenOrdersService().
		Symbol(pairSymbol).Do(context.Background())
	if clientErr != nil {
		// handle error
		if a.isErrorAboutUnknownOrder(clientErr) {
			/* canceling all orders failed,
			let's try to request a list of them and cancel them individually*/
			orders, err := a.GetPairOpenOrders(pairSymbol)
			if err != nil {
				// =(
				return err
			}
			if len(orders) == 0 {
				// orders already cancelled
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
		return errors.New("service request failed: " + clientErr.Error() +
			", stack: " + GetTrace())
	}
	return nil
}

func (a *BinanceSpotAdapter) isErrorAboutUnknownOrder(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "Unknown order sent")
}

// GetPairData - get pair data & limits
func (a *BinanceSpotAdapter) GetPairData(pairSymbol string) (*ExchangePairData, error) {
	exchangeInfo, err := a.binanceAPI.NewExchangeInfoService().Symbol(pairSymbol).Do(context.Background())
	if err != nil {
		return nil, err
	}

	// find pairSymbol
	for _, symbolData := range exchangeInfo.Symbols {
		return a.getExchangePairData(symbolData)
	}

	return nil, errors.New("data for " + pairSymbol + " pair not found")
}

//GetPairOpenOrders - get open orders array
func (a *BinanceSpotAdapter) GetPairOpenOrders(pairSymbol string) ([]*OrderData, error) {
	ordersRaw, err := a.binanceAPI.NewListOpenOrdersService().Symbol(pairSymbol).Do(context.Background())
	if err != nil {
		return nil, err
	}

	return a.convertOrders(ordersRaw)
}

func (a *BinanceSpotAdapter) ping() error {
	err := a.binanceAPI.NewPingService().Do(context.Background())
	if err != nil {
		return errors.New("service disconnected: failed to connect to the exchange, please try again later, stack: " + GetTrace())
	}
	return nil
}

// VerifyAPIKeys - create new exchange client & attempt to get account data
func (a *BinanceSpotAdapter) VerifyAPIKeys(keyPublic, keySecret string) error {
	newClient := binance.NewClient(keyPublic, keySecret)
	accountService, err := newClient.NewGetAccountService().Do(context.Background())
	if err != nil {
		return errors.New("invalid api key: " + err.Error())
	}
	if !accountService.CanTrade {
		return errors.New("service no access: Your API key does not have permission to trade, change its restrictions")
	}
	return nil
}

// convert binance.Symbol to ExchangePairData
func (a *BinanceSpotAdapter) getExchangePairData(symbolData binance.Symbol) (*ExchangePairData, error) {
	pairData := ExchangePairData{
		ExchangeID:     a.ExchangeID,
		BaseAsset:      symbolData.BaseAsset,
		BasePrecision:  symbolData.BaseAssetPrecision,
		QuoteAsset:     symbolData.QuoteAsset,
		QuotePrecision: symbolData.QuotePrecision,
		Status:         symbolData.Status,
		Symbol:         symbolData.Symbol,
		MinQty:         PairDefaultMinQty,
		MaxQty:         PairDefaultMaxQty,
		MinDeposit:     PairMinDeposit,
		MinPrice:       PairDefaultMinPrice,
		QtyStep:        PairDefaultQtyStep,
		PriceStep:      PairDefaultPriceStep,
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

	return &pairData, optionalErr
}

func binanceParseMinNotionalFilter(symbolData *binance.Symbol, pairData *ExchangePairData) error {
	var err error
	minNotionalFilter := symbolData.MinNotionalFilter()
	pairData.OriginalMinDeposit, err = strconv.ParseFloat(minNotionalFilter.MinNotional, 64)
	if err != nil {
		return errors.New("failed to parse float: " + err.Error())
	}
	pairData.MinDeposit = RoundMinDeposit(pairData.OriginalMinDeposit)
	return nil
}

func binanceParsePriceFilter(symbolData *binance.Symbol, pairData *ExchangePairData) error {
	var err error
	priceFilter := symbolData.PriceFilter()
	if priceFilter == nil {
		return errors.New("failed to get price filter for symbol data")
	}
	minPriceRaw := priceFilter.MinPrice
	pairData.MinPrice, err = strconv.ParseFloat(minPriceRaw, 64)
	if err != nil {
		return errors.New("data handle error: " + err.Error())
	}
	if pairData.MinPrice == 0 {
		pairData.MinPrice = PairDefaultMinPrice
	}
	priceStepRaw := priceFilter.TickSize
	pairData.PriceStep, err = strconv.ParseFloat(priceStepRaw, 64)
	if err != nil {
		return errors.New("data handle error: " + err.Error())
	}
	if pairData.PriceStep == 0 {
		pairData.PriceStep = pairData.MinPrice
	}
	return nil
}

func binanceParseLotSizeFilter(symbolData *binance.Symbol, pairData *ExchangePairData) error {
	lotSizeFilter := symbolData.LotSizeFilter()
	if lotSizeFilter == nil {
		return errors.New("failed to get lot size filter for symbol data: " + symbolData.Symbol)
	}
	minQtyRaw := lotSizeFilter.MinQuantity
	maxQtyRaw := lotSizeFilter.MaxQuantity

	var err error
	pairData.MinQty, err = strconv.ParseFloat(minQtyRaw, 64)
	if err != nil {
		return errors.New("failed to parse pair min qty: " + err.Error())
	}
	if pairData.MinQty == 0 {
		pairData.MinQty = PairDefaultMinQty
	}

	pairData.MaxQty, err = strconv.ParseFloat(maxQtyRaw, 64)
	if err != nil {
		return errors.New("failed to parse pair max qty: " + err.Error())
	}

	qtyStepRaw := lotSizeFilter.StepSize
	pairData.QtyStep, err = strconv.ParseFloat(qtyStepRaw, 64)
	if err != nil {
		return errors.New("failed to parse pair qty step: " + err.Error())
	}
	if pairData.QtyStep == 0 {
		pairData.QtyStep = pairData.MinQty
	}
	return nil
}

// GetPairs get all Binance pairs
func (a *BinanceSpotAdapter) GetPairs() ([]*ExchangePairData, error) {
	service := a.binanceAPI.NewExchangeInfoService()
	res, err := service.Do(context.Background())
	if err != nil {
		return nil, errors.New("service disconnected: error while connecting to ExchangeInfoService: " + err.Error() + ", stack: " + GetTrace())
	}

	var lastError error
	pairs := []*ExchangePairData{}
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

func (a *BinanceSpotAdapter) GetPairOrdersHistory(task GetOrdersHistoryTask) ([]*OrderData, error) {
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
		return nil, errors.New("failed to get orders history: " + err.Error())
	}
	if ordersRaw == nil {
		return nil, errors.New("failed to request orders history: orders is nil")
	}

	// convert orders
	return a.convertOrders(ordersRaw)
}

func (a *BinanceSpotAdapter) parseOrderOriginalQty(orderRaw *binance.Order) (float64, error) {
	awaitQty, err := strconv.ParseFloat(orderRaw.OrigQuantity, 32)
	if err != nil {
		return 0, errors.New("failed to parse order original qty: " + err.Error())
	}
	return awaitQty, nil
}

func (a *BinanceSpotAdapter) parseOrderExecutedQty(orderRaw *binance.Order) (float64, error) {
	filledQty, err := strconv.ParseFloat(orderRaw.ExecutedQuantity, 32)
	if err != nil {
		return 0, errors.New("failed to parse order executed qty: " + err.Error())
	}
	return filledQty, nil
}

func (a *BinanceSpotAdapter) parseOrderPrice(orderRaw *binance.Order) (float64, error) {
	price, err := strconv.ParseFloat(orderRaw.Price, 32)
	if err != nil {
		return 0, errors.New("failed to parse order price: " + err.Error())
	}
	return price, nil
}

// converting the order from binance to our format
func (a *BinanceSpotAdapter) convertOrder(orderRaw *binance.Order) (*OrderData, error) {
	awaitQty, err := a.parseOrderOriginalQty(orderRaw)
	if err != nil {
		return nil, err
	}

	filledQty, err := a.parseOrderExecutedQty(orderRaw)
	if err != nil {
		return nil, err
	}

	price, err := a.parseOrderPrice(orderRaw)
	if err != nil {
		return nil, err
	}

	orderType, err := a.convertOrderSide(orderRaw.Side)
	if err != nil {
		return nil, err
	}

	return &OrderData{
		OrderID:       orderRaw.OrderID,
		ClientOrderID: orderRaw.ClientOrderID,
		Status:        string(orderRaw.Status),
		AwaitQty:      awaitQty,
		FilledQty:     filledQty,
		Price:         price,
		Symbol:        orderRaw.Symbol,
		Type:          orderType,
	}, nil
}

func (a *BinanceSpotAdapter) convertOrders(ordersRaw []*binance.Order) ([]*OrderData, error) {
	orders := []*OrderData{}
	for _, orderRaw := range ordersRaw {
		order, err := a.convertOrder(orderRaw)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}
	return orders, nil
}

// GetPairBalance - get pair balance: ticker, quote asset balance for pair symbol
func (a *BinanceSpotAdapter) GetPairBalance(pair PairSymbolData) (*PairBalance, error) {
	accountData, err := a.GetAccountData()
	if err != nil {
		return nil, err
	}

	pairBalanceData := &PairBalance{}
	for _, balanceData := range accountData.Balances {
		if balanceData.Asset == pair.BaseTicker {
			// base asset found
			pairBalanceData.BaseAsset = &AssetBalance{
				Ticker: balanceData.Asset,
				Free:   balanceData.Free,
				Locked: balanceData.Locked,
			}
		}
		if balanceData.Asset == pair.QuoteTicker {
			// quote asset found
			pairBalanceData.QuoteAsset = &AssetBalance{
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
		pairBalanceData.BaseAsset = &AssetBalance{
			Ticker: pair.BaseTicker,
			Free:   0,
			Locked: 0,
		}
	}
	if pairBalanceData.QuoteAsset == nil {
		pairBalanceData.QuoteAsset = &AssetBalance{
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
func (a *BinanceSpotAdapter) GetPriceWorker(callback workers.PriceEventCallback) workers.IPriceWorker {
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
		Symbol: event.Symbol,
		Ask:    eventAsk,
		Bid:    eventBid,
	})
}

// SubscribeToPriceEvents - websocket subscription to change quotes and ask-, bid-qty on the exchange
// returns map[pair symbol] -> worker channels
func (w *PriceWorkerBinance) SubscribeToPriceEvents(
	pairSymbols []string,
	eventCallback workers.PriceEventCallback,
	errorHandler func(err error),
) (map[string]workers.WorkerChannels, error) {
	result := map[string]workers.WorkerChannels{}

	// event handler func
	w.WsChannels = new(workers.WorkerChannels)

	var openWsErr error
	for _, pairSymbol := range pairSymbols {
		newChannels := workers.WorkerChannels{}
		newChannels.WsDone, newChannels.WsStop, openWsErr = binance.WsBookTickerServe(pairSymbol, w.handlePriceEvent, errorHandler)
		if openWsErr != nil {
			return result, errors.New("failed to subscribe to `" + pairSymbol + "` price: " + openWsErr.Error())
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
func (a *BinanceSpotAdapter) GetCandleWorker() workers.ICandleWorker {
	w := CandleWorkerBinance{}
	w.ExchangeTag = a.GetTag()
	return &w
}

// SubscribeToCandleEvents - websocket subscription to change trade candles on the exchange
func (w *CandleWorkerBinance) SubscribeToCandleEvents(
	pairSymbols []string,
	eventCallback func(event workers.CandleEvent),
	errorHandler func(err error),
) error {
	symbolIntervalsMap := make(map[string]string)
	for _, symbol := range pairSymbols {
		symbolIntervalsMap[symbol] = candlesInterval
	}

	wsCandleHandler := func(event *binance.WsKlineEvent) {
		if event != nil {
			// fix endTime
			if strings.HasSuffix(strconv.FormatInt(event.Kline.EndTime, 10), "999") {
				event.Kline.EndTime++
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
			if LogNotNilError(errs) {
				return
			}

			eventCallback(wEvent)
		}
	}
	wsErrHandler := func(err error) {
		errorHandler(errors.New("service request failed: " + err.Error()))
	}
	var openWsErr error
	w.WsChannels = new(workers.WorkerChannels)
	w.WsChannels.WsDone, w.WsChannels.WsStop, openWsErr = binance.WsCombinedKlineServe(symbolIntervalsMap, wsCandleHandler, wsErrHandler)
	if openWsErr != nil {
		return errors.New("service request failed: " + openWsErr.Error())
	}
	return nil
}

// TradeEventWorkerBinance - TradeEventWorker for binance
type TradeEventWorkerBinance struct {
	workers.TradeEventWorker
}

// GetTradeEventsWorker - create new market candle worker
func (a *BinanceSpotAdapter) GetTradeEventsWorker() workers.ITradeEventWorker {
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
		errorHandler(errors.New("service request failed: " + err.Error()))
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
				BuyerOrderID:  event.BuyerOrderID,
				SellerOrderID: event.SellerOrderID,
			}
			errs := make([]error, 2)
			wEvent.Price, errs[0] = strconv.ParseFloat(event.Price, 64)
			wEvent.Quantity, errs[0] = strconv.ParseFloat(event.Quantity, 64)
			if LogNotNilError(errs) {
				return
			}
			eventCallback(wEvent)
		}
	}

	var openWsErr error
	w.WsChannels = new(workers.WorkerChannels)
	w.WsChannels.WsDone, w.WsChannels.WsStop, openWsErr = binance.WsTradeServe(symbol, wsTradeHandler, wsErrHandler)
	if openWsErr != nil {
		return errors.New("failed to subscribe to trade events: " + openWsErr.Error())
	}
	return nil
}
