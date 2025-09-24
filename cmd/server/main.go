package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"captcha-service/internal/domain/entity"
	"captcha-service/internal/infrastructure/balancer"
	"captcha-service/internal/infrastructure/config"
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
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load captcha service config: %v", err)
	}

	if err := logger.Init(cfg.LogLevel); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Get().Sync()

	// Tracing disabled for now
	// if err := tracing.Init(serviceName, cfg.JaegerEndpoint); err != nil {
	// 	logger.Fatal("Failed to initialize tracer", zap.Error(err))
	// }

	// Use memory-optimized repository for better performance
	repo := persistence.NewMemoryOptimizedRepository(10000) // 10k max challenges

	templateEngine := template.NewTemplateEngineService("./templates")

	registry := service.NewGeneratorRegistry()
	registry.Register("slider-puzzle", service.NewSliderPuzzleGenerator(cfg, repo, templateEngine))

	entityConfig := &entity.Config{
		MaxAttempts:      cfg.MaxAttempts,
		BlockDurationMin: cfg.BlockDurationMin,
	}
	captchaService := service.NewCaptchaService(repo, registry, nil, entityConfig)

	portFinder := port.NewPortFinder(8080, 8090)
	availablePort, err := portFinder.FindAvailablePortWithRetry(3, 1*time.Second)
	if err != nil {
		logger.Fatal("Failed to find available port", zap.Error(err))
	}

	logger.Info("Found available port", zap.Int("port", availablePort))

	// Port is set to availablePort

	balancerClient := balancer.NewClient(cfg)

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
		zap.String("balancer_addr", "localhost:9090"))

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	logger.Info("Shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
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
