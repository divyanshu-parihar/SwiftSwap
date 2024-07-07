package pkg

import (
	"bytes"
	helper "crypto-exchange-swap/helper"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type mexcHandler struct {
	api_key    string
	api_secret string
	base_url   string
}

func NewMexcHandler() mexcHandler {
	return mexcHandler{
		api_key:    os.Getenv("MEXC_API_KEY"),
		api_secret: os.Getenv("MEXC_API_SECRET"),
		base_url:   "https://api.mexc.com",
	}
}

// Function to generate HMAC SHA256 signature
type MexcOrderResponse struct {
	Symbol       string `json:"symbol"`
	OrderId      string `json:"orderId"`
	OrderListId  int    `json:"orderListId"`
	Price        string `json:"price"`
	OrigQty      string `json:"origQty"`
	Type         string `json:"type"`
	Side         string `json:"side"`
	TransactTime int    `json:"transactTime"`
}

func (h *mexcHandler) SellForUSDC(coin helper.Coin, amount float64) (helper.OrderResponse, error) {
	tickers, err := helper.GetTickerMexc()
	if err != nil {
		fmt.Println("Error fetching ticker data:", err)
		return helper.OrderResponse{}, err
	}

	ticker, err := helper.FindCoinMexc(tickers, coin.Token+"usdt")
	if err != nil {
		fmt.Println("Error finding coin:", err)
		return helper.OrderResponse{}, err
	}

	price, err := strconv.ParseFloat(ticker.LastPrice, 64)
	fmt.Println("Price:", price)
	if err != nil {
		fmt.Println("Error parsing price:", err)
		return helper.OrderResponse{}, err
	}

	quantity := float64(int((amount/price)*100000)) / 100000
	fmt.Println("Quantity to sell:", quantity)

	body := map[string]string{
		"symbol":     strings.ToUpper(coin.Token) + "USDC",
		"side":       "SELL",
		"type":       "LIMIT",
		"quantity":   fmt.Sprintf("%f", quantity),
		"price":      ticker.LastPrice,
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

	req, err := http.NewRequest("POST", h.base_url+"/api/v3/order", bytes.NewBufferString(payload.Encode()))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return helper.OrderResponse{}, err
	}

	req.Header.Set("X-MEXC-APIKEY", h.api_key)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println("Error creating request:", err)
		return helper.OrderResponse{}, err
	}

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

func (h *mexcHandler) BuyWithUSDC(coin helper.Coin, usdtAmount float64) (helper.OrderResponse, error) {
	tickers, err := helper.GetTickerMexc()
	if err != nil {
		fmt.Println("Error fetching ticker data:", err)
		return helper.OrderResponse{}, err
	}
	ticker, err := helper.FindCoinMexc(tickers, coin.Token+"USDC")

	if err != nil {
		fmt.Println("Error finding coin:", err)
		return helper.OrderResponse{}, err
	}
	fmt.Print("Ticker: ", ticker)
	price, err := strconv.ParseFloat(ticker.LastPrice, 64)
	if err != nil {
		fmt.Println("Error parsing price:", err)
		return helper.OrderResponse{}, err
	}

	quantity := float64(int((usdtAmount/price)*100000)) / 100000
	fmt.Println("Quantity to buy:", quantity)
	fmt.Println("Price:", price)
	body := map[string]string{
		"symbol":        strings.ToUpper(coin.Token) + "USDC",
		"side":          "BUY",
		"type":          "LIMIT",
		"price":         fmt.Sprintf("%f", price),
		"quoteOrderQty": fmt.Sprintf("%f", quantity),
		"quantity":      fmt.Sprintf("%f", quantity),
		"recvWindow":    "10000",
		"timestamp":     strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10),
	}
	fmt.Println("Body:", body)
	payload := url.Values{}
	for key, value := range body {
		payload.Set(key, value)
	}

	signature := generateSignature(payload.Encode(), h.api_secret)
	payload.Set("signature", signature)

	req, err := http.NewRequest("POST", h.base_url+"/api/v3/order", bytes.NewBufferString(payload.Encode()))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return helper.OrderResponse{}, err
	}

	req.Header.Set("X-MEXC-APIKEY", h.api_key)
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
