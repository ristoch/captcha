package service

import (
	"context"
	"sync"

	"captcha-service/internal/domain/entity"
)

type ChallengeGenerator interface {
	Generate(ctx context.Context, complexity int32, userID string) (*entity.Challenge, error)
	Validate(answer interface{}, data interface{}) (bool, int32, error)
}

type GeneratorRegistry struct {
	generators map[string]ChallengeGenerator
	mu         sync.RWMutex
}

func NewGeneratorRegistry() *GeneratorRegistry {
	return &GeneratorRegistry{
		generators: make(map[string]ChallengeGenerator),
	}
}

func (r *GeneratorRegistry) Register(name string, generator ChallengeGenerator) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.generators[name] = generator
}

func (r *GeneratorRegistry) Get(name string) (ChallengeGenerator, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	generator, exists := r.generators[name]
	return generator, exists
}
