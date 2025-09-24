package event_processing

import (
	"context"
	"log"

	"captcha-service/internal/domain/entity"
	"captcha-service/internal/repository"
)

type EventProcessorService struct {
	repository     repository.CaptchaRepository
	eventPublisher *EventPublisherService
}

func NewEventProcessorService(repo repository.CaptchaRepository, eventPublisher *EventPublisherService) *EventProcessorService {
	return &EventProcessorService{
		repository:     repo,
		eventPublisher: eventPublisher,
	}
}

func (e *EventProcessorService) ProcessEvent(ctx context.Context, event *entity.BinaryEvent) error {
	log.Printf("Processing event: %+v", event)
	return nil
}
