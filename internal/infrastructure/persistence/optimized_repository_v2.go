package persistence

import (
	"context"
	"sync"
	"time"

	"captcha-service/internal/domain/entity"
	"captcha-service/internal/domain/interfaces"
)

type OptimizedRepositoryV2 struct {
	challenges map[string]*entity.Challenge
	mu         sync.RWMutex
}

func NewOptimizedRepositoryV2() interfaces.ChallengeRepository {
	return &OptimizedRepositoryV2{
		challenges: make(map[string]*entity.Challenge),
	}
}

func (r *OptimizedRepositoryV2) SaveChallenge(ctx context.Context, challenge *entity.Challenge) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.challenges[challenge.ID] = challenge
	return nil
}

func (r *OptimizedRepositoryV2) GetChallenge(ctx context.Context, challengeID string) (*entity.Challenge, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	challenge, exists := r.challenges[challengeID]
	if !exists {
		return nil, nil
	}

	if time.Now().After(challenge.ExpiresAt) {
		delete(r.challenges, challengeID)
		return nil, nil
	}

	return challenge, nil
}

func (r *OptimizedRepositoryV2) DeleteChallenge(ctx context.Context, challengeID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.challenges, challengeID)
	return nil
}

func (r *OptimizedRepositoryV2) GetChallengesByUser(ctx context.Context, userID string) ([]*entity.Challenge, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var userChallenges []*entity.Challenge
	for _, challenge := range r.challenges {
		if challenge.UserID == userID {
			userChallenges = append(userChallenges, challenge)
		}
	}

	return userChallenges, nil
}
