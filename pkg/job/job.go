package job

import (
	"bytes"
	"encoding/gob"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/khalt00/price_notification/pkg/config"
	"github.com/khalt00/price_notification/pkg/dbsvc"
	"github.com/khalt00/price_notification/pkg/utils"
	"github.com/robfig/cron/v3"
	"github.com/syndtr/goleveldb/leveldb"
)

type BackgroundJob struct {
	config  *config.Configuration
	teleBot *tgbotapi.BotAPI

	cron *cron.Cron
	db   *leveldb.DB
}

func NewBackgroundJob(config *config.Configuration, teleBot *tgbotapi.BotAPI, db *leveldb.DB) *BackgroundJob {
	var cronjob *cron.Cron
	return &BackgroundJob{config, teleBot, cronjob, db}
}

func (job *BackgroundJob) Start() {
	job.cron = cron.New()

	// // Schedule the job to run every 10 seconds
	_, err := job.cron.AddFunc(config.Config.CRON_PRICE_TIME, func() {
		iter := job.db.NewIterator(nil, nil)
		for iter.Next() {
			// Remember that the contents of the returned slice should not be modified, and
			// only valid until the next call to Next.
			key := iter.Key()
			value := iter.Value()
			// chatID, _ := utils.ByteToInt64(key)
			var valueStruct dbsvc.LevelDBValueExample
			buffer := bytes.NewBuffer(value)
			decoder := gob.NewDecoder(buffer)
			if err := decoder.Decode(&valueStruct); err != nil {
				continue
			}

			_, err := utils.ByteToInt64(key)
			if err != nil {
				continue
			}
			// msg := tgbotapi.NewMessage(newKey, data)
			// job.teleBot.Send(msg)

		}
		// fmt.Println("hehe")
	}) // "*/10 * * * * *" means every 10 seconds
	if err != nil {
		fmt.Printf("Error scheduling cron job: %v\n", err)
		return
	}

	job.cron.Start()
}

func (job *BackgroundJob) GetStockPrice() {

}
