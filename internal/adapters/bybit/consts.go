package bybit

import (
	"time"
)

var candleIntervalsToBybit = map[string]intervalData{
	"1m":  {"1", time.Minute},
	"3m":  {"3", time.Minute * 3},
	"5m":  {"5", time.Minute * 5},
	"15m": {"15", time.Minute * 15},
	"30m": {"30", time.Minute * 30},
	"1h":  {"60", time.Hour},
	"2h":  {"120", time.Hour * 2},
	"4h":  {"240", time.Hour * 4},
	"5h":  {"360", time.Hour * 5},
	"6h":  {"720", time.Hour * 6},
	"1d":  {"D", time.Hour * 24},
	"1w":  {"W", time.Hour * 24 * 7},
	"1M":  {"M", time.Hour * 24 * 30},
}

type intervalData struct {
	Code     string
	Duration time.Duration
}
