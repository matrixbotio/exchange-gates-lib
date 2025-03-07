package mappers

import (
	"fmt"

	"github.com/gateio/gateapi-go/v6"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	"github.com/shopspring/decimal"
)

func ConvertBalances(balances []gateapi.SpotAccount) ([]structs.Balance, error) {
	r := []structs.Balance{}
	for _, data := range balances {
		b, err := parseAssetBalance(data)
		if err != nil {
			return nil, fmt.Errorf("parse: %w", err)
		}

		r = append(r, b)
	}
	return r, nil
}

func parseAssetBalance(balance gateapi.SpotAccount) (structs.Balance, error) {
	assetFree, err := decimal.NewFromString(balance.Available)
	if err != nil {
		return structs.Balance{}, fmt.Errorf("available: %w", err)
	}

	assetLocked, err := decimal.NewFromString(balance.Locked)
	if err != nil {
		return structs.Balance{}, fmt.Errorf("locked: %w", err)
	}

	return structs.Balance{
		Asset:  balance.Currency,
		Free:   assetFree.InexactFloat64(),
		Locked: assetLocked.InexactFloat64(),
	}, nil
}
