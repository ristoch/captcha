package config

import (
	"github.com/caarlos0/env/v11"
)

type BalancerConfig struct {
	Host     string `env:"HOST" envDefault:"localhost"`
	Port     string `env:"PORT" envDefault:"8080"`
	GRPCPort string `env:"GRPC_PORT" envDefault:"9090"`

	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`

	CleanupInterval int32 `env:"CLEANUP_INTERVAL" envDefault:"300"`
	StaleThreshold  int32 `env:"STALE_THRESHOLD" envDefault:"600"`

	MaxAttempts      int32 `env:"MAX_ATTEMPTS" envDefault:"3"`
	BlockDurationMin int32 `env:"BLOCK_DURATION_MINUTES" envDefault:"5"`

	MaxTimeoutAttempts int32 `env:"MAX_TIMEOUT_ATTEMPTS" envDefault:"3"`

	MinOverlapPct int32 `env:"MIN_OVERLAP_PCT" envDefault:"20"`

	JaegerEndpoint string `env:"JAEGER_ENDPOINT" envDefault:""`
}

func LoadBalancerConfig() (*BalancerConfig, error) {
	config := &BalancerConfig{}

	if err := env.Parse(config); err != nil {
		return nil, err
	}

	return config, nil
}
