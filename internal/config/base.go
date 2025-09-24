package config

// BaseConfig базовые настройки, общие для всех сервисов
type BaseConfig struct {
	// Logging
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`

	// Tracing
	JaegerEndpoint string `env:"JAEGER_ENDPOINT" envDefault:""`

	// User blocking settings
	MaxAttempts      int32 `env:"MAX_ATTEMPTS" envDefault:"3"`
	BlockDurationMin int32 `env:"BLOCK_DURATION_MINUTES" envDefault:"5"`
}

// Common constants
const (
	DefaultHost = "localhost"
	DefaultPort = "8080"
)
