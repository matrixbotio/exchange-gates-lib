package accessors

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/hirokisan/bybit/v2"
	"github.com/stretchr/testify/require"
)

func TestGetPairSymbolPointerV5(t *testing.T) {
	// given
	symbol := "LTCUSDT"

	// when
	result := GetPairSymbolPointerV5(symbol)

	// then
	require.NotNil(t, result)
	assert.Equal(t, symbol, string(*result))
}

func TestGetOrderIDFromHistoryOrdersParam(t *testing.T) {
	// given
	expectedOrderID := "12345"
	param := bybit.V5GetHistoryOrdersParam{
		OrderID: &expectedOrderID,
	}

	// when
	orderID := GetOrderIDFromHistoryOrdersParam(param)

	// then
	assert.Equal(t, expectedOrderID, orderID)
}

func TestGetOrderIDFromHistoryOrdersParam2(t *testing.T) {
	// given
	expectedClientOrderID := "12345"
	param := bybit.V5GetHistoryOrdersParam{
		OrderLinkID: &expectedClientOrderID,
	}

	// when
	orderID := GetOrderIDFromHistoryOrdersParam(param)

	// then
	assert.Equal(t, expectedClientOrderID, orderID)
}

func TestGetOrderIDFromHistoryUnknown(t *testing.T) {
	// given
	param := bybit.V5GetHistoryOrdersParam{}

	// when
	orderID := GetOrderIDFromHistoryOrdersParam(param)

	// then
	assert.Equal(t, unknownOrderID, orderID)
}

func TestGetOrderSymbolFromHistoryOrdersParam(t *testing.T) {
	// given
	expectedSymbol := "LTCUSDT"
	param := bybit.V5GetHistoryOrdersParam{
		Symbol: GetPairSymbolPointerV5(expectedSymbol),
	}

	// when
	symbol := GetOrderSymbolFromHistoryOrdersParam(param)

	// then
	assert.Equal(t, expectedSymbol, symbol)
}

func TestGetOrderSymbolFromHistoryOUnknown(t *testing.T) {
	// given
	param := bybit.V5GetHistoryOrdersParam{}

	// when
	symbol := GetOrderSymbolFromHistoryOrdersParam(param)

	// then
	assert.Equal(t, unknownPairSymbol, symbol)
}

func TestGetAccountBalanceSpot(t *testing.T) {
	// given
	ticker := bybit.Coin("BTC")
	data := bybit.V5GetWalletBalanceResponse{
		Result: bybit.V5WalletBalanceResult{
			List: []bybit.V5WalletBalanceList{
				{
					AccountType: string(bybit.AccountTypeV5SPOT),
					Coin: []bybit.V5WalletBalanceCoin{
						{
							Coin: ticker,
							Free: "0.01",
						},
					},
				},
				{
					AccountType: string(bybit.AccountTypeV5UNIFIED),
					Coin: []bybit.V5WalletBalanceCoin{
						{
							Coin: "LTC",
							Free: "5",
						},
					},
				},
			},
		},
	}
	accountType := bybit.AccountTypeV5SPOT

	// when
	spotBalance, err := GetAccountBalance(data, accountType)

	// then
	require.NoError(t, err)
	assert.Equal(t, ticker, spotBalance.Coin[0].Coin)
}
