package websocket

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"captcha-service/internal/config"
	"captcha-service/internal/domain/entity"
	"captcha-service/internal/infrastructure/repository"
	"captcha-service/internal/usecase"

	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second

	pongWait = 60 * time.Second

	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 512

	MessageTypeChallengeRequest   = "challenge_request"
	MessageTypeChallengeResponse  = "challenge_response"
	MessageTypeValidateRequest    = "validate_request"
	MessageTypeValidationResponse = "validation_response"
	MessageTypeError              = "error"
	MessageTypeUserBlocked        = "user_blocked"
	MessageTypeCaptchaEvent       = "captcha_event"
)

type DemoWebSocketHandler struct {
	upgrader    websocket.Upgrader
	usecase     *usecase.DemoUsecase
	sessionRepo *repository.InMemorySessionRepository
	config      *config.DemoConfig
	connections map[string]*websocket.Conn
}

type WebSocketMessage struct {
	Type          string                 `json:"type"`
	UserID        string                 `json:"user_id,omitempty"`
	SessionID     string                 `json:"session_id,omitempty"`
	ChallengeID   string                 `json:"challenge_id,omitempty"`
	ChallengeType string                 `json:"challenge_type,omitempty"`
	Complexity    int32                  `json:"complexity,omitempty"`
	Data          map[string]interface{} `json:"data,omitempty"`
}

type ValidationResponse struct {
	ChallengeID string `json:"challenge_id"`
	Valid       bool   `json:"valid"`
	Confidence  int32  `json:"confidence"`
	Success     bool   `json:"success"`
	Blocked     bool   `json:"blocked"`
	Error       string `json:"error,omitempty"`
}

func NewDemoWebSocketHandler(demoConfig *config.DemoConfig, sessionRepo *repository.InMemorySessionRepository) *DemoWebSocketHandler {
	usecase := usecase.NewDemoUsecase(sessionRepo, demoConfig)

	return &DemoWebSocketHandler{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		usecase:     usecase,
		sessionRepo: sessionRepo,
		config:      demoConfig,
		connections: make(map[string]*websocket.Conn),
	}
}

func (h *DemoWebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Printf("WebSocket connection attempt from %s", r.RemoteAddr)

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	sessionID := h.getOrCreateSession(w, r)
	userID := h.getUserIDFromSession(sessionID)
	h.connections[userID] = conn

	log.Printf("WebSocket client connected for user: %s (session: %s)", userID, sessionID)

	log.Printf("Checking if user %s is blocked...", userID)
	blocked, err := h.usecase.IsUserBlocked(userID)
	if err != nil {
		log.Printf("Error checking user block status: %v", err)
		h.sendErrorResponse(conn, userID, sessionID, "", "Internal error")
		return
	}

	log.Printf("User %s blocked status: %v", userID, blocked)
	if blocked {
		log.Printf("User %s is blocked, sending blocked response", userID)
		h.sendBlockedResponse(conn, userID, sessionID, "", "User is blocked due to too many failed attempts")
		return
	}

	h.sendConnectionMessage(conn, userID, sessionID)

	go func() {
		conn.SetReadLimit(maxMessageSize)
		conn.SetReadDeadline(time.Now().Add(pongWait))
		conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
		for {
			messageType, data, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					log.Printf("WebSocket read error for user %s: %v", userID, err)
				}
				break
			}

			if messageType == websocket.BinaryMessage {
				h.handleBinaryMessage(conn, data, userID)
			} else {
				h.handleTextMessage(conn, data, userID)
			}
		}
		log.Printf("WebSocket client disconnected for user: %s", userID)
		delete(h.connections, userID)
	}()

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("WebSocket ping error:", err)
				return
			}
		case <-r.Context().Done():
			return
		}
	}
}

