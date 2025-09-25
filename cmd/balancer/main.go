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

	grpcPort := "9090"
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

	handlers := balancer.NewHandlers(balancerService.(*service.BalancerService))

	grpcServer := grpcLib.NewServer()
	protoBalancer.RegisterBalancerServiceServer(grpcServer, handlers)

	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"balancer"}`))
	})

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: httpMux,
	}

	logger.Info("Balancer server starting",
		zap.String("grpc_port", grpcPort),
		zap.String("http_port", "8080"),
		zap.String("log_level", cfg.LogLevel),
		zap.Duration("cleanup_interval", time.Duration(cfg.CleanupInterval)*time.Second),
		zap.Duration("stale_threshold", time.Duration(cfg.StaleThreshold)*time.Second))

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to serve HTTP: %v", err)
		}
	}()

	logger.Info("Balancer server started successfully",
		zap.String("grpc_port", grpcPort),
		zap.String("http_port", "8080"))

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	logger.Info("Shutting down balancer server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("HTTP server shutdown error", zap.Error(err))
	}

	grpcServer.GracefulStop()
	logger.Info("Balancer server stopped")
}
