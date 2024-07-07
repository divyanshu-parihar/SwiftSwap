package main

import (
	b "crypto-exchange-swap/bot"
	d "crypto-exchange-swap/db"
	wallets "crypto-exchange-swap/wallets"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

// entry point for the application
func main() {
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

	fmt.Println("Hello, Starting the bot")
	// coin := h.Coin{
	// 	Name:  "Bitcoin",
	// 	Token: "LTC",
	// }
	// mexc := w.NewMexcHandler()
	// mexc.BuyWithUSDC(coin, 10)
	wg.Wait()
}
