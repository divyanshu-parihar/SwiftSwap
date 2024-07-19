package walletGen

import (
	"context"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	// "github.com/btcsuite/btcd/chaincfg"
	// "github.com/btcsuite/btcec"
	// "github.com/btcsuite/btcutil"
	// "github.com/ethereum/go-ethereum/common"
	// "github.com/ethereum/go-ethereum/core/types"
	// "github.com/ethereum/go-ethereum/crypto"
	// "github.com/ethereum/go-ethereum/ethclient"
	// "github.com/fbsobreira/gotron-sdk/pkg/common"
	// "github.com/fbsobreira/gotron-sdk/pkg/crypto"
)

// Configuration constants for Ethereum and Coinbase Base Chain
const (
	ethMainnetURL = "https://mainnet.infura.io/v3/f889505d408845eeb60528672e3b61fb"
	baseChainURL  = "https://mainnet.base.org" // Replace with actual URL if different
)

type NWallet struct {
	PrivateKey string
	Chain      string
	Address    string
}

type TronWallet struct {
	PrivateKey string
	Address    string
}

// func CreateTronWallet() (*NWallet, error) {
// 	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to generate private key: %v", err)
// 	}

// 	privKeyBytes := crypto.FromECDSA(privKey)
// 	privKeyHex := hex.EncodeToString(privKeyBytes)

// 	pubKey := privKey.PublicKey
// 	pubKeyBytes := elliptic.Marshal(elliptic.P256(), pubKey.X, pubKey.Y)
// 	addressHex := address.PubkeyToAddress(pubKeyBytes).Hex()

//		return &TronWallet{
//			PrivateKey: privKeyHex,
//			Address:    addressHex,
//		}, nil
//	}
func CreateWallet() (*NWallet, *NWallet, string) {
	// Generate a random 256-bit seed
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		log.Fatalf("Failed to generate entropy: %v", err)
	}

	// Generate a mnemonic for the entropy
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		log.Fatalf("Failed to generate mnemonic: %v", err)
	}

	fmt.Println("Mnemonic:", mnemonic)

	// Generate a seed from the mnemonic
	seed := bip39.NewSeed(mnemonic, "")

	// Generate a master key from the seed
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		log.Fatalf("Failed to generate master key: %v", err)
	}

	// Derive the key for the Ethereum address (using the default path: m/44'/60'/0'/0/0)
	purpose, _ := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)
	coinType, _ := purpose.NewChildKey(bip32.FirstHardenedChild + 60)
	account, _ := coinType.NewChildKey(bip32.FirstHardenedChild + 0)
	change, _ := account.NewChildKey(0)
	addressKey, _ := change.NewChildKey(0)

	// Convert the derived key to an ECDSA private key
	privateKeyECDSA, err := crypto.ToECDSA(addressKey.Key)
	if err != nil {
		log.Fatalf("Failed to convert to ECDSA: %v", err)
	}

	// Generate the public key and Ethereum address
	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalf("Error casting public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	privateKeyBytes := crypto.FromECDSA(privateKeyECDSA)

	fmt.Printf("Wallet address: %s\n", address)
	ethClient, err := ethclient.Dial(ethMainnetURL)

	if err != nil {
		log.Fatalf("Failed to connect to Ethereum mainnet: %v", err)
	}
	fmt.Println("Connected to Ethereum mainnet")
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Printf("Using address: %s\n", fromAddress.Hex())

	// Get the nonce (number of transactions sent from this address)
	nonce, err := ethClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("Failed to get the nonce: %v", err)
	}

	// Get the suggested gas price
	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to suggest gas price: %v", err)
	}

	// Define the recipient address and transaction details
	toAddress := common.HexToAddress("RECIPIENT_ADDRESS")
	value := big.NewInt(0)    // 0 ETH
	gasLimit := uint64(21000) // in units

	// Create a new transaction
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	// Get the network ID for signing the transaction
	chainID, err := ethClient.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get network ID: %v", err)
	}

	// Sign the transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKeyECDSA)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}

	// Send the transaction
	err = ethClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}

	fmt.Printf("Transaction sent: %s\n", signedTx.Hash().Hex())
	// Connect to Coinbase Base Chain

	return &NWallet{
			PrivateKey: hex.EncodeToString(privateKeyBytes),
			Chain:      "Ethereum",
			Address:    address,
		}, &NWallet{
			PrivateKey: hex.EncodeToString(privateKeyBytes),
			Chain:      "Coinbase Base",
			Address:    address,
		},
		mnemonic
}

