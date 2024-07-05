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
	"time"
)

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
	fmt.Println(secret)
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
