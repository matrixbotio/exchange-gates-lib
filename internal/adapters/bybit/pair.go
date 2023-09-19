package bybit

import (
	"fmt"
	"strconv"

	"github.com/hirokisan/bybit/v2"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bybit/helpers/accessors"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/bybit/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
)

func (a *adapter) GetPairData(pairSymbol string) (structs.ExchangePairData, error) {
	response, err := a.getTradePairs(pairSymbol)
	if err != nil {
		return structs.ExchangePairData{}, fmt.Errorf("get instruments info: %w", err)
	}

	pairsData, err := mappers.ConvertPairsData(response.Result.Spot, a.GetID())
	if err != nil {
		return structs.ExchangePairData{}, fmt.Errorf("convert pairs: %w", err)
	}

	if len(pairsData) == 0 {
		return structs.ExchangePairData{}, fmt.Errorf("pair %q data not available", pairSymbol)
	}

	return pairsData[0], nil
}

func (a *adapter) GetPairLastPrice(pairSymbol string) (float64, error) {
	response, err := a.client.V5().Market().GetTickers(bybit.V5GetTickersParam{
		Category: bybit.CategoryV5Spot,
		Symbol:   accessors.GetPairSymbolPointerV5(pairSymbol),
	})
	if err != nil {
		return 0, fmt.Errorf("get %s last price: %w", pairSymbol, err)
	}

	if len(response.Result.Spot.List) == 0 {
		return 0, fmt.Errorf("last price for pair %q not found", pairSymbol)
	}

	lastPrice := response.Result.Spot.List[0].LastPrice

	price, err := strconv.ParseFloat(lastPrice, 64)
	if err != nil {
		return 0, fmt.Errorf(
			"parse last price for pair %q: %q: %w",
			lastPrice, pairSymbol, err,
		)
	}
	return price, nil
}

func (a *adapter) GetPairOpenOrders(pairSymbol string) ([]structs.OrderData, error) {
	response, err := a.client.V5().Order().GetOpenOrders(bybit.V5GetOpenOrdersParam{
		Category: bybit.CategoryV5Spot,
		Symbol:   accessors.GetPairSymbolPointerV5(pairSymbol),
	})
	if err != nil {
		return nil, fmt.Errorf("get open orders: %w", err)
	}

	var result []structs.OrderData
	for _, rawOrderData := range response.Result.List {
		orderData, err := mappers.ConvertOrderData(rawOrderData)
		if err != nil {
			return nil, fmt.Errorf("convert order: %w", err)
		}

		result = append(result, orderData)
	}
	return result, nil
}

// TBD: remove: https://github.com/matrixbotio/exchange-gates-lib/issues/149
func (a *adapter) GetPairOrdersHistory(task structs.GetOrdersHistoryTask) (
	[]structs.OrderData,
	error,
) {
	return nil, nil
}

func (a *adapter) GetPairs() ([]structs.ExchangePairData, error) {
	response, err := a.getTradePairs()
	if err != nil {
		return nil, fmt.Errorf("get instruments info: %w", err)
	}

	pairsData, err := mappers.ConvertPairsData(response.Result.Spot, a.GetID())
	if err != nil {
		return nil, fmt.Errorf("convert pairs: %w", err)
	}

	return pairsData, nil
}

func (a *adapter) GetAccountBalance() ([]structs.Balance, error) {
	response, err := a.client.V5().Account().GetWalletBalance(
		bybit.AccountType(bybit.AccountTypeV5SPOT),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("get wallet balance: %w", err)
	}

	if response == nil {
		return nil, fmt.Errorf("get wallet balance: response is empty")
	}
	return mappers.ConvertAccountBalance(*response)
}

func (a *adapter) GetPairBalance(pair structs.PairSymbolData) (structs.PairBalance, error) {
	baseTickerBalance, err := a.getTickerBalance(pair.BaseTicker)
	if err != nil {
		return structs.PairBalance{}, fmt.Errorf("get base asset balance: %w", err)
	}

	quoteTickerBalance, err := a.getTickerBalance(pair.QuoteTicker)
	if err != nil {
		return structs.PairBalance{}, fmt.Errorf("get quote asset balance: %w", err)
	}

	return structs.PairBalance{
		BaseAsset:  &baseTickerBalance,
		QuoteAsset: &quoteTickerBalance,
	}, nil
}

func (a *adapter) getTradePairs(symbol ...string) (*bybit.V5GetInstrumentsInfoResponse, error) {
	args := bybit.V5GetInstrumentsInfoParam{
		Category: bybit.CategoryV5Spot,
	}
	if len(symbol) > 0 {
		args.Symbol = accessors.GetPairSymbolPointerV5(symbol[0])
	}

	response, err := a.client.V5().Market().GetInstrumentsInfo(args)
	if err != nil {
		return nil, fmt.Errorf("get info: %w", err)
	}
	return response, nil
}
