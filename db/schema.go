package db

import (
	"github.com/google/uuid"
)

type Transaction struct {
	ID              uuid.UUID         `bun:"id,pk,autoincrement,type:uuid,"`
	TxnID           string            `bun:"txn_id,notnull,"`
	PrimaryCurrency string            `bun:"primary,notnull,"`
	FinalCurrency   string            `bun:"final,notnull,"`
	Memo            string            `bun:"memo",notnull`
	Network         string            `bun:"network",notnull`
	MiddleCurrency  string            `bun:"middle"`
	User            string            `bun:"user,notnull,"`
	CreatedAt       string            `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt       string            `bun:"updated_at,notnull,default:current_timestamp"`
	Filled          bool              `bun:"filled,notnull,default:false"`
	PartiallyFilled bool              `bun:"partially_filled,notnull,default:false"`
	InitialQuantity float64           `bun:"initial_quantity,notnull,"`
	FinalQuantity   float64           `bun:"final_quantity,notnull,"`
	platformFees    int               `bun:"platform_fees,notnull,default:0"`
	ExchangeUsed    string            `bun:"exchange_used,notnull,"`
	Tries           int64             `bun:",notnull`
	Status          TransactionStatus `bun:",notnull`
	Address         string            `bun:"address,notnull,"`
}
type TransactionStatus string

const (
	StatusPending    TransactionStatus = "pending"
	StatusCompleted  TransactionStatus = "completed"
	StatusIncomplete TransactionStatus = "incomplete"
	StatusFailed     TransactionStatus = "failed"
	StatusDeposit    TransactionStatus = "deposit"
	StatusWithdrawal TransactionStatus = "pw"
	StatusTransfer   TransactionStatus = "wd"
)

// type Deposits struct {
// 	ID            int    `bun:"id,pk,autoincrement,type:uuid,"`
// 	User          string `bun:"user,notnull,"`
// 	Currency      string `bun:"currency,notnull,"`
// 	Chain         string `bun:"chain,notnull,"`
// 	TransactionID int    `bun:"transaction_id,notnull,"`

// 	confirmed bool    `bun:"confirms,notnull,default:false"`
// 	Amount    float64 `bun:"amount,notnull,"`
// 	Exchange  string  `bun:"exchange,notnull,"`
// 	Addr      string  `bun:"addr,notnull,"`
// 	CreatedAt string  `bun:"created_at,notnull,default:current_timestamp"`
// 	UpdatedAt string  `bun:"updated_at,notnull,default:current_timestamp"`
// }

type UserSavedWalletWithPrivateKey struct {
	ID                uuid.UUID `bun:"type:uuid,default:uuid_generate_v4()" json:"id"`
	Tid               string    `bun:"tid,notnull,"`
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
