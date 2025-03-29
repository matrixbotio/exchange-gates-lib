//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package wrapper

import (
	"context"
	"fmt"
	"time"

	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/mappers"

	"github.com/adshao/go-binance/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
)

type BinanceAPIWrapper interface {
	Sync(context.Context)
	Connect(ctx context.Context, keyPublic, keySecret string) error
	Ping(context.Context) error
	GetAccountData(context.Context) (*binance.Account, error)

	GetPrices(
		ctx context.Context,
		pairSymbol string,
	) ([]*binance.SymbolPrice, error)

	GetOpenOrders(
		ctx context.Context,
		pairSymbol string,
	) ([]*binance.Order, error)

	GetExchangeInfo(
		ctx context.Context,
		pairSymbol string,
	) (*binance.ExchangeInfo, error)

	CancelOrderByID(
		ctx context.Context,
		pairSymbol string,
		orderID int64,
	) error

	CancelOrderByClientOrderID(
		ctx context.Context,
		pairSymbol string,
		clientOrderID string,
	) error

	GetOrderDataByOrderID(
		ctx context.Context,
		pairSymbol string,
		orderID int64,
	) (*binance.Order, error)

	GetOrderDataByClientOrderID(
		ctx context.Context,
		pairSymbol string,
		clientOrderID string,
	) (*binance.Order, error)

	PlaceLimitOrder(
		ctx context.Context,
		pairSymbol string,
		orderSide binance.SideType,
		qty string,
		price string,
		optionalClientOrderID string,
	) (*binance.CreateOrderResponse, error)

	PlaceMarketOrder(
		ctx context.Context,
		pairSymbol string,
		orderSide binance.SideType,
		qty string,
		price string,
		optionalClientOrderID string,
	) (*binance.CreateOrderResponse, error)

	GetKlines(
		ctx context.Context,
		pairSymbol string,
		interval string,
		limit int,
	) ([]*binance.Kline, error)

	SubscribeToCandle(
		pairSymbol string,
		interval string,
		eventCallback func(event workers.CandleEvent),
		errorHandler func(err error),
	) (doneC chan struct{}, stopC chan struct{}, err error)

	SubscribeToCandlesList(
		intervalsPerPair map[string]string,
		eventCallback func(event workers.CandleEvent),
		errorHandler func(err error),
	) (doneC chan struct{}, stopC chan struct{}, err error)

	SubscribeToPriceEvents(
		pairSymbol string,
		eventCallback binance.WsBookTickerHandler,
		errorHandler func(err error),
	) (doneC chan struct{}, stopC chan struct{}, err error)

	SubscribeToTradeEventsPrivate(
		exchangeTag string,
		callback workers.TradeEventPrivateCallback,
		handler func(err error),
	) (doneC chan struct{}, stopC chan struct{}, err error)

	GetOrderTradeHistory(
		ctx context.Context,
		orderID int64,
		pairSymbol string,
	) ([]*binance.TradeV3, error)
}

type BinanceClientWrapper struct {
	*binance.Client
}

func NewWrapper() BinanceAPIWrapper {
	return &BinanceClientWrapper{}
}

func (b *BinanceClientWrapper) Sync(ctx context.Context) {
	//goland:noinspection GoUnhandledErrorResult
	b.NewSetServerTimeService().Do(ctx)
}

func (b *BinanceClientWrapper) Connect(
	ctx context.Context,
	keyPublic,
	keySecret string,
) error {
	b.Client = binance.NewClient(keyPublic, keySecret)
	if err := b.Ping(ctx); err != nil {
		return fmt.Errorf("ping binance: %w", err)
	}
	return nil
}

func (b *BinanceClientWrapper) Ping(ctx context.Context) error {
	var err error
	for attemptNumber := 1; attemptNumber <= consts.PingRetryAttempts; attemptNumber++ {
		if err := b.NewPingService().Do(ctx); err == nil {
			return nil
		}

		time.Sleep(consts.PingRetryWaitTime)
	}

	return fmt.Errorf("ping exchange: %w", err)
}

func (b *BinanceClientWrapper) GetOrderDataByOrderID(
	ctx context.Context,
	pairSymbol string,
	orderID int64,
) (*binance.Order, error) {
	return b.NewGetOrderService().Symbol(pairSymbol).
		OrderID(orderID).Do(ctx)
}

func (b *BinanceClientWrapper) GetOrderDataByClientOrderID(
	ctx context.Context,
	pairSymbol string,
	clientOrderID string,
) (*binance.Order, error) {
	return b.NewGetOrderService().Symbol(pairSymbol).
		OrigClientOrderID(clientOrderID).Do(ctx)
}

func (b *BinanceClientWrapper) GetAccountData(ctx context.Context) (
	*binance.Account,
	error,
) {
	return b.NewGetAccountService().Do(ctx)
}

func (b *BinanceClientWrapper) GetPrices(
	ctx context.Context,
	pairSymbol string,
) ([]*binance.SymbolPrice, error) {
	return b.NewListPricesService().Symbol(pairSymbol).Do(ctx)
}

func (b *BinanceClientWrapper) GetOpenOrders(ctx context.Context, pairSymbol string) (
	[]*binance.Order,
	error,
) {
	return b.NewListOpenOrdersService().Symbol(pairSymbol).Do(ctx)
}

func (b *BinanceClientWrapper) GetExchangeInfo(ctx context.Context, pairSymbol string) (
	*binance.ExchangeInfo,
	error,
) {
	srv := b.NewExchangeInfoService()
	if pairSymbol != "" {
		srv.Symbol(pairSymbol)
	}

	return srv.Do(ctx)
}

