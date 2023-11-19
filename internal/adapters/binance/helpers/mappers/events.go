package mappers

import (
	"strconv"

	"github.com/adshao/go-binance/v2"
)

func ConvertPriceEvent(event binance.WsBookTickerEvent) (ask, bid float64, err error) {
	ask, err = strconv.ParseFloat(event.BestAskPrice, 64)
	if err != nil {
		return
	}

	bid, err = strconv.ParseFloat(event.BestBidPrice, 64)
	return
}
