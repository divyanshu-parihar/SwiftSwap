package main

import (
	"crypto-exchage-swap/pkg"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hello, Starting the bot")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	handler := pkg.NewCoinHandler()

	coin := pkg.Coin{
		Name:  "coin",
		Token: "BTC",
	}
	handler.BuyWithUSDT(coin, 10.0)
}
