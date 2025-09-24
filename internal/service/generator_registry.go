package service

import (
	"sync"

	"captcha-service/internal/domain/interfaces"
)

type GeneratorRegistry struct {
	generators map[string]interfaces.ChallengeGenerator
	mu         sync.RWMutex
}

func NewGeneratorRegistry() *GeneratorRegistry {
	return &GeneratorRegistry{
		generators: make(map[string]interfaces.ChallengeGenerator),
	}
}

func (r *GeneratorRegistry) Register(name string, generator interfaces.ChallengeGenerator) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.generators[name] = generator
}

func (r *GeneratorRegistry) Get(name string) (interfaces.ChallengeGenerator, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	generator, exists := r.generators[name]
	return generator, exists
}
