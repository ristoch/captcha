package http

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	protoBalancer "captcha-service/gen/proto/proto/balancer"
	captchaProto "captcha-service/gen/proto/proto/captcha"
	"captcha-service/internal/domain/entity"
	"captcha-service/internal/service"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func isUserBlockedError(err error) bool {
	return errors.Is(err, entity.ErrUserBlocked)
}

type BalancerProxy struct {
	captchaClients []captchaProto.CaptchaServiceClient
	serviceAddrs   []string
	balancerClient protoBalancer.BalancerServiceClient
	mu             sync.RWMutex
	roundRobin     int
	upgrader       websocket.Upgrader
	config         *entity.Config
	sessions       map[string]*entity.UserSession
	sessionMu      sync.RWMutex
	globalBlocker  *service.GlobalUserBlocker
}

func NewBalancerProxy(config *entity.Config) *BalancerProxy {
	captchaClients := make([]captchaProto.CaptchaServiceClient, 0)
	serviceAddrs := make([]string, 0)
	sessions := make(map[string]*entity.UserSession)

	return &BalancerProxy{
		captchaClients: captchaClients,
		serviceAddrs:   serviceAddrs,
		roundRobin:     0,
		sessions:       sessions,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		config:        config,
		globalBlocker: service.NewGlobalUserBlocker(config),
	}
}

func (bp *BalancerProxy) ConnectToBalancer(balancerAddr string) error {
	conn, err := grpc.Dial(balancerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to balancer: %w", err)
	}

	bp.balancerClient = protoBalancer.NewBalancerServiceClient(conn)
	return nil
}

func (bp *BalancerProxy) AddCaptchaService(addr string) error {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to captcha service %s: %w", addr, err)
	}

	client := captchaProto.NewCaptchaServiceClient(conn)

	bp.mu.Lock()
	bp.captchaClients = append(bp.captchaClients, client)
	bp.serviceAddrs = append(bp.serviceAddrs, addr)
	bp.mu.Unlock()

	log.Printf("Added captcha service: %s", addr)
	return nil
}

func (bp *BalancerProxy) RemoveCaptchaService(addr string) error {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	for i, serviceAddr := range bp.serviceAddrs {
		if serviceAddr == addr {
			bp.captchaClients = append(bp.captchaClients[:i], bp.captchaClients[i+1:]...)
			bp.serviceAddrs = append(bp.serviceAddrs[:i], bp.serviceAddrs[i+1:]...)

			log.Printf("Removed captcha service: %s", addr)
			return nil
		}
	}

	return fmt.Errorf("service not found: %s", addr)
}

func (bp *BalancerProxy) GetNextClient() captchaProto.CaptchaServiceClient {
	bp.mu.RLock()
	defer bp.mu.RUnlock()

	if len(bp.captchaClients) == 0 {
		return nil
	}

	client := bp.captchaClients[bp.roundRobin%len(bp.captchaClients)]
	bp.roundRobin++
	return client
}

