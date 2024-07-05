package pkg

import (
	"bytes"
	helper "crypto-exchange-swap/helper"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

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

func (h *coinDCXHandler) SellforUSDT(coin helper.Coin, amount float64) (helper.OrderResponse, error) {
	tickers, err := helper.GetTicker()
	if err != nil {
		fmt.Println("Error fetching ticker data:", err)
		return helper.OrderResponse{}, err
	}

	// Find the coin in the ticker data
	ticker, err := helper.FindCoin(tickers, coin.Token+"USDT")
	if err != nil {
		fmt.Println("Error finding coin:", err)
		return helper.OrderResponse{}, err
	}

	price, err := strconv.ParseFloat(ticker.LastPrice, 64)
	if err != nil {
		fmt.Println("Error parsing price:", err)
		return helper.OrderResponse{}, err
	}
	// Calculate quantity to buy based on USDT amount
	quantity := float64(int((amount/price)*100000)) / 100000
	fmt.Println("Quantity to buy:", quantity)

	body := map[string]interface{}{
		"side":           "sell",              // Toggle between 'buy' or 'sell'.
		"order_type":     "market_order",      // Toggle between a 'market_order' or 'limit_order'.
		"market":         coin.Token + "USDT", // Replace 'SNTBTC' with your desired market pair.
		"total_quantity": quantity,            // Replace this with the quantity you want
	}

	payload, signature, err := helper.CreateFinalPayload(body)
	fmt.Printf("Payload: %s, signature %s", string(payload), signature)
	if err != nil {
		log.Fatal("Could not create hash")
	}
	req, err := http.NewRequest("POST", h.base_url+"/exchange/v1/orders/create", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return helper.OrderResponse{}, err
	}

	req.Header.Set("X-AUTH-APIKEY", h.api_key)
	req.Header.Set("X-AUTH-SIGNATURE", signature)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return helper.OrderResponse{}, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return helper.OrderResponse{}, err
	}

	fmt.Println("Response:", string(bodyBytes))
	var result helper.OrderResponse
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return helper.OrderResponse{}, err
	}
	return result, nil

}
func (h *coinDCXHandler) BuyWithUSDT(coin helper.Coin, usdtAmount float64) (helper.OrderResponse, error) {
	// Fetch ticker data
	tickers, err := helper.GetTicker()
	if err != nil {
		fmt.Println("Error fetching ticker data:", err)
		return helper.OrderResponse{}, err
	}

	// Find the coin in the ticker data
	ticker, err := helper.FindCoin(tickers, coin.Token+"USDT")
	if err != nil {
		fmt.Println("Error finding coin:", err)
		return helper.OrderResponse{}, err
	}

	price, err := strconv.ParseFloat(ticker.LastPrice, 64)
	if err != nil {
		fmt.Println("Error parsing price:", err)
		return helper.OrderResponse{}, err
	}

	// Calculate quantity to buy based on USDT amount
	quantity := float64(int((usdtAmount/price)*100000)) / 100000
	fmt.Println("Quantity to buy:", quantity)

	body := map[string]interface{}{
		"side":           "buy",               // Buy order
		"order_type":     "market_order",      // Market order
		"market":         coin.Token + "USDT", // Market pair
		"total_quantity": quantity,            // Calculated quantity maximum 5 decimal points
	}

	payload, signature, err := helper.CreateFinalPayload(body)
	if err != nil {
		fmt.Println("Error creating payload:", err)
		return helper.OrderResponse{}, err
	}

	req, err := http.NewRequest("POST", h.base_url+"/exchange/v1/orders/create", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return helper.OrderResponse{}, err
	}

	req.Header.Set("X-AUTH-APIKEY", h.api_key)
	req.Header.Set("X-AUTH-SIGNATURE", signature)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return helper.OrderResponse{}, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return helper.OrderResponse{}, err
	}

	bodyString := string(bodyBytes) // Convert byte[] to string
	fmt.Println("Response:", bodyString)
	//convert bodyString to OrderResponse using encoding/json
	var result helper.OrderResponse
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return helper.OrderResponse{}, err
	}
	return result, nil
}
func exchange_token(initial, final helper.Coin) {
	// sell coin to usdt

	// buy fina coin from the same amount of usdt
}
