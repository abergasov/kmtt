package pricer

import (
	"context"
	"time"
)

func (s *Service) fetchPricer(fetcherName string) {
	ticker := time.NewTicker(s.priceFetchers[fetcherName].FetchInterval())
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.wrapFetcher(fetcherName)
		}
	}
}

func (s *Service) wrapFetcher(fetcherName string) {
	s.wg.Add(1)
	defer s.wg.Done()
	ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)
	defer cancel()
	if err := s.priceFetchers[fetcherName].FetchPrice(ctx); err != nil {
		s.log.Error("failed to fetch price", err)
	}
}
