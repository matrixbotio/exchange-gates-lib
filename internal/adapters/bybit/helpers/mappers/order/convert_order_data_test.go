package order_mappers

import (
	"reflect"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/hirokisan/bybit/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
	"github.com/stretchr/testify/require"
)

func Test_convertOrderData(t *testing.T) {
	type args struct {
		data bybit.V5GetOrder
	}
	tests := []struct {
		name    string
		args    args
		want    structs.OrderData
		wantErr bool
	}{
		{
			name:    "empty ID",
			args:    args{data: bybit.V5GetOrder{}},
			want:    structs.OrderData{},
			wantErr: true,
		},
		{
			name:    "empty qty",
			args:    args{data: bybit.V5GetOrder{OrderID: "X"}},
			want:    structs.OrderData{},
			wantErr: true,
		},
		{
			name: "empty exec qty",
			args: args{data: bybit.V5GetOrder{
				OrderID: "X",
				Qty:     "0.56",
			}},
			want:    structs.OrderData{},
			wantErr: true,
		},
		{
			name: "empty price",
			args: args{data: bybit.V5GetOrder{
				OrderID:    "X",
				Qty:        "0.56",
				CumExecQty: "0",
			}},
			want:    structs.OrderData{},
			wantErr: true,
		},
		{
			name: "empty time",
			args: args{data: bybit.V5GetOrder{
				OrderID:    "X",
				Qty:        "0.56",
				CumExecQty: "0",
				Price:      "80",
			}},
			want:    structs.OrderData{},
			wantErr: true,
		},
		{
			name: "empty side",
			args: args{data: bybit.V5GetOrder{
				OrderID:     "X",
				Qty:         "0.56",
				CumExecQty:  "0",
				Price:       "80",
				UpdatedTime: "12345",
			}},
			want:    structs.OrderData{},
			wantErr: true,
		},
		{
			name: "empty status",
			args: args{data: bybit.V5GetOrder{
				OrderID:     "X",
				Qty:         "0.56",
				CumExecQty:  "0",
				Price:       "80",
				UpdatedTime: "12345",
				Side:        bybit.SideSell,
			}},
			want:    structs.OrderData{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertOrderData(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertOrderData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertOrderData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertOrderStatus(t *testing.T) {
	type args struct {
		status bybit.OrderStatus
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "order cancelled",
			args:    args{status: bybit.OrderStatusCancelled},
			want:    pkgStructs.OrderStatusCancelled,
			wantErr: false,
		},
		{
			name:    "order created",
			args:    args{status: bybit.OrderStatusCreated},
			want:    pkgStructs.OrderStatusNew,
			wantErr: false,
		},
		{
			name:    "order unknown",
			args:    args{status: bybit.OrderStatus("wtf")},
			want:    pkgStructs.OrderStatusUnknown,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertOrderStatus(tt.args.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertOrderStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("convertOrderStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertOrderType(t *testing.T) {
	type args struct {
		side bybit.Side
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "buy",
			args:    args{side: bybit.SideBuy},
			want:    pkgStructs.OrderTypeBuy,
			wantErr: false,
		},
		{
			name:    "sell",
			args:    args{side: bybit.SideSell},
			want:    pkgStructs.OrderTypeSell,
			wantErr: false,
		},
		{
			name:    "unknown",
			args:    args{side: bybit.SideNone},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertOrderType(tt.args.side)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertOrderType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("convertOrderType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertOrderData(t *testing.T) {
	// given
	rawOrder := bybit.V5GetOrder{
		Symbol:      bybit.SymbolV5("LTCUSDT"),
		OrderType:   bybit.OrderTypeLimit,
		OrderID:     "12345",
		Side:        bybit.SideSell,
		Qty:         "0.35",
		Price:       "80.156",
		CumExecQty:  "0.1",
		UpdatedTime: "1692119310600",
		OrderStatus: bybit.OrderStatusActive,
	}

	// when
	orderData, err := ConvertOrderData(rawOrder)

	// then
	require.NoError(t, err)
	assert.Equal(t, string(rawOrder.Symbol), orderData.Symbol)
	assert.Equal(t, pkgStructs.OrderTypeSell, orderData.Type)
	assert.Equal(t, int64(12345), orderData.OrderID)
	assert.Equal(t, float64(0.35), orderData.AwaitQty)
	assert.Equal(t, float64(80.156), orderData.Price)
	assert.Equal(t, float64(0.1), orderData.FilledQty)
	assert.Equal(t, pkgStructs.OrderStatusNew, orderData.Status)
}

func TestConvertOrderSideToBybit(t *testing.T) {
	// given
	side := "buy"
	sideExpected := bybit.Side("Buy")

	// when
	sideFormatted := ConvertOrderSideToBybit(side)

	// then
	assert.Equal(t, sideExpected, sideFormatted)
}
