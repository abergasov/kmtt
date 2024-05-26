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

### Challenge #2 
Cosmos SDK module design
Prepare high level (without much technical details) Cosmos SDK module design in form of
documentation for business requirements mentioned onwards. You can use a structure similar
to the standard Cosmos module structure. Important parts to be included are: State, Keepers,
Messages and any other if required.

Business requirements
* As XYZ token minter I want to have the ability to freeze XYZ tokens on any account.
Frozen tokens can’t be moved outside of account in any way unless unfrozen.
* As XYZ token minter I want to have the ability to unfreeze XYZ token on any
account.

#### Design
in evm chain we can create erc20 compatible smart contract to freeze and unfreeze tokens. 
2 internal storages of the contract will store:
* the frozen amount of tokens for each account
* the allowed amount of tokens for each account

balance check will look at the allowed storage.

transfers of tokens will be done by the contract owner. pseudocode:
```solidity
contract XYZToken is ERC20 {
   mapping(address => uint256) private _frozen;
   mapping(address => uint256) private _allowed;

   function balanceOf(address account) external view returns (uint256) {
      return _allowed[account];
   }
   
   function freeze(address account, uint256 amount) external onlyOwner {
      if (_allowed[account] < amount) {
         revert("insufficient balance");
      }
      _frozen[account] += amount;
      _allowed[account] -= amount;
   }

   function unfreeze(address account, uint256 amount) external onlyOwner {
      if (_frozen[account] < amount) {
         revert("insufficient frozen balance");
      }
      _frozen[account] -= amount;
      _allowed[account] += amount;
   }
}
```
same approach can be applied to cosmos sdk module. 2 maps track frozen and allowed amounts of tokens for each account. 

token keeper will be responsible for storing the frozen amount of tokens for each account
```go
package keeper

import (
   "cosmossdk.io/core/store"
   "github.com/cosmos/cosmos-sdk/codec"
)

type TokenKeeper struct {
   storeKeyFrozen  sdk.StoreKey
   storeKeyAllowed sdk.StoreKey
   storeService    store.KVStoreService
   cdc             codec.BinaryCodec
}

func (k *TokenKeeper) Freeze(ctx sdk.Context, addr sdk.AccAddress, amount sdk.Coin) error {
	// get balance of the account by allowed key
	// verify that the account has enough balance
	// store the frozen amount in the store
	// increment amount by frozen amount, decrement amount by allowed amount
	return nil
}

func (k *TokenKeeper) UnFreeze(ctx sdk.Context, addr sdk.AccAddress, amount sdk.Coin) error {
	// get frozen balance of the account
	// verify that the account has enough balance for unfreeze
	// increment amount by allowed amount, decrement amount by frozen amount
	return nil
}
```

messages will be used to freeze and unfreeze tokens
```go
type MsgManageTokens struct {
	Address sdk.AccAddress `json:"address"`
	Amount  sdk.Coins      `json:"amount"`
}
```

handlers to process message
```go
func HandleMsgFreezeTokens(ctx sdk.Context, keeper TokenKeeper, msg MsgManageTokens) sdk.Result {
    if err := keeper.FreezeTokens(ctx, msg.Address, msg.Amount); err != nil {
        return sdk.ErrUnknownRequest(err.Error()).Result()
    }
    return sdk.Result{Events: ctx.EventManager().Events()}
}

func HandleMsgUnfreezeTokens(ctx sdk.Context, keeper TokenKeeper, msg MsgManageTokens) sdk.Result {
    if err := keeper.UnfreezeTokens(ctx, msg.Address, msg.Amount); err != nil {
        return sdk.ErrUnknownRequest(err.Error()).Result()
    }
    return sdk.Result{Events: ctx.EventManager().Events()}
}
```