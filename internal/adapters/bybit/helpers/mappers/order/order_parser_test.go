package order_mappers

import (
	"testing"

	"github.com/hirokisan/bybit/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseOrderExecFee(t *testing.T) {
	// given
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
	fees, err := ParseOrderExecFee(orderExecData)

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(0.000144), fees)
}

func TestParseOrderExecFeeZero(t *testing.T) {
	// given
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
	fees, err := ParseOrderExecFee(orderExecData)

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(0), fees)
}
