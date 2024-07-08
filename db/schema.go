package db

import (
	"github.com/google/uuid"
)

type Transaction struct {
	ID              int     `bun:"id,pk,autoincrement,type:uuid,"`
	PrimaryCurrency string  `bun:"primary,notnull,"`
	FinalCurrency   string  `bun:"final,notnull,"`
	MiddleCurrency  string  `bun:"middle"`
	User            string  `bun:"user,notnull,"`
	CreatedAt       string  `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt       string  `bun:"updated_at,notnull,default:current_timestamp"`
	Filled          bool    `bun:"filled,notnull,default:false"`
	PartiallyFilled bool    `bun:"partially_filled,notnull,default:false"`
	InitialQuantity float64 `bun:"initial_quantity,notnull,"`
	FinalQuantity   float64 `bun:"final_quantity,notnull,"`
	platformFees    int     `bun:"platform_fees,notnull,default:0"`
	ExchangeUsed    string  `bun:"exchange_used,notnull,"`
}

type Deposits struct {
	ID            int    `bun:"id,pk,autoincrement,type:uuid,"`
	User          string `bun:"user,notnull,"`
	Currency      string `bun:"currency,notnull,"`
	Chain         string `bun:"chain,notnull,"`
	TransactionID int    `bun:"transaction_id,notnull,"`

	confirmed bool    `bun:"confirms,notnull,default:false"`
	Amount    float64 `bun:"amount,notnull,"`
	Exchange  string  `bun:"exchange,notnull,"`
	Addr      string  `bun:"addr,notnull,"`
	CreatedAt string  `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt string  `bun:"updated_at,notnull,default:current_timestamp"`
}

type UserSavedWalletWithPrivateKey struct {
	ID                uuid.UUID `bun:"type:uuid,default:uuid_generate_v4()" json:"id"`
	User              int64     `bun:"user,notnull,"`
	Currency          string    `bun:"currency,notnull,"`
	Chain             string    `bun:"chain,notnull,"`
	Memo              string    `bun:"memo,"`
	Address           string    `bun:"address,notnull,"`
	UseGeneratedFalse bool      `bun:"use_generated,notnull,default:false"`
	PrivateKey        string    `bun:"private_key,notnull"`
}

type UserTrustedWallets struct {
	ID      uuid.UUID `bun:"type:uuid,default:uuid_generate_v4()" json:"id"`
	User    int64     `bun:"user,notnull,"`
	Address string    `bun:"wallets,notnull,"`
}
