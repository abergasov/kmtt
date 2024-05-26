package main

import (
	"context"
	"flag"
	"kmtt/internal/config"
	"kmtt/internal/entities"
	"kmtt/internal/logger"
	"kmtt/internal/service/fetchers/coingecko"
	"kmtt/internal/service/fetchers/currencylayer"
	"kmtt/internal/service/pricer"
	"os"
	"os/signal"
	"syscall"
)

var (
	confFile = flag.String("config", "configs/app_conf.yml", "Configs file path")
	appHash  = os.Getenv("GIT_HASH")
)

func main() {
	flag.Parse()
	appLog := logger.NewAppSLogger(appHash)

	appLog.Info("app starting", logger.WithString("conf", *confFile))
	appConf, err := config.InitConf(*confFile)
	if err != nil {
		appLog.Fatal("unable to init config", err, logger.WithString("config", *confFile))
	}

	appLog.Info("init fetchers")
	fetcherCoinGecko := coingecko.NewFetcher(appConf.CoinGeckoAPI, appConf.CoinGeckoDuration)
	fetcherCurrencyLayer := currencylayer.NewFetcher(appConf.CurrencyLayerAPI, appConf.CurrencyLayerDuration)

	appLog.Info("init services")
	service := pricer.NewService(appLog, entities.AllCoins)
	if err = service.RegisterFetcher("CoinGecko", fetcherCoinGecko); err != nil {
		appLog.Fatal("unable to register CoinGecko fetcher", err)
	}
	if err = service.RegisterFetcher("CurrencyLayer", fetcherCurrencyLayer); err != nil {
		appLog.Fatal("unable to register CurrencyLayer fetcher", err)
	}

	fetcherCurrencyLayer.FetchPrice(context.Background())

	appLog.Info("start services")
	service.Start()

	// register app shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c // This blocks the main thread until an interrupt is received
	if err = service.Stop(); err != nil {
		appLog.Fatal("unable to stop service", err)
	}
	appLog.Info("app stopped")
}
