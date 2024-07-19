package main

import (
	b "crypto-exchange-swap/bot"
	d "crypto-exchange-swap/db"
	ex "crypto-exchange-swap/handler"
	"crypto-exchange-swap/wallets"
	"sync"

	// "crypto-exchange-swap/wallets"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// entry point for the application
func main() {
	// wg.CreateWallet()
	err := godotenv.Load()
	args := os.Args
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db := d.StartServer()
	if len(args) > 1 && args[1] == "migrate" {

		fmt.Println("Running migrations...")
		err := d.RunMigrations(db)
		if err != nil {
			log.Fatalf("Error running migrations: %v", err)
		}
		fmt.Println("migrateions ran successfully")
		return
	}
	// wallet := walletGen.CreateBtcWallet()
	// fmt.Println(wallet)
	var wg sync.WaitGroup
	language := wallets.NewDepositWallets()
	go b.NewBot(&wg, db, language)
	wg.Add(1)
	// go t.Trader(db)
	// wg.Add(1)
	d.UpdateTransactionStatus(db, "5dfde67cd24a80005516bcbca6e7d14e7181f22e2551cf955f3e5464d97d2268", "pending")
	// value,err ""= ("bc1qzfr92rqd4rezwcvfwhf3rw8tuq8y4ve3u34n67")
	// value, err := cw.GetBalance("bc1p8pr29y9ypz9gzs4tgp45yauxtfzdc4yse94e204up20yh77s2lysdrnqw3")
	// fmt.Println(value, err)
	fmt.Println("Hello, Starting the bot")
	// response, err := walletGen.TransferTron("TEPMj7aPCoH6xmf2bjeegUMuoxP3HNGgCh", "TBcGvboH3ptu7GcRcJhYytUkwes6xZgkV7", 1, "bde36c8a93f39dc55e604ce65f0e1f1b6bb88991b77eca1ddd93115c0cbc9a2e")
	// fmt.Println(response, err)
	// coin := h.Coin{
	// 	Name:  "Bitcoin",
	// 	Token: "LTC",
	// }
	//
	// details, err := wgen.GetTransactionDetails("a9e243b07d885e2be7ce66150123e494f2a3e23a5b47671a3cad1caf285d854c")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(details)

	// if len(details.RawData.Contract) > 0 && details.RawData.Contract[0].Parameter.Value.OwnerAddress != "" {
	// 	hexAddress := details.RawData.Contract[0].Parameter.Value.OwnerAddress
	// 	base58Address, err := wgen.HexToBase58Check(hexAddress[2:]) // Remove the "41" prefix
	// 	if err != nil {
	// 		fmt.Println("Error converting address:", err)
	// 		return
	// 	}
	// 	// fmt.Printf("Transaction Details: %+v\n", details)
	// 	fmt.Printf("Sender (Owner) Address: %s\n", base58Address)
	// } else {
	// 	fmt.Println("Owner address not found in the transaction details")
	// }
	// fmt.Println("The balance is: ", value)
	// mexc := ex.NewMexcHandler()
	// return
	mexHandler := ex.NewMexcHandler()
	// coin := h.Coin{
	// 	Name:  "TRON",
	// 	Token: "TRX",
	// }
	//
	// ethMainnetURL := "https://holesky.infura.io/v3/f889505d408845eeb60528672e3b61fb"
	// startTime := time.Now()
	// ethClient, err := ethclient.Dial(ethMainnetURL)
	// if err != nil {
	// 	log.Fatalf("Failed to connect to Ethereum mainnet: %v", err)
	// }
	// fmt.Println("Connected to Ethereum mainnet")
	// middleTime := time.Now()
	// fmt.Println("Time taken to connect to Ethereum mainnet: ", time.Since(startTime).Milliseconds())
	// walletGen.CheckBalance(ethClient, common.HexToAddress("0xdf99A0839818B3f120EBAC9B73f82B617Dc6A555"), "Eth MainChain")

	// fmt.Println("Time take to check the wallet amount", time.Since(middleTime).Milliseconds())
	// fmt.Println("Now starting to withdraw")
	// mexc.WithDrawCryptoToWallet(coin, "TRX", 20, "TBcGvboH3ptu7GcRcJhYytUkwes6xZgkV7", "")
	mexHandler.CheckTheDespositTransaction("0xc32e70752a9d18a941116ea86d79015223d757973de8e93a5ea7884814692715")
	// mexc.BuyWithUSDC(coin, 10)
	// //
	// baseClient, err := ethclient.Dial("https://mainnet.base.org")
	// if err != nil {
	// 	log.Fatalf("Failed to connect to Coinbase Base Chain: %v", err)
	// }
	// fmt.Println("Connected to Coinbase Base Chain")
	// walletGen.CheckBalance(baseClient, common.HexToAddress("0x74293CD5dd47F03A46b686eCCf5bfaE2b6bC2A0F"), "coinbase wallet")
	wg.Wait()
}
