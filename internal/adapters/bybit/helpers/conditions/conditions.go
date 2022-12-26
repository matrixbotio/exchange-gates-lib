package conditions

import "github.com/hirokisan/bybit/v2"

func IsSpotTradingAvailable(tradePairStatus bybit.InstrumentStatus) bool {
	return tradePairStatus == bybit.InstrumentStatusTrading ||
		tradePairStatus == bybit.InstrumentStatusAvailable
}
