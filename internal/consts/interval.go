package consts

import (
	"errors"
	"fmt"
	"reflect"
)

type Interval string

var intervalType = reflect.TypeOf(Interval(""))

var errIntervalNotSet = errors.New("not set")

const (
	Interval1min   Interval = "1m"
	Interval5min   Interval = "5m"
	Interval15min  Interval = "15m"
	Interval30min  Interval = "30m"
	Interval1hour  Interval = "1h"
	Interval6hour  Interval = "6h"
	Interval12hour Interval = "12h"
	Interval1day   Interval = "1d"
)

var allIntervals = GetIntervals()

func GetIntervals() []Interval {
	return []Interval{
		Interval1min,
		Interval5min,
		Interval15min,
		Interval30min,
		Interval1hour,
		Interval6hour,
		Interval1day,
	}
}

func ValidateInterval(interval string) error {
	if interval == "" {
		return errIntervalNotSet
	}

	for _, validInterval := range allIntervals {
		if string(validInterval) == interval {
			return nil
		}
	}

	return fmt.Errorf("invalid interval: %s", interval)
}
