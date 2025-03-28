package binance

import (
	"testing"

	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/wrapper"
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestSubscribeCandle(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	w := wrapper.NewMockBinanceAPIWrapper(ctrl)
	a := New(w)

	eventHandler := func(event workers.CandleEvent) {}
	errHandler := func(err error) {}

	pairSymbol := "LTCUSDT"
	interval := consts.Interval1min

	w.EXPECT().SubscribeToCandle(
		gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any(),
	).
		Return(make(chan struct{}), make(chan struct{}), nil).
		AnyTimes()

	// when
	err := a.SubscribeCandle(pairSymbol, interval, eventHandler, errHandler)

	// then
	require.NoError(t, err)
}

func TestUnsubscribeCandle(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	w := wrapper.NewMockBinanceAPIWrapper(ctrl)
	a := New(w)

	eventHandler := func(event workers.CandleEvent) {}
	errHandler := func(err error) {}

	pairSymbol := "LTCUSDT"
	interval := consts.Interval1min

	w.EXPECT().SubscribeToCandle(
		gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any(),
	).
		Return(make(chan struct{}), make(chan struct{}), nil).
		AnyTimes()

	require.NoError(t, a.SubscribeCandle(
		pairSymbol, interval,
		eventHandler, errHandler,
	))

	// when
	a.UnsubscribeCandle(pairSymbol, interval)
}
