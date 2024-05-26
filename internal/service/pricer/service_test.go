package pricer_test

import (
	"context"
	"kmtt/internal/entities"
	"kmtt/internal/logger"
	"kmtt/internal/service/fetchers"
	"kmtt/internal/service/pricer"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type testFetcher struct {
	fetchers.BaseFetcher
}

func newTestFetcher() *testFetcher {
	return &testFetcher{
		BaseFetcher: fetchers.NewBaseFetcher(100 * time.Millisecond),
	}
}

func (t *testFetcher) FetchInterval() time.Duration {
	return t.BaseFetcher.GetFetchInterval()
}

func (t *testFetcher) SetObserveCoins(_ []entities.Coin) {}

func (t *testFetcher) FetchPrice(_ context.Context) error {
	t.BaseFetcher.AddRates(time.Now(), map[entities.Coin]*big.Float{
		entities.BTC: big.NewFloat(1),
		entities.LTC: big.NewFloat(2),
		entities.ETH: big.NewFloat(3),
	})
	return nil
}

func (t *testFetcher) GetRates(coins []entities.Coin, from, to time.Time) []entities.ServiceRates {
	return t.BaseFetcher.GetRates(coins, from, to)
}

func TestService_GetRates(t *testing.T) {
	// given
	log := logger.NewAppSLogger("test")
	observeCoins := []entities.Coin{entities.BTC, entities.LTC, entities.ETH}
	srv := pricer.NewService(log, observeCoins)

	require.NoError(t, srv.RegisterFetcher(uuid.NewString(), newTestFetcher()))
	require.NoError(t, srv.RegisterFetcher(uuid.NewString(), newTestFetcher()))
	require.NoError(t, srv.RegisterFetcher(uuid.NewString(), newTestFetcher()))
	startTime := time.Now()

	// when
	srv.Start()

	// then
	require.Eventually(t, func() bool {
		rates := srv.GetRates(observeCoins, startTime, time.Now())
		if len(rates) != 3 {
			return false
		}
		for i := range rates {
			if len(rates[i]) == 0 {
				return false
			}
		}
		return true
	}, time.Second*2, time.Millisecond*100)

	t.Run("concurrent fetch", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(100)
		for i := 0; i < 100; i++ {
			go func() {
				defer wg.Done()
				require.Len(t, srv.GetRates(observeCoins, startTime, time.Now()), 3)
			}()
		}
		wg.Wait()
	})
}
