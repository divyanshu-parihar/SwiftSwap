package main

import (
	b "crypto-exchange-swap/bot"
	d "crypto-exchange-swap/db"
	ex "crypto-exchange-swap/handler"

	// h "crypto-exchange-swap/helper"
	"crypto-exchange-swap/walletGen"
	"crypto-exchange-swap/wallets"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

// entry point for the application
func main() {
	// wg.CreateWallet()
	err := godotenv.Load()
	args := os.Args
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db := d.StartServer()
	if len(args) > 1 && args[1] == "migrate" {

		fmt.Println("Running migrations...")
		err := d.RunMigrations(db)
		if err != nil {
			log.Fatalf("Error running migrations: %v", err)
		}
		fmt.Println("migrateions ran successfully")
		return
	}

	var wg sync.WaitGroup
	language := wallets.NewDepositWallets()
	go b.NewBot(&wg, db, language)
	wg.Add(1)

	// fmt.Println("Hello, Starting the bot")
	// coin := h.Coin{
	// 	Name:  "Bitcoin",
	// 	Token: "LTC",
	// }
	mexc := ex.NewMexcHandler()
	// mexHandler := .NewMexcHandler(db, bot, wallets)
	// coin := h.Coin{
	// 	Name:  "TRON",
	// 	Token: "TRX",
	// }
	fmt.Println("Now starting to withdraw")
	// mexc.WithDrawCryptoToWallet(coin, "TRX", 20, "TBcGvboH3ptu7GcRcJhYytUkwes6xZgkV7", "")
	mexc.CheckTheDespositTransaction("0xc32e70752a9d18a941116ea86d79015223d757973de8e93a5ea7884814692715")
	// mexc.BuyWithUSDC(coin, 10)
	//
	baseClient, err := ethclient.Dial("https://mainnet.base.org")
	if err != nil {
		log.Fatalf("Failed to connect to Coinbase Base Chain: %v", err)
	}
	fmt.Println("Connected to Coinbase Base Chain")
	walletGen.CheckBalance(baseClient, common.HexToAddress("0x74293CD5dd47F03A46b686eCCf5bfaE2b6bC2A0F"), "coinbase wallet")
	wg.Wait()
}
