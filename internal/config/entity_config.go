package config

// EntityConfig конфигурация для entity.Config с тегами env
type EntityConfig struct {
	// User blocking settings
	MaxAttempts      int32 `env:"MAX_ATTEMPTS" envDefault:"3"`
	BlockDurationMin int32 `env:"BLOCK_DURATION_MINUTES" envDefault:"5"`

	// Complexity settings
	ComplexityLow    int32 `env:"COMPLEXITY_LOW" envDefault:"30"`
	ComplexityMedium int32 `env:"COMPLEXITY_MEDIUM" envDefault:"50"`
	ComplexityHigh   int32 `env:"COMPLEXITY_HIGH" envDefault:"70"`

	// Puzzle size settings
	PuzzleSizeLow    int32 `env:"PUZZLE_SIZE_LOW" envDefault:"200"`
	PuzzleSizeMedium int32 `env:"PUZZLE_SIZE_MEDIUM" envDefault:"300"`
	PuzzleSizeHigh   int32 `env:"PUZZLE_SIZE_HIGH" envDefault:"400"`

	// Tolerance settings
	ToleranceLow    int32 `env:"TOLERANCE_LOW" envDefault:"10"`
	ToleranceMedium int32 `env:"TOLERANCE_MEDIUM" envDefault:"5"`
	ToleranceHigh   int32 `env:"TOLERANCE_HIGH" envDefault:"3"`

	// Expiration time settings
	ExpirationTimeLow    int32 `env:"EXPIRATION_TIME_LOW" envDefault:"300"`
	ExpirationTimeMedium int32 `env:"EXPIRATION_TIME_MEDIUM" envDefault:"180"`
	ExpirationTimeHigh   int32 `env:"EXPIRATION_TIME_HIGH" envDefault:"120"`

	// Time settings
	MinTimeMs          int32 `env:"MIN_TIME_MS" envDefault:"1000"`
	MaxTimeMs          int32 `env:"MAX_TIME_MS" envDefault:"30000"`
	MaxTimeoutAttempts int32 `env:"MAX_TIMEOUT_ATTEMPTS" envDefault:"3"`

	// Other settings
	MinOverlapPct   int32 `env:"MIN_OVERLAP_PCT" envDefault:"20"`
	CleanupInterval int32 `env:"CLEANUP_INTERVAL" envDefault:"300"`
	StaleThreshold  int32 `env:"STALE_THRESHOLD" envDefault:"600"`
}
