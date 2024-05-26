package fetchers

import (
	"kmtt/internal/entities"
	"kmtt/internal/utils"
	"math/big"
	"sync"
	"time"
)

type BaseFetcher struct {
	fetchInterval time.Duration
	data          map[time.Time]map[entities.Coin]*big.Float
	dataMU        sync.RWMutex
}

func NewBaseFetcher(fetchInterval time.Duration) BaseFetcher {
	return BaseFetcher{
		fetchInterval: fetchInterval,
		data:          make(map[time.Time]map[entities.Coin]*big.Float),
	}

}

func (f *BaseFetcher) GetFetchInterval() time.Duration {
	return f.fetchInterval
}

func (f *BaseFetcher) AddRates(timeEvent time.Time, data map[entities.Coin]*big.Float) {
	timeKey := utils.RoundToNearestInterval(timeEvent, f.fetchInterval)
	f.dataMU.Lock()
	f.data[timeKey] = data
	f.dataMU.Unlock()
}

func (f *BaseFetcher) GetRates(coins []entities.Coin, from, to time.Time) []entities.ServiceRates {
	times := utils.GetRangeBetween(from, to, f.fetchInterval)
	result := make([]entities.ServiceRates, 0, len(times))
	f.dataMU.RLock()
	defer f.dataMU.RUnlock()
	for _, t := range times {
		if _, ok := f.data[t]; !ok {
			continue
		}
		dataRates := make(map[entities.Coin]*big.Float, len(coins))
		for _, coin := range coins {
			if rate, ok := f.data[t][coin]; ok {
				dataRates[coin] = rate
			}
		}
		result = append(result, entities.ServiceRates{
			RateTime: t,
			Rates:    dataRates,
		})
	}
	return result
}
