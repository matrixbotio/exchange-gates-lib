package mappers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/adshao/go-binance/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
)

// ConvertOrderSide - convert order side to bot order type
func ConvertOrderSide(orderSide binance.SideType) string {
	return strings.ToLower(string(orderSide))
}

// GetBinanceOrderSide - convert bot order type to binance order type
func GetBinanceOrderSide(botOrderSide string) (binance.SideType, error) {
	switch botOrderSide {
	default:
		return "", fmt.Errorf("unknown order side: %q", botOrderSide)
	case pkgStructs.OrderTypeBuy:
		return binance.SideTypeBuy, nil
	case pkgStructs.OrderTypeSell:
		return binance.SideTypeSell, nil
	}
}

func ConvertBinanceToBotOrder(orderRes *binance.CreateOrderResponse) (
	structs.CreateOrderResponse,
	error,
) {
	orderResOrigQty, err := strconv.ParseFloat(orderRes.OrigQuantity, 64)
	if err != nil {
		return structs.CreateOrderResponse{},
			fmt.Errorf("parse order origQty: %w", err)
	}

	orderResPrice, err := strconv.ParseFloat(orderRes.Price, 64)
	if err != nil {
		return structs.CreateOrderResponse{},
			fmt.Errorf("parse order price: %w", err)
	}

	return structs.CreateOrderResponse{
		OrderID:       orderRes.OrderID,
		ClientOrderID: orderRes.ClientOrderID,
		OrigQuantity:  orderResOrigQty,
		Price:         orderResPrice,
		Symbol:        orderRes.Symbol,
		Type:          ConvertOrderSide(orderRes.Side),
	}, nil
}
