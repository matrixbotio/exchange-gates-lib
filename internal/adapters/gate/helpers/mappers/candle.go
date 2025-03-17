package mappers

import (
	"fmt"

	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
)

func ConvertCandles(rawData [][]string) ([]workers.CandleData, error) {
	// temp
	for _, v := range rawData {
		fmt.Println(v)
	}
	// TODO
	return nil, nil
}
