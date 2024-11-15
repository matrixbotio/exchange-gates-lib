package helpers

import (
	"fmt"

	"github.com/hirokisan/bybit/v2"

	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bybit/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
)

const defaultCandleInterval = bybit.Interval("1m")

type CandleEventsHandler struct {
	pairSymbol string
	callback   func(event workers.CandleEvent)
}

func (h *CandleEventsHandler) handle(e bybit.V5WebsocketPublicKlineResponse) error {
	for _, eventData := range e.Data {
		event, err := mappers.ConvertWsCandle(h.pairSymbol, eventData)
		if err != nil {
			return fmt.Errorf("convert candle: %w", err)
		}

		h.callback(event)
	}
	return nil
}
