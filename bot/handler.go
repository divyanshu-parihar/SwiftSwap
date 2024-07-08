package bot

import (
	"fmt"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleCommands(bot *tg.BotAPI, update tg.Update, command string, userAddWallet map[int64]*AddWalletsForm, userTrustedWallets map[int64]*UserTrustedWalletForm) {
	// TODO : Handle commands
	msg := tg.NewMessage(update.Message.Chat.ID, "")

	switch command {
	case "start":
		fmt.Println("Start")
		Start(bot, update)
	case "add":
		AddWallet(bot, update, userAddWallet)
	case "deposit":
		Deposit(bot, update)
	case "trust":
		AddUserTrustedWallets(bot, update, userTrustedWallets)
	case "help":
		Help(bot, update)
	case "status":
		Status(bot, update)
	case "verify":
		VerifyAndConvert(bot, update)

	default:
		msg.Text = "I don't know that command"
	}
	bot.Send(msg)

}
