package persistence

import (
	"context"
	"fmt"
	"sync"
	"time"

	"captcha-service/internal/domain/entity"
	"captcha-service/internal/domain/interfaces"
)

// MemoryOptimizedRepository implements memory-optimized storage for challenges
type MemoryOptimizedRepository struct {
	challenges    map[string]*entity.Challenge
	mu            sync.RWMutex
	maxChallenges int
	cleanupTicker *time.Ticker
	stopChan      chan struct{}
}

// NewMemoryOptimizedRepository creates a new memory-optimized repository
func NewMemoryOptimizedRepository(maxChallenges int) interfaces.ChallengeRepository {
	repo := &MemoryOptimizedRepository{
		challenges:    make(map[string]*entity.Challenge),
		maxChallenges: maxChallenges,
		stopChan:      make(chan struct{}),
	}

	// Start cleanup goroutine
	repo.startCleanup()
	return repo
}

func (r *MemoryOptimizedRepository) SaveChallenge(ctx context.Context, challenge *entity.Challenge) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if we need to evict old challenges
	if len(r.challenges) >= r.maxChallenges {
		r.evictOldestChallenges()
	}

	r.challenges[challenge.ID] = challenge
	return nil
}

func (r *MemoryOptimizedRepository) GetChallenge(ctx context.Context, challengeID string) (*entity.Challenge, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	challenge, exists := r.challenges[challengeID]
	if !exists {
		return nil, fmt.Errorf("challenge with ID %s not found", challengeID)
	}

	// Check if challenge is expired
	if challenge.ExpiresAt.Before(time.Now()) {
		// Remove expired challenge
		delete(r.challenges, challengeID)
		return nil, fmt.Errorf("challenge with ID %s has expired", challengeID)
	}

	return challenge, nil
}

func (r *MemoryOptimizedRepository) DeleteChallenge(ctx context.Context, challengeID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.challenges, challengeID)
	return nil
}

// evictOldestChallenges removes the oldest challenges to make room for new ones
func (r *MemoryOptimizedRepository) evictOldestChallenges() {
	// Calculate how many to evict (remove 20% of max capacity)
	evictCount := r.maxChallenges / 5
	if evictCount == 0 {
		evictCount = 1
	}

	// Find oldest challenges
	type challengeWithTime struct {
		id        string
		createdAt time.Time
	}

	var oldest []challengeWithTime
	for id, challenge := range r.challenges {
		oldest = append(oldest, challengeWithTime{
			id:        id,
			createdAt: challenge.CreatedAt,
		})
	}

	// Sort by creation time (oldest first)
	for i := 0; i < len(oldest)-1; i++ {
		for j := i + 1; j < len(oldest); j++ {
			if oldest[i].createdAt.After(oldest[j].createdAt) {
				oldest[i], oldest[j] = oldest[j], oldest[i]
			}
		}
	}

	// Remove oldest challenges
	for i := 0; i < evictCount && i < len(oldest); i++ {
		delete(r.challenges, oldest[i].id)
	}
}

// startCleanup starts a background goroutine to clean up expired challenges
func (r *MemoryOptimizedRepository) startCleanup() {
	r.cleanupTicker = time.NewTicker(5 * time.Minute)

	go func() {
		for {
			select {
			case <-r.cleanupTicker.C:
				r.cleanup()
			case <-r.stopChan:
				return
			}
		}
	}()
}

// cleanup removes expired challenges
func (r *MemoryOptimizedRepository) cleanup() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for id, challenge := range r.challenges {
		if challenge.ExpiresAt.Before(now) {
			delete(r.challenges, id)
		}
	}
}

// Stop stops the cleanup goroutine
func (r *MemoryOptimizedRepository) Stop() {
	if r.cleanupTicker != nil {
		r.cleanupTicker.Stop()
	}
	close(r.stopChan)
}

// GetStats returns repository statistics
func (r *MemoryOptimizedRepository) GetStats() map[string]interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return map[string]interface{}{
		"total_challenges": len(r.challenges),
		"max_capacity":     r.maxChallenges,
		"usage_percent":    float64(len(r.challenges)) / float64(r.maxChallenges) * 100,
	}
}
