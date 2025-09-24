package balancer

import (
	"context"
	"log"

	protoBalancer "captcha-service/gen/proto/proto/balancer"
	"captcha-service/internal/domain/entity"
	"captcha-service/internal/service"
)

type Handlers struct {
	protoBalancer.UnimplementedBalancerServiceServer
	balancerService *service.BalancerService
}

func NewHandlers(balancerService *service.BalancerService) *Handlers {
	return &Handlers{
		balancerService: balancerService,
	}
}

func (h *Handlers) RegisterInstance(stream protoBalancer.BalancerService_RegisterInstanceServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			log.Printf("Failed to receive register instance request: %v", err)
			return err
		}

		log.Printf("Received register instance request: %+v", req)

		// Конвертируем proto в entity
		entityReq := &entity.RegisterInstanceRequest{
			EventType:     req.EventType.String(),
			InstanceID:    req.InstanceId,
			ChallengeType: req.ChallengeType,
			Host:          req.Host,
			PortNumber:    req.PortNumber,
			Timestamp:     req.Timestamp,
		}

		// Регистрируем инстанс в service
		if err := h.balancerService.RegisterInstance(entityReq); err != nil {
			log.Printf("Failed to register instance: %v", err)
			resp := &protoBalancer.RegisterInstanceResponse{
				Status:  protoBalancer.RegisterInstanceResponse_ERROR,
				Message: "Failed to register instance",
			}
			stream.Send(resp)
			continue
		}

		resp := &protoBalancer.RegisterInstanceResponse{
			Status:  protoBalancer.RegisterInstanceResponse_SUCCESS,
			Message: "Instance registered successfully",
		}

		if err := stream.Send(resp); err != nil {
			log.Printf("Failed to send register instance response: %v", err)
			return err
		}
	}
}

func (h *Handlers) CheckUserBlocked(ctx context.Context, req *protoBalancer.CheckUserBlockedRequest) (*protoBalancer.CheckUserBlockedResponse, error) {
	log.Printf("Checking if user is blocked: %s", req.UserId)

	return &protoBalancer.CheckUserBlockedResponse{
		IsBlocked:    false,
		Reason:       "",
		BlockedUntil: 0,
	}, nil
}

func (h *Handlers) BlockUser(ctx context.Context, req *protoBalancer.BlockUserRequest) (*protoBalancer.BlockUserResponse, error) {
	log.Printf("Blocking user: %s for %d minutes", req.UserId, req.DurationMinutes)

	return &protoBalancer.BlockUserResponse{
		Status:  protoBalancer.BlockUserResponse_SUCCESS,
		Message: "User blocked successfully",
	}, nil
}

func (h *Handlers) GetInstances(ctx context.Context, req *protoBalancer.GetInstancesRequest) (*protoBalancer.GetInstancesResponse, error) {
	log.Printf("Getting instances")

	instances, err := h.balancerService.GetInstances()
	if err != nil {
		log.Printf("Failed to get instances: %v", err)
		return &protoBalancer.GetInstancesResponse{
			Instances: []*protoBalancer.InstanceInfo{},
			Count:     0,
		}, err
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
