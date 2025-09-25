package entity

import (
	"errors"
	"time"
)

const (
	MessageTypeChallengeRequest   = "challenge_request"
	MessageTypeChallengeCreated   = "challenge_created"
	MessageTypeChallengeCompleted = "challenge_completed"
	MessageTypeCaptchaEvent       = "captcha_event"
	MessageTypeNewChallenge       = "new_challenge"
	MessageTypeGRPCResponse       = "gRPC_response"
	MessageTypeError              = "error"
	MessageTypeBlocked            = "blocked"
)

const (
	EventTypeCaptchaFailed = "captchaFailed"
)

const (
	FieldHTML           = "html"
	FieldTimestamp      = "timestamp"
	FieldIsCorrect      = "isCorrect"
	FieldDistance       = "distance"
	FieldTargetPosition = "targetPosition"
)

const (
	DefaultTimeoutSeconds      = 5
	DefaultWebSocketTimeout    = 30 * time.Second
	DefaultDiscoveryInterval   = 5 * time.Second
	DefaultCleanupInterval     = 10 * time.Second
	DefaultMaxShutdownInterval = 10 * time.Minute
)

const (
	WebSocketBufferSize = 100
	MinDataLength       = 10
)

const (
	StatusOK            = 200
	StatusBadRequest    = 400
	StatusUnauthorized  = 401
	StatusForbidden     = 403
	StatusNotFound      = 404
	StatusInternalError = 500
)

const (
	ContentTypeHTML = "text/html; charset=utf-8"
	ContentTypeJSON = "application/json"
)

const (
	DefaultUserID = "demo_user"
)

const (
	ErrorWebSocketNotConnected = "WebSocket not connected"
	ErrorUnknownError          = "Unknown error"
	ErrorTooManyAttempts       = "You are blocked due to too many attempts"
	ErrorFailedToConnect       = "Failed to connect to WebSocket"
)

const (
	MessageChallengeCreated     = "Challenge created successfully"
	MessageNewChallengeReceived = "New challenge received!"
	MessageWebSocketConnected   = "WebSocket connected"
	MessageConnecting           = "Connecting to WebSocket..."
)

const (
	GzipMagicByte1 = 0x1f
	GzipMagicByte2 = 0x8b
)

const (
	SliderEventDataLength = 7
	DragEventDataLength   = 8
)

const (
	HealthCheckVersion = "1.0.0"
	HealthCheckStatus  = "healthy"
)

const (
	ChallengeTypeSliderPuzzle = "slider-puzzle"
	ChallengeTypeDragDrop     = "drag-drop"
)

const (
	FieldChallengeID        = "challenge_id"
	FieldUserID             = "user_id"
	FieldSessionID          = "session_id"
	FieldEventType          = "event_type"
	FieldComplexity         = "complexity"
	FieldChallengeType      = "challenge_type"
	FieldAnswer             = "answer"
	FieldStatus             = "status"
	FieldInstanceID         = "instance_id"
	FieldReason             = "reason"
	FieldBlockedUntil       = "blocked_until"
	FieldType               = "type"
	FieldPosition           = "position"
	FieldX                  = "x"
	FieldY                  = "y"
	FieldClick              = "click"
	FieldChallengeCompleted = "challenge_completed"
	FieldChallengeFailed    = "challenge_failed"
)

const (
	EventTypeCaptchaSolved   = "captchaSolved"
	EventTypeCaptchaSendData = "captcha:sendData"
)

const (
	EventTypeSliderMove         = "slider_move"
	EventTypeValidation         = "validation"
	EventTypeSliderMovedStr     = "slider_moved"
	EventTypeFieldEventType     = "eventType"
	EventTypeValidationComplete = "validation_complete"
)

const (
	BalancerEventTypeUserBlocked    = "user_blocked"
	BalancerEventTypeUserUnblocked  = "user_unblocked"
	BalancerEventTypeInstanceStatus = "instance_status"
)

const (
	ChallengeTypeSliderPuzzleReg = "slider-puzzle"
	ChallengeTypeDragDropReg     = "drag-drop"
)

const (
	CanvasWidth          = 400
	CanvasHeight         = 300
	DragDropCanvasWidth  = 400
	DragDropCanvasHeight = 300
)

