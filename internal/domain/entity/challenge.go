package entity

import (
	"captcha-service/internal/domain/dto"
	"errors"
	"fmt"
	"time"
)

type ChallengeData interface {
	GetType() string
}

type SliderPuzzleData struct {
	ChallengeData dto.ChallengeData `json:"challenge"`
	CanvasWidth   int               `json:"canvas_width"`
	CanvasHeight  int               `json:"canvas_height"`
}

func (s SliderPuzzleData) GetType() string {
	return ChallengeTypeSliderPuzzle
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type DragDropChallenge struct {
	UserID         string   `json:"user_id"`
	TargetPosition Position `json:"target_position"`
	ObjectPosition Position `json:"object_position"`
	ObjectSize     Size     `json:"object_size"`
	TargetSize     Size     `json:"target_size"`
	Tolerance      int      `json:"tolerance"`
}

type DragDropData struct {
	ChallengeData DragDropChallenge `json:"challenge"`
	CanvasWidth   int               `json:"canvas_width"`
	CanvasHeight  int               `json:"canvas_height"`
}

func (d DragDropData) GetType() string {
	return ChallengeTypeDragDrop
}

func (c *Challenge) GetSliderPuzzleData() (*SliderPuzzleData, error) {
	data, ok := c.Data.(SliderPuzzleData)
	if !ok {
		return nil, fmt.Errorf("challenge data is not SliderPuzzleData")
	}
	return &data, nil
}

func (c *Challenge) GetDragDropData() (*DragDropData, error) {
	data, ok := c.Data.(DragDropData)
	if !ok {
		return nil, fmt.Errorf("challenge data is not DragDropData")
	}
	return &data, nil
}

type Challenge struct {
	ID                 string
	ChallengeID        string
	UserID             string
	Type               string
	Complexity         int32
	Data               ChallengeData
	Answer             interface{}
	HTML               string
	ExpiresAt          time.Time
	CreatedAt          time.Time
	StartTime          *time.Time
	Attempts           int32
	MaxAttempts        int32
	MinTime            int64
	MaxTime            int64
	IsBlocked          bool
	BlockReason        string
	TimeoutAttempts    int32
	MaxTimeoutAttempts int32
	BlockedUntil       *time.Time
}

var (
	ErrChallengeNotFound  = errors.New("challenge not found")
	ErrChallengeExpired   = errors.New("challenge expired")
	ErrChallengeBlocked   = errors.New("challenge blocked")
	ErrInvalidAnswer      = errors.New("invalid answer")
	ErrMaxAttemptsReached = errors.New("max attempts reached")
)
