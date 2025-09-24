package grpc

import (
	"context"
	"encoding/json"
	"log"

	captchaProto "captcha-service/gen/proto/proto/captcha"


	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EventStreamHandler struct {
	captchaService *service.CaptchaService
}

func NewEventStreamHandler(captchaService *service.CaptchaService) *EventStreamHandler {
	return &EventStreamHandler{
		captchaService: captchaService,
	}
}

func (h *EventStreamHandler) MakeEventStream(stream captchaProto.CaptchaService_MakeEventStreamServer) error {
	log.Println("Event stream started")
	defer log.Println("Event stream ended")

	for {
		// Получаем событие от клиента
		clientEvent, err := stream.Recv()
		if err != nil {
			log.Printf("Error receiving client event: %v", err)
			return err
		}

		log.Printf("Received client event: type=%v, challenge_id=%s, user_id=%s",
			clientEvent.EventType, clientEvent.ChallengeId, clientEvent.UserId)

		// Обрабатываем событие в зависимости от типа
		switch clientEvent.EventType {
		case captchaProto.ClientEvent_FRONTEND_EVENT:
			err = h.handleFrontendEvent(stream, clientEvent)
		case captchaProto.ClientEvent_CONNECTION_CLOSED:
			log.Printf("Connection closed for challenge: %s", clientEvent.ChallengeId)
			return nil
		case captchaProto.ClientEvent_BALANCER_EVENT:
			err = h.handleBalancerEvent(stream, clientEvent)
		default:
			log.Printf("Unknown event type: %v", clientEvent.EventType)
			continue
		}

		if err != nil {
			log.Printf("Error handling event: %v", err)
			return err
		}
	}
}

func (h *EventStreamHandler) handleFrontendEvent(stream captchaProto.CaptchaService_MakeEventStreamServer, event *captchaProto.ClientEvent) error {
	// Парсим данные события
	var eventData map[string]interface{}
	if err := json.Unmarshal(event.Data, &eventData); err != nil {
		log.Printf("Error unmarshaling event data: %v", err)
		return status.Errorf(codes.InvalidArgument, "invalid event data")
	}

	eventType, ok := eventData["eventType"].(string)
	if !ok {
		log.Printf("Missing eventType in event data")
		return status.Errorf(codes.InvalidArgument, "missing eventType")
	}

	switch eventType {
	case "slider_move":
		return h.handleSliderMove(stream, event, eventData)
	case "validation":
		return h.handleValidation(stream, event, eventData)
	default:
		log.Printf("Unknown frontend event type: %s", eventType)
		return nil
	}
}

func (h *EventStreamHandler) handleSliderMove(stream captchaProto.CaptchaService_MakeEventStreamServer, event *captchaProto.ClientEvent, eventData map[string]interface{}) error {
	// Логируем движение слайдера
	log.Printf("Slider move for challenge %s: %+v", event.ChallengeId, eventData["data"])

	// Можно отправить обратную связь клиенту
	response := &captchaProto.ServerEvent{
		Event: &captchaProto.ServerEvent_ClientData{
			ClientData: &captchaProto.ServerEvent_SendClientData{
				ChallengeId: event.ChallengeId,
				Data:        []byte(`{"type":"feedback","message":"slider_moved"}`),
			},
		},
	}

	return stream.Send(response)
}

func (h *EventStreamHandler) handleValidation(stream captchaProto.CaptchaService_MakeEventStreamServer, event *captchaProto.ClientEvent, eventData map[string]interface{}) error {
	// Получаем ответ от клиента
	answerData, ok := eventData["data"].(map[string]interface{})
	if !ok {
		return status.Errorf(codes.InvalidArgument, "invalid answer data")
	}

	// Валидируем challenge
	valid, confidence, err := h.captchaService.ValidateChallenge(context.Background(), event.ChallengeId, answerData)
	if err != nil {
		log.Printf("Error validating challenge: %v", err)
		return status.Errorf(codes.Internal, "validation error")
	}

	// Отправляем результат валидации
	response := &captchaProto.ServerEvent{
		Event: &captchaProto.ServerEvent_Result{
			Result: &captchaProto.ServerEvent_ChallengeResult{
				ChallengeId:       event.ChallengeId,
				ConfidencePercent: int32(confidence),
			},
		},
	}

	if err := stream.Send(response); err != nil {
		return err
	}

	// Отправляем дополнительную информацию клиенту
	clientData := map[string]interface{}{
		"valid":      valid,
		"confidence": confidence,
		"message":    "validation_complete",
	}

	clientDataBytes, _ := json.Marshal(clientData)

	clientResponse := &captchaProto.ServerEvent{
		Event: &captchaProto.ServerEvent_ClientData{
			ClientData: &captchaProto.ServerEvent_SendClientData{
				ChallengeId: event.ChallengeId,
				Data:        clientDataBytes,
			},
		},
	}

	return stream.Send(clientResponse)
}

func (h *EventStreamHandler) handleBalancerEvent(stream captchaProto.CaptchaService_MakeEventStreamServer, event *captchaProto.ClientEvent) error {
	// Обрабатываем события от балансера
	log.Printf("Balancer event for challenge %s: %+v", event.ChallengeId, event.Data)

	// Здесь можно добавить логику обработки событий от балансера
	// Например, принудительное завершение challenge, обновление конфигурации и т.д.

	return nil
}
