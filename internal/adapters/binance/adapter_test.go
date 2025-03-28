package binance

import (
	"context"
	"errors"
	"testing"

	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/errs"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/wrapper"
	"github.com/matrixbotio/exchange-gates-lib/pkg/structs"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestConnectSucess(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	w := wrapper.NewMockBinanceAPIWrapper(ctrl)
	a := New(w)
	credentials := structs.APICredentials{
		Type: structs.APICredentialsTypeKeypair,
	}

	w.EXPECT().Sync(context.Background())

	w.EXPECT().Connect(context.Background(), gomock.Any(), gomock.Any()).
		Return(nil)

	// when
	err := a.Connect(credentials)

	// then
	require.NoError(t, err)
}

func TestConnectErrorInvalidCredentials(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	w := wrapper.NewMockBinanceAPIWrapper(ctrl)
	a := New(w)
	credentials := structs.APICredentials{
		Type: structs.APICredentialsType("wtf"),
	}

	// when
	err := a.Connect(credentials)

	// then
	require.ErrorIs(t, err, errs.ErrInvalidCredentials)
}

func TestConnectErrorPingFailed(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	w := wrapper.NewMockBinanceAPIWrapper(ctrl)
	a := New(w)
	credentials := structs.APICredentials{
		Type: structs.APICredentialsTypeKeypair,
	}

	w.EXPECT().Connect(context.Background(), gomock.Any(), gomock.Any()).
		Return(errors.New("ping: timeout"))

	// when
	err := a.Connect(credentials)

	// then
	require.ErrorContains(t, err, "timeout")
}
