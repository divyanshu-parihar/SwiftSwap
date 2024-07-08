package bot

import (
	d "crypto-exchange-swap/db"
	"fmt"
	"strings"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	Bun "github.com/uptrace/bun"
)

func AddWallet(bot *tgbotapi.BotAPI, update tg.Update, userAddWallet map[int64]*AddWalletsForm) {
	fmt.Println("Add wallet")
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	var markdown strings.Builder
	userAddWallet[update.Message.Chat.ID] = &AddWalletsForm{
		User:    update.Message.Chat.ID,
		Network: "",
		Address: "",
		Coin:    "",
		Memo:    "None",
		Step:    1,
	}

	markdown.WriteString("Let's Get started with your Swap.\n") // TODO : Add the current ping here
	markdown.WriteString("*Please enter the network of the coin you want to add.*\n")
	msg.ParseMode = "markdown"
	msg.Text = markdown.String()
	bot.Send(msg)
}
func AddWalletWithContext(db *Bun.DB, form *AddWalletsForm, bot *tgbotapi.BotAPI, update tgbotapi.Update, userAddWallet map[int64]*AddWalletsForm, userID int64) {
	if form.Network == "" {
		form.Network = update.Message.Text
		userAddWallet[userID] = form
		SendMessage(update, bot, "*Which Coin ? (BTC):*")
	} else if form.Coin == "" {
		form.Coin = update.Message.Text
		userAddWallet[userID] = form
		SendMessage(update, bot, "*Enter your Address :*")
	} else if form.Address == "" {
		form.Address = update.Message.Text
		userAddWallet[userID] = form
		SendMessage(update, bot, "*Enter your Memo (NA if not applicable):*")
	} else if form.Memo == "None" {
		if update.Message.Text == "NA" || update.Message.Text == "na" {
			form.Memo = "None"
		} else {
			form.Memo = update.Message.Text
		}

		userAddWallet[userID] = form
		SendMessage(update, bot, "*You Sure ? (yes/no)*")
	} else {
		if update.Message.Text != "yes" && update.Message.Text != "Yes" && update.Message.Text != "YES" && update.Message.Text != "no" && update.Message.Text != "No" && update.Message.Text != "NO" {
			SendMessage(update, bot, "*I didn't get that. You Sure ? (yes/no)*")

		}

		if update.Message.Text == "no" || update.Message.Text == "No" || update.Message.Text == "NO" {
			userAddWallet[userID] = nil
			SendMessage(update, bot, "*Wallet not saved*")
			return
		}
		// Save the wallet
		userAddWallet[userID] = form
		fmt.Println("form", form)
		// AddWallet(bot *tg.BotAPI, update tg.Update, userAddWallet map[string]*AddWalletsForm)
		userForm := userAddWallet[userID]
		if err := d.AddUserWalletToDB(db, userForm.User, userForm.Network, userForm.Coin, userForm.Address, userForm.Memo, false, ""); err != nil {
			fmt.Println("Error saving wallet", err)
			SendMessage(update, bot, "*Error saving wallet*")
			return
		}

		response := "*Wallet added successfully\n*"
		fmt.Println("userForm", userForm.User)
		response += "*Network: *" + form.Network + "\n" + "*Coin: *" + form.Coin + "\n*Address: *" + form.Address + "\n" + "*Memo: *" + form.Memo

		SendMessage(update, bot, response)

		userAddWallet[userID] = nil
	}
}
