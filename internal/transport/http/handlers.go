package http

import (
	"context"
	"encoding/json"
	"net/http"
	"runtime"
	"sync/atomic"
	"time"

	"captcha-service/internal/domain/entity"
	"captcha-service/internal/infrastructure/cache"
	"captcha-service/internal/infrastructure/persistence"
	"captcha-service/internal/service"
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
	memoryMonitor  *MemoryMonitor

	requestsTotal    int64
	challengesTotal  int64
	validationsTotal int64
	errorsTotal      int64
	startTime        time.Time
}

func NewHandlersWithMemoryMonitor(
	captchaService CaptchaService,
	challengeRepo *persistence.MemoryOptimizedRepository,
	sessionCache *cache.SessionCache,
	globalBlocker *service.GlobalUserBlocker,
) *Handlers {
	return &Handlers{
		captchaService: captchaService,
		memoryMonitor:  NewMemoryMonitor(challengeRepo, sessionCache, globalBlocker),
		startTime:      time.Now(),
	}
}

func (h *Handlers) HandleChallengeRequest(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt64(&h.requestsTotal, 1)

	var req struct {
		ChallengeType string `json:"challenge_type"`
		Complexity    int32  `json:"complexity"`
		UserID        string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		atomic.AddInt64(&h.errorsTotal, 1)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	userID := req.UserID
	if userID == "" {
		userID = "demo_user"
	}

	challenge, err := h.captchaService.CreateChallenge(r.Context(), req.ChallengeType, req.Complexity, userID)
	if err != nil {
		atomic.AddInt64(&h.errorsTotal, 1)
		logger.Error("Failed to create challenge", zap.Error(err))
		http.Error(w, "Failed to create challenge", http.StatusInternalServerError)
		return
	}

	atomic.AddInt64(&h.challengesTotal, 1)
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

func (h *Handlers) HandleValidateRequest(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt64(&h.requestsTotal, 1)

	var req struct {
		ChallengeID string      `json:"challenge_id"`
		Answer      interface{} `json:"answer"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		atomic.AddInt64(&h.errorsTotal, 1)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	valid, confidence, err := h.captchaService.ValidateChallenge(r.Context(), req.ChallengeID, req.Answer)
	if err != nil {
		atomic.AddInt64(&h.errorsTotal, 1)
		logger.Error("Failed to validate challenge", zap.Error(err))
		http.Error(w, "Failed to validate challenge", http.StatusInternalServerError)
		return
	}

	atomic.AddInt64(&h.validationsTotal, 1)
	response := map[string]interface{}{
		"valid":      valid,
		"confidence": confidence,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) HandleMemoryStats(w http.ResponseWriter, r *http.Request) {
	if h.memoryMonitor != nil {
		h.memoryMonitor.HandleMemoryStats(w, r)
		return
	}

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	stats := map[string]interface{}{
		"timestamp": time.Now().Unix(),
		"system": map[string]interface{}{
			"alloc_mb":        memStats.Alloc / 1024 / 1024,
			"total_alloc_mb":  memStats.TotalAlloc / 1024 / 1024,
			"sys_mb":          memStats.Sys / 1024 / 1024,
			"num_gc":          memStats.NumGC,
			"gc_cpu_fraction": memStats.GCCPUFraction,
		},
		"status": "healthy",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (h *Handlers) HandleStats(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(h.startTime)

	stats := map[string]interface{}{
		"timestamp":         time.Now().Unix(),
		"uptime_seconds":    uptime.Seconds(),
		"requests_total":    atomic.LoadInt64(&h.requestsTotal),
		"challenges_total":  atomic.LoadInt64(&h.challengesTotal),
		"validations_total": atomic.LoadInt64(&h.validationsTotal),
		"errors_total":      atomic.LoadInt64(&h.errorsTotal),
		"status":            "healthy",
		"uptime_human":      uptime.String(),
	}

	if uptime.Seconds() > 0 {
		stats["rps"] = float64(atomic.LoadInt64(&h.requestsTotal)) / uptime.Seconds()
	} else {
		stats["rps"] = 0.0
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (h *Handlers) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"version": "1.0.0",
	})
}
