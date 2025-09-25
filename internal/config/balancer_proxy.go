package config

import (
	"github.com/caarlos0/env/v11"
)

type BalancerProxyConfig struct {
	Host string `env:"HOST" envDefault:"localhost"`
	Port string `env:"PORT" envDefault:"8080"`

	BalancerURL string `env:"BALANCER_URL" envDefault:"http://balancer:8080"`

	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`

	MaxAttempts      int32 `env:"MAX_ATTEMPTS" envDefault:"3"`
	BlockDurationMin int32 `env:"BLOCK_DURATION_MINUTES" envDefault:"5"`

	MaxTimeoutAttempts int32 `env:"MAX_TIMEOUT_ATTEMPTS" envDefault:"2"`

	MinOverlapPct int32 `env:"MIN_OVERLAP_PCT" envDefault:"70"`

	JaegerEndpoint string `env:"JAEGER_ENDPOINT" envDefault:""`
}

func LoadBalancerProxyConfig() (*BalancerProxyConfig, error) {
	config := &BalancerProxyConfig{}

	if err := env.Parse(config); err != nil {
		return nil, err
	}

	return config, nil
}