const (
	PuzzleGapTop = 50
)

var BackgroundImages = []string{
	"background1.png",
	"background2.png",
	"background3.png",
	"background4.png",
}

var PuzzleShapes = []string{
	"circle",
	"square",
	"triangle",
}

const (
	SignalBufferSize = 1
)

var ErrWebSocketNotConnected = errors.New("websocket not connected")
var ErrUserBlocked = errors.New("user blocked")

type Config struct {
	Host                 string `env:"HOST" envDefault:"localhost"`
	Port                 string `env:"PORT" envDefault:"8080"`
	BalancerAddress      string `env:"BALANCER_ADDRESS" envDefault:""`
	LogLevel             string `env:"LOG_LEVEL" envDefault:"info"`
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
	MinTimeMs            int32  `env:"MIN_TIME_MS" envDefault:"1000"`
	MaxTimeMs            int32  `env:"MAX_TIME_MS" envDefault:"30000"`
	MaxTimeoutAttempts   int32  `env:"MAX_TIMEOUT_ATTEMPTS" envDefault:"3"`
	MinOverlapPct        int32  `env:"MIN_OVERLAP_PCT" envDefault:"20"`
	CleanupInterval      int32  `env:"CLEANUP_INTERVAL" envDefault:"300"`
	StaleThreshold       int32  `env:"STALE_THRESHOLD" envDefault:"600"`
	MaxAttempts          int32  `env:"MAX_ATTEMPTS" envDefault:"3"`
	BlockDurationMin     int32  `env:"BLOCK_DURATION_MINUTES" envDefault:"5"`
}

type DemoConfig struct {
	Port              string `env:"DEMO_PORT" envDefault:"8082"`
	CaptchaServiceURL string `env:"CAPTCHA_SERVICE_URL" envDefault:"http://localhost:8081"`
	MaxAttempts       int32  `env:"MAX_ATTEMPTS" envDefault:"3"`
	BlockDuration     int32  `env:"BLOCK_DURATION_MINUTES" envDefault:"5"`
}

type Instance struct {
	ID           string    `json:"id"`
	Type         string    `json:"type"`
	Host         string    `json:"host"`
	Port         int32     `json:"port"`
	Status       string    `json:"status"`
	LastSeen     time.Time `json:"last_seen"`
	RegisteredAt time.Time `json:"registered_at"`
}

type BlockedUser struct {
	UserID       string    `json:"user_id"`
	BlockedUntil time.Time `json:"blocked_until"`
	Reason       string    `json:"reason"`
	Attempts     int32     `json:"attempts"`
	LastAttempt  time.Time `json:"last_attempt"`
}

type RegisterInstanceRequest struct {
	EventType     string `json:"event_type"`
	InstanceID    string `json:"instance_id"`
	ChallengeType string `json:"challenge_type"`
	Host          string `json:"host"`
	PortNumber    int32  `json:"port_number"`
	Timestamp     int64  `json:"timestamp"`
}

type WebSocketMessage struct {
	Type   string                 `json:"type"`
	UserID string                 `json:"user_id,omitempty"`
	Data   map[string]interface{} `json:"data,omitempty"`
}

type UserSession struct {
	UserID       string    `json:"user_id"`
	SessionID    string    `json:"session_id"`
	CreatedAt    time.Time `json:"created_at"`
	LastSeen     time.Time `json:"last_seen"`
	Attempts     int32     `json:"attempts"`
	IsBlocked    bool      `json:"is_blocked"`
	BlockedUntil time.Time `json:"blocked_until"`
}

type DemoData struct {
	UserID      string `json:"user_id"`
	SessionID   string `json:"session_id"`
	ChallengeID string `json:"challenge_id"`
	HTML        string `json:"html"`
}

type SessionRepository interface {
	CreateSession(userID string) (*UserSession, error)
	GetSession(sessionID string) (*UserSession, error)
	GetSessionByUserID(userID string) (*UserSession, error)
	UpdateSession(session *UserSession) error
	DeleteSession(sessionID string) error
	GetAllSessions() ([]*UserSession, error)
	CleanupExpired()
}

var ErrSessionNotFound = errors.New("session not found")
