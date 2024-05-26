package entities

import (
	"fmt"
	"strings"
)

type Coin string

const (
	BTC  Coin = "BTC"
	ETH  Coin = "ETH"
	LTC  Coin = "LTC"
	DOT  Coin = "DOT"
	ATOM Coin = "ATOM"
)

var (
	AllCoins = []Coin{BTC, ETH, LTC, DOT, ATOM}
)

func CoinFromString(src string) (Coin, error) {
	src = strings.TrimSpace(src)
	switch strings.ToLower(src) {
	case "btc", "bitcoin":
		return BTC, nil
	case "eth", "ethereum":
		return ETH, nil
	case "ltc", "litecoin":
		return LTC, nil
	case "dot", "polkadot":
		return DOT, nil
	case "atom", "cosmos":
		return ATOM, nil
	default:
		return "", fmt.Errorf("unknown coin: %s", src)
	}
}
