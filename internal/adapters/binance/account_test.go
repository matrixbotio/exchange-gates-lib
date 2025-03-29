package binance

import (
	"errors"
	"testing"

	"github.com/adshao/go-binance/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/wrapper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var errTestException = errors.New("test exception")

func TestCanTradeSuccess(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	w := wrapper.NewMockBinanceAPIWrapper(ctrl)
	a := New(w)

	w.EXPECT().GetAccountData(gomock.Any()).Return(
		&binance.Account{
			CanTrade: true,
		}, nil,
	)

	// when
	isTradingAllowed, err := a.CanTrade()

	// then
	require.NoError(t, err)
	assert.True(t, isTradingAllowed)
}

func TestCanTradeNotAllowed(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	w := wrapper.NewMockBinanceAPIWrapper(ctrl)
	a := New(w)

	w.EXPECT().GetAccountData(gomock.Any()).Return(
		&binance.Account{
			CanTrade: false,
		}, nil,
	)

	// when
	isTradingAllowed, err := a.CanTrade()

	// then
	require.NoError(t, err)
	assert.False(t, isTradingAllowed)
}

func TestCanTradeError(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	w := wrapper.NewMockBinanceAPIWrapper(ctrl)
	a := New(w)

	w.EXPECT().GetAccountData(gomock.Any()).Return(nil, errTestException)

	// when
	_, err := a.CanTrade()

	// then
	require.ErrorIs(t, err, errTestException)
}
