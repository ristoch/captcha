package event_processing

import (
	"context"
	"log"
)

type EventPublisherService struct {
	subscribers []chan interface{}
}

func NewEventPublisherService() *EventPublisherService {
	return &EventPublisherService{
		subscribers: make([]chan interface{}, 0),
	}
}

func (e *EventPublisherService) Publish(ctx context.Context, event interface{}) error {
	log.Printf("Publishing event: %+v", event)

	for _, subscriber := range e.subscribers {
		select {
		case subscriber <- event:
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}

	return nil
}

func (e *EventPublisherService) Subscribe() chan interface{} {
	subscriber := make(chan interface{}, 100)
	e.subscribers = append(e.subscribers, subscriber)
	return subscriber
}