func (b *BinanceClientWrapper) CancelOrderByID(
	ctx context.Context,
	pairSymbol string,
	orderID int64,
) error {
	_, err := b.NewCancelOrderService().Symbol(pairSymbol).
		OrderID(orderID).Do(ctx)
	return err
}

func (b *BinanceClientWrapper) CancelOrderByClientOrderID(
	ctx context.Context,
	pairSymbol string,
	clientOrderID string,
) error {
	_, err := b.NewCancelOrderService().Symbol(pairSymbol).
		OrigClientOrderID(clientOrderID).Do(ctx)
	return err
}

func (b *BinanceClientWrapper) PlaceLimitOrder(
	ctx context.Context,
	pairSymbol string,
	orderSide binance.SideType,
	qty string,
	price string,
	optionalClientOrderID string,
) (*binance.CreateOrderResponse, error) {
	orderService := b.NewCreateOrderService().Symbol(pairSymbol).
		Side(orderSide).Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeGTC).Quantity(qty).
		Price(price)

	if optionalClientOrderID != "" {
		orderService.NewClientOrderID(optionalClientOrderID)
	}

	return orderService.Do(ctx)
}

func (b *BinanceClientWrapper) PlaceMarketOrder(
	ctx context.Context,
	pairSymbol string,
	orderSide binance.SideType,
	qty string,
	price string,
	optionalClientOrderID string,
) (*binance.CreateOrderResponse, error) {
	orderService := b.NewCreateOrderService().Symbol(pairSymbol).
		Side(orderSide).Type(binance.OrderTypeMarket).
		TimeInForce(binance.TimeInForceTypeGTC).Quantity(qty)

	if optionalClientOrderID != "" {
		orderService.NewClientOrderID(optionalClientOrderID)
	}

	return orderService.Do(ctx)
}

func (b *BinanceClientWrapper) GetKlines(
	ctx context.Context,
	pairSymbol string,
	interval string,
	limit int,
) ([]*binance.Kline, error) {
	return b.NewKlinesService().Symbol(pairSymbol).Interval(interval).
		Limit(limit).Do(ctx)
}

func (b *BinanceClientWrapper) SubscribeToCandle(
	pairSymbol string,
	interval string,
	eventCallback func(event workers.CandleEvent),
	errorHandler func(err error),
) (doneC chan struct{}, stopC chan struct{}, err error) {
	return binance.WsKlineServe(
		pairSymbol,
		interval,
		helpers.GetCandleEventsHandler(eventCallback, errorHandler),
		errorHandler,
	)
}

func (b *BinanceClientWrapper) SubscribeToCandlesList(
	intervalsPerPair map[string]string,
	eventCallback func(event workers.CandleEvent),
	errorHandler func(err error),
) (doneC chan struct{}, stopC chan struct{}, err error) {
	return binance.WsCombinedKlineServe(
		intervalsPerPair,
		helpers.GetCandleEventsHandler(eventCallback, errorHandler),
		errorHandler,
	)
}

func (b *BinanceClientWrapper) SubscribeToPriceEvents(
	pairSymbol string,
	eventCallback binance.WsBookTickerHandler,
	errorHandler func(err error),
) (doneC chan struct{}, stopC chan struct{}, err error) {
	return binance.WsBookTickerServe(
		pairSymbol,
		eventCallback,
		errorHandler,
	)
}

func (b *BinanceClientWrapper) SubscribeToTradeEventsPrivate(
	exchangeTag string,
	eventCallback workers.TradeEventPrivateCallback,
	errorHandler func(err error),
) (doneC chan struct{}, stopC chan struct{}, err error) {
	listenKey, err := b.Client.NewStartUserStreamService().Do(context.Background())
	if err != nil {
		return nil, nil, err
	}

	binanceEventCallback := func(event *binance.WsUserDataEvent) {
		if event == nil {
			return
		}

		if event.Event != binance.UserDataEventTypeExecutionReport {
			// ignore non order-related events
			return
		}

		wEvent, err := mappers.ConvertTradeEventPrivate(*event, exchangeTag)
		if err != nil {
			errorHandler(fmt.Errorf("convert trade event: %w", err))
			return
		}

		eventCallback(wEvent)
	}

	doneC, bStopC, err := binance.WsUserDataServe(listenKey, binanceEventCallback, errorHandler)
	if err != nil {
		return nil, nil, fmt.Errorf("serve: %w", err)
	}

	stopC = make(chan struct{})

	go func() {
		for {
			select {
			case <-time.After(time.Minute * 30):
				err := b.Client.NewKeepaliveUserStreamService().ListenKey(listenKey).Do(context.Background())
				if err != nil {
					b.Logger.Println(err)

					doneC, bStopC, err = binance.WsUserDataServe(listenKey, binanceEventCallback, errorHandler)
					if err != nil {
						b.Logger.Println("Failed to reestablish WebSocket connection:", err)
						return
					}
				}
			case <-stopC:
				bStopC <- struct{}{}
				return
			}
		}
	}()

	return doneC, stopC, nil
}

func (b *BinanceClientWrapper) GetOrderTradeHistory(
	ctx context.Context,
	orderID int64,
	pairSymbol string,
) ([]*binance.TradeV3, error) {
	return b.NewListTradesService().OrderId(orderID).
		Symbol(pairSymbol).Do(ctx)
}
