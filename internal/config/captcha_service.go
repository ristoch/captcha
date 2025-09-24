package config

import (
	"github.com/caarlos0/env/v11"
)

// CaptchaServiceConfig конфигурация для captcha-service
type CaptchaServiceConfig struct {
	// Server settings
	Host string `env:"HOST" envDefault:"localhost"`
	Port string `env:"PORT" envDefault:"8080"`

	// Balancer settings
	BalancerAddress string `env:"BALANCER_ADDRESS" envDefault:""`

	// Logging
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`

	// Challenge settings
	ChallengeType        string `env:"CHALLENGE_TYPE" envDefault:"slider-puzzle"`
	ComplexityLow        int32  `env:"COMPLEXITY_LOW" envDefault:"30"`
	ComplexityMedium     int32  `env:"COMPLEXITY_MEDIUM" envDefault:"50"`
	ComplexityHigh       int32  `env:"COMPLEXITY_HIGH" envDefault:"70"`
	PuzzleSizeLow        int32  `env:"PUZZLE_SIZE_LOW" envDefault:"200"`
	PuzzleSizeMedium     int32  `env:"PUZZLE_SIZE_MEDIUM" envDefault:"300"`
	PuzzleSizeHigh       int32  `env:"PUZZLE_SIZE_HIGH" envDefault:"400"`
	ToleranceLow         int32  `env:"TOLERANCE_LOW" envDefault:"10"`
	ToleranceMedium      int32  `env:"TOLERANCE_MEDIUM" envDefault:"5"`
	ToleranceHigh        int32  `env:"TOLERANCE_HIGH" envDefault:"3"`
	ExpirationTimeLow    int32  `env:"EXPIRATION_TIME_LOW" envDefault:"300"`
	ExpirationTimeMedium int32  `env:"EXPIRATION_TIME_MEDIUM" envDefault:"180"`
	ExpirationTimeHigh   int32  `env:"EXPIRATION_TIME_HIGH" envDefault:"120"`

	// Timing settings
	MinTimeMs          int32 `env:"MIN_TIME_MS" envDefault:"1000"`
	MaxTimeMs          int32 `env:"MAX_TIME_MS" envDefault:"30000"`
	MaxTimeoutAttempts int32 `env:"MAX_TIMEOUT_ATTEMPTS" envDefault:"3"`

	// Validation settings
	MinOverlapPct int32 `env:"MIN_OVERLAP_PCT" envDefault:"20"`

	// Cleanup settings
	CleanupInterval int32 `env:"CLEANUP_INTERVAL" envDefault:"300"`
	StaleThreshold  int32 `env:"STALE_THRESHOLD" envDefault:"600"`

	// User blocking settings
	MaxAttempts      int32 `env:"MAX_ATTEMPTS" envDefault:"3"`
	BlockDurationMin int32 `env:"BLOCK_DURATION_MINUTES" envDefault:"5"`
}

// LoadCaptchaServiceConfig загружает конфигурацию для captcha-service
func LoadCaptchaServiceConfig() (*CaptchaServiceConfig, error) {
	config := &CaptchaServiceConfig{}

	if err := env.Parse(config); err != nil {
		return nil, err
	}

	return config, nil
}
