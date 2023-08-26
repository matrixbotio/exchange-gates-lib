package mappers

import (
	"errors"
	"fmt"

	"github.com/hirokisan/bybit/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bybit/helpers/accessors"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bybit/helpers/errs"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
)

func ConvertAccountBalance(data bybit.V5GetWalletBalanceResponse) ([]structs.Balance, error) {
	spotBalance, err := accessors.GetAccountBalanceSpot(data)
	if err != nil {
		if errors.Is(err, errs.ErrSpotBalanceNotFound) {
			return nil, nil
		}
		return nil, err
	}

	var result []structs.Balance
	for _, tickerData := range spotBalance.Coin {
		tickerBalance, err := convertCoinData(tickerData, string(tickerData.Coin))
		if err != nil {
			return nil, fmt.Errorf("convert ticker balance: %w", err)
		}

		result = append(result, structs.Balance{
			Asset:  tickerBalance.Ticker,
			Free:   tickerBalance.Free,
			Locked: tickerBalance.Locked,
		})
	}

	return result, nil
}
