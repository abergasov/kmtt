## Overview
### Challenge #1 
Build a currency conversion gRPC microservice.
1. Prepare .proto file with gRPC service & message type definitions (no implementation
   required). Describe service methods to perform conversion, batch conversion, listing
   of rates with pagination.
   * see `rpc/service.proto`
   * implementation to fetch price from multiply sources: `internal/service/pricer/service.go:GetRates`
2. To perform conversion fast we want to store live exchange rates in local memory and
   sync them with some external source via API e.g (Coingecko API or CurrencyLayer).
   Implement part of exchange rates sync between external API and local memory
   (concurrency should be handled). Also define an interface which provider should
   implement to fetch live rates suitable for different providers
   (Coingecko/CurrencyLayer etc)
   * each fetcher should implement `internal/service/fetchers/abstract.go:PriceFetcher`. implementations:
     * coingeko fetcher: `internal/service/fetchers/coingecko.go`
     * currency layer fetcher: `internal/service/fetchers/currencylayer.go`
   * each fetcher store rates in memory of object. get rates from memory: `internal/service/pricer/service.go:GetRates`

## Run
local run:
1. create `configs/app_conf.yml`
```shell
cp configs/sample.app_conf.yml configs/app_conf.yml1
# fill in the values
make run

# code quality check 
make lint && make test
```
