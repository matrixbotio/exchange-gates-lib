package mappers

import (
	"encoding/json"
	"testing"

	"github.com/adshao/go-binance/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestBinancePairData() binance.Symbol {
	return binance.Symbol{
		Symbol:               "LTCUSDC",
		Status:               string(binance.SymbolStatusTypeTrading),
		BaseAsset:            "LTC",
		QuoteAsset:           "USDC",
		BaseAssetPrecision:   8,
		QuotePrecision:       4,
		IsSpotTradingAllowed: true,
		Filters:              GetTestPairDataFilters(),
	}
}

func TestGetPairPriceSuccess(t *testing.T) {
	// given
	pairSymbol := "LTCUSDT"
	prices := []*binance.SymbolPrice{
		{
			Symbol: "USDCUSDT",
			Price:  "1.001",
		},
		{
			Symbol: pairSymbol,
			Price:  "65",
		},
	}

	// when
	lastPrice, err := GetPairPrice(prices, pairSymbol)

	// then
	require.NoError(t, err)
	assert.Equal(t, float64(65), lastPrice)
}

func TestGetPairPriceParseError(t *testing.T) {
	// given
	pairSymbol := "USDCUSDT"
	prices := []*binance.SymbolPrice{
		{
			Symbol: pairSymbol,
			Price:  "1-001",
		},
	}

	// when
	_, err := GetPairPrice(prices, pairSymbol)

	// then
	require.Error(t, err)
}

func TestGetPairPriceNotFound(t *testing.T) {
	// given
	pairSymbol := "MTXBUSDC"
	prices := []*binance.SymbolPrice{
		{
			Symbol: "USDCUSDT",
			Price:  "1.001",
		},
	}

	// when
	_, err := GetPairPrice(prices, pairSymbol)

	// then
	require.Error(t, err)
}

func TestGetExchangePairData(t *testing.T) {
	// given
	exchangeID := 1
	symbolJSON := `{
		"symbol": "BTCBUSD",
		"status": "TRADING",
		"baseAsset": "BTC",
		"baseAssetPrecision": 8,
		"quoteAsset": "BUSD",
		"quotePrecision": 8,
		"quoteAssetPrecision": 8,
		"orderTypes": [
			"LIMIT",
			"LIMIT_MAKER",
			"MARKET",
			"STOP_LOSS_LIMIT",
			"TAKE_PROFIT_LIMIT"
		],
		"icebergAllowed": true,
		"ocoAllowed": true,
		"isSpotTradingAllowed": true,
		"isMarginTradingAllowed": true,
		"filters": [
			{
			"filterType": "PRICE_FILTER",
			"maxPrice": "1000000.00000000",
			"minPrice": "0.01000000",
			"tickSize": "0.01000000"
			},
			{
			"filterType": "LOT_SIZE",
			"maxQty": "9000.00000000",
			"minQty": "0.00001000",
			"stepSize": "0.00001000"
			},
			{
			"applyToMarket": true,
			"avgPriceMins": 5,
			"filterType": "NOTIONAL",
			"minNotional": "10.00000000"
			},
			{
			"filterType": "ICEBERG_PARTS",
			"limit": 10
			},
			{
			"filterType": "MARKET_LOT_SIZE",
			"maxQty": "83.17102002",
			"minQty": "0.00000000",
			"stepSize": "0.00000000"
			},
			{
			"filterType": "TRAILING_DELTA",
			"maxTrailingAboveDelta": 2000,
			"maxTrailingBelowDelta": 2000,
			"minTrailingAboveDelta": 10,
			"minTrailingBelowDelta": 10
			},
			{
			"askMultiplierDown": "0.2",
			"askMultiplierUp": "5",
			"avgPriceMins": 5,
			"bidMultiplierDown": "0.2",
			"bidMultiplierUp": "5",
			"filterType": "PERCENT_PRICE_BY_SIDE"
			},
			{
			"filterType": "MAX_NUM_ORDERS",
			"maxNumOrders": 200
			},
			{
			"filterType": "MAX_NUM_ALGO_ORDERS",
			"maxNumAlgoOrders": 5
			}
		],
		"permissions": [
			"SPOT",
			"MARGIN",
			"TRD_GRP_004",
			"TRD_GRP_005",
			"TRD_GRP_006"
		]
	}`

	var symbolData binance.Symbol
	require.NoError(t, json.Unmarshal([]byte(symbolJSON), &symbolData))

	// when
	pairData, err := ConvertExchangePairData(symbolData, exchangeID)

	// then
	require.NoError(t, err)
	assert.Greater(t, pairData.MinDeposit, float64(0))
	assert.Greater(t, pairData.MinPrice, float64(0))
	assert.Greater(t, pairData.MinQty, float64(0))
	assert.NotEmpty(t, pairData.Symbol)
}

func TestConvertExchangePairsDataEmpty(t *testing.T) {
	// given
	var pairsResponse = binance.ExchangeInfo{}
	var exchangeID = 1

	// when
	data, err := ConvertExchangePairsData(pairsResponse, exchangeID)

	// then
	require.NoError(t, err)
	assert.Len(t, data, 0)
}

func TestConvertExchangePairsDataSuccess(t *testing.T) {
	// given
	var pairsResponse = binance.ExchangeInfo{
		Symbols: []binance.Symbol{
			getTestBinancePairData(),
		},
	}
	var exchangeID = 1

	// when
	data, err := ConvertExchangePairsData(pairsResponse, exchangeID)

	// then
	require.NoError(t, err)
	require.Len(t, data, 1)
	assert.Equal(t, "LTCUSDC", data[0].Symbol)
	assert.Equal(t, exchangeID, data[0].ExchangeID)
	assert.Equal(t, float64(0.0001), data[0].MinQty)
	assert.Equal(t, float64(0.01), data[0].MinPrice)
	assert.Equal(t, float64(1), data[0].MinDeposit)
}

func TestConvertExchangePairsDataFiltersEmpty(t *testing.T) {
	// given
	var pairsResponse = binance.ExchangeInfo{
		Symbols: []binance.Symbol{
			getTestBinancePairData(),
		},
	}
	var exchangeID = 1

	pairsResponse.Symbols[0].Filters = make([]map[string]interface{}, 0)

	// when
	_, err := ConvertExchangePairsData(pairsResponse, exchangeID)

	// then
	require.ErrorContains(t, err, "filter not available")
}
