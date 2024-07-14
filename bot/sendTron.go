package bot

import (
	d "crypto-exchange-swap/db"
	wg "crypto-exchange-swap/walletGen"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	Bun "github.com/uptrace/bun"
)

type UserTronSendReq struct {
	ToAddress string
	amount    string
}

func SendTron(bot *tgbotapi.BotAPI, update tgbotapi.Update, userTronSendWallet map[int64]*UserTronSendReq) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.Text = "Send me the receiver address :  "
	userTronSendWallet[update.Message.Chat.ID] = &UserTronSendReq{}
	bot.Send(msg)
}

func SendTronWithContext(db *Bun.DB, form *UserTronSendReq, userTronSendWallet map[int64]*UserTronSendReq, bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	if form.ToAddress == "" {
		form.ToAddress = update.Message.Text
		userTronSendWallet[update.Message.Chat.ID] = form
		SendMessage(update, bot, "Send me the amount (Eg. 1 = 1TRX) :  ")
	} else if form.amount == "" {
		form.amount = update.Message.Text
		userTronSendWallet[update.Message.Chat.ID] = form
		SendMessage(update, bot, "*You sure* :  \n*Receiver Address* : *"+form.ToAddress+"*\n*Amount* : *"+form.amount+"*\n\n Send  'yes' to confirm")

	} else {
		fmt.Println("Sending TRX")
		fmt.Println(form)
		if update.Message.Text == "yes" || update.Message.Text == "Yes" {
			SendMessage(update, bot, "Transaction processing...")

			wallets, err := d.GetWallet(db, update.Message.From.ID)
			if err != nil {
				SendMessage(update, bot, "Could not get your wallet.. Aborting")
				return
			}
			fmt.Println(wallets)
			var tronWallet *d.UserSavedWalletWithPrivateKey
			for _, wallet := range wallets {
				fmt.Println(wallet)
				if wallet.Currency == "TRON" {
					tronWallet = wallet
					break
				}
			}
			num, err := strconv.ParseInt(form.amount, 10, 64)
			if err != nil {
				SendMessage(update, bot, "Invalid amount")
				return
			}
			response, err := wg.TransferTron(tronWallet.Address, form.ToAddress, num*1000000, tronWallet.PrivateKey)
			if err != nil {
				SendMessage(update, bot, "Error sending TRX : "+err.Error())
				return
			}
			SendMessage(update, bot, "Transaction DONE :"+response.TxId)
		}
		userTronSendWallet[update.Message.Chat.ID] = nil

	}

}
