package entity

import (
	"errors"
	"time"
)

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

type BinaryEvent struct {
	Type      BinaryEventType
	X         int32
	Y         int32
	Timestamp int64
	Data      []byte
}

func PackSliderEvent(position int32, timestamp int64) []byte {
	if position < 0 || position > 8191 {
		return nil
	}

	packed := (uint64(position) << 38) | (uint64(timestamp) & 0x3FFFFFFFFF)

	result := make([]byte, 7)
	for i := 0; i < 7; i++ {
		result[6-i] = byte((packed >> (i * 8)) & 0xFF)
	}

	return result
}

func PackClickEvent(x, y int32, timestamp int64) []byte {
	if x < 0 || x > 8191 || y < 0 || y > 8191 {
		return nil
	}

	packed := (uint64(x) << 51) | (uint64(y) << 38) | (uint64(timestamp) & 0x3FFFFFFFFF)

	result := make([]byte, 8)
	for i := 0; i < 8; i++ {
		result[7-i] = byte((packed >> (i * 8)) & 0xFF)
	}

	return result
}

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

type EventData struct {
	ChallengeID string                 `json:"challenge_id"`
	UserID      string                 `json:"user_id"`
	EventType   string                 `json:"event_type"`
	Data        map[string]interface{} `json:"data"`
	Timestamp   int64                  `json:"timestamp"`
}

func NewEventData(challengeID, userID, eventType string, data map[string]interface{}) *EventData {
	return &EventData{
		ChallengeID: challengeID,
		UserID:      userID,
		EventType:   eventType,
		Data:        data,
		Timestamp:   time.Now().UnixMilli(),
	}
}

type EventResult struct {
	Success   bool                   `json:"success"`
	Message   string                 `json:"message,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Timestamp int64                  `json:"timestamp"`
}

func NewEventResult(success bool, message string, data map[string]interface{}) *EventResult {
	return &EventResult{
		Success:   success,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UnixMilli(),
	}
}
