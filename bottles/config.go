package bottles

import (
	"time"
)


type Config struct {
	Seed                    int
	SendBottlePeriod        time.Duration
	GenerateBottleCoolTime  time.Duration
}

func NewConfig() *Config {
	return &Config{
		Seed:                    42,
		SendBottlePeriod:        time.Duration(15 * time.Minute),
		GenerateBottleCoolTime:  time.Duration(1 * time.Minute),
	}
}
