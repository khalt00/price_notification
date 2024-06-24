package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/khalt00/price_notification/pkg/config"
	"github.com/khalt00/price_notification/pkg/dbsvc"
	"github.com/syndtr/goleveldb/leveldb"
)

type TelegramBot struct {
	config *config.Configuration
	bot    *tgbotapi.BotAPI
	db     *leveldb.DB
}

func NewTelegramService(config *config.Configuration, bot *tgbotapi.BotAPI, db *leveldb.DB) *TelegramBot {
	return &TelegramBot{config, bot, db}
}

func (bot *TelegramBot) Listener() {
	update := tgbotapi.NewUpdate(0)
	update.Timeout = 60
	bot.bot.Debug = true

	updates := bot.bot.GetUpdatesChan(update)

	for u := range updates {
		if u.Message == nil {
			continue
		}
		if u.Message != nil {
			chatID := u.Message.Chat.ID
			if u.Message.Command() == string(Register) {

				var values []dbsvc.PriceNotification
				for _, v := range bot.config.PRICE_NOTIFICATIONS {
					values = append(values, dbsvc.PriceNotification{
						Symbol: v.Symbol,
						High:   v.High,
						Low:    v.Low,
					})
				}
				dbsvc.PutLevelDB(bot.db, chatID, values)
				msg := tgbotapi.NewMessage(chatID, "Registration completed successfully")
				bot.bot.Send(msg)
			}
		}
	}
}
