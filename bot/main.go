package bot

import (
	"crypto-exchange-swap/lang"
	w "crypto-exchange-swap/wallets"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	bun "github.com/uptrace/bun"
)

func NewBot(wg *sync.WaitGroup, db *bun.DB, wallets []*w.DespositWallet) {
	bot, err := tg.NewBotAPI(os.Getenv("BOT_TOKEN"))
	en := lang.NewLanguageText()
	defer wg.Done()
	if err != nil {
		log.Panic(err)

	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tg.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		msg := tg.NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Command() {
		case "start":
			msg.Text = en.Intro
		case "deposit":
			// wallets.Deposit(db, update.Message.Chat.ID)
			msg.Text = "Deposit Addresses"
			// addresses := wallets

			// for _, address := range addresses {
			// 	msg.Text += `COIN: ` + address.Crypto + `Chain : ` + address.Chain + ` ADDRESS: ` + "`" + address.Addr + "`" + `\n\n`
			// }
			// var markdown strings.Builder

			// Header for wallets
			var markdown strings.Builder
			markdown.WriteString("*Deposit Addresses*\n\n")
			for _, address := range wallets {
				markdown.WriteString(fmt.Sprintf("*%s Wallet*\n", address.Crypto))
				markdown.WriteString(fmt.Sprintf("- *Chain:* %s\n", address.Chain))
				markdown.WriteString(fmt.Sprintf("- *ADDRESS:* `%s`", address.Addr))
				markdown.WriteString("\n\n")

			}
			markdown.WriteString("*⚠️ Please verify the chain before making any transactions to avoid loss of funds.*\n\n⚠️ We'll only consider further communication on this account.")
			msg.ParseMode = "markdown"
			msg.Text = markdown.String()
		case "help":
			msg.Text = "Contact  : @divvy7. (Drop a message, I'll take care of it. Trust me   🤝)"
		case "status":
			msg.Text = "I'm ok." // TODO : Add the current ping here
		default:
			msg.Text = "I don't know that command"
		}
		bot.Send(msg)

	}
}
