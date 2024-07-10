package bot

import (
	d "crypto-exchange-swap/db"
	"fmt"
	"strings"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	Bun "github.com/uptrace/bun"
)

func AddUserTrustedWallets(bot *tgbotapi.BotAPI, update tg.Update, userTrustWallets map[int64]*UserTrustedWalletForm) {
	fmt.Println("Add wallet")
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	var markdown strings.Builder
	userTrustWallets[update.Message.Chat.ID] = &UserTrustedWalletForm{
		User:    update.Message.Chat.ID,
		Address: "",
	}

	markdown.WriteString("*Great let's add a trusted wallet: \n*")
	markdown.WriteString("*Please enter the address: (crypto will only be transfered to these accounts Security++)*\n")
	msg.ParseMode = "markdown"
	msg.Text = markdown.String()
	bot.Send(msg)
}
func AddUserTrustedWalletscontext(db *Bun.DB, form *UserTrustedWalletForm, bot *tgbotapi.BotAPI, update tg.Update, userTrustWallets map[int64]*UserTrustedWalletForm, userid int64) {
	if form.Address == "" {
		form.Address = update.Message.Text
		// userTrustWallets[userid] = form
		SendMessage(update, bot, "*You Sure ? (yes/no)*")
	} else {
		if update.Message.Text != "yes" && update.Message.Text != "Yes" && update.Message.Text != "YES" && update.Message.Text != "no" && update.Message.Text != "No" && update.Message.Text != "NO" {
			SendMessage(update, bot, "*I didn't get that. You Sure ? (yes/no)*")

		}

		if update.Message.Text == "no" || update.Message.Text == "No" || update.Message.Text == "NO" {
			userTrustWallets[userid] = nil
			SendMessage(update, bot, "*Wallet not saved*")
			return
		}
		// Save the wallet
		userTrustWallets[userid] = form
		d.AddUserTrustedWalletToDB(db, userid, form.Address)
		// fmt.Println("form", form)
		userTrustWallets[userid] = nil

	}
}
