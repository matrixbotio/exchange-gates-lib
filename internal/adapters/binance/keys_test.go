package binance

import (
	"testing"

	"github.com/adshao/go-binance/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/errs"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/wrapper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var testAPIPubkey = "test pubkey"
var testAPISecret = "test secret"

func TestVerifyAPIKeysSuccess(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)

	w.EXPECT().Sync(mock.Anything)
	w.EXPECT().Connect(mock.Anything, testAPIPubkey, testAPISecret).Return(nil)
	w.EXPECT().GetAccountData(mock.Anything).Return(
		&binance.Account{
			CanTrade: true,
		}, nil,
	)

	// when
	result, err := a.VerifyAPIKeys(testAPIPubkey, testAPISecret)

	// then
	require.NoError(t, err)
	assert.True(t, result.Active)
}

func TestVerifyAPIKeysError(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)

	w.EXPECT().Sync(mock.Anything)
	w.EXPECT().Connect(mock.Anything, testAPIPubkey, testAPISecret).Return(nil)
	w.EXPECT().GetAccountData(mock.Anything).Return(
		nil, errTestException,
	)

	// when
	result, err := a.VerifyAPIKeys(testAPIPubkey, testAPISecret)

	// then
	require.ErrorIs(t, err, errTestException)
	assert.False(t, result.Active)
}

func TestVerifyAPIKeysTradingNotAllowed(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)

	w.EXPECT().Sync(mock.Anything)
	w.EXPECT().Connect(mock.Anything, testAPIPubkey, testAPISecret).Return(nil)
	w.EXPECT().GetAccountData(mock.Anything).Return(
		&binance.Account{}, nil,
	)

	// when
	_, err := a.VerifyAPIKeys(testAPIPubkey, testAPISecret)

	// then
	require.ErrorIs(t, err, errs.ErrTradingNotAllowed)
}
