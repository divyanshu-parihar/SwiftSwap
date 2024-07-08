package bot

import (
	wallets "crypto-exchange-swap/wallets"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Deposit(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	var markdown strings.Builder
	markdown.WriteString("*Deposit Addresses*\n\n")
	for _, address := range wallets.NewDepositWallets() {
		markdown.WriteString(fmt.Sprintf("*%s Wallet*\n", address.Crypto))
		markdown.WriteString(fmt.Sprintf("- *Chain:* %s\n", address.Chain))
		markdown.WriteString(fmt.Sprintf("- *ADDRESS:* `%s`", address.Addr))
		markdown.WriteString("\n\n")

	}
	markdown.WriteString("*⚠️ Please verify the chain before making any transactions to avoid loss of funds.*\n\n⚠️ We'll only consider further communication on this account.")
	msg.ParseMode = "markdown"
	msg.Text = markdown.String()
	bot.Send(msg)
}
