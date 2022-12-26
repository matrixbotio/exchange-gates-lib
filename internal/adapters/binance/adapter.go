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
	"github.com/matrixbotio/exchange-gates-lib/pkg/errs"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
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

	var logOrderID string
	if orderID > 0 {
		logOrderID = strconv.FormatInt(orderID, 10)
		tradeData.OrderID = orderID
		s.OrderID(orderID)
	} else {
		logOrderID = clientOrderID
		tradeData.ClientOrderID = clientOrderID
		s.OrigClientOrderID(clientOrderID)
	}

	orderResponse, err := s.Do(context.Background())
	if err != nil {
		if strings.Contains(err.Error(), "Order does not exist") {
			tradeData.Status = pkgStructs.OrderStatusUnknown
			return nil, fmt.Errorf(
				"%w %s",
				errs.OrderNotFound,
				logOrderID,
			)
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
func (a *adapter) GetOrderByClientOrderID(pairSymbol string, clientOrderID string) (
	structs.OrderData,
	error,
) {
	return a.getOrderData(pairSymbol, 0, clientOrderID)
}

// PlaceOrder - place order on exchange
func (a *adapter) PlaceOrder(ctx context.Context, order structs.BotOrderAdjusted) (
	structs.CreateOrderResponse,
	error,
) {
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

func (a *adapter) CanTrade() (bool, error) {
	binanceAccountData, clientErr := a.binanceAPI.NewGetAccountService().
		Do(context.Background())
	if clientErr != nil {
		return false, errors.New("send request to trade, " + clientErr.Error())
	}
	return binanceAccountData.CanTrade, nil
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
	exchangeInfo, err := a.binanceAPI.NewExchangeInfoService().
		Symbol(pairSymbol).Do(context.Background())
	if err != nil {
		return structs.ExchangePairData{}, err
	}

	// find pairSymbol
	for _, symbolData := range exchangeInfo.Symbols {
		return getExchangePairData(symbolData, a.ExchangeID)
	}

	return structs.ExchangePairData{}, errors.New("data for " + pairSymbol + " pair not found")
}

// GetPairOpenOrders - get open orders array
func (a *adapter) GetPairOpenOrders(pairSymbol string) ([]structs.OrderData, error) {
	ordersRaw, err := a.binanceAPI.NewListOpenOrdersService().
		Symbol(pairSymbol).Do(context.Background())
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
		return errors.New("your API key does not have permission to trade," +
			" change its restrictions")
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
		pairData, err := getExchangePairData(symbolData, a.ExchangeID)
		if err != nil {
			lastError = err
		} else {
			pairs = append(pairs, pairData)
		}
	}
	return pairs, lastError
}

func (a *adapter) GetPairOrdersHistory(task structs.GetOrdersHistoryTask) (
	[]structs.OrderData,
	error,
) {
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
