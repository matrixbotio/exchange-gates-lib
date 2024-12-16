package utils

import (
	"github.com/matrixbotio/exchange-gates-lib/internal/consts"
	pkgStructs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"
)

func GetTPOrderType(strategy pkgStructs.BotStrategy) pkgStructs.OrderSide {
	if strategy == pkgStructs.BotStrategyLong {
		return consts.OrderSideSell
	}
	return consts.OrderSideBuy
}
