package mappers

import (
	"errors"
	"fmt"

	"github.com/hirokisan/bybit/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bybit/helpers/accessors"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bybit/helpers/errs"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
)

func ConvertAccountBalance(
	data bybit.V5GetWalletBalanceResponse,
	accountType bybit.AccountTypeV5,
) ([]structs.Balance, error) {
	balance, err := accessors.GetAccountBalance(data, accountType)
	if err != nil {
		if errors.Is(err, errs.ErrSpotBalanceNotFound) {
			return nil, nil
		}
		return nil, err
	}

	var result []structs.Balance
	for _, tickerData := range balance.Coin {
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
