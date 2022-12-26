package conditions

import (
	"testing"

	"github.com/hirokisan/bybit/v2"
)

func Test_isSpotTradingAvailable(t *testing.T) {
	type args struct {
		tradePairStatus bybit.InstrumentStatus
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "available",
			args: args{tradePairStatus: bybit.InstrumentStatusAvailable},
			want: true,
		},
		{
			name: "trading",
			args: args{tradePairStatus: bybit.InstrumentStatusTrading},
			want: true,
		},
		{
			name: "closed",
			args: args{tradePairStatus: bybit.InstrumentStatusClosed},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSpotTradingAvailable(tt.args.tradePairStatus); got != tt.want {
				t.Errorf("isSpotTradingAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}
