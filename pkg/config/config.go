package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type PriceNotification struct {
	Symbol string
	High   float64
	Low    float64
}

type Configuration struct {
	CRON_PRICE_TIME      string
	CRON_BOT_STILL_ALIVE string

	TELEGRAM_BOT_TOKEN     string
	LEVEL_DB_PATH          string
	COIN_MARKETCAP_API_KEY string

	PRICE_NOTIFICATIONS []PriceNotification
}

var Config *Configuration

func LoadConfig(path string) (*Configuration, error) {
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Configuration
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config:   %w", err)
	}

	Config = &config
	return &config, nil
}
