package mappers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hirokisan/bybit/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
)

// old -> new
var orderStatusConvertor = map[bybit.OrderStatus]string{
	bybit.OrderStatusCreated:         pkgStructs.OrderStatusNew,
	bybit.OrderStatusRejected:        pkgStructs.OrderStatusRejected,
	bybit.OrderStatusNew:             pkgStructs.OrderStatusNew,
	bybit.OrderStatusPartiallyFilled: pkgStructs.OrderStatusPartiallyFilled,
	bybit.OrderStatusFilled:          pkgStructs.OrderStatusFilled,
	bybit.OrderStatusCancelled:       pkgStructs.OrderStatusCancelled,
	bybit.OrderStatusPendingCancel:   pkgStructs.OrderStatusPendingCancel,
	bybit.OrderStatusUntriggered:     pkgStructs.OrderStatusUntriggered,
	bybit.OrderStatusDeactivated:     pkgStructs.OrderStatusDeactivated,
	bybit.OrderStatusTriggered:       pkgStructs.OrderStatusTriggered,
	bybit.OrderStatusActive:          pkgStructs.OrderStatusNew,
}

func ConvertOrderData(data bybit.V5GetOrder) (structs.OrderData, error) {
	if data.OrderID == "" {
		return structs.OrderData{}, errors.New("order ID is empty")
	}
	orderID, err := strconv.ParseInt(data.OrderID, 10, 64)
	if err != nil {
		return structs.OrderData{}, fmt.Errorf("parse order ID: %w", err)
	}

	if data.Qty == "" {
		return structs.OrderData{}, errors.New("order qty is empty")
	}
	awaitQty, err := strconv.ParseFloat(data.Qty, 64)
	if err != nil {
		return structs.OrderData{}, fmt.Errorf("parse qty: %w", err)
	}

	if data.CumExecQty == "" {
		return structs.OrderData{}, errors.New("order executed qty is empty")
	}
	filledQty, err := strconv.ParseFloat(data.CumExecQty, 64)
	if err != nil {
		return structs.OrderData{}, fmt.Errorf("parse executed qty: %w", err)
	}

	if data.Price == "" {
		return structs.OrderData{}, errors.New("order price is empty")
	}
	price, err := strconv.ParseFloat(data.Price, 64)
	if err != nil {
		return structs.OrderData{}, fmt.Errorf("parse price: %w", err)
	}

	if data.UpdatedTime == "" {
		return structs.OrderData{}, errors.New("order time is empty")
	}
	updatedTime, err := strconv.ParseInt(data.UpdatedTime, 10, 64)
	if err != nil {
		return structs.OrderData{}, fmt.Errorf("parse time: %w", err)
	}

	if data.Side == "" {
		return structs.OrderData{}, errors.New("order side is empty")
	}
	orderType, err := convertOrderType(data.Side)
	if err != nil {
		return structs.OrderData{}, fmt.Errorf("get order type: %w", err)
	}

	if string(data.OrderStatus) == "" {
		return structs.OrderData{}, errors.New("order status is empty")
	}
	orderStatus, err := convertOrderStatus(data.OrderStatus)
	if err != nil {
		return structs.OrderData{}, fmt.Errorf("convert order status: %w", err)
	}

	return structs.OrderData{
		OrderID:       orderID,
		ClientOrderID: data.OrderLinkID,
		Status:        orderStatus,
		AwaitQty:      awaitQty,
		FilledQty:     filledQty,
		Price:         price,
		Symbol:        string(data.Symbol),
		Type:          orderType,
		UpdatedTime:   updatedTime,
	}, nil
}

func convertOrderType(side bybit.Side) (string, error) {
	orderType := strings.ToLower(string(side))

	if orderType != pkgStructs.OrderTypeBuy &&
		orderType != pkgStructs.OrderTypeSell {
		return "", fmt.Errorf("unknown order type: %q", string(side))
	}

	return orderType, nil
}

func convertOrderStatus(status bybit.OrderStatus) (string, error) {
	formattedStatus, isExists := orderStatusConvertor[status]
	if !isExists {
		return pkgStructs.OrderStatusUnknown, fmt.Errorf("uknown status: %q", string(status))
	}
	return formattedStatus, nil
}

func ParseHistoryOrder(
	ordersResponse *bybit.V5GetOrdersResponse,
	orderID string,
	pairSymbol string,
) (structs.OrderData, error) {
	if len(ordersResponse.Result.List) == 0 {
		return structs.OrderData{}, fmt.Errorf(
			"order %q in %q not found",
			orderID, pairSymbol,
		)
	}

	return ConvertOrderData(ordersResponse.Result.List[0])
}

func ConvertOrderSideToBybit(side string) bybit.Side {
	return bybit.Side(cases.Title(language.Und, cases.NoLower).String(side))
}
