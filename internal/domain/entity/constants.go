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
	CanvasWidth  = 400
	CanvasHeight = 300
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

var ErrWebSocketNotConnected = errors.New("websocket not connected")
var ErrUserBlocked = errors.New("user blocked")

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
