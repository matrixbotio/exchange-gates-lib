package binance

import (
	"testing"

	"github.com/adshao/go-binance/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/errs"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/wrapper"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var testAPIPubkey = "test pubkey"
var testAPISecret = "test secret"

func TestVerifyAPIKeysSuccess(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	w := wrapper.NewMockBinanceAPIWrapper(ctrl)
	a := New(w)

	w.EXPECT().Sync(gomock.Any())
	w.EXPECT().Connect(gomock.Any(), testAPIPubkey, testAPISecret).Return(nil)
	w.EXPECT().GetAccountData(gomock.Any()).Return(
		&binance.Account{
			CanTrade: true,
		}, nil,
	)

	// when
	err := a.VerifyAPIKeys(testAPIPubkey, testAPISecret)

	// then
	require.NoError(t, err)
}

func TestVerifyAPIKeysError(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	w := wrapper.NewMockBinanceAPIWrapper(ctrl)
	a := New(w)

	w.EXPECT().Sync(gomock.Any())
	w.EXPECT().Connect(gomock.Any(), testAPIPubkey, testAPISecret).Return(nil)
	w.EXPECT().GetAccountData(gomock.Any()).Return(
		nil, errTestException,
	)

	// when
	err := a.VerifyAPIKeys(testAPIPubkey, testAPISecret)

	// then
	require.ErrorIs(t, err, errTestException)
}

func TestVerifyAPIKeysTradingNotAllowed(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	w := wrapper.NewMockBinanceAPIWrapper(ctrl)
	a := New(w)

	w.EXPECT().Sync(gomock.Any())
	w.EXPECT().Connect(gomock.Any(), testAPIPubkey, testAPISecret).Return(nil)
	w.EXPECT().GetAccountData(gomock.Any()).Return(
		&binance.Account{}, nil,
	)

	// when
	err := a.VerifyAPIKeys(testAPIPubkey, testAPISecret)

	// then
	require.ErrorIs(t, err, errs.ErrTradingNotAllowed)
}
