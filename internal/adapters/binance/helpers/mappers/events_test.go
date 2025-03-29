package mappers

import (
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/adshao/go-binance/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertPriceEventSuccess(t *testing.T) {
	// given
	event := binance.WsBookTickerEvent{
		Symbol:       "BTCUSDT",
		BestBidPrice: "20000",
		BestAskPrice: "20100",
		BestBidQty:   "0.1",
		BestAskQty:   "0,2",
	}

	// when
	ask, bid, err := ConvertPriceEvent(event)

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(20100), ask)
	assert.Equal(t, float64(20000), bid)
}

func TestConvertPriceEventError(t *testing.T) {
	// given
	event := binance.WsBookTickerEvent{
		Symbol:       "BTCUSDT",
		BestBidPrice: "wtf",
		BestAskPrice: "omg",
	}

	// when
	_, _, err := ConvertPriceEvent(event)

	// then
	require.ErrorContains(t, err, "invalid syntax")
}

func TestConvertTradeEventPrivateSuccess(t *testing.T) {
	// given
	exchangeTag := "binance-spot"
	eventTime := time.Now().UnixMilli()
	symbol := "BTCUSDT"
	orderID := int64(11111)
	clientOrderID := uuid.NewString()
	quantity := 10.12
	price := 134.3335
	tradeID := int64(12345)

	event := binance.WsUserDataEvent{
		Event: binance.UserDataEventTypeExecutionReport,
		Time:  eventTime,
		OrderUpdate: binance.WsOrderUpdate{
			Id:            orderID,
			Symbol:        symbol,
			ClientOrderId: clientOrderID,
			Volume:        strconv.FormatFloat(quantity, 'f', -1, 64),
			Price:         strconv.FormatFloat(price, 'f', -1, 64),
			TradeId:       tradeID,
		},
	}

	// when
	result, err := ConvertTradeEventPrivate(event, exchangeTag)

	// then
	require.NoError(t, err)
	assert.Equal(t, strconv.FormatInt(tradeID, 10), result.ID)
	assert.Equal(t, eventTime, result.Time)
	assert.Equal(t, exchangeTag, result.ExchangeTag)
	assert.Equal(t, symbol, result.Symbol)
	assert.Equal(t, strconv.FormatInt(orderID, 10), result.OrderID)
	assert.Equal(t, clientOrderID, result.ClientOrderID)
	assert.Equal(t, price, result.Price)
	assert.Equal(t, quantity, result.Quantity)
}
