package matrixgates

import (
	"github.com/matrixbotio/exchange-gates/workers"
)

// ExchangeAdapter - abstract universal exchange adapter
type ExchangeAdapter struct {
	ExchangeID int
	Name       string
	Tag        string
}

// Placeholders

// Connect to exchange
func (a *ExchangeAdapter) Connect(credentials APICredentials) error {
	return nil
}

// PlaceOrder - place order on exchange
func (a *ExchangeAdapter) PlaceOrder(order BotOrder, pairLimits ExchangePairData) (*CreateOrderResponse, error) {
	return nil, nil
}

// GetAccountData ..
func (a *ExchangeAdapter) GetAccountData() (*AccountData, error) {
	return nil, nil
}

// GetPairLastPrice ..
func (a *ExchangeAdapter) GetPairLastPrice(pairSymbol string) (float64, error) {
	return 0, nil
}

// CancelPairOrder ..
func (a *ExchangeAdapter) CancelPairOrder(pairSymbol string, orderID int64) error {
	return nil
}

// CancelPairOrders ..
func (a *ExchangeAdapter) CancelPairOrders(pairSymbol string) error {
	return nil
}

// GetOrderData ..
func (a *ExchangeAdapter) GetOrderData(pairSymbol string, orderID int64) (*TradeEventData, error) {
	return nil, nil
}

// GetPairOpenOrders ..
func (a *ExchangeAdapter) GetPairOpenOrders(pairSymbol string) ([]*Order, error) {
	// TODO
	return nil, nil
}

// GetPairs get all Binance pairs
func (a *ExchangeAdapter) GetPairs() ([]*ExchangePairData, error) {
	return nil, nil
}

// VerifyAPIKeys ..
func (a *ExchangeAdapter) VerifyAPIKeys(keyPublic, keySecret string) error {
	return nil
}

// GetPriceWorker - create empty market data worker
func (a *ExchangeAdapter) GetPriceWorker() workers.IPriceWorker {
	return nil
}

// GetCandleWorker - create empty market candle worker
func (a *ExchangeAdapter) GetCandleWorker() workers.ICandleWorker {
	return nil
}
