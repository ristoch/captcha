package service

import (
	"context"
	"time"

	protoBalancer "captcha-service/gen/proto/proto/balancer"
	"captcha-service/internal/config"
	"captcha-service/internal/domain/entity"
	"captcha-service/pkg/logger"

	"go.uber.org/zap"
)

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

type BalancerServiceInterface interface {
	RegisterInstance(req *entity.RegisterInstanceRequest) error
	GetInstances() ([]*entity.Instance, error)
	IsUserBlocked(userID string) bool
	BlockUser(userID string, reason string) error
	StartCleanup()
	Stop()
}

type BalancerService struct {
	instanceRepo  InstanceRepository
	userBlockRepo UserBlockRepository
	config        *config.ServiceConfig
}

func NewBalancerService(instanceRepo InstanceRepository, userBlockRepo UserBlockRepository, config *config.ServiceConfig) BalancerServiceInterface {
	return &BalancerService{
		instanceRepo:  instanceRepo,
		userBlockRepo: userBlockRepo,
		config:        config,
	}
}

func (s *BalancerService) RegisterInstance(req *entity.RegisterInstanceRequest) error {
	logger.Info("New instance registration started")

	instance := &entity.Instance{
		ID:           req.InstanceID,
		Type:         req.ChallengeType,
		Host:         req.Host,
		Port:         req.PortNumber,
		LastSeen:     time.Now(),
		Status:       req.EventType,
		RegisteredAt: time.Now(),
	}

	if req.EventType == "STOPPED" {
		logger.Info("Instance stopped", zap.String("instance_id", req.InstanceID))
		s.instanceRepo.RemoveInstance(req.InstanceID)
	} else {
		s.instanceRepo.SaveInstance(instance)
	}

	logger.Info("Instance registered",
		zap.String("instance_id", req.InstanceID),
		zap.String("challenge_type", req.ChallengeType),
		zap.String("host", req.Host),
		zap.Int32("port", req.PortNumber),
		zap.String("event_type", req.EventType))

	return nil
}

func (s *BalancerService) GetInstances() ([]*entity.Instance, error) {
	instances, err := s.instanceRepo.GetAllInstances()
	if err != nil {
		return nil, err
	}
	return instances, nil
}

func (s *BalancerService) GetInstancesGRPC(ctx context.Context, req *protoBalancer.GetInstancesRequest) (*protoBalancer.GetInstancesResponse, error) {
	instances, err := s.instanceRepo.GetAllInstances()
	if err != nil {
		return nil, err
	}

	protoInstances := make([]*protoBalancer.InstanceInfo, 0, len(instances))
	for _, instance := range instances {
		protoInstances = append(protoInstances, &protoBalancer.InstanceInfo{
			InstanceId:    instance.ID,
			ChallengeType: instance.Type,
			Host:          instance.Host,
			PortNumber:    instance.Port,
			Status:        instance.Status,
			LastSeen:      instance.LastSeen.Unix(),
		})
	}

	return &protoBalancer.GetInstancesResponse{
		Instances: protoInstances,
		Count:     int32(len(protoInstances)),
	}, nil
}

func (s *BalancerService) CheckUserBlocked(ctx context.Context, req *protoBalancer.CheckUserBlockedRequest) (*protoBalancer.CheckUserBlockedResponse, error) {
	isBlocked := s.userBlockRepo.IsUserBlocked(req.UserId)
	var blockedUser *entity.BlockedUser
	if isBlocked {
		blockedUser, _ = s.userBlockRepo.GetBlockedUser(req.UserId)
	}
	if !isBlocked {
		return &protoBalancer.CheckUserBlockedResponse{
			IsBlocked: false,
		}, nil
	}

	return &protoBalancer.CheckUserBlockedResponse{
		IsBlocked:    true,
		Reason:       blockedUser.Reason,
		BlockedUntil: blockedUser.BlockedUntil.Unix(),
	}, nil
}

func (s *BalancerService) IsUserBlocked(userID string) bool {
	return s.userBlockRepo.IsUserBlocked(userID)
}

func (s *BalancerService) BlockUser(userID, reason string) error {
	return s.userBlockRepo.BlockUser(userID, reason)
}

func (s *BalancerService) BlockUserGRPC(ctx context.Context, req *protoBalancer.BlockUserRequest) (*protoBalancer.BlockUserResponse, error) {
	err := s.BlockUser(req.UserId, req.Reason)
	if err != nil {
		return &protoBalancer.BlockUserResponse{
			Status:  protoBalancer.BlockUserResponse_ERROR,
			Message: err.Error(),
		}, nil
	}

	return &protoBalancer.BlockUserResponse{
		Status:  protoBalancer.BlockUserResponse_SUCCESS,
		Message: "User blocked successfully",
	}, nil
}

func (s *BalancerService) StartCleanup() {
	ticker := time.NewTicker(time.Duration(s.config.CleanupInterval) * time.Second)
	go func() {
		for range ticker.C {
			s.userBlockRepo.CleanupExpiredBlocks()
		}
	}()
}

func (s *BalancerService) Stop() {
}
