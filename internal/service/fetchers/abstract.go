package fetchers

import (
	"context"
	"kmtt/internal/entities"
	"time"
)

type PriceFetcher interface {
	FetchInterval() time.Duration
	SetObserveCoins(coins []entities.Coin)
	FetchPrice(ctx context.Context) error
	GetRates(coins []entities.Coin, from, to time.Time) []entities.ServiceRates
}
