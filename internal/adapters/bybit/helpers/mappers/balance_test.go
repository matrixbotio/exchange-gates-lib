package mappers

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/hirokisan/bybit/v2"
	"github.com/stretchr/testify/require"
)

func TestConvertAccountBalance(t *testing.T) {
	// given
	data := bybit.V5GetWalletBalanceResponse{
		Result: bybit.V5WalletBalanceResult{
			List: []bybit.V5WalletBalanceList{
				{
					AccountType: string(bybit.AccountTypeV5SPOT),
					Coin: []bybit.V5WalletBalanceCoin{
						{
							Coin: "BTC",
							Free: "0.01",
						},
					},
				},
				{
					AccountType: string(bybit.AccountTypeV5OPTION),
					Coin: []bybit.V5WalletBalanceCoin{
						{
							Coin: "USDT",
							Free: "150",
						},
					},
				},
			},
		},
	}

	// when
	balances, err := ConvertAccountBalance(data)

	// then
	require.NoError(t, err)
	require.Len(t, balances, 1)
	assert.Equal(t, float64(0.01), balances[0].Free)
}
