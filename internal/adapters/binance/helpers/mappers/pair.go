package mappers

import (
	"fmt"
	"strconv"

	"github.com/adshao/go-binance/v2"
)

func GetPairPrice(prices []*binance.SymbolPrice, pairSymbol string) (float64, error) {
	for _, p := range prices {
		if p.Symbol == pairSymbol {
			price, err := strconv.ParseFloat(p.Price, 64)
			if err != nil {
				return 0, fmt.Errorf("parse price %q: %w", p.Price, err)
			}
			return price, nil
		}
	}
	return 0, fmt.Errorf("last price not found for pair %q", pairSymbol)
}
