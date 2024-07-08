package pkg

import (
	"bytes"
	helper "crypto-exchange-swap/helper"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type wazirXHandler struct {
	api_key    string
	api_secret string
	base_url   string
}

func NewWazirXHandler() wazirXHandler {
	return wazirXHandler{
		api_key:    os.Getenv("WAZIR_API_KEY"),
		api_secret: os.Getenv("WAZIR_API_SECRET"),
		base_url:   "https://api.wazirx.com",
	}
}

// Function to generate HMAC SHA256 signature
//
//
//

func generateSignature(data, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func (h *wazirXHandler) SellForUSDT(coin helper.Coin, amount float64) (helper.OrderResponse, error) {
	tickers, err := helper.GetTicker()
	if err != nil {
		fmt.Println("Error fetching ticker data:", err)
		return helper.OrderResponse{}, err
	}

	ticker, err := helper.FindCoin(tickers, coin.Token+"usdt")
	if err != nil {
		fmt.Println("Error finding coin:", err)
		return helper.OrderResponse{}, err
	}

	price, err := strconv.ParseFloat(ticker.LastPrice, 64)
	if err != nil {
		fmt.Println("Error parsing price:", err)
		return helper.OrderResponse{}, err
	}

	quantity := float64(int((amount/price)*100000)) / 100000
	fmt.Println("Quantity to sell:", quantity)

	body := map[string]string{
		"symbol":     coin.Token + "usdt",
		"side":       "sell",
		"type":       "market",
		"quantity":   fmt.Sprintf("%f", quantity),
		"recvWindow": "10000",
		"timestamp":  strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10),
	}

	payload := url.Values{}
	for key, value := range body {
		payload.Set(key, value)
	}

	signature := generateSignature(payload.Encode(), h.api_secret)
	payload.Set("signature", signature)

	req, err := http.NewRequest("POST", h.base_url+"/sapi/v1/order", bytes.NewBufferString(payload.Encode()))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return helper.OrderResponse{}, err
	}

	req.Header.Set("X-Api-Key", h.api_key)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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

func (h *wazirXHandler) BuyWithUSDT(coin helper.Coin, usdtAmount float64) (helper.OrderResponse, error) {
	tickers, err := helper.GetTickerWazir()
	if err != nil {
		fmt.Println("Error fetching ticker data:", err)
		return helper.OrderResponse{}, err
	}
	ticker, err := helper.FindCoinWazirx(tickers, coin.Token+"usdt")
	if err != nil {
		fmt.Println("Error finding coin:", err)
		return helper.OrderResponse{}, err
	}
	price, err := strconv.ParseFloat(ticker.LastPrice, 64)
	if err != nil {
		fmt.Println("Error parsing price:", err)
		return helper.OrderResponse{}, err
	}

	quantity := float64(int((usdtAmount/price)*100000)) / 100000
	fmt.Println("Quantity to buy:", quantity)
	fmt.Println("Price:", price)
	body := map[string]string{
		"symbol":     coin.Token + "usdt",
		"side":       "buy",
		"type":       "limit",
		"price":      fmt.Sprintf("%f", price),
		"quantity":   fmt.Sprintf("%f", quantity),
		"recvWindow": "10000",
		"timestamp":  strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10),
	}
	fmt.Println("Body:", body)
	payload := url.Values{}
	for key, value := range body {
		payload.Set(key, value)
	}

	signature := generateSignature(payload.Encode(), h.api_secret)
	payload.Set("signature", signature)

	req, err := http.NewRequest("POST", h.base_url+"/sapi/v1/order", bytes.NewBufferString(payload.Encode()))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return helper.OrderResponse{}, err
	}

	req.Header.Set("X-Api-Key", h.api_key)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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

func exchangeToken(initial, final helper.Coin, amount float64) {
	handler := NewWazirXHandler()

	// Sell initial coin to get USDT
	sellResponse, err := handler.SellForUSDT(initial, amount)
	if err != nil {
		log.Fatalf("Error selling %s: %v", initial.Name, err)
	}
	fmt.Printf("Sold %s for USDT: %+v\n", initial.Name, sellResponse)

	// Use the received USDT amount to buy the final coin
	buyResponse, err := handler.BuyWithUSDT(final, amount)
	if err != nil {
		log.Fatalf("Error buying %s: %v", final.Name, err)
	}
	fmt.Printf("Bought %s with USDT: %+v\n", final.Name, buyResponse)
}
