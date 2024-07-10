package bot

import (
	d "crypto-exchange-swap/db"
	"crypto-exchange-swap/lang"
	wg "crypto-exchange-swap/walletGen"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	Bun "github.com/uptrace/bun"
)

func Start(db *Bun.DB, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	en := lang.NewLanguageText()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.Text = en.Intro

	// Check if user has a wallet
	wallet, err := d.GetWallet(db, update.Message.Chat.ID)
	exists := len(wallet) > 0
	if err != nil {
		fmt.Println("Error checking wallet", err)
		SendMessage(update, bot, "*Error checking wallet*")
		return
	}
	if exists {
		fmt.Println("Wallet exists")
		SendMessage(update, bot, "*You already have a wallet*")
		return
	}
	ethwallet, baseWallet := wg.CreateWallet()
	if err := d.AddUserWalletToDB(db, update.Message.Chat.ID, ethwallet.Chain, "ETH", ethwallet.Address, "", true, ethwallet.PrivateKey); err != nil {
		fmt.Println("Error saving wallet", err)
		SendMessage(update, bot, "*Error saving eth wallet*")
		return
	}

	msg.Text += "\n\n*You have been assigned a wallet*"
	msg.Text += "\n═ Your Wallets ═"
	msg.Text += "\nw1: [ETH](https://etherscan.io/address/" + ethwallet.Address + "): ⟠ | [BASE](https://basescan.org/address/" + baseWallet.Address + "): ⟠ "
	msg.Text += "\n\n*ETH Address | BASE Address* 	```" + ethwallet.Address + "```"

	msg.ParseMode = "Markdown"
	bot.Send(msg)
}
