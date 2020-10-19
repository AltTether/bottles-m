package bottles

import (
	"time"
)

func NewTestConfig() *Config {
	return &Config{
		Seed:                   42,
		TokenSize:              10,
		TokenExpiration:        time.Duration(10 * time.Millisecond),
		SendBottleDelay:        time.Duration(15 * time.Millisecond),
		GenerateBottleCoolTime: time.Duration(10 * time.Millisecond),
	}
}
