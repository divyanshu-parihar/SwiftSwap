package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func Status(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.Text = "I'm ok." // TODO : Add the current ping here
	bot.Send(msg)
}
