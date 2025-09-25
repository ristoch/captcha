package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"captcha-service/internal/config"
	"captcha-service/internal/domain/entity"
	"captcha-service/internal/infrastructure/balancer"
	"captcha-service/internal/infrastructure/persistence"
	"captcha-service/internal/infrastructure/port"
	"captcha-service/internal/infrastructure/template"
	"captcha-service/internal/service"
	"captcha-service/internal/transport/grpc"
	"captcha-service/internal/transport/http"
	"captcha-service/pkg/logger"

	"go.uber.org/zap"
)

const serviceName = "captcha-service"

func main() {
	cfg, err := config.LoadCaptchaServiceConfig()
	if err != nil {
		log.Fatalf("Failed to load captcha service config: %v", err)
	}

	if err := logger.Init(cfg.LogLevel); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Get().Sync()

	repo := persistence.NewMemoryOptimizedRepository(int(cfg.MaxChallenges))

	templateEngine := template.NewTemplateEngineService("./templates")

	entityConfig := cfg

	registry := service.NewGeneratorRegistry()
	registry.Register(entity.ChallengeTypeSliderPuzzle, service.NewSliderPuzzleGenerator(entityConfig, repo, templateEngine))

	captchaService := service.NewCaptchaService(repo, registry, nil, entityConfig)

	portFinder := port.NewPortFinder(int(cfg.MinPort), int(cfg.MaxPort))
	availablePort, err := portFinder.FindAvailablePortWithRetry(3, 1*time.Second)
	if err != nil {
		logger.Fatal("Failed to find available port", zap.Error(err))
	}

	logger.Info("Found available port", zap.Int("port", availablePort))

	balancerClient := balancer.NewClient(entityConfig)

	ctx := context.Background()
	if err := balancerClient.Connect(ctx); err != nil {
		logger.Error("Failed to connect to balancer", zap.Error(err))
	}

	grpcHandlers := grpc.NewHandlers(captchaService)
	grpcServer := grpc.NewServer(grpcHandlers, availablePort)

	httpHandlers := http.NewHandlers(captchaService)
	httpServer := http.NewServer(httpHandlers, cfg.Port)

	go func() {
		if err := grpcServer.Start(); err != nil {
			logger.Error("gRPC server error", zap.Error(err))
		}
	}()

	go func() {
		if err := httpServer.Start(); err != nil {
			logger.Error("HTTP server error", zap.Error(err))
		}
	}()

	logger.Info("Captcha service started",
		zap.Int("port", availablePort),
		zap.String("log_level", cfg.LogLevel),
		zap.String("balancer_addr", cfg.BalancerAddr))

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	logger.Info("Shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Duration(cfg.ShutdownTimeoutSec)*time.Second)
	defer shutdownCancel()

	if err := balancerClient.Stop(shutdownCtx); err != nil {
		logger.Error("Failed to stop balancer client", zap.Error(err))
	}
	if err := grpcServer.Stop(shutdownCtx); err != nil {
		logger.Error("Failed to stop gRPC server", zap.Error(err))
	}

	if err := httpServer.Stop(shutdownCtx); err != nil {
		logger.Error("Failed to stop HTTP server", zap.Error(err))
	}

	logger.Info("Server stopped")
}
