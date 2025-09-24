package persistence

import (
	"fmt"
	"sync"
	"time"

	"captcha-service/internal/domain/entity"
)

type InstanceInfo struct {
	ID            string
	ChallengeType string
	Host          string
	Port          int32
	Status        string
	LastSeen      time.Time
}

type UserBlockInfo struct {
	UserID       string
	BlockedUntil time.Time
	Reason       string
}

type MemoryInstanceRepository struct {
	instances map[string]*entity.Instance
	mu        sync.RWMutex
}

func NewMemoryInstanceRepository() *MemoryInstanceRepository {
	return &MemoryInstanceRepository{
		instances: make(map[string]*entity.Instance),
	}
}

func (r *MemoryInstanceRepository) SaveInstance(instance *entity.Instance) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.instances[instance.ID] = instance
	return nil
}

func (r *MemoryInstanceRepository) GetInstance(id string) (*entity.Instance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	instance, exists := r.instances[id]
	if !exists {
		return nil, fmt.Errorf("instance not found")
	}
	return instance, nil
}

func (r *MemoryInstanceRepository) GetAllInstances() ([]*entity.Instance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	instances := make([]*entity.Instance, 0, len(r.instances))
	for _, instance := range r.instances {
		instances = append(instances, instance)
	}
	return instances, nil
}

func (r *MemoryInstanceRepository) RemoveInstance(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.instances, id)
	return nil
}

type MemoryUserBlockRepository struct {
	blocks map[string]*UserBlockInfo
	mu     sync.RWMutex
}

func NewMemoryUserBlockRepository() *MemoryUserBlockRepository {
	return &MemoryUserBlockRepository{
		blocks: make(map[string]*UserBlockInfo),
	}
}

func (r *MemoryUserBlockRepository) Save(block *UserBlockInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.blocks[block.UserID] = block
	return nil
}

func (r *MemoryUserBlockRepository) Get(userID string) (*UserBlockInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	block, exists := r.blocks[userID]
	if !exists {
		return nil, nil
	}
	return block, nil
}

func (r *MemoryUserBlockRepository) SaveBlockedUser(blockedUser *entity.BlockedUser) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.blocks[blockedUser.UserID] = &UserBlockInfo{
		UserID:       blockedUser.UserID,
		BlockedUntil: blockedUser.BlockedUntil,
		Reason:       blockedUser.Reason,
	}
	return nil
}

func (r *MemoryUserBlockRepository) GetBlockedUser(userID string) (*entity.BlockedUser, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	block, exists := r.blocks[userID]
	if !exists {
		return nil, fmt.Errorf("user not blocked")
	}
	return &entity.BlockedUser{
		UserID:       block.UserID,
		BlockedUntil: block.BlockedUntil,
		Reason:       block.Reason,
		Attempts:     0, // Default value
	}, nil
}

func (r *MemoryUserBlockRepository) RemoveBlockedUser(userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.blocks, userID)
	return nil
}

func (r *MemoryUserBlockRepository) GetAllBlockedUsers() ([]*entity.BlockedUser, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	blocks := make([]*entity.BlockedUser, 0, len(r.blocks))
	for _, block := range r.blocks {
		blocks = append(blocks, &entity.BlockedUser{
			UserID:       block.UserID,
			BlockedUntil: block.BlockedUntil,
			Reason:       block.Reason,
			Attempts:     0, // Default value
		})
	}
	return blocks, nil
}

func (r *MemoryUserBlockRepository) IsUserBlocked(userID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	block, exists := r.blocks[userID]
	if !exists {
		return false
	}
	return time.Now().Before(block.BlockedUntil)
}

func (r *MemoryUserBlockRepository) BlockUser(userID, reason string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.blocks[userID] = &UserBlockInfo{
		UserID:       userID,
		BlockedUntil: time.Now().Add(5 * time.Minute), // Default block duration
		Reason:       reason,
	}
	return nil
}

func (r *MemoryUserBlockRepository) CleanupExpiredBlocks() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	for userID, block := range r.blocks {
		if now.After(block.BlockedUntil) {
			delete(r.blocks, userID)
		}
	}
	return nil
}

func (r *MemoryUserBlockRepository) Delete(userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.blocks, userID)
	return nil
}