func (h *DemoWebSocketHandler) handleBinaryMessage(conn *websocket.Conn, data []byte, userID string) {
	if len(data) < 4 {
		log.Printf("Invalid binary message length: %d", len(data))
		return
	}

	jsonLength := binary.LittleEndian.Uint32(data[:4])
	if len(data) < int(4+jsonLength) {
		log.Printf("Invalid binary message format")
		return
	}

	jsonData := data[4 : 4+jsonLength]
	var msg WebSocketMessage
	if err := json.Unmarshal(jsonData, &msg); err != nil {
		log.Printf("Failed to unmarshal JSON from binary message: %v", err)
		return
	}

	binaryPayload := data[4+jsonLength:]
	msg.UserID = userID

	h.processMessage(conn, msg, binaryPayload)
}

func (h *DemoWebSocketHandler) handleTextMessage(conn *websocket.Conn, data []byte, userID string) {
	var msg WebSocketMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return
	}

	msg.UserID = userID
	h.processMessage(conn, msg, nil)
}

func (h *DemoWebSocketHandler) processMessage(conn *websocket.Conn, msg WebSocketMessage, binaryData []byte) {
	switch msg.Type {
	case MessageTypeChallengeRequest:
		h.handleChallengeRequest(conn, msg)
	case MessageTypeValidateRequest:
		h.handleValidateRequest(conn, msg, binaryData)
	case MessageTypeCaptchaEvent:
		h.handleCaptchaEvent(conn, msg)
	default:
		h.sendErrorResponse(conn, msg.UserID, msg.SessionID, msg.ChallengeID, fmt.Sprintf("Unknown message type: %s", msg.Type))
	}
}

func (h *DemoWebSocketHandler) handleChallengeRequest(conn *websocket.Conn, msg WebSocketMessage) {
	userID := msg.UserID
	sessionID := msg.SessionID

	blocked, err := h.usecase.IsUserBlocked(userID)
	if err != nil {
		log.Printf("Error checking user block status: %v", err)
		h.sendErrorResponse(conn, userID, sessionID, "", "Internal error")
		return
	}

	if blocked {
		h.sendBlockedResponse(conn, userID, sessionID, "", "User is blocked due to too many failed attempts")
		return
	}

	complexity := h.config.DefaultComplexity
	if msg.Data != nil {
		if comp, ok := msg.Data["complexity"].(float64); ok {
			complexity = int32(comp)
		}
	}

	challenge, err := h.usecase.CreateChallenge(userID, "slider-puzzle", complexity)
	if err != nil {
		log.Printf("Error creating challenge: %v", err)
		h.sendErrorResponse(conn, userID, sessionID, "", "Failed to create challenge")
		return
	}

	response := WebSocketMessage{
		Type:          MessageTypeChallengeResponse,
		UserID:        userID,
		SessionID:     sessionID,
		ChallengeID:   challenge.ID,
		ChallengeType: challenge.Type,
		Complexity:    challenge.Complexity,
		Data: map[string]interface{}{
			"html": challenge.HTML,
		},
	}

	h.sendMessage(conn, response)
	log.Printf("Challenge created for user %s: %s", userID, challenge.ID)
}

