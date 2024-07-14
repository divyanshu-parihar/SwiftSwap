package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func StartServer() *bun.DB {

	dsn := os.Getenv("DATABASE_URL")
	// dsn := "unix://user:pass@dbname/var/run/postgresql/.s.PGSQL.5432"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	return bun.NewDB(sqldb, pgdialect.New())
}

func CreateTransation(db *bun.DB, txnid string, primary, secondary string, userid string, initialQuantiy int, memo string) error {
	fmt.Println("Adding Transaction to db : ")
	transaction := &Transaction{
		ID:              uuid.New(),
		TxnID:           txnid,
		Memo:            memo,
		PrimaryCurrency: strings.ToUpper(primary),
		FinalCurrency:   strings.ToLower(secondary),
		MiddleCurrency:  strings.ToUpper(primary),
		User:            userid,
		Filled:          false,
		PartiallyFilled: false,
		InitialQuantity: float64(initialQuantiy),
		FinalQuantity:   0,
		ExchangeUsed:    "MEXC",
		platformFees:    0,
		Tries:           0,
	}
	_, err := db.NewInsert().Model(transaction).Exec(context.Background())
	// return err
	// _, err := db.NewCreateTable().Model((*Transaction)(nil)).Exec(context.Background())
	return err
}

// func GetTransaction(db *bun.DB) ([]*Transaction, error) {
// 	transactions := []*Transaction{}
// 	err := db.NewSelect().Model(&transactions).Where("tid = ?", string(strconv.FormatInt(userid, 10))).Scan(context.Background())

//		fmt.Println("Wallet : ", wallets)
//		if len(wallets) < 0 && err != nil {
//			return []*UserSavedWalletWithPrivateKey{}, nil
//		}
//		return wallets, nil
//	}
func GetWallet(db *bun.DB, userid int64) ([]*UserSavedWalletWithPrivateKey, error) {
	wallets := []*UserSavedWalletWithPrivateKey{}
	err := db.NewSelect().Model(&wallets).Where("tid = ?", string(strconv.FormatInt(userid, 10))).Scan(context.Background())

	fmt.Println("Wallet : ", wallets)
	if len(wallets) < 0 && err != nil {
		return []*UserSavedWalletWithPrivateKey{}, nil
	}
	return wallets, nil
}
func AddUserWalletToDB(db *bun.DB, userid int64, network, coin, address, memo string, useGeneratedFalse bool, privateKey string) error {
	fmt.Println("Adding wallet to db : ", userid, network, coin, address, memo)
	wallet := &UserSavedWalletWithPrivateKey{
		ID:                uuid.New(),
		Tid:               string(strconv.FormatInt(userid, 10)),
		Currency:          strings.ToUpper(coin),
		Chain:             strings.ToUpper(network),
		Address:           address,
		Memo:              memo,
		PrivateKey:        privateKey,
		UseGeneratedFalse: false,
	}
	_, err := db.NewInsert().Model(wallet).Exec(context.Background())
	return err
}

// func GetWallet(db *bun.DB, userid int64) ([]*UserSavedWalletWithPrivateKey, error) {
// 	wallets := []*UserSavedWalletWithPrivateKey{}
// 	err := db.NewSelect().Model(&wallets).Where("tid = ?", string(strconv.FormatInt(userid, 10))).Scan(context.Background())

// 	fmt.Println("Wallet : ", wallets)
// 	if len(wallets) < 0 && err != nil {
// 		return []*UserSavedWalletWithPrivateKey{}, nil
// 	}
// 	return wallets, nil
// }

func AddUserTrustedWalletToDB(db *bun.DB, userid int64, address string) error {

	wallet := &UserTrustedWallets{
		ID:      uuid.New(),
		User:    userid,
		Address: strings.ToUpper(address),
	}
	_, err := db.NewInsert().Model(wallet).Exec(context.Background())
	return err
}
