package bot

import (
	"crypto-exchange-swap/lang"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Start(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	en := lang.NewLanguageText()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.Text = en.Intro
	bot.Send(msg)
}
