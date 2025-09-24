package event_processing

import (
	"sync"
)

type EventStreamManagerService struct {
	streams map[string]chan interface{}
	mu      sync.RWMutex
}

func NewEventStreamManagerService() *EventStreamManagerService {
	return &EventStreamManagerService{
		streams: make(map[string]chan interface{}),
	}
}

func (e *EventStreamManagerService) CreateStream(streamID string) chan interface{} {
	e.mu.Lock()
	defer e.mu.Unlock()

	stream := make(chan interface{}, 100)
	e.streams[streamID] = stream
	return stream
}

func (e *EventStreamManagerService) GetStream(streamID string) (chan interface{}, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	stream, exists := e.streams[streamID]
	return stream, exists
}

func (e *EventStreamManagerService) CloseStream(streamID string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if stream, exists := e.streams[streamID]; exists {
		close(stream)
		delete(e.streams, streamID)
	}
}
