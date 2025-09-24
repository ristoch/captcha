package interfaces

import (
	"captcha-service/internal/domain/entity"
	"context"
)

type CaptchaService interface {
	CreateChallenge(ctx context.Context, challengeType string, complexity int32, userID string) (*entity.Challenge, error)
	ValidateChallenge(ctx context.Context, challengeID string, answer interface{}) (bool, int32, error)
	GetChallenge(ctx context.Context, challengeID string) (*entity.Challenge, error)
}

type GeneratorRegistry interface {
	Get(challengeType string) (ChallengeGenerator, bool)
	Register(challengeType string, generator ChallengeGenerator)
}

type ChallengeGenerator interface {
	Generate(ctx context.Context, complexity int32, userID string) (*entity.Challenge, error)
	Validate(answer interface{}, data interface{}) (bool, int32, error)
}

type EventProcessor interface {
	ProcessEvent(event *entity.BinaryEvent) (*entity.EventResult, error)
}

type TemplateEngine interface {
	Render(templateName string, data interface{}) (string, error)
}

type ChallengeRepository interface {
	SaveChallenge(ctx context.Context, challenge *entity.Challenge) error
	GetChallenge(ctx context.Context, challengeID string) (*entity.Challenge, error)
	DeleteChallenge(ctx context.Context, challengeID string) error
}

type WebSocketSender interface {
	SendMessage(userID string, message interface{}) error
}

type InstanceRepository interface {
	SaveInstance(instance *entity.Instance) error
	GetInstance(id string) (*entity.Instance, error)
	GetAllInstances() ([]*entity.Instance, error)
	RemoveInstance(id string) error
}

type UserBlockRepository interface {
	SaveBlockedUser(blockedUser *entity.BlockedUser) error
	GetBlockedUser(userID string) (*entity.BlockedUser, error)
	RemoveBlockedUser(userID string) error
	GetAllBlockedUsers() ([]*entity.BlockedUser, error)
	IsUserBlocked(userID string) bool
	BlockUser(userID string, reason string) error
	CleanupExpiredBlocks() error
}

type EventStreamManager interface {
	CreateStream(userID string) (EventStream, error)
	CloseStream(userID string) error
	GetStream(userID string) (EventStream, error)
}

type EventPublisher interface {
	Publish(event *entity.BinaryEvent) error
	Subscribe(userID string, handler EventHandler) error
	Unsubscribe(userID string) error
}

type EventStream interface {
	Send(event *entity.BinaryEvent) error
	Receive() (*entity.BinaryEvent, error)
	Close() error
}

type EventHandler interface {
	Handle(event *entity.BinaryEvent) error
}

type NewTemplateData interface {
	GetData() map[string]interface{}
}

type BalancerService interface {
	RegisterInstance(req *entity.RegisterInstanceRequest) error
	GetInstances() ([]*entity.Instance, error)
	IsUserBlocked(userID string) bool
	BlockUser(userID string, reason string) error
	StartCleanup()
	Stop()
}
