package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"captcha-service/internal/config"
	httpDelivery "captcha-service/internal/transport/http"
)

func main() {
	cfg, err := config.LoadBalancerProxyConfig()
	if err != nil {
		log.Fatalf("Failed to load balancer-proxy config: %v", err)
	}

	entityConfig := &config.ServiceConfig{
		MaxAttempts:      cfg.MaxAttempts,
		BlockDurationMin: cfg.BlockDurationMin,
		ComplexityMedium: cfg.MinOverlapPct,
	}
	proxy := httpDelivery.NewBalancerProxy(entityConfig)

	if err := proxy.ConnectToBalancer(cfg.BalancerAddress); err != nil {
		log.Fatalf("Failed to connect to balancer: %v", err)
	}

	go proxy.StartServiceDiscovery()

	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			proxy.CleanupSessions()
		}
	}()

	mux := httpDelivery.SetupBalancerProxyRoutes(proxy, cfg)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	go func() {
		log.Printf("Balancer proxy started on http://%s:%s", cfg.Host, cfg.Port)
		log.Printf("Open %s in your browser", cfg.DemoURL)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
