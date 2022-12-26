package bybit

import (
	"errors"
	"fmt"

	"github.com/hirokisan/bybit/v2"

	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bybit/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
)

func (a *adapter) getTickersBalance(tickers []bybit.Coin) ([]bybit.V5WalletBalanceList, error) {
	balanceData, err := a.client.V5().Account().
		GetWalletBalance(bybit.AccountType(bybit.AccountTypeV5SPOT), tickers)
	if err != nil {
		return nil, fmt.Errorf("get ticker balance: %w", err)
	}

	if len(balanceData.Result.List) == 0 {
		return nil, errors.New("balance data not found: list is empty")
	}

	return balanceData.Result.List, nil
}

func (a *adapter) getTickerBalance(tickerTag string) (structs.AssetBalance, error) {
	balanceData, err := a.getTickersBalance([]bybit.Coin{bybit.Coin(tickerTag)})
	if err != nil {
		return structs.AssetBalance{}, fmt.Errorf("get tickers balance: %w", err)
	}

	if len(balanceData[0].Coin) == 0 {
		return structs.AssetBalance{
			Ticker: tickerTag,
		}, nil
	}

	return mappers.ConvertBalances(balanceData[0].Coin, tickerTag)
}
