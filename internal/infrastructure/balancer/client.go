package balancer

import (
	"context"
	"fmt"
	"log"
	"time"

	protoBalancer "captcha-service/gen/proto/proto/balancer"
	"captcha-service/internal/domain/entity"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	config         *entity.Config
	conn           *grpc.ClientConn
	balancerClient protoBalancer.BalancerServiceClient
	stream         protoBalancer.BalancerService_RegisterInstanceClient
	instanceID     string
	host           string
	port           int32
}

func NewClient(cfg *entity.Config) *Client {
	return &Client{
		config:     cfg,
		instanceID: generateInstanceID(),
		host:       cfg.Host,
		port:       8080, // Default port
	}
}

func (c *Client) Connect(ctx context.Context) error {
	balancerAddr := c.config.BalancerAddress
	if balancerAddr == "" {
		balancerAddr = fmt.Sprintf("%s:9090", c.config.Host)
	}
	log.Printf("Connecting to balancer at %s", balancerAddr)
	conn, err := grpc.Dial(balancerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	c.conn = conn
	c.balancerClient = protoBalancer.NewBalancerServiceClient(conn)

	stream, err := c.balancerClient.RegisterInstance(ctx)
	if err != nil {
		return err
	}
	c.stream = stream

	err = c.sendReadyEvent()
	if err != nil {
		return err
	}

	go c.keepAlive(ctx)

	log.Printf("Successfully connected to balancer")
	return nil
}

func (c *Client) sendReadyEvent() error {
	req := &protoBalancer.RegisterInstanceRequest{
		EventType:     protoBalancer.RegisterInstanceRequest_READY,
		InstanceId:    c.instanceID,
		ChallengeType: entity.ChallengeTypeSliderPuzzleReg,
		Host:          c.host,
		PortNumber:    c.port,
		Timestamp:     time.Now().Unix(),
	}

	return c.stream.Send(req)
}

func (c *Client) sendStoppedEvent() error {
	req := &protoBalancer.RegisterInstanceRequest{
		EventType:     protoBalancer.RegisterInstanceRequest_STOPPED,
		InstanceId:    c.instanceID,
		ChallengeType: entity.ChallengeTypeSliderPuzzleReg,
		Host:          c.host,
		PortNumber:    c.port,
		Timestamp:     time.Now().Unix(),
	}

	return c.stream.Send(req)
}

func (c *Client) keepAlive(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := c.sendReadyEvent(); err != nil {
				log.Printf("Failed to send keepalive: %v", err)
				return
			}
		}
	}
}

func (c *Client) Stop(ctx context.Context) error {
	log.Printf("Stopping balancer client")

	if c.stream != nil {
		c.sendStoppedEvent()
		c.stream.CloseSend()
	}

	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}

func generateInstanceID() string {
	return "captcha-instance-" + time.Now().Format("20060102150405")
}
