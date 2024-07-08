package bot

import (
	"strings"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendMessage(update tgbotapi.Update, bot *tgbotapi.BotAPI, content string) {
	msg := tg.NewMessage(update.Message.Chat.ID, "")
	var markdown strings.Builder
	markdown.WriteString(content)
	msg.ParseMode = "markdown"
	msg.Text = markdown.String()

	bot.Send(msg)

}
