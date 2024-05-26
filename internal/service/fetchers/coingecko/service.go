package coingecko

import (
	"context"
	"fmt"
	"kmtt/internal/entities"
	"kmtt/internal/service/fetchers"
	"kmtt/internal/utils"
	"math/big"
	"net/http"
	"strings"
	"time"
)

type Fetcher struct {
	apiToken    string
	castedCoins string
	mapCoins    map[entities.Coin]string
	fetchers.BaseFetcher
}

func NewFetcher(apiToken string, fetchInterval time.Duration) *Fetcher {
	return &Fetcher{
		apiToken: apiToken,
		mapCoins: map[entities.Coin]string{
			entities.BTC:  "bitcoin",
			entities.ETH:  "ethereum",
			entities.LTC:  "litecoin",
			entities.DOT:  "polkadot",
			entities.ATOM: "cosmos",
		},
		BaseFetcher: fetchers.NewBaseFetcher(fetchInterval),
	}
}

func (f *Fetcher) SetObserveCoins(coins []entities.Coin) {
	castCoins := make([]string, 0, len(coins))
	for _, coin := range coins {
		coinStr, ok := f.mapCoins[coin]
		if !ok {
			continue
		}
		castCoins = append(castCoins, coinStr)
	}
	f.castedCoins = strings.Join(castCoins, ",")
}

func (f *Fetcher) FetchInterval() time.Duration {
	return f.BaseFetcher.GetFetchInterval()
}

func (f *Fetcher) FetchPrice(ctx context.Context) error {
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", f.castedCoins)
	res, code, err := utils.GetCurl[map[string]struct {
		USD float64 `json:"usd"`
	}](ctx, url, map[string]string{
		"accept":            "application/json",
		"x-cg-demo-api-key": f.apiToken,
	})
	if err != nil {
		return fmt.Errorf("unable to fetch price: %w", err)
	}
	if code != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", code)
	}
	data := make(map[entities.Coin]*big.Float, len(*res))
	for coin := range f.mapCoins {
		coinStr := f.mapCoins[coin]
		if price, ok := (*res)[coinStr]; ok {
			data[coin] = big.NewFloat(price.USD)
		}
	}
	f.BaseFetcher.AddRates(time.Now(), data)
	return nil
}

func (f *Fetcher) GetRates(coins []entities.Coin, from, to time.Time) []entities.ServiceRates {
	return f.BaseFetcher.GetRates(coins, from, to)
}
