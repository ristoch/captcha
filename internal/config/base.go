package config

import (
	"github.com/caarlos0/env/v11"
)

// BaseConfig базовые настройки, общие для всех сервисов
type BaseConfig struct {
	DefaultHost string `env:"DEFAULT_HOST" envDefault:"localhost"`
	DefaultPort string `env:"DEFAULT_PORT" envDefault:"8080"`

	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`

	JaegerEndpoint string `env:"JAEGER_ENDPOINT" envDefault:""`

	MaxAttempts      int32 `env:"MAX_ATTEMPTS" envDefault:"3"`
	BlockDurationMin int32 `env:"BLOCK_DURATION_MINUTES" envDefault:"5"`
}

func LoadConfigFromEnv[T any]() (*T, error) {
	config := new(T)
	if err := env.Parse(config); err != nil {
		return nil, err
	}
	return config, nil
}
