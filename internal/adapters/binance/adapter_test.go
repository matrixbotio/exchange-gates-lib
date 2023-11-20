package binance

import (
	"context"
	"testing"

	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/pkg/structs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBinanceAdapter(t *testing.T) {
	// given
	a := New()

	// when
	exchangeID := a.GetID()

	// then
	assert.Equal(t, exchangeID, consts.ExchangeIDbinanceSpot)
}

func TestConnect(t *testing.T) {
	// given
	w := NewMockBinanceAPIWrapper(t)
	a := createAdapter(w)
	credentials := structs.APICredentials{
		Type: structs.APICredentialsTypeKeypair,
	}

	w.EXPECT().Sync(context.Background())

	w.EXPECT().Connect(mock.Anything, mock.Anything, context.Background()).
		Return(nil)

	// when
	err := a.Connect(credentials)

	// then
	require.NoError(t, err)
}
