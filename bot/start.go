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

	ethWallet, _, _ := wg.CreateWallet()
	if err := d.AddUserWalletToDB(db, update.Message.From.ID, ethWallet.Chain, "ETH", ethWallet.Address, "NONE", true, ethWallet.PrivateKey); err != nil {
		fmt.Println("Error saving wallet", err)

		SendMessage(update, bot, "*Error saving eth wallet*")
		return
	}

	tronWallet := wg.GenerateTronwallet()
	if err := d.AddUserWalletToDB(db, update.Message.From.ID, "TRC20", "TRON", tronWallet.AddressBase58, "NONE", true, tronWallet.PrivateKey); err != nil {
		fmt.Println("Error saving wallet", err)

		SendMessage(update, bot, "*Error saving eth wallet*")
		return
	}

	msg.Text += "\n\n*You have been assigned a wallet*"
	msg.Text += "\n═ Your Wallets ═"
	msg.Text += "\nw1: [ETH](https://ethscan.io/address/" + ethWallet.Address + "): ⟠ "
	msg.Text += "\n\n*ETH Address* 	```" + ethWallet.Address + "```  \n*ETH Private Key(Don't share this, We don't store it.)* 	```" + ethWallet.PrivateKey + "```"
	msg.Text += "\nw1: [TRON](https://tronscan.org/#/address/" + tronWallet.AddressBase58 + "): ⟠ "
	msg.Text += "\n\n*TRON Address* 	```" + tronWallet.AddressBase58 + "```  \n*TRON Private Key(Don't share this, We don't store it.)* 	```" + tronWallet.PrivateKey + "```"
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}
