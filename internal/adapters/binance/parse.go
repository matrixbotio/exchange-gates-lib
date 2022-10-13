package binance

import (
	"errors"
	"strconv"

	"github.com/adshao/go-binance/v2"
)

func parseOrderOriginalQty(orderRaw *binance.Order) (float64, error) {
	awaitQty, err := strconv.ParseFloat(orderRaw.OrigQuantity, 64)
	if err != nil {
		return 0, errors.New("failed to parse order original qty: " + err.Error())
	}
	return awaitQty, nil
}

func parseOrderExecutedQty(orderRaw *binance.Order) (float64, error) {
	filledQty, err := strconv.ParseFloat(orderRaw.ExecutedQuantity, 64)
	if err != nil {
		return 0, errors.New("failed to parse order executed qty: " + err.Error())
	}
	return filledQty, nil
}

func parseOrderPrice(orderRaw *binance.Order) (float64, error) {
	price, err := strconv.ParseFloat(orderRaw.Price, 64)
	if err != nil {
		return 0, errors.New("failed to parse order price: " + err.Error())
	}
	return price, nil
}
