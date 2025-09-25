package config

import (
	"github.com/caarlos0/env/v11"
)

type DemoConfig struct {
	Port string `env:"DEMO_PORT" envDefault:"8082"`

	CaptchaServiceURL string `env:"CAPTCHA_SERVICE_URL" envDefault:"http://localhost:8081"`

	MaxAttempts   int32 `env:"MAX_ATTEMPTS" envDefault:"3"`
	BlockDuration int32 `env:"BLOCK_DURATION_MINUTES" envDefault:"5"`

	DefaultComplexity int32 `env:"DEFAULT_COMPLEXITY" envDefault:"50"`

	SessionTimeout string `env:"SESSION_TIMEOUT" envDefault:"24h"`

	MaxSessions int `env:"MAX_SESSIONS" envDefault:"5000"`

	MaxChallenges int `env:"MAX_CHALLENGES" envDefault:"1000"`

	ShutdownTimeoutSeconds int `env:"SHUTDOWN_TIMEOUT_SECONDS" envDefault:"10"`
	CleanupIntervalHours   int `env:"CLEANUP_INTERVAL_HOURS" envDefault:"1"`

	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`

	DefaultTargetX    int32 `env:"DEFAULT_TARGET_X" envDefault:"200"`
	DefaultTargetY    int32 `env:"DEFAULT_TARGET_Y" envDefault:"150"`
	DefaultTolerance  int32 `env:"DEFAULT_TOLERANCE" envDefault:"10"`
	DefaultConfidence int32 `env:"DEFAULT_CONFIDENCE" envDefault:"85"`
}

func LoadDemoConfig() (*DemoConfig, error) {
	config := &DemoConfig{}

	if err := env.Parse(config); err != nil {
		return nil, err
	}

	return config, nil
}
