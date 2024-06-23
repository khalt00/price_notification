package main

import (
	"log"

	"github.com/khalt00/price_notification/pkg/config"
	"github.com/khalt00/price_notification/pkg/job"
	"github.com/khalt00/price_notification/pkg/telegram"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	config, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal("ERROR - Failed to load config")
	}

	db, err := leveldb.OpenFile(config.LEVEL_DB_PATH, &opt.Options{})
	if err != nil {
		log.Fatal("ERROR - failed to load config")
	}
	defer db.Close()

	teleBot, err := tgbotapi.NewBotAPI(config.TELEGRAM_BOT_TOKEN)
	if err != nil {
		log.Fatal("ERROR - Invalid configuration for telegram bot")
	}
	teleService := telegram.NewTelegramService(config, teleBot, db)

	jobs := job.NewBackgroundJob(config, teleBot, db)
	jobs.Start()

	go teleService.Listener()

	for {
	}

}
