package balancer

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	protoBalancer "captcha-service/gen/proto/proto/balancer"
	"captcha-service/internal/config"
	"captcha-service/internal/domain/entity"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	config         *config.CaptchaConfig
	conn           *grpc.ClientConn
	balancerClient protoBalancer.BalancerServiceClient
	stream         protoBalancer.BalancerService_RegisterInstanceClient
	instanceID     string
	host           string
	port           int32
}

func NewClient(cfg *config.CaptchaConfig) *Client {
	// Используем порт из конфигурации, если указан
	port := int32(8080)
	if cfg.Port != "" {
		if p, err := strconv.Atoi(cfg.Port); err == nil {
			port = int32(p)
		}
	}

	return &Client{
		config:     cfg,
		instanceID: generateInstanceID(),
		host:       cfg.Host,
		port:       port,
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

func (c *Client) SetPort(port int32) {
	c.port = port
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
	return "captcha-instance-" + uuid.New().String()
}
