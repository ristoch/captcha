package grpc

import (
	"context"
	"encoding/json"
	"log"

	captchaProto "captcha-service/gen/proto/proto/captcha"
	"captcha-service/internal/domain/entity"
	"captcha-service/internal/service"

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
		clientEvent, err := stream.Recv()
		if err != nil {
			log.Printf("Error receiving client event: %v", err)
			return err
		}

		log.Printf("Received client event: type=%v, challenge_id=%s, user_id=%s",
			clientEvent.EventType, clientEvent.ChallengeId, clientEvent.UserId)

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
	var eventData map[string]interface{}
	if err := json.Unmarshal(event.Data, &eventData); err != nil {
		log.Printf("Error unmarshaling event data: %v", err)
		return status.Errorf(codes.InvalidArgument, "invalid event data")
	}

	eventType, ok := eventData[entity.EventTypeFieldEventType].(string)
	if !ok {
		log.Printf("Missing eventType in event data")
		return status.Errorf(codes.InvalidArgument, "missing eventType")
	}

	switch eventType {
	case entity.EventTypeSliderMove:
		return h.handleSliderMove(stream, event, eventData)
	case entity.EventTypeValidation:
		return h.handleValidation(stream, event, eventData)
	default:
		log.Printf("Unknown frontend event type: %s", eventType)
		return nil
	}
}

func (h *EventStreamHandler) handleSliderMove(stream captchaProto.CaptchaService_MakeEventStreamServer, event *captchaProto.ClientEvent, eventData map[string]interface{}) error {
	log.Printf("Slider move for challenge %s: %+v", event.ChallengeId, eventData["data"])

	response := &captchaProto.ServerEvent{
		Event: &captchaProto.ServerEvent_ClientData{
			ClientData: &captchaProto.ServerEvent_SendClientData{
				ChallengeId: event.ChallengeId,
				Data:        []byte(`{"type":"feedback","message":"` + entity.EventTypeSliderMovedStr + `"}`),
			},
		},
	}

	return stream.Send(response)
}

func (h *EventStreamHandler) handleValidation(stream captchaProto.CaptchaService_MakeEventStreamServer, event *captchaProto.ClientEvent, eventData map[string]interface{}) error {
	answerData, ok := eventData["data"].(map[string]interface{})
	if !ok {
		return status.Errorf(codes.InvalidArgument, "invalid answer data")
	}

	valid, confidence, err := h.captchaService.ValidateChallenge(context.Background(), event.ChallengeId, answerData)
	if err != nil {
		log.Printf("Error validating challenge: %v", err)
		return status.Errorf(codes.Internal, "validation error")
	}

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

	clientData := map[string]interface{}{
		"valid":      valid,
		"confidence": confidence,
		"message":    entity.EventTypeValidationComplete,
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
	log.Printf("Balancer event for challenge %s: %+v", event.ChallengeId, event.Data)

	// Парсим данные balancer события
	var balancerData map[string]interface{}
	if err := json.Unmarshal(event.Data, &balancerData); err != nil {
		log.Printf("Error unmarshaling balancer event data: %v", err)
		return status.Errorf(codes.InvalidArgument, "invalid balancer event data")
	}

	// Определяем тип balancer собычытиячыччы
	eventType, ok := balancerData[entity.FieldType].(string)
	if !ok {
		log.Printf("Missing event type in balancer event data")
		return status.Errorf(codes.InvalidArgument, "missing event type")
	}

	switch eventType {
	case entity.BalancerEventTypeUserBlocked:
		return h.handleUserBlockedEvent(stream, event, balancerData)
	case entity.BalancerEventTypeUserUnblocked:
		return h.handleUserUnblockedEvent(stream, event, balancerData)
	case entity.BalancerEventTypeInstanceStatus:
		return h.handleInstanceStatusEvent(stream, event, balancerData)
	default:
		log.Printf("Unknown balancer event type: %s", eventType)
		return nil
	}
}

func (h *EventStreamHandler) handleUserBlockedEvent(stream captchaProto.CaptchaService_MakeEventStreamServer, event *captchaProto.ClientEvent, data map[string]interface{}) error {
	userID, ok := data[entity.FieldUserID].(string)
	if !ok {
		return status.Errorf(codes.InvalidArgument, "missing user_id in user_blocked event")
	}

	reason, _ := data[entity.FieldReason].(string)
	blockedUntil, _ := data[entity.FieldBlockedUntil].(float64)

	log.Printf("User %s blocked: %s (until: %v)", userID, reason, blockedUntil)

	response := &captchaProto.ServerEvent{
		Event: &captchaProto.ServerEvent_ClientData{
			ClientData: &captchaProto.ServerEvent_SendClientData{
				ChallengeId: event.ChallengeId,
				Data:        []byte(`{"` + entity.FieldType + `":"user_blocked","` + entity.FieldUserID + `":"` + userID + `","` + entity.FieldReason + `":"` + reason + `"}`),
			},
		},
	}

	return stream.Send(response)
}

func (h *EventStreamHandler) handleUserUnblockedEvent(stream captchaProto.CaptchaService_MakeEventStreamServer, event *captchaProto.ClientEvent, data map[string]interface{}) error {
	userID, ok := data[entity.FieldUserID].(string)
	if !ok {
		return status.Errorf(codes.InvalidArgument, "missing user_id in user_unblocked event")
	}

	log.Printf("User %s unblocked", userID)

	response := &captchaProto.ServerEvent{
		Event: &captchaProto.ServerEvent_ClientData{
			ClientData: &captchaProto.ServerEvent_SendClientData{
				ChallengeId: event.ChallengeId,
				Data:        []byte(`{"` + entity.FieldType + `":"user_unblocked","` + entity.FieldUserID + `":"` + userID + `"}`),
			},
		},
	}

	return stream.Send(response)
}

func (h *EventStreamHandler) handleInstanceStatusEvent(stream captchaProto.CaptchaService_MakeEventStreamServer, event *captchaProto.ClientEvent, data map[string]interface{}) error {
	instanceID, ok := data[entity.FieldInstanceID].(string)
	if !ok {
		return status.Errorf(codes.InvalidArgument, "missing instance_id in instance_status event")
	}

	status, _ := data[entity.FieldStatus].(string)
	challengeType, _ := data[entity.FieldChallengeType].(string)

	log.Printf("Instance %s status changed: %s (type: %s)", instanceID, status, challengeType)

	response := &captchaProto.ServerEvent{
		Event: &captchaProto.ServerEvent_ClientData{
			ClientData: &captchaProto.ServerEvent_SendClientData{
				ChallengeId: event.ChallengeId,
				Data:        []byte(`{"` + entity.FieldType + `":"instance_status","` + entity.FieldInstanceID + `":"` + instanceID + `","` + entity.FieldStatus + `":"` + status + `","` + entity.FieldChallengeType + `":"` + challengeType + `"}`),
			},
		},
	}

	return stream.Send(response)
}
