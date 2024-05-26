package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	CoinGeckoAPI      string        `yaml:"coingecko_api"`
	CoinGeckoDuration time.Duration `yaml:"coingecko_fetch_duration"`

	CurrencyLayerAPI      string        `yaml:"currencylayer_api"`
	CurrencyLayerDuration time.Duration `yaml:"currencylayer_fetch_duration"`
}

func InitConf(confFile string) (*AppConfig, error) {
	file, err := os.Open(filepath.Clean(confFile))
	if err != nil {
		return nil, fmt.Errorf("error open config file: %w", err)
	}
	defer func() {
		if e := file.Close(); e != nil {
			log.Fatal("Error close config file", e)
		}
	}()

	var cfg AppConfig
	if err = yaml.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("error decode config file: %w", err)
	}

	return &cfg, nil
}