func (h *DemoWebSocketHandler) handleValidateRequest(conn *websocket.Conn, msg WebSocketMessage, binaryData []byte) {
	userID := msg.UserID
	sessionID := msg.SessionID
	challengeID := msg.ChallengeID

	blocked, err := h.usecase.IsUserBlocked(userID)
	if err != nil {
		log.Printf("Error checking user block status: %v", err)
		h.sendErrorResponse(conn, userID, sessionID, challengeID, "Internal error")
		return
	}

	if blocked {
		h.sendBlockedResponse(conn, userID, sessionID, challengeID, "User is blocked due to too many failed attempts")
		return
	}

	var answer map[string]interface{}
	if binaryData != nil && len(binaryData) > 0 {
		if len(binaryData) >= 4 {
			packed := binary.LittleEndian.Uint32(binaryData)
			answer = map[string]interface{}{
				"x": int(packed & 0x1FFF),         // 13 bits for x
				"y": int((packed >> 13) & 0x1FFF), // 13 bits for y
			}
		}
	} else if msg.Data != nil {
		if ans, ok := msg.Data["answer"].(map[string]interface{}); ok {
			answer = ans
		}
	}

	if answer == nil {
		h.sendErrorResponse(conn, userID, sessionID, challengeID, "No answer provided")
		return
	}

	valid, confidence, err := h.usecase.ValidateChallenge(userID, challengeID, answer)
	if err != nil {
		log.Printf("Error validating challenge: %v", err)
		h.sendErrorResponse(conn, userID, sessionID, challengeID, "Validation error")
		return
	}

	if valid {
		h.sendValidationResponse(conn, userID, sessionID, challengeID, true, confidence, false, "")
		log.Printf("Challenge %s validated successfully for user %s", challengeID, userID)
	} else {
		shouldBlock, err := h.usecase.ShouldBlockUser(userID)
		if err != nil {
			log.Printf("Error checking if user should be blocked: %v", err)
		}

		if shouldBlock {
			err := h.usecase.BlockUser(userID)
			if err != nil {
				log.Printf("Error blocking user: %v", err)
			}
			h.sendValidationResponse(conn, userID, sessionID, challengeID, false, confidence, true, "User blocked due to too many failed attempts")
			log.Printf("User %s blocked after failed validation", userID)
		} else {
			newChallenge, err := h.usecase.CreateChallenge(userID, "slider-puzzle", 50)
			if err != nil {
				log.Printf("Error creating new challenge: %v", err)
				h.sendErrorResponse(conn, userID, sessionID, challengeID, "Failed to create new challenge")
				return
			}

			response := WebSocketMessage{
				Type:          MessageTypeValidationResponse,
				UserID:        userID,
				SessionID:     sessionID,
				ChallengeID:   challengeID,
				ChallengeType: newChallenge.Type,
				Complexity:    newChallenge.Complexity,
				Data: map[string]interface{}{
					"valid":      false,
					"confidence": confidence,
					"success":    false,
					"blocked":    false,
					"new_challenge": map[string]interface{}{
						entity.FieldChallengeID: newChallenge.ID,
						"html":                  newChallenge.HTML,
					},
				},
			}
			h.sendMessage(conn, response)
			log.Printf("New challenge created for user %s after failed validation: %s", userID, newChallenge.ID)
		}
	}
}

func (h *DemoWebSocketHandler) getOrCreateSession(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		sessionID := h.generateSessionID()
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			HttpOnly: true,
			Secure:   false, // Set to true in production with HTTPS
		})
		return sessionID
	}
	return cookie.Value
}

func (h *DemoWebSocketHandler) getUserIDFromSession(sessionID string) string {
	return "demo_user"
}

func (h *DemoWebSocketHandler) generateSessionID() string {
	return fmt.Sprintf("session_%d", time.Now().UnixNano())
}

