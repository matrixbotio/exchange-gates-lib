package consts

import "time"

const (
	TestExchangeID        = 0
	PairDefaultExchangeID = 1
	PairDefaultMinQty     = 0.001
	PairDefaultMaxQty     = 99999.99
	PairDefaultMinPrice   = 0.000001
	PairDefaultQtyStep    = 0.001
	PairDefaultPriceStep  = 0.01
	PairMinDeposit        = 10
	PairDefaultBaseAsset  = "BTC"
	PairDefaultQuoteAsset = "BUSD"
	PairDefaultStatus     = "TRADING"
	PairDefaultAsset      = PairDefaultBaseAsset + PairDefaultQuoteAsset
)

type OrderStatus string

const (
	OrderStatusNew                      OrderStatus = "NEW"
	OrderStatusPartiallyFilled          OrderStatus = "PARTIALLY_FILLED"
	OrderStatusPartiallyFilledCancelled OrderStatus = "PARTIALLY_FILLED_CANCELLED"
	OrderStatusFilled                   OrderStatus = "FILLED"
	OrderStatusCancelled                OrderStatus = "CANCELED"
	OrderStatusPendingCancel            OrderStatus = "PENDING_CANCEL"
	OrderStatusRejected                 OrderStatus = "REJECTED"
	OrderStatusExpired                  OrderStatus = "EXPIRED"
	OrderStatusUnknown                  OrderStatus = "UNKNOWN"
	OrderStatusUntriggered              OrderStatus = "UNTRIGGERED"
	OrderStatusTriggered                OrderStatus = "TRIGGERED"
	OrderStatusDeactivated              OrderStatus = "DEACTIVATED"
)

const (
	ExchangeSetupConnTimeout = 3500 // ms
	ReadTimeout              = time.Second * 5
)

const (
	PingRetryAttempts = 3
	PingRetryWaitTime = time.Second * 2
)

const (
	ExchangeIDbinanceSpot = 1
	ExchangeIDbybitSpot   = 2
	ExchangeIDbingx       = 3
	ExchangeIDgateSpot    = 4
)

const (
	CheckOrdersTimeoutBinance = time.Second * 30
	CheckOrdersTimeoutBybit   = time.Second * 15
	CheckOrdersTimeoutGate    = time.Second * 15
)

type OrderSide string

const (
	OrderSideBuy  OrderSide = "buy"
	OrderSideSell OrderSide = "sell"
)

const (
	PairStatusTrading   = "TRADING"
	PairStatusOffline   = "OFFLINE"
	PairStatusPreOpen   = "PRE-OPEN"
	PairStatusSuspended = "SUSPENDED"
	PairStatusUnknown   = "UNKNOWN"
)

const BingXAdapterTag = "bingx-spot"
const BinanceAdapterTag = "binance-spot"
const GateAdapterTag = "gate-spot"