func (bp *BalancerProxy) NewChallengeHandler(w http.ResponseWriter, r *http.Request) {
	complexityStr := r.URL.Query().Get("complexity")
	complexity := bp.config.ComplexityMedium

	if complexityStr != "" {
		if parsed, err := strconv.ParseInt(complexityStr, 10, 32); err == nil {
			complexity = int32(parsed)
		} else {
			http.Error(w, "Invalid complexity parameter", http.StatusBadRequest)
			return
		}
	}

	session := bp.getOrCreateSession(r)
	userID := session.UserID
	log.Printf("Using userID from session: %s", userID)

	if isBlocked, blockDuration := bp.isUserBlockedInSession(userID); isBlocked {
		log.Printf("User %s is blocked in session, showing blocked page", userID)

		displayUserID := userID
		if displayUserID == "anonymous" {
			displayUserID = "неизвестный пользователь"
		}

		blockedHTML := bp.generateBlockedPage(displayUserID, blockDuration)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(blockedHTML))
		return
	}

	if bp.balancerClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), entity.DefaultTimeoutSeconds*time.Second)
		defer cancel()

		checkResp, err := bp.balancerClient.CheckUserBlocked(ctx, &protoBalancer.CheckUserBlockedRequest{
			UserId: userID,
		})
		if err != nil {
			log.Printf("Failed to check user blocked status on balancer: %v", err)
		} else if checkResp.IsBlocked {
			log.Printf("User %s is blocked on balancer, showing blocked page", userID)

			displayUserID := userID
			if displayUserID == "anonymous" {
				displayUserID = "неизвестный пользователь"
			}

			blockedUntil := time.Unix(checkResp.BlockedUntil, 0)
			remainingMinutes := int(time.Until(blockedUntil).Minutes())
			if remainingMinutes < 0 {
				remainingMinutes = 0
			}

			blockedHTML := bp.generateBlockedPage(displayUserID, fmt.Sprintf("%d", remainingMinutes))
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(blockedHTML))
			return
		}
	}

	client := bp.GetNextClient()
	if client == nil {
		http.Error(w, "No captcha services available", http.StatusServiceUnavailable)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.NewChallenge(ctx, &captchaProto.ChallengeRequest{
		Complexity: complexity,
		UserId:     userID,
	})
	if err != nil {
		log.Printf("Failed to create challenge: %v", err)

		if isUserBlockedError(err) {
			blockDuration := fmt.Sprintf("%d", bp.config.BlockDurationMin)

			displayUserID := userID
			if displayUserID == "anonymous" {
				displayUserID = "неизвестный пользователь"
			}

			blockedHTML := bp.generateBlockedPage(displayUserID, blockDuration)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(blockedHTML))
			return
		}

		http.Error(w, "Failed to create challenge", http.StatusInternalServerError)
		return
	}

	htmlWithWebSocket := bp.addWebSocketCode(resp.Html, userID)

	http.SetCookie(w, &http.Cookie{
		Name:     "captcha_user_id",
		Value:    userID,
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(htmlWithWebSocket))
}

func (bp *BalancerProxy) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := bp.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	log.Println("WebSocket client connected")

	responseChan := make(chan entity.WebSocketMessage, 100)
	defer close(responseChan)

	go func() {
		for response := range responseChan {
			if err := conn.WriteJSON(response); err != nil {
				log.Printf("Failed to send response to WebSocket: %v", err)
				return
			}
		}
	}()

	for {
		var msg entity.WebSocketMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		if msg.UserID != "" {
			if isBlocked, blockDuration := bp.isUserBlockedInSession(msg.UserID); isBlocked {
				log.Printf("User %s is blocked in WebSocket, sending blocked response", msg.UserID)

				blockedResponse := entity.WebSocketMessage{
					Type:   "grpc_response",
					UserID: msg.UserID,
					Data: map[string]interface{}{
						"error":          "User is blocked due to too many attempts",
						"blocked":        true,
						"block_duration": blockDuration,
					},
				}

				responseChan <- blockedResponse
				continue
			}

			if msg.Type == "captcha_event" {
				if isNowBlocked := bp.incrementAttempts(msg.UserID); isNowBlocked {
					log.Printf("User %s blocked after incrementing attempts", msg.UserID)

					blockedResponse := entity.WebSocketMessage{
						Type:   "grpc_response",
						UserID: msg.UserID,
						Data: map[string]interface{}{
							"error":          "User is blocked due to too many attempts",
							"blocked":        true,
							"block_duration": fmt.Sprintf("%d", bp.config.BlockDurationMin),
						},
					}

					responseChan <- blockedResponse
					continue
				}
			}
		}

		log.Printf("Received WebSocket message: %+v", msg)

		switch msg.Type {
		case "challenge_request":
			if isBlocked, blockDuration := bp.isUserBlockedInSession(msg.UserID); isBlocked {
				log.Printf("User %s is blocked, cannot create challenge via WebSocket", msg.UserID)
				blockedResponse := entity.WebSocketMessage{
					Type:   "grpc_response",
					UserID: msg.UserID,
					Data: map[string]interface{}{
						"error":          "User is blocked due to too many attempts",
						"blocked":        true,
						"block_duration": blockDuration,
					},
				}
				responseChan <- blockedResponse
				continue
			}

			challengeID := fmt.Sprintf("challenge_%d", time.Now().UnixNano())
			responseChan <- entity.WebSocketMessage{
				Type:   "challenge_created",
				UserID: msg.UserID,
				Data: map[string]interface{}{
					entity.FieldChallengeID: challengeID,
					"html":                  "<div>Mock challenge HTML</div>",
				},
			}

		case "captcha_event":
			log.Printf("Received captcha event from user %s", msg.UserID)

		default:
			log.Printf("Unknown message type: %s", msg.Type)
		}
	}

	log.Println("WebSocket client disconnected")
}

