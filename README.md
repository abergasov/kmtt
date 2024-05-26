## Overview
### Challenge #1 
Build a currency conversion gRPC microservice.
1. Prepare .proto file with gRPC service & message type definitions (no implementation
   required). Describe service methods to perform conversion, batch conversion, listing
   of rates with pagination.
2. To perform conversion fast we want to store live exchange rates in local memory and
   sync them with some external source via API e.g (Coingecko API or CurrencyLayer).
   Implement part of exchange rates sync between external API and local memory
   (concurrency should be handled). Also define an interface which provider should
   implement to fetch live rates suitable for different providers
   (Coingecko/CurrencyLayer etc)

## Run
local run:
1. create `configs/app_conf.yml`
```shell
cp configs/sample.app_conf.yml configs/app_conf.yml1
# fill in the values
make run
```
