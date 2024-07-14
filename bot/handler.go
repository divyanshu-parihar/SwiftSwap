package bot

import (
	"fmt"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	Bun "github.com/uptrace/bun"
)

func handleCommands(db *Bun.DB, bot *tg.BotAPI, update tg.Update, command string, userAddWallet map[int64]*AddWalletsForm, userTrustedWallets map[int64]*UserTrustedWalletForm, userTronSendWallet map[int64]*UserTronSendReq, userConvertForm map[int64]*ConvertForm) {
	// TODO : Handle commands
	msg := tg.NewMessage(update.Message.Chat.ID, "")

	switch command {
	case "start":
		fmt.Println("Start")
		Start(db, bot, update)
	case "add":
		AddWallet(bot, update, userAddWallet)
	case "convert":
		Convert(update, bot, userConvertForm)
	case "deposit":
		Deposit(bot, update)
	case "balance":
		GetUserBalance(bot, update, db)
	case "wallet":
		GetUserWallets(bot, update, db)
	case "trust":
		AddUserTrustedWallets(bot, update, userTrustedWallets)
	case "send":
		SendTron(bot, update, userTronSendWallet)
	case "help":
		Help(bot, update)
	case "status":
		Status(bot, update)
	// case "verify":
	// 	VerifyAndConvert(bot, update)

	default:
		msg.Text = "I don't know that command"
	}
	bot.Send(msg)

}
