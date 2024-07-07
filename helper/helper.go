package pkg

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Coin struct {
	Name  string
	Token string
}
type WazirxTicker struct {
	Symbol     string `json:"symbol"`
	BaseAsset  string `json:"baseAsset"`
	QouteAsset string `json:"quoteAsset"`
	OpenPrice  string `json:"openPrice"`
	LowPrice   string `json:"lowPrice"`
	HighPrice  string `json:"highPrice"`
	LastPrice  string `json:"lastPrice"`
	Volume     string `json:"volume"`
	BidPrice   string `json:"bidPrice"`
	AskPrice   string `json:"askPrice"`
	At         int64  `json:"at"`
}
type Ticker struct {
	Market       string          `json:"market"`
	Change24Hour string          `json:"change_24_hour"`
	High         string          `json:"high"`
	Low          string          `json:"low"`
	Volume       string          `json:"volume"`
	LastPrice    string          `json:"last_price"`
	Bid          json.RawMessage `json:"bid"`
	Ask          json.RawMessage `json:"ask"`
	Timestamp    int64           `json:"timestamp"`
}

type Order struct {
	Id                 string          `json:"id"`
	Client_order_id    string          `json:"client_order_id"`
	Market_order       string          `json:"market_order"`
	Side               string          `json:"side"`
	Status             string          `json:"status"`
	Fee_amount         json.RawMessage `json:"fee_amount"`
	Fee                json.RawMessage `json:"fee"`
	Total_quantity     json.RawMessage `json:"total_quantity"`
	Remaining_quantity json.RawMessage `json:"remaining_quantity"`
	Avg_price          json.RawMessage `json:"avg_price"`
	Price_per_unit     json.RawMessage `json:"price_per_unit"`
	Created_at         int64           `json:"created_at"`
	Update_at          int64           `json:"update_at"`
}

type MexcTicker struct {
	Symbol             string `json:"symbol"`
	PriceChange        string `json:"priceChange"`
	PriceChangePercent string `json:"priceChangePercent"`
	PrevClosePrice     string `json:"prevClosePrice"`
	LastPrice          string `json:"lastPrice"`
	BidPrice           string `json:"bidPrice"`
	BidQty             string `json:"bidQty"`
	AskPrice           string `json:"askPrice"`
	AskQty             string `json:"askQty"`
	OpenPrice          string `json:"openPrice"`
	HighPrice          string `json:"highPrice"`
	LowPrice           string `json:"lowPrice"`
	Volume             string `json:"volume"`
	QuoteVolume        string `json:"quoteVolume"`
	OpenTime           int64  `json:"openTime"`
	CloseTime          int64  `json:"closeTime"`
	Count              string `json:"count"`
}
type OrderResponse struct {
	Orders []Order `json:"orders"`
}

func FindCoinMexc(tickers []MexcTicker, market string) (MexcTicker, error) {
	for _, ticker := range tickers {
		if ticker.Symbol == strings.ToUpper(market) {
			return ticker, nil
		}
	}
	return MexcTicker{}, errors.New("Coin not found")
}

// no need currently

func GetTickerMexc() ([]MexcTicker, error) {
	resp, err := http.Get("https://api.mexc.com/api/v3/ticker/24hr")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tickers []MexcTicker

	if err := json.NewDecoder(resp.Body).Decode(&tickers); err != nil {
		return nil, err
	}

	return tickers, nil
}
func GetTickerWazir() ([]WazirxTicker, error) {
	resp, err := http.Get("https://api.wazirx.com/sapi/v1/tickers/24hr")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var tickers []WazirxTicker
	if err := json.NewDecoder(resp.Body).Decode(&tickers); err != nil {
		return nil, err
	}
	return tickers, nil
}
func GetTicker() ([]Ticker, error) {
	resp, err := http.Get("https://api.coindcx.com/exchange/ticker")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tickers []Ticker
	if err := json.NewDecoder(resp.Body).Decode(&tickers); err != nil {
		return nil, err
	}

	return tickers, nil
}
func FindCoinWazirx(tickers []WazirxTicker, market string) (*WazirxTicker, error) {
	for _, ticker := range tickers {
		if ticker.Symbol == market {
			return &ticker, nil
		}
	}
	return nil, errors.New("coin not found")
}
func FindCoin(tickers []Ticker, market string) (*Ticker, error) {
	for _, ticker := range tickers {
		if ticker.Market == market {
			return &ticker, nil
		}
	}
	return nil, errors.New("coin not found")
}
func CreateFinalPayload(body map[string]any) ([]byte, string, error) {

	var secret string = os.Getenv("API_SECRET")
	timeStamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	body["timestamp"] = timeStamp
	payload, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return []byte{}, "", errors.New("cannot convert to json")
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	signature := hex.EncodeToString(mac.Sum(nil))
	return payload, signature, nil
}
