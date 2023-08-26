package accessors

import (
	"fmt"

	"github.com/hirokisan/bybit/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bybit/helpers/errs"
)

const (
	unknownPairSymbol = "UNKNOWN"
	unknownOrderID    = "unknown-id"
)

func GetAccountBalanceSpot(data bybit.V5GetWalletBalanceResponse) (
	bybit.V5WalletBalanceList, error,
) {
	if len(data.Result.List) == 0 {
		return bybit.V5WalletBalanceList{}, fmt.Errorf("balance data not available")
	}

	for _, tickerData := range data.Result.List {
		if tickerData.AccountType == string(bybit.AccountTypeV5SPOT) {
			return tickerData, nil
		}
	}

	return bybit.V5WalletBalanceList{}, errs.ErrSpotBalanceNotFound
}

func GetPairSymbolPointerV5(pairSymbol string) *bybit.SymbolV5 {
	s := bybit.SymbolV5(pairSymbol)
	return &s
}

func GetOrderIDFromHistoryOrdersParam(param bybit.V5GetHistoryOrdersParam) string {
	if param.OrderID != nil {
		return *param.OrderID
	}
	if param.OrderLinkID != nil {
		return *param.OrderLinkID
	}

	return unknownOrderID
}

func GetOrderSymbolFromHistoryOrdersParam(param bybit.V5GetHistoryOrdersParam) string {
	if param.Symbol != nil {
		return string(*param.Symbol)
	}
	return unknownPairSymbol
}
