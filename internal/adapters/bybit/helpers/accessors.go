package helpers

import "github.com/hirokisan/bybit/v2"

const (
	unknownPairSymbol = "UNKNOWN"
	unknownOrderID    = "unknown-id"
)

func GetPairSymbolPointerV5(pairSymbol string) *bybit.SymbolV5 {
	s := bybit.SymbolV5(pairSymbol)
	return &s
}

func GetOrderIDFromHistoryOrdersParam(param bybit.V5GetHistoryOrdersParam) string {
	if param.OrderID != nil {
		return *param.OrderID
	}
	if param.OrderLinkID != nil {
		return *param.OrderLinkID
	}

	return unknownOrderID
}

func GetOrderSymbolFromHistoryOrdersParam(param bybit.V5GetHistoryOrdersParam) string {
	if param.Symbol != nil {
		return string(*param.Symbol)
	}
	return unknownPairSymbol
}
