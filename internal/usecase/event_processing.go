package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"captcha-service/internal/domain/entity"
	"captcha-service/pkg/logger"

	"go.uber.org/zap"
)

type EventProcessor interface {
	ProcessEvent(event *entity.BinaryEvent) (*entity.EventResult, error)
}

type EventStreamManager interface {
	CreateStream(userID string) (EventStream, error)
	CloseStream(userID string) error
	GetStream(userID string) (EventStream, error)
}

type EventPublisher interface {
	Publish(event *entity.BinaryEvent) error
	Subscribe(userID string, handler EventHandler) error
	Unsubscribe(userID string) error
}

type EventStream interface {
	Send(event *entity.BinaryEvent) error
	Receive() (*entity.BinaryEvent, error)
	Close() error
}

type EventHandler interface {
	Handle(event *entity.BinaryEvent) error
}

type ChallengeRepository interface {
	SaveChallenge(ctx context.Context, challenge *entity.Challenge) error
	GetChallenge(ctx context.Context, challengeID string) (*entity.Challenge, error)
	DeleteChallenge(ctx context.Context, challengeID string) error
}

type GeneratorRegistry interface {
	Get(challengeType string) (ChallengeGenerator, bool)
	Register(challengeType string, generator ChallengeGenerator)
}

type ChallengeGenerator interface {
	Generate(ctx context.Context, complexity int32, userID string) (*entity.Challenge, error)
	Validate(answer interface{}, data interface{}) (bool, int32, error)
}

type EventProcessingUseCase struct {
	eventProcessor    EventProcessor
	eventStreamMgr    EventStreamManager
	eventPublisher    EventPublisher
	challengeRepo     ChallengeRepository
	generatorRegistry GeneratorRegistry
}

func NewEventProcessingUseCase(
	eventProcessor EventProcessor,
	eventStreamMgr EventStreamManager,
	eventPublisher EventPublisher,
	challengeRepo ChallengeRepository,
	generatorRegistry GeneratorRegistry,
) *EventProcessingUseCase {
	return &EventProcessingUseCase{
		eventProcessor:    eventProcessor,
		eventStreamMgr:    eventStreamMgr,
		eventPublisher:    eventPublisher,
		challengeRepo:     challengeRepo,
		generatorRegistry: generatorRegistry,
	}
}

func (uc *EventProcessingUseCase) ProcessBinaryEvent(ctx context.Context, challengeID string, eventType entity.BinaryEventType, data []byte) (*entity.EventResult, error) {
	logger.Debug("Processing binary event",
		zap.String(entity.FieldChallengeID, challengeID),
		zap.Int("event_type", int(eventType)),
		zap.Int("data_length", len(data)))

	challenge, err := uc.challengeRepo.GetChallenge(ctx, challengeID)
	if err != nil {
		return entity.NewEventResult(false, "Challenge not found", nil), err
	}

	if challenge.ExpiresAt.Before(time.Now()) {
		return entity.NewEventResult(false, "Challenge expired", nil), fmt.Errorf("challenge expired")
	}

	binaryEvent := &entity.BinaryEvent{
		Type:      eventType,
		Timestamp: time.Now().UnixNano(),
	}

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

func (uc *EventProcessingUseCase) ProcessJSONEvent(ctx context.Context, eventData *entity.WebSocketMessage) (*entity.EventResult, error) {
	logger.Debug("Processing JSON event",
		zap.String("user_id", eventData.UserID),
		zap.String("type", eventData.Type))

	challengeID, ok := eventData.Data[entity.FieldChallengeID].(string)
	if !ok {
		return entity.NewEventResult(false, "Challenge ID not found", nil), fmt.Errorf("challenge ID not found")
	}

	challenge, err := uc.challengeRepo.GetChallenge(ctx, challengeID)
	if err != nil {
		return entity.NewEventResult(false, "Challenge not found", nil), err
	}

	if challenge.ExpiresAt.Before(time.Now()) {
		return entity.NewEventResult(false, "Challenge expired", nil), fmt.Errorf("challenge expired")
	}

	binaryEvent := &entity.BinaryEvent{
		Timestamp: time.Now().UnixNano(),
	}

	switch eventData.Type {
	case entity.EventTypeSliderMovedStr:
		position, ok := eventData.Data[entity.FieldPosition].(float64)
		if !ok {
			return entity.NewEventResult(false, "Invalid position data", nil), fmt.Errorf("invalid position data")
		}
		binaryEvent.Type = entity.EventTypeSliderMoved
		binaryEvent.X = int32(position)
		binaryEvent.Y = 0
		return uc.eventProcessor.ProcessEvent(binaryEvent)

	case entity.FieldClick:
		x, xOk := eventData.Data[entity.FieldX].(float64)
		y, yOk := eventData.Data[entity.FieldY].(float64)
		if !xOk || !yOk {
			return entity.NewEventResult(false, "Invalid click coordinates", nil), fmt.Errorf("invalid click coordinates")
		}
		binaryEvent.Type = entity.EventTypeClick
		binaryEvent.X = int32(x)
		binaryEvent.Y = int32(y)
		return uc.eventProcessor.ProcessEvent(binaryEvent)

	case entity.FieldChallengeCompleted:
		binaryEvent.Type = entity.EventTypeChallengeCompleted
		return uc.eventProcessor.ProcessEvent(binaryEvent)

	case entity.FieldChallengeFailed:
		binaryEvent.Type = entity.EventTypeChallengeFailed
		return uc.eventProcessor.ProcessEvent(binaryEvent)

	default:
		return entity.NewEventResult(false, "Unknown event type", nil), fmt.Errorf("unknown event type: %s", eventData.Type)
	}
}

func (uc *EventProcessingUseCase) CreateEventStream(ctx context.Context, challengeID string) (EventStream, error) {
	logger.Debug("Creating event stream", zap.String(entity.FieldChallengeID, challengeID))

	_, err := uc.challengeRepo.GetChallenge(ctx, challengeID)
	if err != nil {
		return nil, fmt.Errorf("challenge not found: %w", err)
	}

	stream, err := uc.eventStreamMgr.CreateStream(challengeID)
	if err != nil {
		return nil, fmt.Errorf("failed to create stream: %w", err)
	}

	return stream, nil
}

func (uc *EventProcessingUseCase) PublishEvent(ctx context.Context, event *entity.BinaryEvent) error {
	logger.Debug("Publishing event",
		zap.Int("event_type", int(event.Type)),
		zap.Int32("x", event.X),
		zap.Int32("y", event.Y))

	return uc.eventPublisher.Publish(event)
}
