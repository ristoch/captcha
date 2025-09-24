package entity

import (
	"errors"
	"time"
)

// BinaryEventType represents the type of binary event
type BinaryEventType int

const (
	EventTypeSliderMoved BinaryEventType = iota
	EventTypeClick
	EventTypeDragStart
	EventTypeDragEnd
	EventTypeInteractionStarted
	EventTypeChallengeCompleted
	EventTypeChallengeFailed
)

var ErrInvalidBinaryData = errors.New("invalid binary data")

// BinaryEvent represents a packed binary event
type BinaryEvent struct {
	Type      BinaryEventType
	X         int32
	Y         int32
	Timestamp int64
	Data      []byte
}

// PackSliderEvent packs slider position and timestamp into binary data
// Coordinates are limited to 8192x8192 as per requirements
func PackSliderEvent(position int32, timestamp int64) []byte {
	if position < 0 || position > 8191 {
		return nil
	}

	// Pack: 13 bits for position (0-8191), 38 bits for timestamp
	packed := (uint64(position) << 38) | (uint64(timestamp) & 0x3FFFFFFFFF)

	result := make([]byte, 7)
	for i := 0; i < 7; i++ {
		result[6-i] = byte((packed >> (i * 8)) & 0xFF)
	}

	return result
}

// PackClickEvent packs click coordinates and timestamp into binary data
func PackClickEvent(x, y int32, timestamp int64) []byte {
	if x < 0 || x > 8191 || y < 0 || y > 8191 {
		return nil
	}

	// Pack: 13 bits for x, 13 bits for y, 38 bits for timestamp
	packed := (uint64(x) << 51) | (uint64(y) << 38) | (uint64(timestamp) & 0x3FFFFFFFFF)

	result := make([]byte, 8)
	for i := 0; i < 8; i++ {
		result[7-i] = byte((packed >> (i * 8)) & 0xFF)
	}

	return result
}

// UnpackSliderEvent unpacks binary data to slider position and timestamp
func UnpackSliderEvent(data []byte) (position int32, timestamp int64, err error) {
	if len(data) != 7 {
		return 0, 0, ErrInvalidBinaryData
	}

	var packed uint64
	for i := 0; i < 7; i++ {
		packed |= uint64(data[6-i]) << (i * 8)
	}

	position = int32((packed >> 38) & 0x1FFF)
	timestamp = int64(packed & 0x3FFFFFFFFF)

	return position, timestamp, nil
}

// UnpackClickEvent unpacks binary data to click coordinates and timestamp
func UnpackClickEvent(data []byte) (x, y int32, timestamp int64, err error) {
	if len(data) != 8 {
		return 0, 0, 0, ErrInvalidBinaryData
	}

	var packed uint64
	for i := 0; i < 8; i++ {
		packed |= uint64(data[7-i]) << (i * 8)
	}

	x = int32((packed >> 51) & 0x1FFF)
	y = int32((packed >> 38) & 0x1FFF)
	timestamp = int64(packed & 0x3FFFFFFFFF)

	return x, y, timestamp, nil
}

// EventData represents event data for processing
type EventData struct {
	ChallengeID string                 `json:"challenge_id"`
	UserID      string                 `json:"user_id"`
	EventType   string                 `json:"event_type"`
	Data        map[string]interface{} `json:"data"`
	Timestamp   int64                  `json:"timestamp"`
}

// NewEventData creates a new event data
func NewEventData(challengeID, userID, eventType string, data map[string]interface{}) *EventData {
	return &EventData{
		ChallengeID: challengeID,
		UserID:      userID,
		EventType:   eventType,
		Data:        data,
		Timestamp:   time.Now().UnixMilli(),
	}
}

// EventResult represents the result of event processing
type EventResult struct {
	Success   bool                   `json:"success"`
	Message   string                 `json:"message,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Timestamp int64                  `json:"timestamp"`
}

// NewEventResult creates a new event result
func NewEventResult(success bool, message string, data map[string]interface{}) *EventResult {
	return &EventResult{
		Success:   success,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UnixMilli(),
	}
}
