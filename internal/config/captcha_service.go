package config

import (
	"github.com/caarlos0/env/v11"
)

type CaptchaConfig struct {
	Host string `env:"HOST" envDefault:"localhost"`
	Port string `env:"PORT" envDefault:"8080"`

	BalancerAddress string `env:"BALANCER_ADDRESS" envDefault:""`

	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`

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

	MinTimeMs          int32 `env:"MIN_TIME_MS" envDefault:"1000"`
	MaxTimeMs          int32 `env:"MAX_TIME_MS" envDefault:"30000"`
	MaxTimeoutAttempts int32 `env:"MAX_TIMEOUT_ATTEMPTS" envDefault:"3"`

	MinOverlapPct int32 `env:"MIN_OVERLAP_PCT" envDefault:"20"`

	CleanupInterval int32 `env:"CLEANUP_INTERVAL" envDefault:"300"`
	StaleThreshold  int32 `env:"STALE_THRESHOLD" envDefault:"600"`

	MaxAttempts      int32 `env:"MAX_ATTEMPTS" envDefault:"3"`
	BlockDurationMin int32 `env:"BLOCK_DURATION_MINUTES" envDefault:"5"`

	MinPort int32 `env:"MIN_PORT" envDefault:"38000"`
	MaxPort int32 `env:"MAX_PORT" envDefault:"40000"`

	MaxChallenges      int32  `env:"MAX_CHALLENGES" envDefault:"10000"`
	ShutdownTimeoutSec int32  `env:"SHUTDOWN_TIMEOUT_SEC" envDefault:"30"`
	BalancerAddr       string `env:"BALANCER_ADDR" envDefault:"localhost:9090"`

	DefaultTargetX    int32 `env:"DEFAULT_TARGET_X" envDefault:"200"`
	DefaultTargetY    int32 `env:"DEFAULT_TARGET_Y" envDefault:"150"`
	DefaultTolerance  int32 `env:"DEFAULT_TOLERANCE" envDefault:"10"`
	DefaultConfidence int32 `env:"DEFAULT_CONFIDENCE" envDefault:"85"`
}

func LoadCaptchaServiceConfig() (*CaptchaConfig, error) {
	config := &CaptchaConfig{}

	if err := env.Parse(config); err != nil {
		return nil, err
	}

	return config, nil
}
