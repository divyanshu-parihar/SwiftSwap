package walletGen

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Configuration constants for Ethereum and Coinbase Base Chain
const (
	ethMainnetURL = "https://mainnet.infura.io/v3/f889505d408845eeb60528672e3b61fb"
	baseChainURL  = "https://mainnet.base.org" // Replace with actual URL if different
)

func CreateWallet() {
	// Load the private key
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	// Get the public address
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalf("Error casting public key to ECDSA")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Printf("Wallet address: %s\n", address.Hex())

	// Connect to Ethereum mainnet
	ethClient, err := ethclient.Dial(ethMainnetURL)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum mainnet: %v", err)
	}
	fmt.Println("Connected to Ethereum mainnet")

	// Connect to Coinbase Base Chain
	baseClient, err := ethclient.Dial(baseChainURL)
	if err != nil {
		log.Fatalf("Failed to connect to Coinbase Base Chain: %v", err)
	}
	fmt.Println("Connected to Coinbase Base Chain")

	// Check balance on Ethereum mainnet
	checkBalance(ethClient, address, "Ethereum")

	// Check balance on Coinbase Base Chain
	checkBalance(baseClient, address, "Coinbase Base Chain")

	// Transfer tokens on Ethereum mainnet
	// transferTokens(ethClient, privateKey, recipientAddress, transferAmountWei, "Ethereum")

	// Transfer tokens on Coinbase Base Chain
	// transferTokens(baseClient, privateKey, recipientAddress, transferAmountWei, "Coinbase Base Chain")
}

func checkBalance(client *ethclient.Client, address common.Address, networkName string) {
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
