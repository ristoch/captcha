package config

import (
	"github.com/caarlos0/env/v11"
)

// BalancerConfig конфигурация для balancer
type BalancerConfig struct {
	// Server settings
	Host string `env:"HOST" envDefault:"localhost"`
	Port string `env:"PORT" envDefault:"8080"`

	// Logging
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`

	// Cleanup settings
	CleanupInterval int32 `env:"CLEANUP_INTERVAL" envDefault:"300"`
	StaleThreshold  int32 `env:"STALE_THRESHOLD" envDefault:"600"`

	// User blocking settings
	MaxAttempts      int32 `env:"MAX_ATTEMPTS" envDefault:"3"`
	BlockDurationMin int32 `env:"BLOCK_DURATION_MINUTES" envDefault:"5"`

	// Timeout settings
	MaxTimeoutAttempts int32 `env:"MAX_TIMEOUT_ATTEMPTS" envDefault:"3"`

	// Validation settings
	MinOverlapPct int32 `env:"MIN_OVERLAP_PCT" envDefault:"20"`

	// Tracing
	JaegerEndpoint string `env:"JAEGER_ENDPOINT" envDefault:""`
}

// LoadBalancerConfig загружает конфигурацию для balancer
func LoadBalancerConfig() (*BalancerConfig, error) {
	config := &BalancerConfig{}

	if err := env.Parse(config); err != nil {
		return nil, err
	}

	return config, nil
}