// Connect to Ethereum mainnet

// Transfer tokens on Ethereum mainnet
// transferTokens(ethClient, privateKey, recipientAddress, transferAmountWei, "Ethereum")

// Transfer tokens on Coinbase Base Chain
// transferTokens(baseClient, privateKey, recipientAddress, transferAmountWei, "Coinbase Base Chain")

func CheckBalance(client *ethclient.Client, address common.Address, networkName string) {
	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		log.Fatalf("Failed to retrieve balance from %s: %v", networkName, err)
	}
	fmt.Printf("Balance on %s: %s\n", networkName, balance.String())
}

func transferTokens(client *ethclient.Client, privateKey *ecdsa.PrivateKey, recipient string, amount int64, networkName string) {
	// Get the sender address
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalf("Error casting public key to ECDSA")
	}
	senderAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Get the nonce
	nonce, err := client.PendingNonceAt(context.Background(), senderAddress)
	if err != nil {
		log.Fatalf("Failed to get nonce from %s: %v", networkName, err)
	}

	// Get the gas price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to get gas price from %s: %v", networkName, err)
	}

	// Create the transaction
	toAddress := common.HexToAddress(recipient)
	tx := types.NewTransaction(nonce, toAddress, big.NewInt(amount), uint64(21000), gasPrice, nil)

	// Sign the transaction
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get chain ID from %s: %v", networkName, err)
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatalf("Failed to sign transaction for %s: %v", networkName, err)
	}

	// Send the transaction
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("Failed to send transaction to %s: %v", networkName, err)
	}

	fmt.Printf("Transaction sent to %s: %s\n", networkName, signedTx.Hash().Hex())
}

// btc wallet

func CreateBtcWallet() *NWallet {

	privKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	// Convert the private key to bytes
	// privKeyBytes := privKey.D.Bytes()

	// // Create a new WIF (Wallet Import Format) for the private key
	// wif, err := btcutil.NewWIF((*btcec.PrivateKey)(privKey), &chaincfg.MainNetParams, true)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Derive the public key from the private key
	pubKey := privKey.PublicKey()
	pubKeyBytes := pubKey.Bytes()

	// Generate a public key hash
	pubKeyHash := btcutil.Hash160(pubKeyBytes)

	// Create a Bitcoin address
	address, err := btcutil.NewAddressPubKeyHash(pubKeyHash, &chaincfg.MainNetParams)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Private Key (Hex):", hex.EncodeToString(privKey.Bytes()))

	// Print the Bitcoin address
	fmt.Println("Bitcoin Address:", address.EncodeAddress())

	// Print the public key in hexadecimal format
	fmt.Println("Public Key (Hex):", hex.EncodeToString(pubKeyBytes))

	return &NWallet{PrivateKey: hex.EncodeToString(pubKeyBytes), Chain: "Bitcoin", Address: address.EncodeAddress()}

}

type BalanceResponse struct {
	FinalBalance int64 `json:"final_balance"`
}

func GetBalance(address string) (int64, error) {
	url := fmt.Sprintf("https://blockchain.info/balance?active=%s", address)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to get balance: %s", resp.Status)
	}
	var result map[string]BalanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	balance, exists := result[address]

	if !exists {
		return 0, fmt.Errorf("address not found in response")
	}

	return int64(balance.FinalBalance), nil
}
