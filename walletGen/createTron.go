package walletGen

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/mr-tron/base58"
	w "github.com/ranjbar-dev/tron-wallet"
	"github.com/ranjbar-dev/tron-wallet/enums"
)

const (
	apiKey  = "8e87c19e-99be-4d18-8a87-2baf277eec59"
	baseURL = "https://api.trongrid.io"
)

// Account represents the structure of the account response
type Account struct {
	Balance int64 `json:"balance"`
}

// GetAccountBalance fetches the balance of a TRON account
func GetAccountBalance(address string) (int64, error) {
	url := fmt.Sprintf("%s/v1/accounts/%s", baseURL, address)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("TRON-PRO-API-KEY", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to fetch account info: %s", resp.Status)
	}

	var result struct {
		Data []Account `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	if len(result.Data) == 0 {
		return 0, fmt.Errorf("account not found")
	}

	return result.Data[0].Balance, nil
}
func GenerateTronwallet() *w.TronWallet {
	tronWallet := w.GenerateTronWallet("8e87c19e-99be-4d18-8a87-2baf277eec59")

	return tronWallet
}

type TronTransferResponse struct {
	Response string
	TxId     string
}

func CheckTronBalance(wallet *w.TronWallet) (int64, error) {
	currentBalance, err := wallet.Balance()
	if err != nil {
		return 0, err
	}
	return currentBalance, nil

}

func GenerateWalletWithHex(privateKeyHex string) *w.TronWallet {
	wallet, _ := w.CreateTronWallet(enums.MAIN_NODE, privateKeyHex)

	return wallet
}
func TransferTron(from, toAddr string, amount int64, privateKeyHex string) (*TronTransferResponse, error) {
	wallet, _ := w.CreateTronWallet(enums.MAIN_NODE, privateKeyHex)
	currentBalance, err := CheckTronBalance(wallet)
	if err != nil {
		return &TronTransferResponse{
			Response: "Failed to check balance",
			TxId:     "",
		}, err
	}

	if currentBalance < amount+1 {
		fmt.Println("Insufficient balance")

		return &TronTransferResponse{
			Response: "Insuffcient balance",
			TxId:     "",
		}, errors.New("Insufficient balance")
	}
	txId, err := wallet.Transfer(toAddr, amount)
	if err != nil {
		fmt.Println(err)
		return &TronTransferResponse{
			Response: "Failed to transfer",
			TxId:     "",
		}, err
	}
	return &TronTransferResponse{
		Response: "Success",
		TxId:     txId,
	}, nil
}

type TransactionDetails struct {
	TxID       string   `json:"txID"`
	RawData    RawData  `json:"raw_data"`
	Ret        []Ret    `json:"ret"`
	Signature  []string `json:"signature"`
	RawDataHex string   `json:"raw_data_hex"`
}

// RawData represents the raw_data part of the transaction
type RawData struct {
	Contract      []Contract `json:"contract"`
	Expiration    int64      `json:"expiration"`
	RefBlockBytes string     `json:"ref_block_bytes"`
	RefBlockHash  string     `json:"ref_block_hash"`
	Timestamp     int64      `json:"timestamp"`
}

// Contract represents the contract part of the transaction
type Contract struct {
	Parameter Parameter `json:"parameter"`
	Type      string    `json:"type"`
}

// Parameter represents the parameter part of the contract
type Parameter struct {
	Value Value `json:"value"`
}

// Value represents the value part of the parameter
type Value struct {
	Amount       int    `json:"amount"`
	OwnerAddress string `json:"owner_address"`
	ToAddress    string `json:"to_address"`
}

// Ret represents the return status of the transaction
type Ret struct {
	ContractRet string `json:"contractRet"`
}

func GetTransactionDetails(txID string) (*TransactionDetails, error) {
	url := "https://api.trongrid.io/walletsolidity/gettransactionbyid"

	// Create the JSON payload
	payload := map[string]string{"value": txID}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("TRON-PRO-API-KEY", apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var transactionDetails TransactionDetails
	err = json.Unmarshal(body, &transactionDetails)
	if err != nil {
		return nil, err
	}

	return &transactionDetails, nil
}
func doubleSHA256(input []byte) []byte {
	hash1 := sha256.Sum256(input)
	hash2 := sha256.Sum256(hash1[:])
	return hash2[:]
}

func HexToBase58Check(hexAddress string) (string, error) {
	bytes, err := hex.DecodeString(hexAddress)
	if err != nil {
		return "", err
	}

	// Base58Check encoding: version byte 0x41 + payload + checksum
	payload := append([]byte{0x41}, bytes...)
	checksum := doubleSHA256(payload)[:4]
	addressBytes := append(payload, checksum...)

	return base58.Encode(addressBytes), nil
}
