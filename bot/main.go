package bot

import (
	w "crypto-exchange-swap/wallets"
	"log"
	"os"
	"sync"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	bun "github.com/uptrace/bun"
)

func NewBot(wg *sync.WaitGroup, db *bun.DB, wallets []*w.DespositWallet) {
	bot, err := tg.NewBotAPI(os.Getenv("BOT_TOKEN"))
	// en := lang.NewLanguageText()
	defer wg.Done()
	if err != nil {
		log.Panic(err)

	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tg.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	// Bot memory

	userAddWallet := make(map[int64]*AddWalletsForm)
	userTrustedWallet := make(map[int64]*UserTrustedWalletForm)
	for update := range updates {
		// Very imporant id( Primary key for most DS)
		userID := update.Message.Chat.ID
		if update.Message.IsCommand() {
			handleCommands(bot, update, update.Message.Command(), userAddWallet, userTrustedWallet)
			continue
		}
		if form, ok := userTrustedWallet[userID]; ok {
			AddUserTrustedWalletscontext(db, form, bot, update, userTrustedWallet, userID)
		}
		if form, ok := userAddWallet[userID]; ok {
			AddWalletWithContext(db, form, bot, update, userAddWallet, userID)
		}
	}
}
