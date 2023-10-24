package order_mappers

import (
	"testing"

	"github.com/hirokisan/bybit/v2"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseBuyOrderExecFee(t *testing.T) {
	// given
	orderSide := pkgStructs.OrderTypeBuy
	orderExecData := bybit.V5GetExecutionListResult{
		List: []bybit.V5GetExecutionListItem{
			{
				OrderQty: "0.124",
				ExecFee:  "0.000124",
			},
			{
				ExecFee: "0.00002",
			},
		},
	}

	// when
	fees, err := ParseOrderExecFee(orderExecData, orderSide)

	// then
	require.NoError(t, err)
	assert.Equal(t, decimal.NewFromFloat(float64(0.000144)), fees.BaseAsset)
	assert.Equal(t, decimal.NewFromFloat(0), fees.QuoteAsset)
}

func TestParseSellOrderExecFee(t *testing.T) {
	// given
	orderSide := pkgStructs.OrderTypeSell
	orderExecData := bybit.V5GetExecutionListResult{
		List: []bybit.V5GetExecutionListItem{
			{
				OrderQty: "0.124",
				ExecFee:  "0.000124",
			},
			{
				ExecFee: "0.00002",
			},
		},
	}

	// when
	fees, err := ParseOrderExecFee(orderExecData, orderSide)

	// then
	require.NoError(t, err)
	assert.Equal(t, decimal.NewFromFloat(float64(0.000144)), fees.QuoteAsset)
	assert.Equal(t, decimal.NewFromFloat(0), fees.BaseAsset)
}

func TestParseOrderExecFeeZero(t *testing.T) {
	// given
	orderSide := pkgStructs.OrderTypeBuy
	orderExecData := bybit.V5GetExecutionListResult{
		List: []bybit.V5GetExecutionListItem{
			{
				ExecFee: "0",
			},
			{},
			{},
		},
	}

	// when
	fees, err := ParseOrderExecFee(orderExecData, orderSide)

	// then
	require.NoError(t, err)
	assert.Equal(t, decimal.NewFromFloat(float64(0)), fees.BaseAsset)
	assert.Equal(t, decimal.NewFromFloat(float64(0)), fees.QuoteAsset)
}
