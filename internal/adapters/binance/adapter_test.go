package binance

import (
	"context"
	"errors"
	"testing"

	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/helpers/errs"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/binance/wrapper"
	"github.com/matrixbotio/exchange-gates-lib/pkg/structs"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestConnectSucess(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)
	credentials := structs.APICredentials{
		Type: structs.APICredentialsTypeKeypair,
	}

	w.EXPECT().Sync(context.Background())

	w.EXPECT().Connect(context.Background(), mock.Anything, mock.Anything).
		Return(nil)

	// when
	err := a.Connect(credentials)

	// then
	require.NoError(t, err)
}

func TestConnectErrorInvalidCredentials(t *testing.T) {
	// given
	w := wrapper.NewMockBinanceAPIWrapper(t)
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
	w := wrapper.NewMockBinanceAPIWrapper(t)
	a := New(w)
	credentials := structs.APICredentials{
		Type: structs.APICredentialsTypeKeypair,
	}

	w.EXPECT().Connect(context.Background(), mock.Anything, mock.Anything).
		Return(errors.New("ping: timeout"))

	// when
	err := a.Connect(credentials)

	// then
	require.ErrorContains(t, err, "timeout")
}
