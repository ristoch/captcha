package config

import (
	"github.com/caarlos0/env/v11"
)

type ServiceConfig struct {
	MaxAttempts      int32 `env:"MAX_ATTEMPTS" envDefault:"3"`
	BlockDurationMin int32 `env:"BLOCK_DURATION_MINUTES" envDefault:"5"`
	CleanupInterval  int32 `env:"CLEANUP_INTERVAL" envDefault:"300"`
	StaleThreshold   int32 `env:"STALE_THRESHOLD" envDefault:"600"`

	ComplexityMedium int32 `env:"COMPLEXITY_MEDIUM" envDefault:"50"`
}

func LoadServiceConfig() (*ServiceConfig, error) {
	config := &ServiceConfig{}

	if err := env.Parse(config); err != nil {
		return nil, err
	}

	return config, nil
}