func (bp *BalancerProxy) StartServiceDiscovery() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	bp.discoverServices()

	for range ticker.C {
		bp.discoverServices()
	}
}

func (bp *BalancerProxy) discoverServices() {
	if bp.balancerClient == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), entity.DefaultTimeoutSeconds*time.Second)
	defer cancel()

	resp, err := bp.balancerClient.GetInstances(ctx, &protoBalancer.GetInstancesRequest{})
	if err != nil {
		log.Printf("Failed to get instances from balancer: %v", err)
		return
	}

	currentServices := make(map[string]bool)
	for _, instance := range resp.Instances {
		address := fmt.Sprintf("%s:%d", instance.Host, instance.PortNumber)
		currentServices[address] = true

		bp.mu.RLock()
		exists := false
		for _, addr := range bp.serviceAddrs {
			if addr == address {
				exists = true
				break
			}
		}
		bp.mu.RUnlock()

		if !exists {
			log.Printf("Discovered new captcha service: %s", address)
			if err := bp.AddCaptchaService(address); err != nil {
				log.Printf("Failed to add discovered service %s: %v", address, err)
			}
		}
	}

	bp.mu.Lock()
	var toRemove []int
	for i, addr := range bp.serviceAddrs {
		if !currentServices[addr] {
			toRemove = append(toRemove, i)
		}
	}

	for i := len(toRemove) - 1; i >= 0; i-- {
		idx := toRemove[i]
		addr := bp.serviceAddrs[idx]
		bp.captchaClients = append(bp.captchaClients[:idx], bp.captchaClients[idx+1:]...)
		bp.serviceAddrs = append(bp.serviceAddrs[:idx], bp.serviceAddrs[idx+1:]...)
		log.Printf("Removed stale captcha service: %s", addr)
	}
	bp.mu.Unlock()

}

func (bp *BalancerProxy) addWebSocketCode(html, userID string) string {
	tmpl, err := template.ParseFiles("templates/websocket.js")
	if err != nil {
		log.Printf("Failed to parse websocket template: %v", err)
		return html
	}

	var websocketCode bytes.Buffer
	data := struct {
		UserID string
	}{
		UserID: userID,
	}

	if err := tmpl.Execute(&websocketCode, data); err != nil {
		log.Printf("Failed to execute websocket template: %v", err)
		return html
	}

	scriptCode := fmt.Sprintf("<script>\n%s\n</script>", websocketCode.String())

	if len(html) > 0 {
		lastBodyIndex := -1
		for i := len(html) - 1; i >= 0; i-- {
			if i+6 <= len(html) && html[i:i+6] == "</body>" {
				lastBodyIndex = i
				break
			}
		}

		if lastBodyIndex != -1 {
			return html[:lastBodyIndex] + scriptCode + html[lastBodyIndex:]
		}
	}

	return html + scriptCode
}

