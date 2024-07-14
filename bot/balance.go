package bot

import (
	d "crypto-exchange-swap/db"
	"crypto-exchange-swap/walletGen"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	Bun "github.com/uptrace/bun"
)

func GetUserBalance(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *Bun.DB) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	userid := update.Message.From.ID
	wallets, err := d.GetWallet(db, userid)

	if len(wallets) < 2 {
		SendMessage(update, bot, "* You Don't have any wallets*")
	}
	if err != nil {
		SendMessage(update, bot, "*Error checking wallet*")
		return
	}
	// addr := wallets[1].Address
	balance, err := walletGen.GetAccountBalance(wallets[1].Address)
	if err != nil {
		SendMessage(update, bot, "*Error checking balance*")
		return
	}
	balanceTRX := float64(balance) / 1_000_000

	msg.Text = fmt.Sprintf("Your balance %.2f TRX", balanceTRX)

	bot.Send(msg)
}
