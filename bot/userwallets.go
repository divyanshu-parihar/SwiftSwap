package bot

import (
	d "crypto-exchange-swap/db"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	Bun "github.com/uptrace/bun"
)

func GetUserWallets(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *Bun.DB) {
	wallets, err := d.GetWallet(db, update.Message.From.ID)
	if len(wallets) == 0 {
		SendMessage(update, bot, "* You Don't have any wallets*")
	}
	if err != nil {
		SendMessage(update, bot, "*Error checking wallet*")
	}
	SendMessage(update, bot, "*YOUR WALLETS*\n You can check your balance with /balance command\n"+"*ETH|BASE*: ```"+wallets[0].Address+"```"+"*TRON*: ```"+wallets[1].Address+"```")
}
