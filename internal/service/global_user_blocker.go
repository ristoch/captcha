package service

import (
	"fmt"
	"sync"
	"time"

	"captcha-service/internal/domain/entity"
	"captcha-service/pkg/logger"

	"go.uber.org/zap"
)

// GlobalUserBlocker manages user blocking across all services
type GlobalUserBlocker struct {
	blockedUsers  map[string]*entity.BlockedUser
	mu            sync.RWMutex
	config        *entity.Config
	cleanupTicker *time.Ticker
	stopChan      chan struct{}
}

// NewGlobalUserBlocker creates a new global user blocker
func NewGlobalUserBlocker(config *entity.Config) *GlobalUserBlocker {
	blocker := &GlobalUserBlocker{
		blockedUsers: make(map[string]*entity.BlockedUser),
		config:       config,
		stopChan:     make(chan struct{}),
	}

	// Start cleanup goroutine
	blocker.startCleanup()
	return blocker
}

// IsUserBlocked checks if a user is currently blocked
func (b *GlobalUserBlocker) IsUserBlocked(userID string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	blockedUser, exists := b.blockedUsers[userID]
	if !exists {
		return false
	}

	// Check if block has expired
	if time.Now().After(blockedUser.BlockedUntil) {
		return false
	}

	return true
}

// BlockUser blocks a user for the specified duration
func (b *GlobalUserBlocker) BlockUser(userID, reason string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	blockedUntil := time.Now().Add(time.Duration(b.config.BlockDurationMin) * time.Minute)

	blockedUser := &entity.BlockedUser{
		UserID:       userID,
		BlockedUntil: blockedUntil,
		Reason:       reason,
		Attempts:     0,
	}

	b.blockedUsers[userID] = blockedUser

	logger.Warn("User blocked globally",
		zap.String("userID", userID),
		zap.String("reason", reason),
		zap.Time("blockedUntil", blockedUntil))

	return nil
}

// UnblockUser removes a user from the blocked list
func (b *GlobalUserBlocker) UnblockUser(userID string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	delete(b.blockedUsers, userID)

	logger.Info("User unblocked globally", zap.String("userID", userID))
	return nil
}

// GetBlockedUser returns information about a blocked user
func (b *GlobalUserBlocker) GetBlockedUser(userID string) (*entity.BlockedUser, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	blockedUser, exists := b.blockedUsers[userID]
	if !exists {
		return nil, fmt.Errorf("user %s is not blocked", userID)
	}

	// Check if block has expired
	if time.Now().After(blockedUser.BlockedUntil) {
		return nil, fmt.Errorf("user %s block has expired", userID)
	}

	return blockedUser, nil
}

// RecordAttempt records a failed attempt for a user
func (b *GlobalUserBlocker) RecordAttempt(userID, challengeID string) (isBlocked bool, remainingAttempts int32) {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	blockedUser, exists := b.blockedUsers[userID]

	if !exists {
		// First attempt
		blockedUser = &entity.BlockedUser{
			UserID:       userID,
			Attempts:     1,
			LastAttempt:  now,
			BlockedUntil: time.Time{},
		}
		b.blockedUsers[userID] = blockedUser

		remainingAttempts = b.config.MaxAttempts - 1
		return false, remainingAttempts
	}

	// Check if user is already blocked
	if !blockedUser.BlockedUntil.IsZero() && now.Before(blockedUser.BlockedUntil) {
		return true, 0
	}

	// Reset if block has expired
	if !blockedUser.BlockedUntil.IsZero() && now.After(blockedUser.BlockedUntil) {
		blockedUser.Attempts = 0
		blockedUser.BlockedUntil = time.Time{}
	}

	blockedUser.Attempts++
	blockedUser.LastAttempt = now

	// Check if user should be blocked
	if blockedUser.Attempts >= b.config.MaxAttempts {
		blockedUser.BlockedUntil = now.Add(time.Duration(b.config.BlockDurationMin) * time.Minute)
		blockedUser.Reason = "Too many failed attempts"

		logger.Warn("User blocked due to max attempts",
			zap.String("userID", userID),
			zap.Int32("attempts", blockedUser.Attempts),
			zap.Time("blockedUntil", blockedUser.BlockedUntil))

		return true, 0
	}

	remainingAttempts = b.config.MaxAttempts - blockedUser.Attempts
	return false, remainingAttempts
}

// ResetAttempts resets the attempt count for a user
func (b *GlobalUserBlocker) ResetAttempts(userID string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if blockedUser, exists := b.blockedUsers[userID]; exists {
		blockedUser.Attempts = 0
		blockedUser.BlockedUntil = time.Time{}
		blockedUser.Reason = ""

		logger.Info("User attempts reset", zap.String("userID", userID))
	}
}

// startCleanup starts a background goroutine to clean up expired blocks
func (b *GlobalUserBlocker) startCleanup() {
	b.cleanupTicker = time.NewTicker(5 * time.Minute)

	go func() {
		for {
			select {
			case <-b.cleanupTicker.C:
				b.cleanup()
			case <-b.stopChan:
				return
			}
		}
	}()
}

// cleanup removes expired blocks
func (b *GlobalUserBlocker) cleanup() {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	for userID, blockedUser := range b.blockedUsers {
		if !blockedUser.BlockedUntil.IsZero() && now.After(blockedUser.BlockedUntil) {
			delete(b.blockedUsers, userID)
			logger.Debug("Removed expired block", zap.String("userID", userID))
		}
	}
}

// Stop stops the cleanup goroutine
func (b *GlobalUserBlocker) Stop() {
	if b.cleanupTicker != nil {
		b.cleanupTicker.Stop()
	}
	close(b.stopChan)
}

// GetStats returns blocker statistics
func (b *GlobalUserBlocker) GetStats() map[string]interface{} {
	b.mu.RLock()
	defer b.mu.RUnlock()

	activeBlocks := 0
	now := time.Now()
	for _, blockedUser := range b.blockedUsers {
		if !blockedUser.BlockedUntil.IsZero() && now.Before(blockedUser.BlockedUntil) {
			activeBlocks++
		}
	}

	return map[string]interface{}{
		"total_blocked_users": len(b.blockedUsers),
		"active_blocks":       activeBlocks,
		"max_attempts":        b.config.MaxAttempts,
		"block_duration_min":  b.config.BlockDurationMin,
	}
}
