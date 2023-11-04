package utils

import (
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
)

func GetTPOrderType(strategy pkgStructs.BotStrategy) string {
	if strategy == pkgStructs.BotStrategyLong {
		return pkgStructs.OrderTypeSell
	}
	return pkgStructs.OrderTypeBuy
}
