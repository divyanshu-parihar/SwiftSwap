package trader

import (
	d "crypto-exchange-swap/db"
	handler "crypto-exchange-swap/handler"
	helper "crypto-exchange-swap/helper"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	Bun "github.com/uptrace/bun"
)

type SoldResponse struct {
	order helper.OrderResponse
	txId  string
}
type WithdrawResponse struct {
	order helper.OrderResponse
	txId  string
}

func Trader(db *Bun.DB) {
	fmt.Println("Helle, Trader's executor is live")
	// init mexc handler : handling all the interactions with mexc
	mexcHandler := handler.NewMexcHandler()
	orderTicker := time.NewTicker(2.0 * time.Second)

	newDesposit := make(chan handler.DepostTransaction)
	soldDesposit := make(chan SoldResponse)
	withdrawHash := make(chan WithdrawResponse)
	for {
		select {

		case responsewithid := <-withdrawHash:
			go func() {
				// get the transaction
				order := responsewithid.order
				txnid := responsewithid.txId
				transaction, err := d.GetTransactionWithTxnID(db, txnid)

				if len(transaction) <= 0 {
					fmt.Println("No transaction found")
					return
				}
				if err != nil {
					fmt.Println("Error getting transaction by hash")
					return
				}

				const PlatformFeesPercentage = 0
				var totalqty, fees float64
				json.Unmarshal(order.Orders[0].Total_quantity, &totalqty)
				err = json.Unmarshal(order.Orders[0].Fee, &fees)
				if err != nil {
					fmt.Printf("Error unmarshaling fee: %v\n", err)
					return
				}
				response, err := mexcHandler.WithDrawCryptoToWallet(helper.Coin{Name: transaction[0].FinalCurrency}, transaction[0].Network, totalqty*0.999, transaction[0].Address, transaction[0].Memo)
				if err != nil {
					fmt.Println("withdraw error")
				}

				fmt.Println("Withdraw response ")
				fmt.Println(response)
			}()
		case event := <-soldDesposit:
			go func() {
				sold := event.order
				const PlatformFeesPercentage = 0
				var totalqty, fees float64
				err := json.Unmarshal(sold.Orders[0].Total_quantity, &totalqty)
				if err != nil {
					fmt.Printf("Error unmarshaling total quantity: %v\n", err)
					return
				}
				err = json.Unmarshal(sold.Orders[0].Fee, &fees)
				if err != nil {
					fmt.Printf("Error unmarshaling fee: %v\n", err)
					return
				}
				response, err := mexcHandler.BuyWithUSDC(helper.Coin{Name: "ETH", Token: "ETH"}, (totalqty-fees)*0.99)
				if err != nil {
					fmt.Println("Error buying with USDC")
				}
				d.UpdateTransactionStatus(db, event.txId, string(d.StatusCompleted))
				fmt.Println(response)

				withdrawHash <- WithdrawResponse{response, event.txId}

			}()

		case deposit := <-newDesposit:
			go func() {
				coin := helper.Coin{
					Name:  deposit.Coin,
					Token: strings.ToUpper(deposit.Coin),
				}
				num, err := strconv.ParseFloat(deposit.Amount, 64)
				if err != nil {
					fmt.Println("Error converting amount to int")
					return
				}
				response, err := mexcHandler.SellForUSDC(coin, num)
				if err != nil {
					fmt.Println("Error selling for USDC")
					/// TODO: update the transaction status to failed
				}
				d.UpdateTransactionStatus(db, deposit.TxId, string(d.StatusIncomplete))
				soldDesposit <- SoldResponse{response, deposit.TxId}
			}()
		case <-orderTicker.C:
			go func() {
				fmt.Println("Checking for orders")
				deposits, err := mexcHandler.GetAllDeposits()
				if err != nil {
					fmt.Printf("Couldn't deposits : %s", time.Now().UTC())
				}

				for _, deposit := range deposits {
					fmt.Println(deposit)

					// check if deposit is already in d
				}
				if len(deposits) == 0 {
					fmt.Println("No deposits made in last 7 days")
				}
				transactions, err := d.GetPendingTransaction(db)
				if err != nil {
					fmt.Printf("Couldn't get transactions : %s", time.Now().UTC())
				}
				fmt.Println(transactions)
				// TODO: try making this GODDAAM thing faster;
				for _, transaction := range transactions {
					for _, deposit := range deposits {
						if transaction.TxnID == deposit.TxId {
							d.UpdateTransactionStatus(db, transaction.TxnID, string(d.StatusDeposit))
							newDesposit <- deposit
						}
					}
				}
			}()
		}
	}
}

func handleSellingForUSDT() {

}