func (bp *BalancerProxy) HealthHandler(w http.ResponseWriter, r *http.Request) {
	bp.mu.RLock()
	serviceCount := len(bp.captchaClients)
	bp.mu.RUnlock()

	response := map[string]interface{}{
		"status":        "healthy",
		"service_count": serviceCount,
		"timestamp":     time.Now().Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (bp *BalancerProxy) ListServicesHandler(w http.ResponseWriter, r *http.Request) {
	bp.mu.RLock()
	services := make([]entity.Instance, len(bp.serviceAddrs))
	for i, addr := range bp.serviceAddrs {
		services[i] = entity.Instance{
			ID:     fmt.Sprintf("service_%d", i),
			Type:   "captcha",
			Host:   addr,
			Port:   8080,
			Status: "active",
		}
	}
	bp.mu.RUnlock()

	response := map[string]interface{}{
		"services": services,
		"count":    len(services),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (bp *BalancerProxy) AddServiceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Address string `json:"address"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	bp.mu.RLock()
	for _, addr := range bp.serviceAddrs {
		if addr == req.Address {
			http.Error(w, "Service already exists", http.StatusConflict)
			bp.mu.RUnlock()
			return
		}
	}
	bp.mu.RUnlock()

	if err := bp.AddCaptchaService(req.Address); err != nil {
		http.Error(w, "Failed to add service", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (bp *BalancerProxy) RemoveServiceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Address string `json:"address"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := bp.RemoveCaptchaService(req.Address); err != nil {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (bp *BalancerProxy) ChallengeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Complexity int    `json:"complexity"`
		UserID     string `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), entity.DefaultTimeoutSeconds*time.Second)
	defer cancel()

	instances, err := bp.balancerClient.GetInstances(ctx, &protoBalancer.GetInstancesRequest{})
	if err != nil {
		http.Error(w, "Failed to get instances: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(instances.Instances) == 0 {
		http.Error(w, "No instances available", http.StatusServiceUnavailable)
		return
	}

	instance := instances.Instances[0]
	address := fmt.Sprintf("%s:%d", instance.Host, instance.PortNumber)
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		http.Error(w, "Failed to connect to captcha service: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := captchaProto.NewCaptchaServiceClient(conn)
	resp, err := client.NewChallenge(ctx, &captchaProto.ChallengeRequest{
		Complexity: int32(req.Complexity),
		UserId:     req.UserID,
	})
	if err != nil {
		http.Error(w, "Failed to create challenge: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		entity.FieldChallengeID: resp.ChallengeId,
		"html":                  resp.Html,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (bp *BalancerProxy) ValidateChallengeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ChallengeID string      `json:entity.FieldChallengeID`
		Answer      interface{} `json:"answer"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.ChallengeID == "" {
		http.Error(w, "challenge_id is required", http.StatusBadRequest)
		return
	}

	client := bp.GetNextClient()
	if client == nil {
		http.Error(w, "No captcha services available", http.StatusServiceUnavailable)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	answerJSON, err := json.Marshal(req.Answer)
	if err != nil {
		http.Error(w, "Failed to marshal answer", http.StatusBadRequest)
		return
	}

	resp, err := client.ValidateChallenge(ctx, &captchaProto.ValidateRequest{
		ChallengeId: req.ChallengeID,
		Answer:      string(answerJSON),
	})
	if err != nil {
		log.Printf("Failed to validate challenge: %v", err)
		http.Error(w, "Failed to validate challenge: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"valid":      resp.Valid,
		"confidence": resp.Confidence,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (bp *BalancerProxy) BlockedPageHandler(w http.ResponseWriter, r *http.Request) {
	duration := r.URL.Query().Get("duration")
	userID := r.URL.Query().Get("user_id")

	if duration == "" {
		duration = "5"
	}
	if userID == "" {
		userID = "unknown"
	}

	durationInt, err := strconv.Atoi(duration)
	if err != nil {
		durationInt = int(bp.config.BlockDurationMin)
	}

	templatePath := "./templates/blocked.html"
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Доступ заблокирован</title>
</head>
<body>
    <h1>Доступ заблокирован</h1>
    <p>Превышено максимальное количество попыток решения капчи.</p>
    <p>Попробуйте снова через %d минут</p>
    <p>ID сессии: %s</p>
    <button onclick="window.location.reload()">Обновить страницу</button>
</body>
</html>`, durationInt, userID)
		return
	}

	data := struct {
		BlockDurationMinutes int
		UserID               string
	}{
		BlockDurationMinutes: durationInt,
		UserID:               userID,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusTooManyRequests)
	tmpl.Execute(w, data)
}

func (bp *BalancerProxy) generateBlockedPage(userID, blockDuration string) string {
	durationInt, err := strconv.Atoi(blockDuration)
	if err != nil {
		durationInt = int(bp.config.BlockDurationMin)
	}

	templatePath := "./templates/blocked.html"
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Доступ заблокирован</title>
</head>
<body>
    <h1>Доступ заблокирован</h1>
    <p>Превышено максимальное количество попыток решения капчи.</p>
    <p>Попробуйте снова через %d минут</p>
    <p>ID сессии: %s</p>
    <button onclick="window.location.reload()">Обновить страницу</button>
</body>
</html>`, durationInt, userID)
	}

	data := struct {
		BlockDurationMinutes int
		UserID               string
	}{
		BlockDurationMinutes: durationInt,
		UserID:               userID,
	}

	var buf bytes.Buffer
	tmpl.Execute(&buf, data)

	return buf.String()
}

func (bp *BalancerProxy) generateSecureUserID(r *http.Request) string {
	ip := r.RemoteAddr
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ip = strings.Split(forwarded, ",")[0]
	} else if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		ip = realIP
	}

	userAgent := r.Header.Get("User-Agent")
	if userAgent == "" {
		userAgent = "unknown"
	}

	identifier := fmt.Sprintf("%s|%s", ip, userAgent)

	hash := sha256.Sum256([]byte(identifier))
	userID := hex.EncodeToString(hash[:])[:16]

	return fmt.Sprintf("secure-%s", userID)
}

func (bp *BalancerProxy) getOrCreateSession(r *http.Request) *entity.UserSession {
	var sessionID string
	if cookie, err := r.Cookie("captcha_user_id"); err == nil {
		sessionID = cookie.Value
		log.Printf("Using userID from cookie: %s", sessionID)
	} else {
		sessionID = bp.generateSecureUserID(r)
		log.Printf("Generated new userID: %s", sessionID)
	}

	bp.sessionMu.Lock()
	defer bp.sessionMu.Unlock()

	if session, exists := bp.sessions[sessionID]; exists {
		session.LastSeen = time.Now()

		if session.IsBlocked && time.Now().After(session.BlockedUntil) {
			log.Printf("Block expired for userID: %s", sessionID)
			session.IsBlocked = false
			session.BlockedUntil = time.Time{}
		}

		return session
	}

	session := &entity.UserSession{
		UserID:       sessionID,
		SessionID:    sessionID,
		CreatedAt:    time.Now(),
		LastSeen:     time.Now(),
		IsBlocked:    false,
		BlockedUntil: time.Time{},
		Attempts:     0,
	}

	bp.sessions[sessionID] = session
	log.Printf("Created new session for userID: %s", sessionID)

	return session
}

func (bp *BalancerProxy) CleanupSessions() {
	bp.sessionMu.Lock()
	defer bp.sessionMu.Unlock()

	now := time.Now()
	cutoff := now.Add(-24 * time.Hour)

	for sessionID, session := range bp.sessions {
		if session.LastSeen.Before(cutoff) {
			delete(bp.sessions, sessionID)
			log.Printf("Cleaned up old session: %s", sessionID)
		}
	}
}

func (bp *BalancerProxy) blockUser(userID string, durationMinutes int) {
	bp.sessionMu.Lock()
	defer bp.sessionMu.Unlock()

	if session, exists := bp.sessions[userID]; exists {
		session.IsBlocked = true
		session.BlockedUntil = time.Now().Add(time.Duration(durationMinutes) * time.Minute)
		log.Printf("User %s blocked for %d minutes until %v", userID, durationMinutes, session.BlockedUntil)
	}
}

func (bp *BalancerProxy) incrementAttempts(userID string) bool {
	bp.sessionMu.Lock()
	defer bp.sessionMu.Unlock()

	if session, exists := bp.sessions[userID]; exists {
		if session.Attempts == 0 {
			session.Attempts = 0
		}
		session.Attempts++
		log.Printf("User %s attempt %d/%d", userID, session.Attempts, bp.config.MaxAttempts)

		if int32(session.Attempts) > bp.config.MaxAttempts {
			session.IsBlocked = true
			session.BlockedUntil = time.Now().Add(time.Duration(bp.config.BlockDurationMin) * time.Minute)
			log.Printf("User %s blocked after %d attempts", userID, session.Attempts)

			if bp.balancerClient != nil {
				go func() {
					ctx, cancel := context.WithTimeout(context.Background(), entity.DefaultTimeoutSeconds*time.Second)
					defer cancel()

					_, err := bp.balancerClient.BlockUser(ctx, &protoBalancer.BlockUserRequest{
						UserId:          userID,
						DurationMinutes: int32(bp.config.BlockDurationMin),
						Reason:          "Too many failed attempts",
					})
					if err != nil {
						log.Printf("Failed to block user on balancer: %v", err)
					} else {
						log.Printf("User %s blocked on balancer for %d minutes", userID, bp.config.BlockDurationMin)
					}
				}()
			}

			return true
		}
	}

	return false
}

func (bp *BalancerProxy) resetAttempts(userID string) {
	bp.sessionMu.Lock()
	defer bp.sessionMu.Unlock()

	if session, exists := bp.sessions[userID]; exists {
		session.Attempts = 0
		log.Printf("User %s attempts reset to 0", userID)
	}
}

func (bp *BalancerProxy) isUserBlockedInSession(userID string) (bool, string) {
	if bp.globalBlocker.IsUserBlocked(userID) {
		blockedUser, err := bp.globalBlocker.GetBlockedUser(userID)
		if err == nil {
			remainingMinutes := int(time.Until(blockedUser.BlockedUntil).Minutes())
			return true, fmt.Sprintf("%d", remainingMinutes)
		}
		return true, "Unknown"
	}

	return false, ""
}
