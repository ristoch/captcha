package config

import (
	"github.com/caarlos0/env/v11"
)

// BalancerProxyConfig конфигурация для balancer-proxy
type BalancerProxyConfig struct {
	// Server settings
	Host string `env:"HOST" envDefault:"localhost"`
	Port string `env:"PORT" envDefault:"8080"`

	// Balancer settings
	BalancerURL string `env:"BALANCER_URL" envDefault:"http://balancer:8080"`

	// Logging
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`

	// User blocking settings
	MaxAttempts      int32 `env:"MAX_ATTEMPTS" envDefault:"3"`
	BlockDurationMin int32 `env:"BLOCK_DURATION_MINUTES" envDefault:"5"`

	// Timeout settings
	MaxTimeoutAttempts int32 `env:"MAX_TIMEOUT_ATTEMPTS" envDefault:"2"`

	// Validation settings
	MinOverlapPct int32 `env:"MIN_OVERLAP_PCT" envDefault:"70"`

	// Tracing
	JaegerEndpoint string `env:"JAEGER_ENDPOINT" envDefault:""`
}

// LoadBalancerProxyConfig загружает конфигурацию для balancer-proxy
func LoadBalancerProxyConfig() (*BalancerProxyConfig, error) {
	config := &BalancerProxyConfig{}

	if err := env.Parse(config); err != nil {
		return nil, err
	}

	return config, nil
}
