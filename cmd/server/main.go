package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
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
	grpc_gateway "captcha-service/internal/transport/grpc_gateway"
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

	// Используем порт из конфигурации, если задан
	var availablePort int
	if cfg.Port != "" {
		if p, err := strconv.Atoi(cfg.Port); err == nil {
			availablePort = p
			logger.Info("Using configured port", zap.Int("port", availablePort))
		} else {
			logger.Fatal("Invalid port configuration", zap.String("port", cfg.Port), zap.Error(err))
		}
	} else {
		// Ищем свободный порт только если не задан в конфигурации
		portFinder := port.NewPortFinder(int(cfg.MinPort), int(cfg.MaxPort))
		var err error
		availablePort, err = portFinder.FindAvailablePortWithRetry(3, 1*time.Second)
		if err != nil {
			logger.Fatal("Failed to find available port", zap.Error(err))
		}
		logger.Info("Found available port", zap.Int("port", availablePort))
	}

	balancerClient := balancer.NewClient(entityConfig)
	// Обновляем порт в клиенте балансера
	balancerClient.SetPort(int32(availablePort))

	ctx := context.Background()
	if err := balancerClient.Connect(ctx); err != nil {
		logger.Error("Failed to connect to balancer", zap.Error(err))
	}

	grpcHandlers := grpc.NewHandlers(captchaService)
	httpHandlers := http.NewHandlers(captchaService)

	gatewayServer := grpc_gateway.NewServer(grpcHandlers, httpHandlers, availablePort)

	go func() {
		if err := gatewayServer.Start(); err != nil {
			logger.Error("Gateway server error", zap.Error(err))
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
	if err := gatewayServer.Stop(shutdownCtx); err != nil {
		logger.Error("Failed to stop gateway server", zap.Error(err))
	}

	logger.Info("Server stopped")
}
