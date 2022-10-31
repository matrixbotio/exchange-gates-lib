package structs

import "github.com/matrixbotio/exchange-gates-lib/internal/structs"

// APICredentialsTypeKeypair - public and private key pair
var APICredentialsTypeKeypair APICredentialsType = "keypair"

// APICredentialsType - API credentials type ^ↀᴥↀ^
type APICredentialsType string

// APICredentials - data for authorization to the exchange API
type APICredentials struct {
	Type APICredentialsType `json:"type"`

	Keypair  APIKeypair          `json:"keypair"`
	Password structs.APIPassword `json:"password"`
	Email    structs.APIEmail    `json:"email"`
}

// WorkerChannels - channels container to control the worker
type WorkerChannels struct {
	WsDone chan struct{}
	WsStop chan struct{}
}

// CheckOrdersResponse - data on checked and restored orders
type CheckOrdersResponse struct {
	ExecutedOrders  []structs.OrderData
	CancelledOrders []structs.OrderData
}

/*
BotOrder - structure containing information about the order placed by the bot.
Used when auto-resuming trades
*/
type BotOrder struct {
	// required
	PairSymbol string  `json:"pair"`
	Type       string  `json:"type"`
	Qty        float64 `json:"qty"`
	Price      float64 `json:"price"`
	Deposit    float64 `json:"deposit"`

	// optional
	ClientOrderID string `json:"clientOrderID"`
}

// APIKeypair - data for authorization via public and private keys
type APIKeypair struct {
	Public string `json:"public"`
	Secret string `json:"secret"`
}

// BotOrderAdjusted - the same as BotOrder, only with the given values for the trading pair
type BotOrderAdjusted struct {
	// required
	PairSymbol string `json:"pair"`
	Type       string `json:"type"`
	Qty        string `json:"qty"`
	Price      string `json:"price"`
	Deposit    string `json:"deposit"`

	// optional
	ClientOrderID string `json:"clientOrderID"`

	// calculated
	MinQty           float64 `json:"minQty"`
	MinQtyPassed     bool    `json:"minQtyPassed"`
	MinDeposit       float64 `json:"minDeposit"`
	MinDepositPassed bool    `json:"minDepositPassed"`
}
