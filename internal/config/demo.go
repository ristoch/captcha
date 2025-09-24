package config

import (
	"github.com/caarlos0/env/v11"
)

// DemoConfig конфигурация для demo
type DemoConfig struct {
	// Server settings
	Port string `env:"DEMO_PORT" envDefault:"8082"`

	// Captcha service settings
	CaptchaServiceURL string `env:"CAPTCHA_SERVICE_URL" envDefault:"http://localhost:8081"`

	// User blocking settings
	MaxAttempts   int32 `env:"MAX_ATTEMPTS" envDefault:"3"`
	BlockDuration int32 `env:"BLOCK_DURATION_MINUTES" envDefault:"5"`

	// Session settings
	SessionTimeout string `env:"SESSION_TIMEOUT" envDefault:"24h"`

	// Logging
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
}

// LoadDemoConfig загружает конфигурацию для demo
func LoadDemoConfig() (*DemoConfig, error) {
	config := &DemoConfig{}

	if err := env.Parse(config); err != nil {
		return nil, err
	}

	return config, nil
}
