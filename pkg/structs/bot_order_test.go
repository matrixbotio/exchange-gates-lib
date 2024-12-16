package structs

import (
	"testing"

	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/stretchr/testify/assert"
)

func TestBotOrderToString(t *testing.T) {
	// given
	order := BotOrder{
		PairSymbol: "LTCUSDC",
		Type:       consts.OrderSideBuy,
		Qty:        0.1,
		Price:      65,
	}

	// when
	result := order.String()

	// then
	assert.NotEmpty(t, result)
	assert.Contains(t, result, order.PairSymbol)
}
