package persistence

import protoBalancer "captcha-service/gen/proto/proto/balancer"

type InstanceRepository interface {
	SaveInstance(instance *protoBalancer.InstanceInfo) error
	GetInstance(id string) (*protoBalancer.InstanceInfo, error)
	GetAllInstances() ([]*protoBalancer.InstanceInfo, error)
	RemoveInstance(id string) error
}

type UserBlockRepository interface {
	Save(block *UserBlockInfo) error
	Get(userID string) (*UserBlockInfo, error)
	Delete(userID string) error
}
