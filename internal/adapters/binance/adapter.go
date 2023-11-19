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
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/mappers"
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
				errs.ErrOrderNotFound,
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

	orderResponse, err := a.getOrderFromService(pairSymbol, orderID, clientOrderID)
	if err != nil {
		return structs.OrderData{}, err
	}

	return mappers.ConvertOrderData(orderResponse)
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
	orderSide, err := mappers.GetBinanceOrderSide(order.Type)
	if err != nil {
		return structs.CreateOrderResponse{}, fmt.Errorf("get order side: %w", err)
	}

	orderService := a.binanceAPI.NewCreateOrderService().Symbol(order.PairSymbol).
		Side(orderSide).Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeGTC).Quantity(order.Qty).
		Price(order.Price)

	if order.ClientOrderID != "" {
		orderService.NewClientOrderID(order.ClientOrderID)
	}

	orderResponse, err := orderService.Do(ctx)
	if err != nil {
		return structs.CreateOrderResponse{}, fmt.Errorf("create order: %w", err)
	}

	if orderResponse == nil {
		return structs.CreateOrderResponse{}, errors.New("order response is empty")
	}

	return mappers.ConvertPlacedOrder(*orderResponse)
}

func (a *adapter) ping() error {
	var err error
	for attemptNumber := 1; attemptNumber <= consts.PingRetryAttempts; attemptNumber++ {
		if err := a.binanceAPI.NewPingService().Do(context.Background()); err == nil {
			return nil
		}

		time.Sleep(consts.PingRetryWaitTime)
	}

	return fmt.Errorf("ping exchange: %w", err)
}
