package config

import (
	"os"
	"strconv"
)

type Limits struct {
	MaxReserve int64
}

type Config struct {
	Workers   int
	BatchSize int
	Limits    *Limits
}

func Load() *Config {
	return &Config{
		Workers:   getInt("INV_WORKERS", 2),
		BatchSize: getInt("INV_BATCH_SIZE", 2),
		Limits:    &Limits{MaxReserve: 1000},
	}
}

func getInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return def
	}
	return n
}
