package binance

import (
	"testing"

	"github.com/adshao/go-binance/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/errs"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/wrapper"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestVerifyAPIKeysSuccess(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)

	w.EXPECT().GetAccountData(mock.Anything).Return(
		&binance.Account{
			CanTrade: true,
		}, nil,
	)

	// when
	err := a.VerifyAPIKeys("test", "test")

	// then
	require.NoError(t, err)
}

func TestVerifyAPIKeysError(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)

	w.EXPECT().GetAccountData(mock.Anything).Return(
		nil, errTestException,
	)

	// when
	err := a.VerifyAPIKeys("test", "test")

	// then
	require.ErrorIs(t, err, errTestException)
}

func TestVerifyAPIKeysTradingNotAllowed(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)

	w.EXPECT().GetAccountData(mock.Anything).Return(
		&binance.Account{}, nil,
	)

	// when
	err := a.VerifyAPIKeys("test", "test")

	// then
	require.ErrorIs(t, err, errs.ErrTradingNotAllowed)
}
