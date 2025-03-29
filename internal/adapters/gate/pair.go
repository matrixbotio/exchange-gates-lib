package gate

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/antihax/optional"
	"github.com/gateio/gateapi-go/v6"
	"github.com/matrixbotio/exchange-gates-lib/internal/adapters/gate/helpers/mappers"
	"github.com/matrixbotio/exchange-gates-lib/internal/structs"
	"github.com/matrixbotio/exchange-gates-lib/pkg/errs"
	"github.com/matrixbotio/exchange-gates-lib/pkg/utils"
	"github.com/shopspring/decimal"
)

func (a *adapter) GetPairData(pairSymbol string) (structs.ExchangePairData, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), requestTimeout)
	defer ctxCancel()

	data, _, err := a.client.SpotApi.GetCurrencyPair(ctx, pairSymbol)
	if err != nil {
		return structs.ExchangePairData{},
			fmt.Errorf("get: %w", err)
	}

	result, err := mappers.ConvertPair(data)
	if err != nil {
		return structs.ExchangePairData{},
			fmt.Errorf("convert: %w", err)
	}
	return result, nil
}

func (a *adapter) GetPairLastPrice(pairSymbol string) (float64, error) {
	tickers, _, err := a.client.SpotApi.ListTickers(
		context.Background(),
		&gateapi.ListTickersOpts{
			CurrencyPair: optional.NewString(pairSymbol),
		},
	)
	if err != nil {
		return 0, fmt.Errorf("get ticker: %w", err)
	}

	if len(tickers) == 0 {
		return 0, fmt.Errorf("ticker %q price not found", pairSymbol)
	}

	lastPrice, err := decimal.NewFromString(tickers[0].Last)
	if err != nil {
		return 0, fmt.Errorf("parse price: %w", err)
	}

	return lastPrice.InexactFloat64(), nil
}

func (a *adapter) CancelPairOrder(
	pairSymbol string,
	orderID int64,
	_ context.Context,
) error {
	if !a.creds.Keypair.IsSet() {
		return errs.ErrAPIKeyNotSet
	}

	ctx, ctxCancel := context.WithTimeout(a.auth, requestTimeout)
	defer ctxCancel()

	_, _, err := a.client.SpotApi.CancelOrder(
		ctx,
		strconv.FormatInt(orderID, 10),
		pairSymbol,
		nil,
	)
	return mappers.MapCancelOrderErr(err)
}

func (a *adapter) CancelPairOrderByClientOrderID(
	pairSymbol string,
	clientOrderID string,
	_ context.Context,
) error {
	if !a.creds.Keypair.IsSet() {
		return errs.ErrAPIKeyNotSet
	}

	ctx, ctxCancel := context.WithTimeout(a.auth, requestTimeout)
	defer ctxCancel()

	_, _, err := a.client.SpotApi.CancelOrder(
		ctx,
		clientOrderID,
		pairSymbol,
		nil,
	)
	return mappers.MapCancelOrderErr(err)
}

func (a *adapter) GetPairs() ([]structs.ExchangePairData, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), requestTimeout)
	defer ctxCancel()

	pairs, _, err := a.client.SpotApi.ListCurrencyPairs(ctx)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	result, err := mappers.ConvertPairs(pairs)
	if err != nil {
		return nil, fmt.Errorf("convert: %w", err)
	}
	return result, nil
}

func (a *adapter) GetPairBalance(pair structs.PairSymbolData) (structs.PairBalance, error) {
	if !a.creds.Keypair.IsSet() {
		return structs.PairBalance{}, errs.ErrAPIKeyNotSet
	}

	balances, err := a.GetAccountBalance()
	if err != nil {
		if strings.Contains(err.Error(), "Invalid key provided") {
			return structs.PairBalance{}, errs.ErrAPIKeyInvalid
		}

		return structs.PairBalance{}, fmt.Errorf("get: %w", err)
	}

	return utils.FindPairBalance(balances, pair), nil
}
