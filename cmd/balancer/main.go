package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	protoBalancer "captcha-service/gen/proto/proto/balancer"
	"captcha-service/internal/config"
	"captcha-service/internal/infrastructure/persistence"
	"captcha-service/internal/service"
	"captcha-service/internal/transport/grpc/balancer"
	httpTransport "captcha-service/internal/transport/http"
	"captcha-service/pkg/logger"

	"go.uber.org/zap"
	grpcLib "google.golang.org/grpc"
)

func main() {
	cfg, err := config.LoadBalancerConfig()
	if err != nil {
		log.Fatalf("Failed to load balancer config: %v", err)
	}

	logger.Init(cfg.LogLevel)
	defer logger.Sync()

	grpcPort := cfg.GRPCPort
	lis, err := net.Listen("tcp", "0.0.0.0:"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	instanceRepo := persistence.NewMemoryInstanceRepository()
	userBlockRepo := persistence.NewMemoryUserBlockRepository()

	entityConfig := &config.ServiceConfig{
		MaxAttempts:      cfg.MaxAttempts,
		BlockDurationMin: cfg.BlockDurationMin,
		CleanupInterval:  cfg.CleanupInterval,
		StaleThreshold:   cfg.StaleThreshold,
	}
	balancerService := service.NewBalancerService(instanceRepo, userBlockRepo, entityConfig)

	balancerService.StartCleanup()

	grpcHandlers := balancer.NewHandlers(balancerService.(*service.BalancerService))
	httpHandlers := httpTransport.NewBalancerHandlers(balancerService.(*service.BalancerService))

	grpcServer := grpcLib.NewServer()
	protoBalancer.RegisterBalancerServiceServer(grpcServer, grpcHandlers)

	httpServer := httpTransport.NewBalancerServer(httpHandlers, cfg.Port)

	logger.Info("Balancer server starting",
		zap.String("grpc_port", grpcPort),
		zap.String("http_port", cfg.Port),
		zap.String("log_level", cfg.LogLevel),
		zap.Duration("cleanup_interval", time.Duration(cfg.CleanupInterval)*time.Second),
		zap.Duration("stale_threshold", time.Duration(cfg.StaleThreshold)*time.Second))

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	go func() {
		if err := httpServer.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to serve HTTP: %v", err)
		}
	}()

	logger.Info("Balancer server started successfully",
		zap.String("grpc_port", grpcPort),
		zap.String("http_port", cfg.Port))

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	logger.Info("Shutting down balancer server...")

	balancerService.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Stop(ctx); err != nil {
		logger.Error("HTTP server shutdown error", zap.Error(err))
	}

	grpcServer.GracefulStop()
	logger.Info("Balancer server stopped")
}
