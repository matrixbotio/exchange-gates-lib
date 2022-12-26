package mappers

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hirokisan/bybit/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
)

func ParsePriceEvent(
	r bybit.V5WebsocketPublicTickerResponse,
	exchangeTag string,
) (workers.PriceEvent, error) {
	if r.Data.Spot == nil {
		return workers.PriceEvent{}, errors.New("price event data is empty")
	}

	if r.Data.Spot.LastPrice == "" {
		return workers.PriceEvent{}, errors.New("last price is empty")
	}
	lastPrice, err := strconv.ParseFloat(r.Data.Spot.LastPrice, 64)
	if err != nil {
		return workers.PriceEvent{}, fmt.Errorf("parse last price: %w", err)
	}

	return workers.PriceEvent{
		ExchangeTag: exchangeTag,
		Symbol:      string(r.Data.Spot.Symbol),
		Ask:         lastPrice,
		Bid:         lastPrice,
	}, nil
}