func (h *DemoWebSocketHandler) sendMessage(conn *websocket.Conn, msg WebSocketMessage) {
	conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

func (h *DemoWebSocketHandler) sendErrorResponse(conn *websocket.Conn, userID, sessionID, challengeID, errorMsg string) {
	response := WebSocketMessage{
		Type:        MessageTypeError,
		UserID:      userID,
		SessionID:   sessionID,
		ChallengeID: challengeID,
		Data: map[string]interface{}{
			"error": errorMsg,
		},
	}
	h.sendMessage(conn, response)
}

func (h *DemoWebSocketHandler) sendBlockedResponse(conn *websocket.Conn, userID, sessionID, challengeID, reason string) {
	response := WebSocketMessage{
		Type:        MessageTypeUserBlocked,
		UserID:      userID,
		SessionID:   sessionID,
		ChallengeID: challengeID,
		Data: map[string]interface{}{
			"reason": reason,
		},
	}
	h.sendMessage(conn, response)
}

func (h *DemoWebSocketHandler) sendValidationResponse(conn *websocket.Conn, userID, sessionID, challengeID string, valid bool, confidence int32, blocked bool, errorMsg string) {
	response := WebSocketMessage{
		Type:        MessageTypeValidationResponse,
		UserID:      userID,
		SessionID:   sessionID,
		ChallengeID: challengeID,
		Data: map[string]interface{}{
			"valid":      valid,
			"confidence": confidence,
			"success":    valid,
			"blocked":    blocked,
			"error":      errorMsg,
		},
	}
	h.sendMessage(conn, response)
}

func (h *DemoWebSocketHandler) sendConnectionMessage(conn *websocket.Conn, userID, sessionID string) {
	response := WebSocketMessage{
		Type:      "connected",
		UserID:    userID,
		SessionID: sessionID,
		Data: map[string]interface{}{
			"user_id":    userID,
			"session_id": sessionID,
			"message":    "Connected to WebSocket",
		},
	}
	h.sendMessage(conn, response)
	log.Printf("Sent connection message for user: %s", userID)
}

func (h *DemoWebSocketHandler) handleCaptchaEvent(conn *websocket.Conn, msg WebSocketMessage) {
	userID := msg.UserID
	sessionID := msg.SessionID

	log.Printf("Received captcha event from user %s: %+v", userID, msg.Data)

	eventType, ok := msg.Data["eventType"].(string)
	if !ok {
		log.Printf("Invalid event type in data: %+v", msg.Data)
		h.sendErrorResponse(conn, userID, sessionID, msg.ChallengeID, "Invalid event type")
		return
	}

	log.Printf("Processing captcha event type: %s", eventType)

	switch eventType {
	case "userBlocked":
		err := h.usecase.BlockUser(userID)
		if err != nil {
			log.Printf("Error blocking user %s: %v", userID, err)
		} else {
			log.Printf("User %s blocked due to captcha failures", userID)
		}

		h.sendBlockedResponse(conn, userID, sessionID, msg.ChallengeID, "User blocked due to too many failed attempts")

	case "newChallenge":
		log.Printf("User %s requested new challenge", userID)

		rand.Seed(time.Now().UnixNano())
		backgroundImage := entity.BackgroundImages[rand.Intn(len(entity.BackgroundImages))]
		puzzleShape := entity.PuzzleShapes[rand.Intn(len(entity.PuzzleShapes))]
		newTargetX := rand.Intn(340) + 30 // Random position between 30 and 370

		response := WebSocketMessage{
			Type:      "new_challenge_data",
			UserID:    userID,
			SessionID: sessionID,
			Data: map[string]interface{}{
				"background_image": fmt.Sprintf("http://localhost:8081/backgrounds/%s", backgroundImage),
				"puzzle_shape":     puzzleShape,
				"target_x":         newTargetX,
				"challenge_id":     fmt.Sprintf("mock_challenge_%d", time.Now().UnixNano()),
			},
		}
		h.sendMessage(conn, response)
		log.Printf("Sent new challenge data: background=%s, shape=%s, targetX=%d", backgroundImage, puzzleShape, newTargetX)

	case "captchaFailed":
		log.Printf("User %s failed captcha attempt", userID)

		err := h.usecase.IncrementAttempts(userID)
		if err != nil {
			log.Printf("Error incrementing attempts for user %s: %v", userID, err)
		} else {
			log.Printf("Incremented attempts for user %s", userID)
		}

		shouldBlock, err := h.usecase.ShouldBlockUser(userID)
		if err != nil {
			log.Printf("Error checking if user should be blocked: %v", err)
		} else {
			log.Printf("Should block user %s: %v", userID, shouldBlock)
			if shouldBlock {
				err := h.usecase.BlockUser(userID)
				if err != nil {
					log.Printf("Error blocking user: %v", err)
				} else {
					log.Printf("User %s blocked after failed attempt", userID)
					h.sendBlockedResponse(conn, userID, sessionID, msg.ChallengeID, "User blocked due to too many failed attempts")
					return
				}
			}
		}
	}

	response := WebSocketMessage{
		Type:      "captcha_event_ack",
		UserID:    userID,
		SessionID: sessionID,
		Data: map[string]interface{}{
			"message": "Captcha event received",
			"event":   msg.Data,
		},
	}

	h.sendMessage(conn, response)
	log.Printf("Sent captcha event acknowledgment for user: %s", userID)
}
