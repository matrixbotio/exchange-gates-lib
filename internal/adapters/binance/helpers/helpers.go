package helpers

import (
	"fmt"
)

func GetTradeEventCacheKey(buyerOrderID, sellerOrderID int64) string {
	return fmt.Sprintf(
		"%v-%v",
		buyerOrderID,
		sellerOrderID,
	)
}
