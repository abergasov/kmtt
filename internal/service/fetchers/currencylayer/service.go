package currencylayer

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
			entities.BTC: "BTC",
			entities.LTC: "LTC",
		},
		BaseFetcher: fetchers.NewBaseFetcher(fetchInterval),
	}
}

func (f *Fetcher) FetchInterval() time.Duration {
	return f.BaseFetcher.GetFetchInterval()
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

func (f *Fetcher) FetchPrice(ctx context.Context) error {
	url := fmt.Sprintf("http://api.currencylayer.com/live?access_key=%s&source=USD&currencies=%s", f.apiToken, f.castedCoins) //nolint: gosec
	res, code, err := utils.GetCurl[struct {
		Success   bool               `json:"success"`
		Terms     string             `json:"terms"`
		Privacy   string             `json:"privacy"`
		Timestamp int64              `json:"timestamp"`
		Source    string             `json:"source"`
		Quotes    map[string]float64 `json:"quotes"`
	}](ctx, url, map[string]string{
		"accept": "application/json",
	})
	if err != nil {
		return fmt.Errorf("unable to fetch data: %w", err)
	}
	if code != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", code)
	}

	data := make(map[entities.Coin]*big.Float, len(res.Quotes))
	for i := range res.Quotes {
		coin, errC := entities.CoinFromString(strings.TrimPrefix(i, "USD"))
		if errC != nil {
			continue
		}
		data[coin] = calculateReversalQuote(res.Quotes[i])
	}
	f.BaseFetcher.AddRates(time.Now(), data)
	return nil
}

func calculateReversalQuote(quote float64) *big.Float {
	one := big.NewFloat(1.0)
	rate := new(big.Float).Quo(one, big.NewFloat(quote))
	return rate
}

func (f *Fetcher) GetRates(coins []entities.Coin, from, to time.Time) []entities.ServiceRates {
	return f.BaseFetcher.GetRates(coins, from, to)
}
