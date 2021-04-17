package matrixgates

//ExchangeAdapter - abstract universal exchange adapter
type ExchangeAdapter struct {
	ExchangeID int
	Name       string
}

func newExchangeAdapter(name string, exchangeID int) *ExchangeAdapter {
	return &ExchangeAdapter{
		ExchangeID: exchangeID,
		Name:       name,
	}
}

//ExchangeAdapters - map of all supported exchanges
var ExchangeAdapters map[int]*ExchangeAdapter = map[int]*ExchangeAdapter{
	1: NewBinanceSpotAdapter(),
}
