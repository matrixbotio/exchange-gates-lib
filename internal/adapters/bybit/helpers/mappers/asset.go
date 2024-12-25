package mappers

import (
	"fmt"
	"strconv"

	"github.com/hirokisan/bybit/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
)

func ConvertBalances(
	coins []bybit.V5WalletBalanceCoin,
	tickerTag string,
) (structs.AssetBalance, error) {
	for _, coinData := range coins {
		if tickerTag != string(coinData.Coin) {
			continue
		}

		return convertCoinData(coinData, tickerTag)
	}

	return structs.AssetBalance{}, fmt.Errorf(
		"balance not found for ticker %q",
		tickerTag,
	)
}

func convertCoinData(
	coinData bybit.V5WalletBalanceCoin,
	tickerTag string,
) (structs.AssetBalance, error) {
	var tickerFree float64
	var tickerLocked float64
	var err error

	if coinData.Locked != "" {
		tickerLocked, err = strconv.ParseFloat(coinData.Locked, 64)
		if err != nil {
			return structs.AssetBalance{}, fmt.Errorf("parse locked balance: %w", err)
		}
	}

	if coinData.Free != "" {
		tickerFree, err = strconv.ParseFloat(coinData.Free, 64)
	} else if coinData.Equity != "" {
		equity, err := strconv.ParseFloat(coinData.Equity, 64)
		if err != nil {
			return structs.AssetBalance{}, fmt.Errorf("parse equity: %w", err)
		}

		tickerFree = equity - tickerLocked
	}
	if err != nil {
		return structs.AssetBalance{}, fmt.Errorf("parse free balance: %w", err)
	}

	return structs.AssetBalance{
		Ticker: tickerTag,
		Free:   tickerFree,
		Locked: tickerLocked,
	}, nil
}
