package bot

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func VerifyAndConvert(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	var markdown strings.Builder
	markdown.WriteString("*Deposit Addresses*\n\n")
	markdown.WriteString("Let's Get started with your Conversion. \n * Send me your transaction Id\n") // TODO : Add the current ping here
	msg.ParseMode = "markdown"
	msg.Text = markdown.String()
	bot.Send(msg)
}
