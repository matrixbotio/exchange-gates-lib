package mappers

import (
	"fmt"
	"strconv"

	"github.com/adshao/go-binance/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
)

// ConvertOrderSide - convert order side from binance format to bot order side
func ConvertOrderSide(orderSide binance.SideType) (string, error) {
	switch orderSide {
	default:
		return "", fmt.Errorf("unknown order side: %q", orderSide)
	case binance.SideTypeBuy:
		return pkgStructs.OrderTypeBuy, nil
	case binance.SideTypeSell:
		return pkgStructs.OrderTypeSell, nil
	}
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

func ConvertPlacedOrder(orderResponse binance.CreateOrderResponse) (
	structs.CreateOrderResponse,
	error,
) {
	orderResOrigQty, err := strconv.ParseFloat(orderResponse.OrigQuantity, 64)
	if err != nil {
		return structs.CreateOrderResponse{},
			fmt.Errorf("parse order origQty: %w", err)
	}

	orderResPrice, err := strconv.ParseFloat(orderResponse.Price, 64)
	if err != nil {
		return structs.CreateOrderResponse{},
			fmt.Errorf("parse order price: %w", err)
	}

	orderSide, err := ConvertOrderSide(orderResponse.Side)
	if err != nil {
		return structs.CreateOrderResponse{},
			fmt.Errorf("convert order side: %w", err)
	}

	return structs.CreateOrderResponse{
		OrderID:       orderResponse.OrderID,
		ClientOrderID: orderResponse.ClientOrderID,
		OrigQuantity:  orderResOrigQty,
		Price:         orderResPrice,
		Symbol:        orderResponse.Symbol,
		Type:          orderSide,
	}, nil
}

// ConvertOrderData converting the order data from binance to our format
func ConvertOrderData(orderResponse *binance.Order) (structs.OrderData, error) {
	awaitQty, err := strconv.ParseFloat(orderResponse.OrigQuantity, 64)
	if err != nil {
		return structs.OrderData{}, fmt.Errorf("parse await qty: %w", err)
	}

	filledQty, err := strconv.ParseFloat(orderResponse.ExecutedQuantity, 64)
	if err != nil {
		return structs.OrderData{}, fmt.Errorf("parse executed qty: %w", err)
	}

	price, err := strconv.ParseFloat(orderResponse.Price, 64)
	if err != nil {
		return structs.OrderData{}, fmt.Errorf("parse price: %w", err)
	}

	orderSide, err := ConvertOrderSide(orderResponse.Side)
	if err != nil {
		return structs.OrderData{},
			fmt.Errorf("convert order side: %w", err)
	}

	return structs.OrderData{
		OrderID:       orderResponse.OrderID,
		ClientOrderID: orderResponse.ClientOrderID,
		Status:        string(orderResponse.Status),
		AwaitQty:      awaitQty,
		FilledQty:     filledQty,
		Price:         price,
		Symbol:        orderResponse.Symbol,
		Type:          orderSide,
		CreatedTime:   orderResponse.Time,
		UpdatedTime:   orderResponse.UpdateTime,
	}, nil
}

func ConvertOrders(ordersRaw []*binance.Order) ([]structs.OrderData, error) {
	orders := []structs.OrderData{}
	for _, orderRaw := range ordersRaw {
		order, err := ConvertOrderData(orderRaw)
		if err != nil {
			return nil, fmt.Errorf("convert order: %w", err)
		}

		orders = append(orders, order)
	}
	return orders, nil
}

func GetTestPairDataFilters() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"filterType": string(binance.SymbolFilterTypeLotSize),
			"maxQty":     "999999",
			"minQty":     "0.0001",
			"stepSize":   "0.0001",
		},
		{
			"filterType": string(binance.SymbolFilterTypePriceFilter),
			"maxPrice":   "1000000",
			"minPrice":   "0.01",
			"tickSize":   "0.01",
		},
		{
			"filterType":    string(binance.SymbolFilterTypeNotional),
			"minNotional":   "1",
			"avgPriceMins":  0.01,
			"applyToMarket": true,
		},
		{
			"filterType": string(binance.SymbolFilterTypeMarketLotSize),
			"maxQty":     "999999",
			"minQty":     "0.0001",
			"stepSize":   "0.0001",
		},
	}
}
