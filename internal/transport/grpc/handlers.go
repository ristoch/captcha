package grpc

import (
	"context"
	"encoding/json"
	"strconv"

	captchav1 "captcha-service/gen/proto/proto/captcha"
	"captcha-service/internal/domain/entity"
	"captcha-service/internal/service"
	"captcha-service/pkg/logger"

	"go.uber.org/zap"
)

type Handlers struct {
	captchav1.UnimplementedCaptchaServiceServer
	captchaService     *service.CaptchaService
	eventStreamHandler *EventStreamHandler
}

func NewHandlers(captchaService *service.CaptchaService) *Handlers {
	return &Handlers{
		captchaService:     captchaService,
		eventStreamHandler: NewEventStreamHandler(captchaService),
	}
}

func (h *Handlers) NewChallenge(ctx context.Context, req *captchav1.ChallengeRequest) (*captchav1.ChallengeResponse, error) {
	challenge, err := h.captchaService.CreateChallenge(ctx, "slider-puzzle", req.Complexity, req.UserId)
	if err != nil {
		logger.Error("Failed to create challenge", zap.Error(err))
		return nil, err
	}

	// Генерируем HTML для challenge
	html := h.generateChallengeHTML(challenge)
	if html == "" {
		logger.Error("Failed to generate HTML for challenge")
		return nil, err
	}

	return &captchav1.ChallengeResponse{
		ChallengeId: challenge.ID,
		Html:        html,
	}, nil
}

func (h *Handlers) generateChallengeHTML(challenge *entity.Challenge) string {
	// Здесь должна быть генерация HTML на основе challenge.Data
	// Пока возвращаем простой HTML
	return `<div class="captcha-container">
		<h3>Slider Puzzle Captcha</h3>
		<p>Challenge ID: ` + challenge.ID + `</p>
		<p>Complexity: ` + strconv.Itoa(int(challenge.Complexity)) + `</p>
		<div class="slider-area">
			<input type="range" id="xSlider" min="0" max="380" value="0">
			<input type="range" id="ySlider" min="0" max="180" value="0">
		</div>
		<button id="validateBtn">Validate</button>
	</div>`
}

func (h *Handlers) ValidateChallenge(ctx context.Context, req *captchav1.ValidateRequest) (*captchav1.ValidateResponse, error) {
	var answer interface{}
	if err := json.Unmarshal([]byte(req.Answer), &answer); err != nil {
		logger.Error("Failed to parse answer", zap.Error(err))
		return nil, err
	}

	valid, confidence, err := h.captchaService.ValidateChallenge(ctx, req.ChallengeId, answer)
	if err != nil {
		logger.Error("Failed to validate challenge", zap.Error(err))
		return nil, err
	}

	return &captchav1.ValidateResponse{
		Valid:      valid,
		Confidence: confidence,
	}, nil
}

func (h *Handlers) MakeEventStream(stream captchav1.CaptchaService_MakeEventStreamServer) error {
	return h.eventStreamHandler.MakeEventStream(stream)
}
