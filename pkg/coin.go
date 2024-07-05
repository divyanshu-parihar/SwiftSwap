package pkg

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Coin struct {
	Name  string
	Token string
}

type coinDCXHandler struct {
	api_key    string
	api_secret string
	base_url   string // Ensure consistency in naming
}

func NewCoinHandler() coinDCXHandler {
	return coinDCXHandler{
		api_key:    os.Getenv("API_KEY"),
		api_secret: os.Getenv("API_SECRET"),
		base_url:   "https://api.coindcx.com", // Match the struct field name
	}
}

func (h *coinDCXHandler) SellforUSDT(coin Coin, amount float32) {
	body := map[string]interface{}{
		"side":           "sell",         // Toggle between 'buy' or 'sell'.
		"order_type":     "market_order", // Toggle between a 'market_order' or 'limit_order'.
		"market":         coin.Token,     // Replace 'SNTBTC' with your desired market pair.
		"total_quantity": amount,         // Replace this with the quantity you want
	}

	payload, signature, err := CreateFinalPayload(body)
	if err != nil {
		log.Fatal("Could not create hash")
	}
	req, err := http.NewRequest("POST", h.base_url+"/exchange/v1/orders/create", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("X-AUTH-APIKEY", string(payload))
	req.Header.Set("X-AUTH-SIGNATURE", signature)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println("Response:", string(bodyBytes))

}
func (h *coinDCXHandler) BuyWithUSDT(coin Coin, usdtAmount float64) {
	// Fetch ticker data
	tickers, err := GetTicker()
	if err != nil {
		fmt.Println("Error fetching ticker data:", err)
		return
	}

	// Find the coin in the ticker data
	ticker, err := FindCoin(tickers, coin.Token+"USDT")
	if err != nil {
		fmt.Println("Error finding coin:", err)
		return
	}

	// Convert the last price to a float64
	price, err := strconv.ParseFloat(ticker.LastPrice, 64)
	if err != nil {
		fmt.Println("Error parsing price:", err)
		return
	}

	// Calculate quantity to buy based on USDT amount
	quantity := usdtAmount / price

	body := map[string]interface{}{
		"side":           "buy",               // Buy order
		"order_type":     "market_order",      // Market order
		"market":         coin.Token + "USDT", // Market pair
		"total_quantity": quantity,            // Calculated quantity
	}

	payload, signature, err := CreateFinalPayload(body)
	if err != nil {
		fmt.Println("Error creating payload:", err)
		return
	}

	req, err := http.NewRequest("POST", h.base_url+"/exchange/v1/orders/create", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("X-AUTH-APIKEY", h.api_key)
	req.Header.Set("X-AUTH-SIGNATURE", signature)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	bodyString := string(bodyBytes) // Convert byte[] to string
	fmt.Println("Response:", bodyString)
}
func exchange_token(initial, final Coin) {
	// sell coin to usdt

	// buy fina coin from the same amount of usdt
}
