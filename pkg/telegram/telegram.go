package telegram

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
				msg := tgbotapi.NewMessage(chatID, "Registration completed successfully")
				bot.bot.Send(msg)
				// stocks := []string{"BTC-USD", "FPT.VN"}
				dbsvc.PutLevelDB(bot.db, chatID, dbsvc.LevelDBValueExample{
					// Datas: stocks,
				})
			}
			if u.Message.Command() == "get" {
				val, _ := dbsvc.GetLevelDB(bot.db, chatID)
				// var prices []float64
				for _, stock := range val.Datas {
					price, err := GetStockPrice(stock)
					if err != nil {
						log.Panic(err)
					}
					// prices = append(prices, price)
					msg := tgbotapi.NewMessage(chatID, fmt.Sprint(price))
					bot.bot.Send(msg)
				}
			}
		}
	}
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

// func FetchPrice(){
// 	for i := range marketRequests {
// 		request, _ := http.NewRequest("GET", fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s?range=1mo&interval=1d", marketRequests[i].Symbol), nil)
// 		requests = append(requests, request)
// 	}

// }
