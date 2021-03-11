package bottles

import (
	"time"
)


type Config struct {
	Seed                   int
	SendBottleDelay        time.Duration
	GenerateBottleCoolTime time.Duration
}

func NewConfig() *Config {
	return &Config{
		Seed:                   42,
		SendBottleDelay:        time.Duration(15 * time.Minute),
		GenerateBottleCoolTime: time.Duration(1 * time.Minute),
	}
}
