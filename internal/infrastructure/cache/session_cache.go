package cache

import (
	"sync"
	"time"

	"captcha-service/internal/domain/entity"
)

type SessionCache struct {
	sessions      map[string]*entity.UserSession
	mu            sync.RWMutex
	maxSessions   int
	cleanupTicker *time.Ticker
	stopChan      chan struct{}
}

func NewSessionCache(maxSessions int) *SessionCache {
	cache := &SessionCache{
		sessions:    make(map[string]*entity.UserSession),
		maxSessions: maxSessions,
		stopChan:    make(chan struct{}),
	}

	cache.startCleanup()
	return cache
}

func (c *SessionCache) Get(sessionID string) (*entity.UserSession, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	session, exists := c.sessions[sessionID]
	if !exists {
		return nil, false
	}

	if time.Since(session.LastSeen) > 30*time.Minute {
		return nil, false
	}

	return session, true
}

func (c *SessionCache) Set(sessionID string, session *entity.UserSession) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.sessions) >= c.maxSessions {
		c.evictOldestSessions()
	}

	c.sessions[sessionID] = session
}

func (c *SessionCache) Delete(sessionID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.sessions, sessionID)
}

func (c *SessionCache) evictOldestSessions() {
	evictCount := c.maxSessions / 5
	if evictCount == 0 {
		evictCount = 1
	}

	type sessionWithTime struct {
		id       string
		lastSeen time.Time
	}

	var oldest []sessionWithTime
	for id, session := range c.sessions {
		oldest = append(oldest, sessionWithTime{
			id:       id,
			lastSeen: session.LastSeen,
		})
	}

	for i := 0; i < len(oldest)-1; i++ {
		for j := i + 1; j < len(oldest); j++ {
			if oldest[i].lastSeen.After(oldest[j].lastSeen) {
				oldest[i], oldest[j] = oldest[j], oldest[i]
			}
		}
	}

	for i := 0; i < evictCount && i < len(oldest); i++ {
		delete(c.sessions, oldest[i].id)
	}
}

func (c *SessionCache) startCleanup() {
	c.cleanupTicker = time.NewTicker(5 * time.Minute)

	go func() {
		for {
			select {
			case <-c.cleanupTicker.C:
				c.cleanup()
			case <-c.stopChan:
				return
			}
		}
	}()
}

func (c *SessionCache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for id, session := range c.sessions {
		if now.Sub(session.LastSeen) > 30*time.Minute {
			delete(c.sessions, id)
		}
	}
}

func (c *SessionCache) Stop() {
	if c.cleanupTicker != nil {
		c.cleanupTicker.Stop()
	}
	close(c.stopChan)
}

func (c *SessionCache) GetStats() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return map[string]interface{}{
		"total_sessions": len(c.sessions),
		"max_capacity":   c.maxSessions,
		"usage_percent":  float64(len(c.sessions)) / float64(c.maxSessions) * 100,
	}
}
