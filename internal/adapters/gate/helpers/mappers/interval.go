package mappers

import (
	"fmt"

	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
)

var intervalGateToOur = map[string]consts.Interval{
	"1m":  consts.Interval1min,
	"5m":  consts.Interval5min,
	"15m": consts.Interval15min,
	"30m": consts.Interval30min,
	"1h":  consts.Interval1hour,
	"4h":  consts.Interval4hour,
	"6h":  consts.Interval6hour,
	"12h": consts.Interval12hour,
	"1d":  consts.Interval1day,
}

var ourIntervalToGate = func() map[consts.Interval]string {
	r := map[consts.Interval]string{}
	for gateFormat, ourFormat := range intervalGateToOur {
		r[ourFormat] = gateFormat
	}
	return r
}()

func ConvertIntervalToGate(interval consts.Interval) (string, error) {
	result, isExists := ourIntervalToGate[interval]
	if !isExists {
		return "", fmt.Errorf("unknown interval: %q", interval)
	}

	return result, nil
}

func ConvertGateInterval(interval string) (consts.Interval, error) {
	result, isExists := intervalGateToOur[interval]
	if !isExists {
		return "", fmt.Errorf("unknown interval: %q", interval)
	}

	return result, nil
}
