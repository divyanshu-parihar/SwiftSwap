package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func Help(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.Text = "Contact  : @divvy7. (Drop a message, I'll take care of it. Trust me   🤝)"
	bot.Send(msg)
}
