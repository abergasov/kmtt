package entities

import (
	"math/big"
	"time"
)

type ServiceRates struct {
	RateTime time.Time           `json:"rate_time"`
	Rates    map[Coin]*big.Float `json:"rates"`
}
