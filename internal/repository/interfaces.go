package repository

import (
	"captcha-service/internal/domain/entity"
	"context"
	"time"
)

type ChallengeRepository interface {
	Create(ctx context.Context, challenge *entity.Challenge) error
	GetByID(ctx context.Context, id string) (*entity.Challenge, error)
	Update(ctx context.Context, challenge *entity.Challenge) error
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context) ([]*entity.Challenge, error)
	CleanupExpired(ctx context.Context) error
}

type InstanceRepository interface {
	AddInstance(instance *entity.Instance)
	RemoveInstance(instanceID string)
	GetInstance(instanceID string) (*entity.Instance, bool)
	GetAllInstances() []*entity.Instance
	CleanupStaleInstances(staleThreshold time.Duration)
}

type UserBlockRepository interface {
	BlockUser(userID string, duration time.Duration, reason string) error
	IsUserBlocked(userID string) (bool, *entity.BlockedUser, error)
	UnblockUser(userID string) error
	CleanupExpiredBlocks() error
}
