package pricer

import (
	"context"
	"kmtt/internal/entities"
	"kmtt/internal/logger"
	"kmtt/internal/service/fetchers"
	"sync"
	"time"
)

type Service struct {
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	log           logger.AppLogger
	observeCoins  []entities.Coin
	priceFetchers map[string]fetchers.PriceFetcher
}

func NewService(log logger.AppLogger, observeCoins []entities.Coin) *Service {
	srv := &Service{
		log:           log.With(logger.WithString("service", "pricer")),
		priceFetchers: make(map[string]fetchers.PriceFetcher),
		observeCoins:  observeCoins,
	}
	srv.ctx, srv.cancel = context.WithCancel(context.Background())
	return srv
}

func (s *Service) RegisterFetcher(serviceName string, f fetchers.PriceFetcher) error {
	f.SetObserveCoins(s.observeCoins)
	s.priceFetchers[serviceName] = f
	return nil
}

func (s *Service) Start() {
	for fetcherName := range s.priceFetchers {
		go s.fetchPricer(fetcherName)
	}
}

func (s *Service) Stop() error {
	s.cancel()
	return nil
}

func (s *Service) GetRates(coins []entities.Coin, from, to time.Time) map[string][]entities.ServiceRates {
	var (
		wg     sync.WaitGroup
		mu     sync.Mutex
		result = make(map[string][]entities.ServiceRates, len(s.priceFetchers))
	)
	wg.Add(len(s.priceFetchers))
	for i := range s.priceFetchers {
		go func(j string) {
			defer wg.Done()
			data := s.priceFetchers[j].GetRates(coins, from, to)
			mu.Lock()
			result[j] = data
			mu.Unlock()
		}(i)
	}
	wg.Wait()
	return result
}
