package job

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

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
		log.Printf("Doing cron at: %s", time.Now().Format("2006-01-02 15:04:05"))
		iter := job.db.NewIterator(nil, nil)
		for iter.Next() {
			// Remember that the contents of the returned slice should not be modified, and
			// only valid until the next call to Next.
			key := iter.Key()
			value := iter.Value()
			var valueStruct []dbsvc.PriceNotification
			buffer := bytes.NewBuffer(value)
			decoder := gob.NewDecoder(buffer)
			if err := decoder.Decode(&valueStruct); err != nil {
				continue
			}

			newKey, err := utils.ByteToInt64(key)
			if err != nil {
				continue
			}

			for _, val := range valueStruct {
				price, err := GetStockPrice(val.Symbol)
				if err != nil {
					log.Println("err: ", err)
					continue
				}
				if price >= float64(val.High) {
					msg := tgbotapi.NewMessage(newKey, fmt.Sprintf("The price of %s is higher than %f", val.Symbol, val.High))
					job.teleBot.Send(msg)
				}
				if price <= float64(val.Low) {
					msg := tgbotapi.NewMessage(newKey, fmt.Sprintf("The price of %s is lower than %f", val.Symbol, val.Low))
					job.teleBot.Send(msg)
				}
			}
		}

	})
	if err != nil {
		fmt.Printf("Error scheduling cron job: %v\n", err)
		return
	}

	_, err = job.cron.AddFunc(config.Config.CRON_BOT_STILL_ALIVE, func() {
		iter := job.db.NewIterator(nil, nil)
		for iter.Next() {
			// Remember that the contents of the returned slice should not be modified, and
			// only valid until the next call to Next.
			key := iter.Key()
			value := iter.Value()
			var valueStruct []dbsvc.PriceNotification
			buffer := bytes.NewBuffer(value)
			decoder := gob.NewDecoder(buffer)
			if err := decoder.Decode(&valueStruct); err != nil {
				continue
			}

			newKey, err := utils.ByteToInt64(key)
			if err != nil {
				continue
			}
			msg := tgbotapi.NewMessage(newKey, "Bot still alive")
			job.teleBot.Send(msg)
		}
	})
	if err != nil {
		fmt.Printf("Error scheduling cron job: %v\n", err)
		return
	}

	job.cron.Start()
}

func (job *BackgroundJob) GetStockPrice() {

}

type StockData struct {
	Chart struct {
		Result []struct {
			Meta struct {
				RegularMarketPrice float64 `json:"regularMarketPrice"`
			} `json:"meta"`
			Timestamp  []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Open   []float64 `json:"open"`
					Low    []float64 `json:"low"`
					Close  []float64 `json:"close"`
					Volume []int64   `json:"volume"`
					High   []float64 `json:"high"`
				} `json:"quote"`
				Adjclose []struct {
					Adjclose []float64 `json:"adjclose"`
				} `json:"adjclose"`
			} `json:"indicators"`
		} `json:"result"`
		Error interface{} `json:"error"`
	} `json:"chart"`
}

func GetStockPrice(symbol string) (float64, error) {
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s?range=1d&interval=1d", symbol)
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to get stock price: %v", err)
	}
	defer resp.Body.Close()

	var stockPrice StockData
	if err := json.NewDecoder(resp.Body).Decode(&stockPrice); err != nil {
		return 0, fmt.Errorf("failed to decode stock price response: %v", err)
	}

	if len(stockPrice.Chart.Result) == 0 {
		return 0, fmt.Errorf("no stock price found for symbol: %s", symbol)
	}

	return stockPrice.Chart.Result[0].Meta.RegularMarketPrice, nil
}
