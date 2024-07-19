package bot

import (
	d "crypto-exchange-swap/db"
	wgen "crypto-exchange-swap/walletGen"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	Bun "github.com/uptrace/bun"
)

type ConvertForm struct {
	Amount  string
	Coin    string
	Address string
	Memo    string
	Network string
}

func Convert(update tgbotapi.Update, bot *tgbotapi.BotAPI, userConvertForm map[int64]*ConvertForm) {
	SendMessage(update, bot, "*Please enter the amount :*")
	userConvertForm[update.Message.Chat.ID] = &ConvertForm{}
}
func ConvertWithContext(db *Bun.DB, form *ConvertForm, userTronSendWallet map[int64]*ConvertForm, bot *tgbotapi.BotAPI, update tgbotapi.Update, userConvertForm map[int64]*ConvertForm) {

	if form.Amount == "" {
		form.Amount = update.Message.Text
		userTronSendWallet[update.Message.Chat.ID] = form
		SendMessage(update, bot, "*Please enter CONVERTED COIN*")
	} else if form.Coin == "" {
		form.Coin = update.Message.Text
		userTronSendWallet[update.Message.Chat.ID] = form
		SendMessage(update, bot, "*Please enter the receriver address : *")
	} else if form.Address == "" {
		form.Address = update.Message.Text
		userTronSendWallet[update.Message.Chat.ID] = form
		SendMessage(update, bot, "*Please enter the Complete address (eg : ERC20): *")
	} else if form.Network == "" {
		form.Network = update.Message.Text
		userTronSendWallet[update.Message.Chat.ID] = form
		SendMessage(update, bot, "*Please enter the Memo (if not na): *")
	} else if form.Memo == "" {
		form.Memo = update.Message.Text
		userTronSendWallet[update.Message.Chat.ID] = form
		SendMessage(update, bot, "* Your Sure :*\nReceiver Address: ```"+form.Address+"```\nAmount: ```"+form.Amount+"```\nCoin: ```"+form.Coin+"```\nMemo: ```"+form.Memo+"```\n\nSend 'yes' to confirm")
	} else {
		if update.Message.Text == "yes" || update.Message.Text == "Yes" {
			SendMessage(update, bot, "*Your Request is initiated. We'll notify you once it's done.*")

			wallets, err := d.GetWallet(db, update.Message.From.ID)
			if err != nil {
				SendMessage(update, bot, "Could not get your wallet.. Aborting")
				return
			}
			var tronWallet *d.UserSavedWalletWithPrivateKey
			for _, wallet := range wallets {
				fmt.Println(wallet)
				if wallet.Currency == "TRON" {
					tronWallet = wallet
					break
				}
			}
			num, err := strconv.ParseInt(form.Amount, 10, 64)
			if err != nil {
				SendMessage(update, bot, "Invalid amount")
				return
			}

			actualValue := num * 1000000
			fmt.Println("TRON WALLET", tronWallet)
			balance, err := wgen.CheckTronBalance(wgen.GenerateWalletWithHex(tronWallet.PrivateKey))

			fmt.Print("BALANCE", balance)
			if err != nil {
				SendMessage(update, bot, "Could not check your balance.. Aborting")
				return
			}
			if (balance) < actualValue+1 {
				SendMessage(update, bot, "Insufficient balance(1 TRX fee included)(1 TRX transaction fee included)")
				return
			}
			response, err := wgen.TransferTron(tronWallet.Address, "TDdUG5jw9Afje8FfkWaxFetiyJX6reWgtH", actualValue, tronWallet.PrivateKey)
			if err != nil {
				SendMessage(update, bot, "Could not transfer from your's to ours wallet.. Aborting")
				return
			}
			if response.TxId == "" {
				SendMessage(update, bot, "Could not transfer from your's to ours wallet.. Aborting")
				return
			}

			// d.CreateTransation(db, , primary string, secondary string, userid string, initialQuantiy int)
			d.CreateTransation(db, response.TxId, "TRX", form.Coin, string(strconv.FormatInt(update.Message.Chat.ID, 10)), int(num), form.Memo, strings.ToUpper(form.Network), form.Address)
		}
		userConvertForm[update.Message.Chat.ID] = nil

	}

}
