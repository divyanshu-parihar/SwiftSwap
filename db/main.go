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

func AddUserWalletToDB(db *bun.DB, userid int64, network, coin, address, memo string, useGeneratedFalse bool, privateKey string) error {
	fmt.Println("Adding wallet to db : ", userid, network, coin, address, memo)
	wallet := &UserSavedWalletWithPrivateKey{
		ID:                uuid.New(),
		Tid:               string(strconv.FormatInt(userid, 10)),
		Currency:          strings.ToUpper(coin),
		Chain:             strings.ToUpper(network),
		Address:           strings.ToUpper(address),
		Memo:              strings.ToUpper(memo),
		PrivateKey:        strings.ToUpper(privateKey),
		UseGeneratedFalse: false,
	}
	_, err := db.NewInsert().Model(wallet).Exec(context.Background())
	return err
}

func GetWallet(db *bun.DB, userid int64) ([]*UserSavedWalletWithPrivateKey, error) {
	wallets := []*UserSavedWalletWithPrivateKey{}
	err := db.NewSelect().Model(&wallets).Where("tid = ?", string(strconv.FormatInt(userid, 10))).Scan(context.Background())

	fmt.Println("Wallet : ", wallets)
	if len(wallets) < 0 && err != nil {
		return []*UserSavedWalletWithPrivateKey{}, nil
	}
	return wallets, nil
}

func AddUserTrustedWalletToDB(db *bun.DB, userid int64, address string) error {

	wallet := &UserTrustedWallets{
		ID:      uuid.New(),
		User:    userid,
		Address: strings.ToUpper(address),
	}
	_, err := db.NewInsert().Model(wallet).Exec(context.Background())
	return err
}
