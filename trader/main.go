package trader

import (
	handler "crypto-exchange-swap/handler"
	"fmt"
	"time"

	Bun "github.com/uptrace/bun"
)

func main() {
	fmt.Println("Helle, Trader's executor is live")
	// init mexc handler : handling all the interactions with mexc
	mexcHandler := handler.NewMexcHandler()
	orderTicker := time.NewTicker(2.0 * time.Second)
	for {
		select {
		case <-orderTicker.C:
			fmt.Println("Checking for orders")
			deposits, err := mexcHandler.GetAllDeposits()
			if err != nil {
				fmt.Printf("Couldn't deposits : %s", time.Now().UTC())
			}

			fmt.Println(deposits)
			if len(deposits) == 0 {
				fmt.Println("No deposits made in last 7 days")
			}



		}
	}
}

func checkOrderPayment(db *Bun.DB)
