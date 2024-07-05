package main

import (
	w "crypto-exchange-swap/handler"
	h "crypto-exchange-swap/helper"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

// entry point for the application
func main() {
	fmt.Println("Hello, Starting the bot")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// _ := pkg.NewCoinHandler()

	coin := h.Coin{
		Name:  "Bitcoin",
		Token: "btc",
	}
	// handler.SellforUSDT(coin, 9.979)
	wazir := w.NewWazirXHandler()

	wazir.BuyWithUSDT(coin, 10)
}
