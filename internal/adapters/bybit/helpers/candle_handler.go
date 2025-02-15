package helpers

import (
	"fmt"

	"github.com/hirokisan/bybit/v2"

	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bybit/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
)

const defaultCandleInterval = "1m"

type CandleEventsHandler struct {
	symbols  symbolPerTopic
	callback func(event workers.CandleEvent)
}

type symbolPerTopic map[string]string

func (h *CandleEventsHandler) handle(e bybit.V5WebsocketPublicKlineResponse) error {
	for _, eventData := range e.Data {
		pairSymbol, isExists := h.symbols[e.Topic]
		if !isExists {
			return nil
		}

		event, err := mappers.ConvertWsCandle(pairSymbol, eventData)
		if err != nil {
			return fmt.Errorf("convert candle: %w", err)
		}

		h.callback(event)
	}
	return nil
}
