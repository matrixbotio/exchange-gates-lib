package mappers

import (
	"fmt"
	"strings"

	"github.com/adshao/go-binance/v2"
	"github.com/matrixbotio/exchange-gates-lib/pkg/structs"
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
	case structs.OrderTypeBuy:
		return binance.SideTypeBuy, nil
	case structs.OrderTypeSell:
		return binance.SideTypeSell, nil
	}
}
