package binance

import (
	"context"
	"errors"
	"fmt"

	"github.com/adshao/go-binance/v2"
)

func (a *adapter) VerifyAPIKeys(keyPublic, keySecret string) error {
	newClient := binance.NewClient(keyPublic, keySecret)
	accountService, err := newClient.NewGetAccountService().Do(context.Background())
	if err != nil {
		return fmt.Errorf("invalid api key: %w", err)
	}
	if !accountService.CanTrade {
		return errors.New("your API key does not have permission to trade," +
			" change its restrictions")
	}
	return nil
}
