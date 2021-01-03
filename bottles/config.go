package bottles

import (
	"time"
)


type Config struct {
	Seed                   int
	TokenSize              int
	TokenExpiration        time.Duration
	SendBottleDelay        time.Duration
	GenerateBottleCoolTime time.Duration
}

func NewConfig() *Config {
	return &Config{
		Seed:                   42,
		TokenSize:              10,
		TokenExpiration:        time.Duration(1 * time.Minute),
		SendBottleDelay:        time.Duration(15 * time.Minute),
		GenerateBottleCoolTime: time.Duration(1 * time.Minute),
	}
}
