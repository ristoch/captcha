package usecase

import (
	"captcha-service/internal/config"
	"captcha-service/internal/domain/entity"
	"fmt"
	"time"
)

type DemoUsecase struct {
	sessionRepo entity.SessionRepository
	config      *config.DemoConfig
}

func NewDemoUsecase(sessionRepo entity.SessionRepository, config *config.DemoConfig) *DemoUsecase {
	return &DemoUsecase{
		sessionRepo: sessionRepo,
		config:      config,
	}
}

func (u *DemoUsecase) CreateSession(userID string) (*entity.UserSession, error) {
	session := &entity.UserSession{
		UserID:       userID,
		SessionID:    u.generateSessionID(),
		CreatedAt:    time.Now(),
		LastSeen:     time.Now(),
		Attempts:     0,
		IsBlocked:    false,
		BlockedUntil: time.Time{},
	}

	session, err := u.sessionRepo.CreateSession(userID)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (u *DemoUsecase) GetOrCreateSession(userID string) (*entity.UserSession, error) {
	session, err := u.sessionRepo.GetSessionByUserID(userID)
	if err != nil {
		fmt.Printf("GetOrCreateSession: creating new session for user %s\n", userID)
		return u.CreateSession(userID)
	}

	if u.isSessionExpired(session) {
		fmt.Printf("GetOrCreateSession: session expired for user %s, creating new one\n", userID)
		u.sessionRepo.DeleteSession(session.SessionID)
		return u.CreateSession(userID)
	}

	fmt.Printf("GetOrCreateSession: found existing session for user %s, attempts: %d\n", userID, session.Attempts)
	session.LastSeen = time.Now()
	u.sessionRepo.UpdateSession(session)
	return session, nil
}

func (u *DemoUsecase) IsUserBlocked(userID string) (bool, error) {
	session, err := u.GetOrCreateSession(userID)
	if err != nil {
		return false, err
	}

	if !session.IsBlocked {
		return false, nil
	}

	if time.Now().After(session.BlockedUntil) {
		session.IsBlocked = false
		session.BlockedUntil = time.Time{}
		session.Attempts = 0
		u.sessionRepo.UpdateSession(session)
		return false, nil
	}

	return true, nil
}

func (u *DemoUsecase) BlockUser(userID string) error {
	session, err := u.GetOrCreateSession(userID)
	if err != nil {
		return err
	}

	session.IsBlocked = true
	session.BlockedUntil = time.Now().Add(time.Duration(u.config.BlockDuration) * time.Minute)
	session.Attempts = 0

	return u.sessionRepo.UpdateSession(session)
}

func (u *DemoUsecase) IncrementAttempts(userID string) error {
	session, err := u.GetOrCreateSession(userID)
	if err != nil {
		return err
	}

	oldAttempts := session.Attempts
	session.Attempts++
	session.LastSeen = time.Now()

	fmt.Printf("IncrementAttempts: user %s attempts %d -> %d (max: %d)\n", userID, oldAttempts, session.Attempts, u.config.MaxAttempts)

	return u.sessionRepo.UpdateSession(session)
}

func (u *DemoUsecase) ShouldBlockUser(userID string) (bool, error) {
	session, err := u.GetOrCreateSession(userID)
	if err != nil {
		return false, err
	}

	shouldBlock := session.Attempts >= u.config.MaxAttempts
	fmt.Printf("ShouldBlockUser: user %s attempts %d >= %d = %v\n", userID, session.Attempts, u.config.MaxAttempts, shouldBlock)

	return shouldBlock, nil
}

func (u *DemoUsecase) isSessionExpired(session *entity.UserSession) bool {
	return time.Since(session.LastSeen) > 30*time.Minute // Default session timeout
}

func (u *DemoUsecase) generateSessionID() string {
	sessions, _ := u.sessionRepo.GetAllSessions()
	return fmt.Sprintf("session_%d_%d", time.Now().UnixNano(), len(sessions))
}

func (u *DemoUsecase) CreateChallenge(userID, challengeType string, complexity int32) (*entity.Challenge, error) {
	blocked, err := u.IsUserBlocked(userID)
	if err != nil {
		return nil, err
	}
	if blocked {
		return nil, fmt.Errorf("user is blocked")
	}

	challenge := &entity.Challenge{
		ID:         fmt.Sprintf("demo_%d", time.Now().UnixNano()),
		Type:       challengeType,
		Complexity: complexity,
		UserID:     userID,
		CreatedAt:  time.Now(),
		ExpiresAt:  time.Now().Add(5 * time.Minute),
		HTML:       u.generateChallengeHTML(challengeType, complexity),
		Data:       entity.SliderPuzzleData{},
	}

	return challenge, nil
}

func (u *DemoUsecase) ValidateChallenge(userID, challengeID string, answer map[string]interface{}) (bool, int32, error) {
	session, err := u.sessionRepo.GetSessionByUserID(userID)
	if err != nil {
		session, err = u.CreateSession(userID)
		if err != nil {
			return false, 0, err
		}
	}

	blocked, err := u.IsUserBlocked(userID)
	if err != nil {
		return false, 0, err
	}
	if blocked {
		return false, 0, fmt.Errorf("user is blocked")
	}

	x, xOk := answer["x"].(float64)
	y, yOk := answer["y"].(float64)

	if !xOk || !yOk {
		u.recordFailedAttempt(session)
		return false, 0, fmt.Errorf("invalid answer format")
	}

	targetX, targetY := float64(u.config.DefaultTargetX), float64(u.config.DefaultTargetY)
	tolerance := 20.0

	diffX := x - targetX
	diffY := y - targetY

	if diffX*diffX+diffY*diffY <= tolerance*tolerance {
		session.Attempts = 0
		session.IsBlocked = false
		session.BlockedUntil = time.Time{}
		u.sessionRepo.UpdateSession(session)
		return true, 85, nil
	}

	u.recordFailedAttempt(session)
	return false, 0, nil
}

func (u *DemoUsecase) recordFailedAttempt(session *entity.UserSession) {
	session.Attempts++
	session.LastSeen = time.Now()

	if session.Attempts >= u.config.MaxAttempts {
		session.IsBlocked = true
		session.BlockedUntil = time.Now().Add(time.Duration(u.config.BlockDuration) * time.Minute)
	}

	u.sessionRepo.UpdateSession(session)
}

func (u *DemoUsecase) generateChallengeHTML(challengeType string, complexity int32) string {
	return `<div class="captcha-container">
		<h3>Slider Puzzle Captcha</h3>
		<p>Challenge ID: demo_challenge</p>
		<p>Complexity: ` + fmt.Sprintf("%d", complexity) + `</p>
		<div class="slider-area">
			<input type="range" id="xSlider" min="0" max="380" value="0" oninput="updatePosition()">
			<input type="range" id="ySlider" min="0" max="180" value="0" oninput="updatePosition()">
			<p>Position: <span id="position">(0, 0)</span></p>
		</div>
		<button id="validateBtn" onclick="validateChallenge()">Validate</button>
		<script>
			function updatePosition() {
				const x = document.getElementById('xSlider').value;
				const y = document.getElementById('ySlider').value;
				document.getElementById('position').textContent = '(' + x + ', ' + y + ')';
			}
			
			function validateChallenge() {
				const x = parseInt(document.getElementById('xSlider').value);
				const y = parseInt(document.getElementById('ySlider').value);
				
				if (window.parent && window.parent.validateChallengeViaWebSocket) {
					window.parent.validateChallengeViaWebSocket(x, y);
				} else {
					alert('WebSocket validation not available');
				}
			}
		</script>
	</div>`
}
