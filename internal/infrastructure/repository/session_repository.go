package repository

import (
	"captcha-service/internal/domain/entity"
	"fmt"
	"sync"
	"time"
)

type InMemorySessionRepository struct {
	sessions     map[string]*entity.UserSession
	userSessions map[string]string
	mutex        sync.RWMutex
}

func NewInMemorySessionRepository() *InMemorySessionRepository {
	return &InMemorySessionRepository{
		sessions:     make(map[string]*entity.UserSession),
		userSessions: make(map[string]string),
	}
}

func (r *InMemorySessionRepository) CreateSession(userID string) (*entity.UserSession, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	sessionID := fmt.Sprintf("session_%d_%s", time.Now().UnixNano(), userID)
	session := &entity.UserSession{
		UserID:       userID,
		SessionID:    sessionID,
		CreatedAt:    time.Now(),
		LastSeen:     time.Now(),
		Attempts:     0,
		IsBlocked:    false,
		BlockedUntil: time.Time{},
	}

	r.sessions[sessionID] = session
	r.userSessions[userID] = sessionID

	return session, nil
}

func (r *InMemorySessionRepository) GetSession(sessionID string) (*entity.UserSession, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	session, exists := r.sessions[sessionID]
	if !exists {
		return nil, entity.ErrSessionNotFound
	}

	return session, nil
}

func (r *InMemorySessionRepository) GetSessionByUserID(userID string) (*entity.UserSession, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, session := range r.sessions {
		if session.UserID == userID {
			return session, nil
		}
	}
	return nil, entity.ErrSessionNotFound
}

func (r *InMemorySessionRepository) UpdateSession(session *entity.UserSession) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.sessions[session.SessionID]; !exists {
		return entity.ErrSessionNotFound
	}

	r.sessions[session.SessionID] = session
	return nil
}

func (r *InMemorySessionRepository) DeleteSession(sessionID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	session, exists := r.sessions[sessionID]
	if !exists {
		return entity.ErrSessionNotFound
	}

	delete(r.sessions, sessionID)
	delete(r.userSessions, session.UserID)

	return nil
}

func (r *InMemorySessionRepository) GetAllSessions() ([]*entity.UserSession, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	sessions := make([]*entity.UserSession, 0, len(r.sessions))
	for _, session := range r.sessions {
		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (r *InMemorySessionRepository) CleanupExpired() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	now := time.Now()
	for sessionID, session := range r.sessions {
		if now.Sub(session.LastSeen) > 24*time.Hour {
			delete(r.sessions, sessionID)
			delete(r.userSessions, session.UserID)
		}
	}
}
