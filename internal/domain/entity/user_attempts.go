package entity

import (
	"sync"
	"time"
)

type UserAttempts struct {
	mu            sync.RWMutex
	attempts      map[string]*UserAttempt
	config        *DemoConfig
	cleanupTicker *time.Ticker
	stopChan      chan struct{}
}

type UserAttempt struct {
	UserID       string
	Attempts     int32
	LastAttempt  time.Time
	BlockedUntil *time.Time
	ChallengeID  string
}

func NewUserAttempts(config *DemoConfig) *UserAttempts {
	ua := &UserAttempts{
		attempts: make(map[string]*UserAttempt),
		config:   config,
		stopChan: make(chan struct{}),
	}

	ua.startCleanup()
	return ua
}

func (ua *UserAttempts) RecordAttempt(userID, challengeID string) (isBlocked bool, remainingAttempts int32) {
	ua.mu.Lock()
	defer ua.mu.Unlock()

	now := time.Now()
	attempt, exists := ua.attempts[userID]

	if !exists {
		attempt = &UserAttempt{
			UserID:      userID,
			Attempts:    0,
			LastAttempt: now,
			ChallengeID: challengeID,
		}
		ua.attempts[userID] = attempt
	}

	if attempt.BlockedUntil != nil && now.Before(*attempt.BlockedUntil) {
		return true, 0
	}

	if attempt.BlockedUntil != nil && now.After(*attempt.BlockedUntil) {
		attempt.BlockedUntil = nil
		attempt.Attempts = 0
	}

	attempt.Attempts++
	attempt.LastAttempt = now
	attempt.ChallengeID = challengeID

	if attempt.Attempts >= ua.config.MaxAttempts {
		blockedUntil := now.Add(time.Duration(ua.config.BlockDuration) * time.Minute)
		attempt.BlockedUntil = &blockedUntil
		return true, 0
	}

	remainingAttempts = ua.config.MaxAttempts - attempt.Attempts
	return false, remainingAttempts
}

func (ua *UserAttempts) ResetAttempts(userID string) {
	ua.mu.Lock()
	defer ua.mu.Unlock()

	if attempt, exists := ua.attempts[userID]; exists {
		attempt.Attempts = 0
		attempt.BlockedUntil = nil
	}
}

func (ua *UserAttempts) IsBlocked(userID string) bool {
	ua.mu.RLock()
	defer ua.mu.RUnlock()

	attempt, exists := ua.attempts[userID]
	if !exists {
		return false
	}

	if attempt.BlockedUntil == nil {
		return false
	}

	return time.Now().Before(*attempt.BlockedUntil)
}

func (ua *UserAttempts) GetRemainingAttempts(userID string) int32 {
	ua.mu.RLock()
	defer ua.mu.RUnlock()

	attempt, exists := ua.attempts[userID]
	if !exists {
		return ua.config.MaxAttempts
	}

	if attempt.BlockedUntil != nil && time.Now().Before(*attempt.BlockedUntil) {
		return 0
	}

	return ua.config.MaxAttempts - attempt.Attempts
}

func (ua *UserAttempts) startCleanup() {
	ua.cleanupTicker = time.NewTicker(5 * time.Minute) // Default cleanup interval

	go func() {
		for {
			select {
			case <-ua.cleanupTicker.C:
				ua.cleanup()
			case <-ua.stopChan:
				return
			}
		}
	}()
}

func (ua *UserAttempts) cleanup() {
	ua.mu.Lock()
	defer ua.mu.Unlock()

	now := time.Now()
	expiredTime := now.Add(-time.Duration(ua.config.BlockDuration) * time.Minute)

	for userID, attempt := range ua.attempts {
		if attempt.BlockedUntil != nil && attempt.BlockedUntil.Before(now) {
			delete(ua.attempts, userID)
		} else if attempt.LastAttempt.Before(expiredTime) && attempt.BlockedUntil == nil {
			delete(ua.attempts, userID)
		}
	}
}

func (ua *UserAttempts) Stop() {
	if ua.cleanupTicker != nil {
		ua.cleanupTicker.Stop()
	}
	close(ua.stopChan)
}
