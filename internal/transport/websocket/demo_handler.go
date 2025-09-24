package websocket

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

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
)

type DemoWebSocketHandler struct {
	upgrader     websocket.Upgrader
	usecase      *usecase.DemoUsecase
	sessionRepo  *repository.InMemorySessionRepository
	config       *entity.DemoConfig
	userSessions map[string]*UserSession
	connections  map[string]*websocket.Conn
}

type UserSession struct {
	UserID       string    `json:"user_id"`
	SessionID    string    `json:"session_id"`
	CreatedAt    time.Time `json:"created_at"`
	BlockedUntil time.Time `json:"blocked_until"`
	Attempts     int32     `json:"attempts"`
	MaxAttempts  int32     `json:"max_attempts"`
	IsBlocked    bool      `json:"is_blocked"`
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

func NewDemoWebSocketHandler(demoConfig *entity.DemoConfig, sessionRepo *repository.InMemorySessionRepository) *DemoWebSocketHandler {
	usecase := usecase.NewDemoUsecase(sessionRepo, demoConfig)

	return &DemoWebSocketHandler{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for demo purposes
			},
		},
		usecase:      usecase,
		sessionRepo:  sessionRepo,
		config:       demoConfig,
		userSessions: make(map[string]*UserSession),
		connections:  make(map[string]*websocket.Conn),
	}
}

func (h *DemoWebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	sessionID := h.getOrCreateSession(w, r)
	userID := h.getUserIDFromSession(sessionID)
	h.connections[userID] = conn // Store connection for user

	log.Printf("WebSocket client connected for user: %s", userID)

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
		delete(h.connections, userID) // Remove connection on disconnect
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
	default:
		h.sendErrorResponse(conn, msg.UserID, msg.SessionID, msg.ChallengeID, fmt.Sprintf("Unknown message type: %s", msg.Type))
	}
}

func (h *DemoWebSocketHandler) handleChallengeRequest(conn *websocket.Conn, msg WebSocketMessage) {
	userID := msg.UserID
	sessionID := msg.SessionID

	if h.isUserBlocked(userID) {
		h.sendBlockedResponse(conn, userID, sessionID, "", "User is blocked due to too many failed attempts")
		return
	}

	challenge := &entity.Challenge{
		ID:         fmt.Sprintf("challenge_%d", time.Now().UnixNano()),
		Type:       "slider_puzzle",
		Complexity: 50,
		HTML:       h.generateSliderPuzzleHTML(),
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

	if h.isUserBlocked(userID) {
		h.sendBlockedResponse(conn, userID, sessionID, challengeID, "User is blocked due to too many failed attempts")
		return
	}

	var answer interface{}
	if binaryData != nil && len(binaryData) > 0 {
		if len(binaryData) >= 4 {
			packed := binary.LittleEndian.Uint32(binaryData)
			answer = map[string]interface{}{
				"x": int(packed & 0x1FFF),         // 13 bits for x
				"y": int((packed >> 13) & 0x1FFF), // 13 bits for y
			}
		}
	} else if msg.Data != nil {
		answer = msg.Data["answer"]
	}

	if answer == nil {
		h.sendErrorResponse(conn, userID, sessionID, challengeID, "No answer provided")
		return
	}

	valid := h.mockValidateAnswer(answer)
	confidence := int32(85)

	if valid {
		h.resetUserAttempts(userID)
		h.sendValidationResponse(conn, userID, sessionID, challengeID, true, confidence, false, "")
		log.Printf("Challenge %s validated successfully for user %s", challengeID, userID)
	} else {
		shouldBlock := h.incrementUserAttempts(userID)

		if shouldBlock {
			h.blockUser(userID)
			h.sendValidationResponse(conn, userID, sessionID, challengeID, false, confidence, true, "User blocked due to too many failed attempts")
			log.Printf("User %s blocked after failed validation", userID)
		} else {
			newChallenge := &entity.Challenge{
				ID:         fmt.Sprintf("challenge_%d", time.Now().UnixNano()),
				Type:       "slider_puzzle",
				Complexity: 50,
				HTML:       h.generateSliderPuzzleHTML(),
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
						"challenge_id": newChallenge.ID,
						"html":         newChallenge.HTML,
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
	return sessionID
}

func (h *DemoWebSocketHandler) generateSessionID() string {
	return fmt.Sprintf("session_%d", time.Now().UnixNano())
}

func (h *DemoWebSocketHandler) isUserBlocked(userID string) bool {
	session, exists := h.userSessions[userID]
	if !exists {
		return false
	}
	return session.IsBlocked && time.Now().Before(session.BlockedUntil)
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

func (h *DemoWebSocketHandler) mockValidateAnswer(answer interface{}) bool {
	return false
}

func (h *DemoWebSocketHandler) generateSliderPuzzleHTML() string {
	return `
<!DOCTYPE html>
<html>
<head>
    <title>Slider Puzzle Captcha</title>
    <style>
        body { font-family: Arial, sans-serif; text-align: center; padding: 20px; }
        .puzzle-container { margin: 20px auto; width: 300px; }
        .slider { width: 100%; margin: 10px 0; }
        .message { margin: 10px 0; font-weight: bold; }
    </style>
</head>
<body>
    <h3>Complete the puzzle</h3>
    <div class="puzzle-container">
        <canvas id="puzzle" width="300" height="200" style="border: 1px solid #ccc;"></canvas>
        <input type="range" class="slider" min="0" max="100" value="0" id="slider">
        <div class="message" id="message">Move the slider to complete the puzzle</div>
    </div>
    <script>
        const canvas = document.getElementById('puzzle');
        const ctx = canvas.getContext('2d');
        const slider = document.getElementById('slider');
        const message = document.getElementById('message');
        
        ctx.fillStyle = '#f0f0f0';
        ctx.fillRect(0, 0, 300, 200);
        ctx.fillStyle = '#333';
        ctx.fillRect(50, 50, 200, 100);
        
        slider.addEventListener('input', function() {
            const value = this.value;
            window.top.postMessage({
                type: 'captcha:sendData',
                eventType: 'slider_move',
                data: { x: parseInt(value), y: 100 }
            }, '*');
        });
    </script>
</body>
</html>`
}

func (h *DemoWebSocketHandler) resetUserAttempts(userID string) {
	if session, exists := h.userSessions[userID]; exists {
		session.Attempts = 0
		session.IsBlocked = false
		session.BlockedUntil = time.Time{}
	}
}

func (h *DemoWebSocketHandler) incrementUserAttempts(userID string) bool {
	session, exists := h.userSessions[userID]
	if !exists {
		session = &UserSession{
			UserID:      userID,
			SessionID:   h.generateSessionID(),
			CreatedAt:   time.Now(),
			Attempts:    0,
			MaxAttempts: 3, // Demo limit
			IsBlocked:   false,
		}
		h.userSessions[userID] = session
	}

	session.Attempts++
	return session.Attempts >= session.MaxAttempts
}

func (h *DemoWebSocketHandler) blockUser(userID string) {
	if session, exists := h.userSessions[userID]; exists {
		session.IsBlocked = true
		session.BlockedUntil = time.Now().Add(5 * time.Minute) // Block for 5 minutes
	}
}
