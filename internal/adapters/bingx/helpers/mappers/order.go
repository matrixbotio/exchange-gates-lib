package mappers

import (
	"errors"
	"fmt"
	"strconv"

	bingxgo "github.com/Sagleft/go-bingx"

	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
)

const (
	sideBuy  = "BUY"
	sideSell = "SELL"
)

// exchange const -> our const
var orderStatusConvertor = map[string]consts.OrderStatus{
	"NEW":              pkgStructs.OrderStatusNew,
	"PENDING":          pkgStructs.OrderStatusNew,
	"PARTIALLY_FILLED": pkgStructs.OrderStatusPartiallyFilled,
	"FILLED":           pkgStructs.OrderStatusFilled,
	"CANCELED":         pkgStructs.OrderStatusCancelled,
	"FAILED":           pkgStructs.OrderStatusRejected,
}

func ConvertBingXStatus(status string) (consts.OrderStatus, error) {
	result, isExists := orderStatusConvertor[status]
	if !isExists {
		return "", fmt.Errorf("unknown status: %q", status)
	}
	return result, nil
}

func ConvertBingXSide(side string) (consts.OrderSide, error) {
	if side == "" {
		return "", errors.New("not set")
	}

	switch side {
	default:
		return "", fmt.Errorf("unknown: %q", side)
	case sideBuy:
		return consts.OrderSideBuy, nil
	case sideSell:
		return consts.OrderSideSell, nil
	}
}

func ConvertBingXOrderData(data *bingxgo.SpotOrder) (structs.OrderData, error) {
	if data == nil {
		return structs.OrderData{}, errors.New("order data not set")
	}

	orderPrice, err := strconv.ParseFloat(data.Price, 64)
	if err != nil {
		return structs.OrderData{}, fmt.Errorf("price: %w", err)
	}

	orderStatus, err := ConvertBingXStatus(data.Status)
	if err != nil {
		return structs.OrderData{}, fmt.Errorf("convert status: %w", err)
	}

	orderQty, err := strconv.ParseFloat(data.OrigQty, 64)
	if err != nil {
		return structs.OrderData{}, fmt.Errorf("qty: %w", err)
	}

	orderFilledQty, err := strconv.ParseFloat(data.ExecutedQty, 64)
	if err != nil {
		return structs.OrderData{}, fmt.Errorf("filled qty: %w", err)
	}

	orderSide, err := ConvertBingXSide(data.Side)
	if err != nil {
		return structs.OrderData{}, fmt.Errorf("convert side: %w", err)
	}

	return structs.OrderData{
		OrderID:       int64(data.OrderID),
		ClientOrderID: data.ClientOrderID,
		Status:        orderStatus,
		AwaitQty:      orderQty,
		FilledQty:     orderFilledQty,
		Price:         orderPrice,
		Symbol:        data.Symbol,
		Side:          orderSide,
		CreatedTime:   data.Time,
		UpdatedTime:   data.UpdateTime,
	}, nil
}
