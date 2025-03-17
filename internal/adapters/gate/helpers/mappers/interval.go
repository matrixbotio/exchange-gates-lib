package mappers

import (
	"fmt"

	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
)

var ourIntervalToGate = map[consts.Interval]string{
	consts.Interval1min:   "1m",
	consts.Interval5min:   "5m",
	consts.Interval15min:  "15m",
	consts.Interval30min:  "30m",
	consts.Interval1hour:  "1h",
	consts.Interval4hour:  "4h",
	consts.Interval6hour:  "6h",  // TODO: test it
	consts.Interval12hour: "12h", // TODO: test it
	consts.Interval1day:   "1d",
}

func ConvertIntervalToGate(interval consts.Interval) (string, error) {
	result, isExists := ourIntervalToGate[interval]
	if !isExists {
		return "", fmt.Errorf("unknown interval: %q", interval)
	}
	return result, nil
}
