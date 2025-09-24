package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"captcha-service/internal/domain/entity"
	"captcha-service/internal/domain/interfaces"
	"captcha-service/pkg/logger"

	"go.uber.org/zap"
)

// EventProcessingUseCase handles event processing business logic
type EventProcessingUseCase struct {
	eventProcessor    interfaces.EventProcessor
	eventStreamMgr    interfaces.EventStreamManager
	eventPublisher    interfaces.EventPublisher
	challengeRepo     interfaces.ChallengeRepository
	generatorRegistry interfaces.GeneratorRegistry
}

// NewEventProcessingUseCase creates a new event processing use case
func NewEventProcessingUseCase(
	eventProcessor interfaces.EventProcessor,
	eventStreamMgr interfaces.EventStreamManager,
	eventPublisher interfaces.EventPublisher,
	challengeRepo interfaces.ChallengeRepository,
	generatorRegistry interfaces.GeneratorRegistry,
) *EventProcessingUseCase {
	return &EventProcessingUseCase{
		eventProcessor:    eventProcessor,
		eventStreamMgr:    eventStreamMgr,
		eventPublisher:    eventPublisher,
		challengeRepo:     challengeRepo,
		generatorRegistry: generatorRegistry,
	}
}

// ProcessBinaryEvent processes a binary event from the client
func (uc *EventProcessingUseCase) ProcessBinaryEvent(ctx context.Context, challengeID string, eventType entity.BinaryEventType, data []byte) (*entity.EventResult, error) {
	logger.Debug("Processing binary event",
		zap.String("challenge_id", challengeID),
		zap.Int("event_type", int(eventType)),
		zap.Int("data_length", len(data)))

	// Validate challenge exists and is not expired
	challenge, err := uc.challengeRepo.GetChallenge(ctx, challengeID)
	if err != nil {
		return entity.NewEventResult(false, "Challenge not found", nil), err
	}

	if challenge.ExpiresAt.Before(time.Now()) {
		return entity.NewEventResult(false, "Challenge expired", nil), fmt.Errorf("challenge expired")
	}

	// Create binary event
	binaryEvent := &entity.BinaryEvent{
		Type:      eventType,
		Timestamp: time.Now().UnixNano(),
	}

	// Process based on event type
	switch eventType {
	case entity.EventTypeSliderMoved:
		position, _, err := entity.UnpackSliderEvent(data)
		if err != nil {
			return entity.NewEventResult(false, "Invalid slider event data", nil), err
		}
		binaryEvent.X = position
		binaryEvent.Y = 0
		return uc.eventProcessor.ProcessEvent(binaryEvent)

	case entity.EventTypeClick:
		x, y, _, err := entity.UnpackClickEvent(data)
		if err != nil {
			return entity.NewEventResult(false, "Invalid click event data", nil), err
		}
		binaryEvent.X = x
		binaryEvent.Y = y
		return uc.eventProcessor.ProcessEvent(binaryEvent)

	case entity.EventTypeChallengeCompleted:
		var eventData map[string]interface{}
		if err := json.Unmarshal(data, &eventData); err != nil {
			return entity.NewEventResult(false, "Invalid completion data", nil), err
		}
		return uc.eventProcessor.ProcessEvent(binaryEvent)

	case entity.EventTypeChallengeFailed:
		var eventData map[string]interface{}
		if err := json.Unmarshal(data, &eventData); err != nil {
			return entity.NewEventResult(false, "Invalid failure data", nil), err
		}
		return uc.eventProcessor.ProcessEvent(binaryEvent)

	default:
		return entity.NewEventResult(false, "Unknown event type", nil), fmt.Errorf("unknown event type: %d", eventType)
	}
}

// ProcessJSONEvent processes a JSON event from the client
func (uc *EventProcessingUseCase) ProcessJSONEvent(ctx context.Context, eventData *entity.WebSocketMessage) (*entity.EventResult, error) {
	logger.Debug("Processing JSON event",
		zap.String("user_id", eventData.UserID),
		zap.String("type", eventData.Type))

	// Extract challenge ID from data
	challengeID, ok := eventData.Data["challenge_id"].(string)
	if !ok {
		return entity.NewEventResult(false, "Challenge ID not found", nil), fmt.Errorf("challenge ID not found")
	}

	// Validate challenge exists and is not expired
	challenge, err := uc.challengeRepo.GetChallenge(ctx, challengeID)
	if err != nil {
		return entity.NewEventResult(false, "Challenge not found", nil), err
	}

	if challenge.ExpiresAt.Before(time.Now()) {
		return entity.NewEventResult(false, "Challenge expired", nil), fmt.Errorf("challenge expired")
	}

	// Create binary event
	binaryEvent := &entity.BinaryEvent{
		Timestamp: time.Now().UnixNano(),
	}

	// Process based on event type
	switch eventData.Type {
	case "slider_moved":
		position, ok := eventData.Data["position"].(float64)
		if !ok {
			return entity.NewEventResult(false, "Invalid position data", nil), fmt.Errorf("invalid position data")
		}
		binaryEvent.Type = entity.EventTypeSliderMoved
		binaryEvent.X = int32(position)
		binaryEvent.Y = 0
		return uc.eventProcessor.ProcessEvent(binaryEvent)

	case "click":
		x, xOk := eventData.Data["x"].(float64)
		y, yOk := eventData.Data["y"].(float64)
		if !xOk || !yOk {
			return entity.NewEventResult(false, "Invalid click coordinates", nil), fmt.Errorf("invalid click coordinates")
		}
		binaryEvent.Type = entity.EventTypeClick
		binaryEvent.X = int32(x)
		binaryEvent.Y = int32(y)
		return uc.eventProcessor.ProcessEvent(binaryEvent)

	case "challenge_completed":
		binaryEvent.Type = entity.EventTypeChallengeCompleted
		return uc.eventProcessor.ProcessEvent(binaryEvent)

	case "challenge_failed":
		binaryEvent.Type = entity.EventTypeChallengeFailed
		return uc.eventProcessor.ProcessEvent(binaryEvent)

	default:
		return entity.NewEventResult(false, "Unknown event type", nil), fmt.Errorf("unknown event type: %s", eventData.Type)
	}
}

// CreateEventStream creates a new event stream for a challenge
func (uc *EventProcessingUseCase) CreateEventStream(ctx context.Context, challengeID string) (interfaces.EventStream, error) {
	logger.Debug("Creating event stream", zap.String("challenge_id", challengeID))

	// Validate challenge exists
	_, err := uc.challengeRepo.GetChallenge(ctx, challengeID)
	if err != nil {
		return nil, fmt.Errorf("challenge not found: %w", err)
	}

	// Create stream
	stream, err := uc.eventStreamMgr.CreateStream(challengeID)
	if err != nil {
		return nil, fmt.Errorf("failed to create stream: %w", err)
	}

	return stream, nil
}

// PublishEvent publishes an event
func (uc *EventProcessingUseCase) PublishEvent(ctx context.Context, event *entity.BinaryEvent) error {
	logger.Debug("Publishing event",
		zap.Int("event_type", int(event.Type)),
		zap.Int32("x", event.X),
		zap.Int32("y", event.Y))

	return uc.eventPublisher.Publish(event)
}
