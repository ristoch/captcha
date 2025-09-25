package http

import (
	"context"
	"encoding/json"
	"net/http"

	"captcha-service/internal/domain/entity"
	"captcha-service/pkg/logger"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type CaptchaService interface {
	CreateChallenge(ctx context.Context, challengeType string, complexity int32, userID string) (*entity.Challenge, error)
	ValidateChallenge(ctx context.Context, challengeID string, answer interface{}) (bool, int32, error)
	GetChallenge(ctx context.Context, challengeID string) (*entity.Challenge, error)
}

type Handlers struct {
	captchaService CaptchaService
}

func NewHandlers(captchaService CaptchaService) *Handlers {
	return &Handlers{
		captchaService: captchaService,
	}
}

func (h *Handlers) HandleChallengeRequest(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ChallengeType string `json:"challenge_type"`
		Complexity    int32  `json:"complexity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	challenge, err := h.captchaService.CreateChallenge(r.Context(), req.ChallengeType, req.Complexity, "demo_user")
	if err != nil {
		logger.Error("Failed to create challenge", zap.Error(err))
		http.Error(w, "Failed to create challenge", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(challenge)
}

func (h *Handlers) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		userID = "anonymous"
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("Failed to upgrade connection", zap.Error(err))
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			logger.Error("WebSocket read error", zap.Error(err))
			break
		}

		if err := conn.WriteMessage(messageType, message); err != nil {
			logger.Error("WebSocket write error", zap.Error(err))
			break
		}
	}
}

func (h *Handlers) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"version": "1.0.0",
	})
}
