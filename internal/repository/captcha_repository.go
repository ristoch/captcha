package repository

import (
	"context"
)

type CaptchaPort interface {
	CreateChallenge(ctx context.Context, challengeType string, complexity int32) (interface{}, error)
	ValidateAnswer(ctx context.Context, challengeID string, answer interface{}) (bool, int32, error)
	GetChallenge(ctx context.Context, challengeID string) (interface{}, error)
	GetHealth(ctx context.Context) (interface{}, error)
}
